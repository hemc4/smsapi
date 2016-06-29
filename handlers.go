package main

import (
	"encoding/json"
	_"fmt"
	"net/http"
	"strings"
	_"log"
)



func (env *Env) InboundSms(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	to := strings.TrimSpace(r.FormValue("to"))
	from := strings.TrimSpace(r.FormValue("from"))
	text:= strings.TrimSpace(r.FormValue("text"))

	//validate the form data
	validateError :=ValidateFormData(from, to, text)
	if validateError !=""{
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(`{"message": "", "error":"`+validateError+`"}`); err != nil {
			panic(err)
		}
		return
	}


	// check if the to number exists for the authorized user
	if numberExists(env.db,to){

		if text=="STOP" || text=="STOP\n" || text=="STOP\r" || text=="STOP\r\n" {
			//save to redis
			if cacheSms(env.client,from,to) {
				successMessage:=`{"message": "inbound sms ok", "error":""}`
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(successMessage); err != nil {
					panic(err)
				}
				return
			}


		}
	}else{
		errorMessage :=`{"message": "","error": "to parameter not found"}`
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(errorMessage); err != nil {
			panic(err)
		}
	}




}

func (env *Env) OutboundSms(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	to := strings.TrimSpace(r.FormValue("to"))
	from := strings.TrimSpace(r.FormValue("from"))
	text := strings.TrimSpace(r.FormValue("text"))


	//validate the formdata
	validateError :=ValidateFormData(from, to, text)
	if validateError !=""{
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(`{"message": "", "error":"`+validateError+`"}`); err != nil {
			panic(err)
		}
		return
	}
	//check the redis cache
	if cacheExists(env.client, from, to ){
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(`{"message": "", "error":"sms from `+from+` to `+to+` blocked by STOP request"}`); err != nil {
			panic(err)
		}
		return
	}

	//check limit

	if limitExceed(from){
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(`{"message": "", "error":"limit reached for from  `+from+` "}`); err != nil {
			panic(err)
		}
		return
	}



	// check if the to number exists for the authorized user
	if !numberExists(env.db, from) {

		errorMessage :=`{"message": "","error": "from parameter not found"}`
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(errorMessage); err != nil {
			panic(err)
		}
		return
	}


	successMessage :=`{"message": "outbound sms ok","error": ""}`
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(successMessage); err != nil {
		panic(err)
	}
}
