package main

import (
	"insurance-otp-service/database"
	"insurance-otp-service/handlers"
	"insurance-otp-service/logger"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// @title Insurance OTP service
// @version 1.0
// @description Insurance OTP service is used for sending and varifying OTPs

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /

var (
	port string
	log  = logger.GetLogger()
	db   = database.GetDB()
)

func init() {
	port = ":" + os.Getenv("NOMAD_HOST_PORT_http")
	if port == "" {
		port = ":3001"
	}

	db.DB.Exec(database.OTPStructSchema)
}

func main() {
	log.Info("Service starting...")

	r := mux.NewRouter()
	r.HandleFunc("/otp", handlers.OtpGetHandler).Methods("GET")
	r.HandleFunc("/otp", handlers.OtpPostHandler).Methods("POST")
	http.Handle("/", r)

	defer db.Close()

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Error("Service didn't started!", err)
		os.Exit(1)
	}
}
