package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"gopkg.in/redis.v3"
	"log"
	"time"
)

var db *sql.DB
var client *redis.Client

var userId int

func init() {

	var err error

	db, err = sql.Open("postgres", "user=postgres password=root  dbname=hemc sslmode=disable")

	if err != nil {
		log.Fatalf("Error on initializing database connection: %s", err.Error())
	}

	//defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error: Could not establish a connection with the database : %s ", err.Error())
	}

	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err = client.Ping().Result()
	//fmt.Println(pong, err)
	if err != nil {
		log.Fatalf("Error: Could not establish a connection with the redis db : %s ", err.Error())
	}

}

func userExists(username, auth_id string) bool {

	var err error

	err = db.QueryRow("select id from account where username=$1 and auth_id=$2", username, auth_id).Scan(&userId)

	if userId != 0 {
		return true
	}

	if err != nil {
		return false
		//log.Fatalf("Error: Could not establish a connection with the database : %s ", err.Error())
	}

	return false
}

func numberExists(number string) bool {

	//fmt.Println("userid :", userId)
	//fmt.Println("number : ", number)

	var id int
	var err error

	err = db.QueryRow("select id from phone_number where number=$1 and account_id=$2", number, userId).Scan(&id)

	if id != 0 {
		return true
	}

	if err != nil {
		return false
	}

	return false
}

func cacheSms(from, to string) bool {
	//save to redis
	//fmt.Println("saving data to redis")
	err := client.Set(from, to, 4*60*60*60*time.Second).Err()
	if err != nil {
		panic(err)
	}

	//val, err := client.Get(from).Result()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("key : ", val)

	return true
}

func validateFormData(from, to, text string) {
	var errorMessage string

	if len(from < 6) || len(from > 16) {
		return errorMessage
	}

	if len(to < 6) || len(to > 16) {

	}
	return errorMessage

}
