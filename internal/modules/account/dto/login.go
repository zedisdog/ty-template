package dto

type LoginByMiniCodeRequest struct {
	Code string `binding:"required"`
}

type SEventCode struct {
	Code string `json:"code"`
}

type Token struct {
	Token string `json:"token"`
}
