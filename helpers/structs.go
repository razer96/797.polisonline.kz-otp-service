package helpers

type TokenResponseStruct struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

type UCSRequest struct {
	Message     string `json:"Message"`
	ChannelName string `json:"ChannelName"`
	Contact     string `json:"Contact"`
	Priority    int    `json:"Priotiry"`
	Source      string `json:"Source"`
	UserID      string `json:"UserID"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
