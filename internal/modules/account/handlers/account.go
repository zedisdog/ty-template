package handlers

import (
	"fback/internal/modules/account/consts"
	"fback/internal/modules/account/dto"
	"fback/internal/modules/account/services"
	"github.com/gin-gonic/gin"
	"github.com/zedisdog/ty/sdk/net/http"
	"github.com/zedisdog/ty/sdk/net/http/response"
	"time"
)

func NewAccount(accountSvc *services.Account) *Account {
	return &Account{
		accountSvc: accountSvc,
	}
}

type Account struct {
	accountSvc *services.Account
}

func (a Account) LoginByMiniCode(ctx *gin.Context) {
	var req dto.LoginByMiniCodeRequest
	if err := http.ValidateJSON(ctx, &req); err != nil {
		response.Error(ctx, err)
		return
	}

	token, err := a.accountSvc.LoginByMiniCode(req.Code)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, gin.H{
		"token": token,
	})
}

func (a Account) Self(ctx *gin.Context) {
	info, err := a.accountSvc.GetWechatInfoByOpenID(ctx.MustGet("open_id").(string), consts.WechatOfficial)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, info)
}

func (a Account) UpdateWechatMiniInfo(ctx *gin.Context) {
	var req dto.WechatInfo
	if err := http.ValidateJSON(ctx, &req); err != nil {
		response.Error(ctx, err)
		return
	}

	err := a.accountSvc.UpdateWechatInfo(req, ctx.MustGet("open_id").(string))
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

func (a Account) SyncMiniInfo(ctx *gin.Context) {
	var req dto.WechatInfo
	if err := http.ValidateJSON(ctx, &req); err != nil {
		response.Error(ctx, err)
		return
	}

	err := a.accountSvc.UpdateInfo(
		req.Nickname,
		req.AvatarUrl,
		req.Mobile,
		ctx.MustGet("open_id").(string),
	)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx)
}

func (a Account) LoginByCode(ctx *gin.Context) {
	a.accountSvc.ServeLoginByCode(ctx)
}

func (a Account) ReceiveWechatMessage(ctx *gin.Context) {
	if str := ctx.Query("echostr"); str != "" {
		ctx.Data(200, "text/plain", []byte(str))
		return
	}
	var req dto.WechatMessage
	err := ctx.ShouldBindXML(&req)
	if err != nil {
		return
	}

	result := a.accountSvc.LoginByCode(req.Content.Text, req.FromUserName.Text)
	if result != "" {
		resp := dto.WechatPassiveReplyMessage{
			WechatMessageCommonFields: dto.WechatMessageCommonFields{
				ToUserName:   dto.CDATA{Text: req.FromUserName.Text},
				FromUserName: dto.CDATA{Text: req.ToUserName.Text},
				CreateTime:   time.Now().Unix(),
				MsgType:      dto.CDATA{Text: string(dto.TEXT)},
			},
			Content: dto.CDATA{Text: result},
		}
		ctx.XML(200, resp)
	}

	return
}
