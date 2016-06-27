package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"encoding/base64"
	_"log"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to sms api !\n")
}


func InboundSms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	to := strings.TrimSpace(r.FormValue("to"))
	from := strings.TrimSpace(r.FormValue("from"))
	text:= strings.TrimSpace(r.FormValue("text"))

	//validate the form data

	//

	// check if the to number exists for the authorized user
	if numberExists(to){

		if text=="STOP" || text=="STOP\n" || text=="STOP\r" || text=="STOP\r\n" {
			//save to redis
			if cacheSms(from,to) {
				successMessage:=`{"message": "inbound sms ok", "error":""}`
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(successMessage); err != nil {
					panic(err)
				}
			}


		}
	}else{
		errorMessage :=`{"message": "","error": "to parameter not found"}`
		w.WriteHeader(401)
		if err := json.NewEncoder(w).Encode(errorMessage); err != nil {
			panic(err)
		}
	}




}

func OutboundSms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	to := strings.TrimSpace(r.FormValue("to"))
	from := strings.TrimSpace(r.FormValue("from"))
	text := strings.TrimSpace(r.FormValue("text"))


	//validate the formdata

	//check the redis cache


	//check limit



	// check if the to number exists for the authorized user
	if !numberExists(from) {

		errorMessage :=`{"message": "","error": "from parameter not found"}`
		w.WriteHeader(401)
		if err := json.NewEncoder(w).Encode(errorMessage); err != nil {
			panic(err)
		}
	}
}

func basicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}


		if !userExists(pair[0],pair[1])  {
			http.Error(w, "Not authorized", 401)
			return
		}

		h.ServeHTTP(w, r)
	}
}