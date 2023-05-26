package models

import (
	"fback/internal/modules/account/consts"
	"github.com/zedisdog/ty/database"
)

type Wechat struct {
	database.CommonField
	AccountID uint64            //账号id
	Type      consts.WechatType //账号类型
	OpenID    string            //openid
	AvatarUrl string            //头像url
	Nickname  string            //昵称
	Mobile    string            //手机号
}

func (w Wechat) TableName() string {
	return "wechat"
}
