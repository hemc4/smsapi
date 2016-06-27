package main

type Sms struct {
	from	string	`json:"from"`
	to	string	`json:"to"`
	text	string	`json:text`
}

type Smss []Sms
