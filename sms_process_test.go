package main


import (
	"testing"
)

type mockDB struct{

}

type mockClient struct{

}


var smstests = []struct {
	from string
	to   string
	text string
	out  string
}{
	{"12345678", "87654321", "hello", ""},
	{"", "87654321", "hello", "from is missing"},
	{"12345678", "", "hello", "to is missing"},
	{"12345678", "87654321", "", "text is missing"},
	{"1234", "87654321", "hello", "from is invalid"},
	{"12341234123412341234", "87654321", "hello", "from is invalid"},
	{"12345678", "87654", "hello", "to is invalid"},
	{"12345678", "8765487654876548765487654", "hello", "to is invalid"},
	{"12345678", "87654321", "hellofromtheoutsidehellofromtheoutsidehellofromtheoutsidehellofromtheoutsidehellofromtheoutsidehellofromtheoutsidehellofromtheoutside", "text is invalid"},

}

func TestValidateFormData(t *testing.T) {

	for _, smstest := range smstests {
		s := ValidateFormData(smstest.from, smstest.to, smstest.text)
		if s != smstest.out {
			t.Errorf("validateformdata(%q,%q,%q) => %q want %q", smstest.from, smstest.to, smstest.text, s, smstest.out)
		}
	}
}

func (mdb *mockDB) UserExists(username, auth_id string) bool {

	var users =[]struct {
		username string
		auth_id string
		id int
	}{
		{"plivo1","20S0KPNOIM",1},
		{"plivo2","54P2EOKQ47",2},

	}

	for _, v  := range users{
		if username==v.username && auth_id==v.auth_id {
			userId=v.id
			return true
		}
	}

	return false
}


func (mdb *mockDB) NumberExists(number string) bool  {
	//fmt.Println("Userid: ", userId)
	//fmt.Println("number: ", number)

	allowednumbers := map[string]int{ "31297728125":1, "441224459571":2 };
	for k, v  := range allowednumbers {
		if number==k && userId==v {
			return true
		}
	}

	return false
}


func (mclient *mockClient) CacheSms(from, to string) bool{

	//fmt.Println("to value in mock cachesms: ",to)
	if to=="31297728125" {
		return true
	}

	return false
}

func (mclient *mockClient) CacheExists(from, to string) bool{

	//fmt.Println("to value in mock cachesms: ",to)

	if to =="31297728125" {
		return true
	}

	return false

}


