package main

import (
	"database/sql"
	_"fmt"
	_ "github.com/lib/pq"
	timerate "golang.org/x/time/rate"
	"gopkg.in/go-redis/rate.v4"
	"gopkg.in/redis.v4"

	_ "strconv"
	"time"
)


var userId int

func NewDB(dataSourceName string) (*sql.DB, error) {
	var err error

	db, err := sql.Open("postgres", dataSourceName)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func NewCache(address, password string, dbcount int) (*redis.Client, error) {
	var err error
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password, // no password set
		DB:       dbcount,  // use default DB
	})

	_, err = client.Ping().Result()

	if err != nil {
		return nil, err
	}

	return client, nil
}

func userExists(db *sql.DB,username, auth_id string) bool {

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

func numberExists(db *sql.DB,number string) bool {

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

func cacheSms(client *redis.Client, from, to string) bool {
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

func cacheExists(client *redis.Client, from, to string) bool {

	val, err := client.Get(from).Result()
	if err != nil {
		//log.Fatalf("Pair not found %s", err.Error())
	}
	//fmt.Println("key : ", val)
	if val == to {
		return true
	}

	return false
}

func ValidateFormData(from, to, text string) string{

	//fmt.Println("from : %q, to = %q , text = %q", from, to, text)
	var errorMessage string

	if len(from) < 6 || len(from) > 16 {
		errorMessage = "from is invalid"
	} else if len(to) < 6 || len(to) > 16 {
		errorMessage = "to is invalid"

	} else if len(text) < 1 || len(text) > 120 {
		errorMessage = "text is invalid"

	}

	if len(from) == 0 {
		errorMessage = "from is missing"
	} else if len(to) == 0 {
		errorMessage = "to is missing"
	} else if len(text) == 0 {
		errorMessage = "text is missing"
	}

	return errorMessage

}

func limitExceed(from string) bool {

	fromID := "from-" + from
	limit := int64(50)
	duration := time.Duration(24) * time.Hour

	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": "localhost:6379",
		},
	})
	fallbackLimiter := timerate.NewLimiter(timerate.Every(time.Second), 100)
	rateLimiter := rate.NewLimiter(ring, fallbackLimiter)

	_, _, allowed := rateLimiter.Allow(fromID, limit, duration)

	if !allowed {
		//fmt.Println("limit exceed")
		return true
	}

	//fmt.Println("Rate limit remaining: ", strconv.FormatInt(limit-rate, 10))

	return false
}
