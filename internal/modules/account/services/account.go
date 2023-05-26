package services

import (
	"errors"
	"fback/internal/modules/account/consts"
	"fback/internal/modules/account/dto"
	"fback/internal/modules/account/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/zedisdog/ty/application"
	"github.com/zedisdog/ty/auth"
	"github.com/zedisdog/ty/errx"
	"github.com/zedisdog/ty/log"
	"github.com/zedisdog/ty/sdk/wechat/mini"
	"github.com/zedisdog/ty/strings"
	"gorm.io/gorm"
	"time"
)

func NewAccount(app application.IApplication) *Account {
	client := mini.NewClient(
		app.Config().GetString("modules.account.wechat.mini.appid"),
		app.Config().GetString("modules.account.wechat.mini.secret"),
	)
	return &Account{
		db:           app.Database("default").(*gorm.DB),
		logger:       app.Logger(),
		client:       client,
		codeSessions: make(map[string]uint64),
	}
}

type Account struct {
	db           *gorm.DB
	logger       log.ILog
	client       *mini.Client
	codeSessions map[string]uint64
}

func (a *Account) FindByID(id any) (account models.Account, err error) {
	err = a.db.First(&account, id).Error
	return
}

func (a *Account) UserExists(id uint64) (exists bool, err error) {
	_, err = a.FindByID(id)
	if err != nil {
		return
	}

	exists = true
	return
}

func (a *Account) LoginByMiniCode(code string) (token string, err error) {
	r, err := a.client.Code2Session(code)
	if err != nil {
		return
	}

	if r.ErrCode != 0 {
		err = errx.New(fmt.Sprintf("code2session failed<%d>: %s", r.ErrCode, r.ErrMsg))
		return
	}

	var account models.Account
	err = a.db.Exec(
		"SELECT accounts.* FROM accounts RIGHT JOIN wechat ON wechat.account_id = accounts.id WHERE wechat.open_id = ? AND wechat.type = ?",
		r.OpenID,
		consts.WechatMini,
	).First(&account).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}

	if err == gorm.ErrRecordNotFound {
		account, err = a.CreateByOpenID(r.OpenID, consts.WechatMini)
		if err != nil {
			return
		}
	}

	token, err = auth.NewJwtTokenBuilder().WithClaims(map[string]interface{}{
		"open_id":       r.OpenID,
		"type":          consts.WechatMini,
		auth.JwtSubject: fmt.Sprintf("%d", account.ID),
	}).BuildToken()

	return
}

func (a *Account) CreateByOpenID(openID string, t consts.WechatType) (account models.Account, err error) {
	err = a.db.Transaction(func(tx *gorm.DB) (err error) {
		err = tx.Create(&account).Error
		if err != nil {
			return
		}

		wechat := models.Wechat{
			AccountID: account.ID,
			OpenID:    openID,
			Type:      t,
			Nickname:  "微信用户",
		}
		err = tx.Create(&wechat).Error
		return
	})
	return
}

func (a *Account) UpdateInfo(nickname string, avatarUrl string, mobile string, openID string) (err error) {
	var wechat models.Wechat
	err = a.db.First(&wechat, "open_id = ?", openID).Error
	if err != nil {
		return
	}

	wechat.Nickname = nickname
	wechat.AvatarUrl = avatarUrl
	wechat.Mobile = mobile

	return a.db.Updates(&wechat).Error
}

func (a *Account) GetWechatInfoByOpenID(openID string, t consts.WechatType) (info dto.WechatInfo, err error) {
	var wechat models.Wechat
	err = a.db.Where("open_id = ? AND type = ?", openID, t).First(&wechat).Error
	if err != nil {
		return
	}

	err = copier.Copy(&info, wechat)
	return
}

func (a *Account) UpdateWechatInfo(info dto.WechatInfo, openID string) (err error) {
	var wechat models.Wechat
	err = a.db.Where("open_id = ? AND type = ?", openID, consts.WechatMini).First(&wechat).Error
	if err != nil {
		return
	}

	err = copier.Copy(&wechat, info)
	if err != nil {
		return
	}

	err = a.db.Select("nickname", "avatar_url").Updates(&wechat).Error
	return
}

func (a *Account) ServeLoginByCode(ctx *gin.Context) {
	var code string
	for {
		code = strings.RandNumeric(4)
		if _, ok := a.codeSessions[code]; !ok {
			break
		}
	}

	a.codeSessions[code] = 0
	defer delete(a.codeSessions, code)

	ctx.SSEvent("message", dto.SEventCode{
		Code: code,
	})
	ctx.Writer.Flush()

	notify := ctx.Writer.CloseNotify()

	for {
		select {
		case <-notify:
			return
		default:
			if a.codeSessions[code] != 0 {
				account, err := a.FindByID(a.codeSessions[code])
				if err != nil {
					ctx.SSEvent("error", gin.H{
						"message": "登录错误,请重试",
					})
					ctx.Writer.Flush()
					return
				}

				var wechat models.Wechat
				err = a.db.Where("account_id = ?", account.ID).First(&wechat).Error
				if err != nil {
					ctx.SSEvent("error", gin.H{
						"message": "登录错误,请重试",
					})
					ctx.Writer.Flush()
					return
				}

				var token string
				token, err = auth.NewJwtTokenBuilder().WithClaims(map[string]interface{}{
					"open_id":       wechat.OpenID,
					"type":          consts.WechatMini,
					auth.JwtSubject: fmt.Sprintf("%d", account.ID),
				}).BuildToken()
				if err != nil {
					ctx.SSEvent("error", gin.H{
						"message": "登录错误,请重试",
					})
					return
				}

				ctx.SSEvent("message", dto.Token{
					Token: token,
				})
				ctx.Writer.Flush()
				return
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func (a *Account) LoginByCode(code string, openID string) (result string) {
	if _, ok := a.codeSessions[code]; !ok {
		return
	}

	var (
		wechat  models.Wechat
		account models.Account
		err     error
	)
	err = a.db.Where("open_id = ? AND type = ?", openID, consts.WechatOfficial).First(&wechat).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	} else if err != nil {
		account, err = a.CreateByOpenID(openID, consts.WechatOfficial)
		if err != nil {
			return
		}
	} else {
		account, err = a.FindByID(wechat.AccountID)
		if err != nil {
			return
		}
	}

	a.codeSessions[code] = account.ID
	result = "登录成功"
	return
}
