package main

import (
	"log"
	"net/http"
	"gopkg.in/redis.v4"
	"database/sql"
	"strings"
	"encoding/base64"
)


type Env struct {
	db *sql.DB
	client *redis.Client
}




func main() {

	db, err := NewDB("user=postgres password=root  dbname=hemc sslmode=disable")
	if err != nil {
		log.Panic(err)
	}

	client, err := NewCache("localhost:6379","",0)
	if err != nil {
		log.Panic(err)
	}
	env := &Env{db: db, client:client}



	http.HandleFunc("/inbound/sms", env.BasicAuth(env.InboundSms))

	http.HandleFunc("/outbound/sms", env.BasicAuth(env.OutboundSms))

	http.ListenAndServe(":8080", nil)


}


//middleware for authentication
func (env *Env)  BasicAuth( h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			http.Error(w, "Not authorized", 403)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, err.Error(), 403)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "Not authorized", 403)
			return
		}


		if !userExists(env.db, pair[0],pair[1])  {
			http.Error(w, "Not authorized", 403)
			return
		}

		h.ServeHTTP(w, r)
	}
}