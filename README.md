# SMS API


## Install


```
go get github.com/hemendra619/smsapi
```


## Dependencies
Install PostgreSQL
https://www.postgresql.org/download/

Install redis
http://redis.io/download

Install go
https://golang.org/doc/install



## Test
Run tests
```
go test
```

No PostgreSQL database or Reidis db required to test. All db methods has been mocked.

#Run
run this project

```
go build && ./smsapi

```


## Libraries Used

Go sql driver for PostgreSQL
https://github.com/lib/pq


Go redis
https://github.com/go-redis/redis

Go redis rate limitor
https://github.com/go-redis/rate




