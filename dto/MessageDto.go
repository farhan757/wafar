package dto

type ReqSend struct {
	No   string `json:"no"`
	Text string `json:"text"`
}

type ReqSendPict struct {
	No   string `form:"no" binding:"required"`
	Pict []byte `form:"pict" binding:"required"`
	Text string `form:"text" binding:"required"`
}

type ReqSendPdf struct {
	No   string `form:"no" binding:"required"`
	Pdf  []byte `form:"pdf" binding:"required"`
	Text string `form:"text" binding:"required"`
}
