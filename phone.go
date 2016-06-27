package main

type Phone struct {
	id	int	`json:"id"`
	number	string	`json:"number"`
	account_id	int	`json:"account_id"`
}

type Phones []Phone
