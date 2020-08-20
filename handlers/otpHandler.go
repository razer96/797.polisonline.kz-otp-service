package handlers

import (
	"encoding/json"
	"insurance-otp-service/database"
	"insurance-otp-service/elastic"
	"insurance-otp-service/helpers"
	"insurance-otp-service/logger"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/xlzd/gotp"
)

var (
	log       = logger.GetLogger()
	otpSecret = os.Getenv("OTP_SECRET")
	db        = database.GetDB()
	es        = elastic.GetESClient()
)

// OtpGetHanler godoc
// @Summary Send a OTP
// @Description Send a OTP to client by sms
// @ID get-otp-to-phon-number
// @Tags otp
// @Accept json
// @Produce json
// @Param phone path string true "Phone number to send OTP"
// @Success 200 {object} GetOtpRespObj
// @Failure 500 {object} helpers.ErrorResponse
// @Failure 422 {object} helpers.ErrorResponse
// @Router /otp [get]
func OtpGetHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("/otp Get method called")

	newUuid, err := uuid.NewUUID()
	if err != nil {
		log.Error("Error getiin uuid", err)
		helpers.ErrorJsonResponse(w, http.StatusInternalServerError, "Error getting uuid")
		return
	}
	randomSecret := gotp.RandomSecret(16)

	phoneNumber := r.URL.Query().Get("phone")

	totp := gotp.NewTOTP(randomSecret, 4, 60, nil)
	otp := totp.Now()

	otpToDb := &database.OTP{
		ID:          newUuid.String(),
		Secret:      randomSecret,
		PhoneNumber: phoneNumber,
		Status:      "1",
		SendAt:      time.Now().Unix(),
	}
	err = db.DB.Save(&otpToDb).Error
	if err != nil {
		log.Error("Error to save record to db", err)
		helpers.ErrorJsonResponse(w, http.StatusInternalServerError, "Error saving record into db")
		return
	}

	smsText := otp + " - код подтверждения. НИКОМУ НЕ ГОВОРИТЕ КОД! ДАЖЕ СОТРУДНИКУ БАНКА!"

	status, err := helpers.SendSms(smsText, phoneNumber)
	if err != nil && status != http.StatusOK {
		log.Error("Error sending sms", err)
		helpers.ErrorJsonResponse(w, status, err.Error())
		return
	}
	err = es.Insert("insurance_otp", newUuid.String(), map[string]string{"test": "test"})
	helpers.JsonResponse(w, http.StatusOK, &GetOtpRespObj{Key: otpToDb.ID})
	if err != nil {
		log.Error("Error insetring into ES", err)
		helpers.ErrorJsonResponse(w, http.StatusInternalServerError, "Error inserting log into ES")
		return
	}
	log.Info("/otp Get method finished")
}

// OtpGetHanler godoc
// @Summary Send a OTP
// @Description Send a OTP to client by sms
// @ID post-otp-to-phon-number
// @Tags otp
// @Accept json
// @Produce json
// @Param validate_otp_req_body body ValidateOtpReqBody true "Body should contain phone number, key, and otp"
// @Success 204
// @Failure 500 {object} helpers.ErrorResponse
// @Failure 400 {object} helpers.ErrorResponse "This status is returned if wrong otp has been sent"
// @Failure 403 {object} helpers.ErrorResponse "This status is returned if key status is no more valid"
// @Failure 404 {object} helpers.ErrorResponse "This status is returned if otp sent in more than 60 sec"
// @Failure 410 {object} helpers.ErrorResponse "This status is returned if otp reached 3 attemps of validation"
// @Router /otp [post]
func OtpPostHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("/otp Post method called")

	requestTimestamp := time.Now().Unix()

	bodyContent := &ValidateOtpReqBody{}

	err := json.NewDecoder(r.Body).Decode(bodyContent)
	if err != nil {
		log.Error("Error parsing request body", err)
		helpers.ErrorJsonResponse(w, http.StatusInternalServerError, "Error parsing request body")
		return
	}

	otpInfo := &database.OTP{}

	err = db.DB.Raw("SELECT * FROM otps WHERE id = ?", bodyContent.Key).Scan(otpInfo).Error
	if err != nil {
		log.Error("Error getting otp info from db", err)
		helpers.ErrorJsonResponse(w, http.StatusInternalServerError, "Error getting otp info from db")
		return
	}
	if otpInfo.Status == "1" {
		if requestTimestamp-otpInfo.SendAt <= 60 {
			totp := gotp.NewTOTP(otpInfo.Secret, 4, 60, nil)

			verification := totp.Verify(bodyContent.Otp, int(otpInfo.SendAt))
			if verification {
				log.Info("Valid otp")
				otpInfo.Status = "0"
				otpInfo.Attempts = otpInfo.Attempts + 1
				err := db.DB.Save(otpInfo).Error
				if err != nil {
					log.Error("Unable to update record", err)
					helpers.ErrorJsonResponse(w, http.StatusInternalServerError, "Unable to update record")
					return
				}
				helpers.JsonResponse(w, http.StatusNoContent, nil)
				return
			} else {
				log.Info("Invalid otp")
				otpInfo.Attempts = otpInfo.Attempts + 1
				if otpInfo.Attempts == 3 {
					otpInfo.Status = "0"
					err := db.DB.Save(otpInfo).Error
					if err != nil {
						log.Error("Unable to update record", err)
						helpers.ErrorJsonResponse(w, http.StatusInternalServerError, "Unable to update record")
						return
					}
					helpers.ErrorJsonResponse(w, http.StatusGone, "You reached maximum attempts")
					return
				} else {
					helpers.ErrorJsonResponse(w, http.StatusBadRequest, "Code is invalid")
					return
				}
			}
		} else {
			log.Info("OTP is expired")
			otpInfo.Status = "0"
			err := db.DB.Save(otpInfo).Error
			if err != nil {
				log.Error("Unable to update record", err)
				helpers.ErrorJsonResponse(w, http.StatusInternalServerError, "Unable to update record")
				return
			}
			helpers.ErrorJsonResponse(w, http.StatusNotFound, "OTP is expired")
			return
		}
	} else {
		log.Info("Otp is no longer valid")
		helpers.ErrorJsonResponse(w, http.StatusForbidden, "OTP is no longer valid")
		return
	}
}
