package dto

import "encoding/xml"

type WechatMessageType string

const (
	TEXT WechatMessageType = "text"
)

type CDATA struct {
	Text string `xml:",cdata"`
}

type WechatMessageCommonFields struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATA
	FromUserName CDATA
	CreateTime   int64
	MsgType      CDATA
}

type WechatMessage struct {
	WechatMessageCommonFields
	Content CDATA
	MsgId   string
}

type WechatPassiveReplyMessage struct {
	WechatMessageCommonFields
	Content CDATA
}
