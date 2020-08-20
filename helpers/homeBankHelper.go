package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"insurance-otp-service/logger"
	"net/http"
	"net/url"
	"os"
)

var (
	homeBankOAuth = os.Getenv("HOMEBANK_OAUTH")
	grantType     = os.Getenv("OAUTH_GRANT_TYPE")
	clientID      = os.Getenv("OAUTH_CLIENT_ID")
	clientSecret  = os.Getenv("OAUTH_CLIENT_SECRET")
	ucsService    = os.Getenv("UCS_SERVICE")
	log           = logger.GetLogger()
)

func getToken() (*TokenResponseStruct, error) {
	formData := url.Values{
		"grant_type":    []string{grantType},
		"scope":         []string{"sms_send"},
		"client_id":     []string{clientID},
		"client_secret": []string{clientSecret},
	}
	resp, err := http.PostForm(homeBankOAuth, formData)
	if err != nil {
		log.Error("Request to "+homeBankOAuth+" error", err)
		return nil, err
	}
	tokenResp := &TokenResponseStruct{}
	err = json.NewDecoder(resp.Body).Decode(tokenResp)
	if err != nil {
		log.Error("Parsing response error", err)
		return nil, err
	}
	return tokenResp, nil
}

func SendSms(text string, phoneNumber string) (int, error) {
	token, err := getToken()
	if err != nil {
		log.Error("Get homebank token error", err)
		return http.StatusInternalServerError, err
	}
	ucsReq := &UCSRequest{
		Message:     text,
		ChannelName: "SMS",
		Contact:     phoneNumber,
		Priority:    0,
		Source:      "7111",
		UserID:      "",
	}
	jsonUCSReq, err := json.Marshal(ucsReq)
	if err != nil {
		log.Error("Parsing ucsRequest struct to json", err)
		return http.StatusInternalServerError, err
	}
	httpClient := http.Client{}
	request, err := http.NewRequest("POST", ucsService, bytes.NewBuffer(jsonUCSReq))
	request.Header.Set("Authorization", "Bearer "+token.AccessToken)
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Error("Creating new requset to UCS error", err)
		return http.StatusInternalServerError, err
	}
	resp, err := httpClient.Do(request)
	if err != nil {
		log.Error("Error requset to UCS", err)
		return http.StatusInternalServerError, err
	}
	var ucsReqBody map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&ucsReqBody)
	if err != nil {
		log.Error("Decoding error", err)
		return http.StatusInternalServerError, err
	}

	if ucsReqBody["Error"].(string) != "" {
		log.Error("UCS status error", ucsReqBody["Error"].(string))
		if ucsReqBody["Status"].(string) == "-2" {
			return http.StatusBadRequest, errors.New("Invalid phone number")
		}
		return http.StatusServiceUnavailable, errors.New(ucsReqBody["Error"].(string))
	}
	return http.StatusOK, nil
}
