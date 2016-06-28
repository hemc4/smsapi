package main

import (
	"log"
	"net/http"

)

func main() {

	InitDb("user=postgres password=root  dbname=hemc sslmode=disable")
	InitRedis("localhost:6379","",0)
	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
