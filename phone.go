package main

type Phone struct {
	Id	int	`json:"id"`
	Number	string	`json:"number"`
	Account_id	int	`json:"account_id"`
}

type Phones []Phone
