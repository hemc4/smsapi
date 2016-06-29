package main

import (
	"encoding/json"
	_ "fmt"
	_ "log"
	"net/http"
	"strings"


)

func (env *Env) InboundSms(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")


	to := strings.TrimSpace(r.FormValue("to"))
	from := strings.TrimSpace(r.FormValue("from"))
	text := strings.TrimSpace(r.FormValue("text"))

	//fmt.Println(from)

	//validate the form data
	validateError := ValidateFormData(from, to, text)
	if validateError != "" {
		w.WriteHeader(http.StatusOK)
		out:=jsonOutput{Message:"",Error:validateError}
		if err := json.NewEncoder(w).Encode(out); err != nil {
			panic(err)
		}
		return
	}

	//fmt.Println(text)
	// check if the to number exists for the authorized user
	if !env.db.NumberExists(to) {
		out:=jsonOutput{Message:"",Error:"to parameter not found"}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(out); err != nil {
			panic(err)

		}
		return

	}

	//save to redis
	if text == `STOP` || text == `STOP\n`|| text == `STOP\r` || text == `STOP\r\n` {
		//save to redis
		if !env.client.CacheSms(from, to) {
			out:=jsonOutput{Message:"",Error:"unknown failure"}
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(out); err != nil {
				panic(err)
			}
		}

	}

	out:=jsonOutput{Message:"inbound sms ok",Error:""}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(out); err != nil {
		panic(err)
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

	//fmt.Println(from)
	//validate the formdata
	validateError := ValidateFormData(from, to, text)
	if validateError != "" {
		w.WriteHeader(http.StatusOK)
		out:=jsonOutput{Message:"",Error:validateError}
		if err := json.NewEncoder(w).Encode(out); err != nil {
			panic(err)
		}
		return
	}

	//check the redis cache
	if env.client.CacheExists(from, to) {
		w.WriteHeader(http.StatusOK)
		out:=jsonOutput{Message:"",Error:`sms from ` + from + ` to ` + to + ` blocked by STOP request`}
		if err := json.NewEncoder(w).Encode(out); err != nil {
			panic(err)
		}
		return
	}

	//check limit
	if limitExceed(from) {
		w.WriteHeader(http.StatusOK)
		out:=jsonOutput{Message:"",Error:`limit reached for from  ` + from}
		if err := json.NewEncoder(w).Encode(out); err != nil {
			panic(err)
		}
		return
	}

	// check if the to number exists for the authorized user
	if !env.db.NumberExists(from) {
		out:=jsonOutput{Message:"",Error:"from parameter not found"}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(out); err != nil {
			panic(err)
		}
		return
	}

	out:=jsonOutput{Message:"outbound sms ok",Error:""}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(out); err != nil {
		panic(err)
	}
}
