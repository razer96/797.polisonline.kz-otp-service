package handlers

type SmsStruct struct {
	Message     string `json:"Message"`
	ChannelName string `json:"ChannelName"`
	Contact     string `json:"Contact"`
	Priority    int    `json:"Priority"`
	Source      string `json:"Source"`
	UserID      string `json:"UserID"`
}

type ValidateOtpReqBody struct {
	Key   string `json:"key"`
	Phone string `json:"phone"`
	Otp   string `json:"otp"`
}

type GetOtpRespObj struct {
	Key string `json:"key"`
}
