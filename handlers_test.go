package  main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"net/url"




	"strings"
)

var inboundtests = []struct {
	username string
	auth_id string
	from string
	to   string
	text string
	outmessage  string
	outerror string
}{
	{"plivo1","20S0KPNOIM","12345678", "87654321", "hello", "","to parameter not found"},
	{"plivo1","20S0KPNOIM","12345678", "31297728125", "hello", "inbound sms ok",""},
	{"plivo1","20S0KPNOIM","12345678", "31297728125", "STOP", "inbound sms ok",""},
	{"plivo1","20S0KPNOIM","12345678", "31297728125", `STOP\n`, "inbound sms ok",""},
	{"plivo1","20S0KPNOIM","12345678", "31297728125", `STOP\r`, "inbound sms ok",""},
	{"plivo1","20S0KPNOIM","12345678", "31297728125", `STOP\r\n`, "inbound sms ok",""},
	{"plivo1","20S0KPNOIM","12345678", "", "hello","","to is missing"},
	{"plivo1","20S0KPNOIM","12345678", "87654321", "","", "text is missing"},
	{"plivo1","20S0KPNOIM","1234", "87654321", "hello", "","from is invalid"},
	{"plivo1","20S0KPNOIM","12341234123412341234", "87654321", "hello","", "from is invalid"},
	{"plivo1","20S0KPNOIM","12345678", "87654", "hello","", "to is invalid"},
	{"plivo1","20S0KPNOIM","12345678", "8765487654876548765487654", "hello","", "to is invalid"},
	{"plivo1","20S0KPNOIM","12345678", "87654321", "hellofromtheoutsidehellofromtheoutsidehellofromtheoutsidehellofromtheoutsidehellofromtheoutsidehellofromtheoutsidehellofromtheoutside","", "text is invalid"},

}


var outboundtests = []struct {
	username string
	auth_id string
	from string
	to   string
	text string
	outmessage  string
	outerror string
}{

	{"plivo1","20S0KPNOIM","12345678", "87654321", "hello", "","from parameter not found"},
	{"plivo1","20S0KPNOIM","12345678", "31297728125", "hello", "","sms from 12345678 to 31297728125 blocked by STOP request"},
	{"plivo1","20S0KPNOIM","31297728125", "31297728190", "STOP", "outbound sms ok",""},
	{"plivo1","20S0KPNOIM","12345678", "", "hello","","to is missing"},
	{"plivo1","20S0KPNOIM","12345678", "87654321", "","", "text is missing"},
	{"plivo1","20S0KPNOIM","1234", "87654321", "hello", "","from is invalid"},
	{"plivo1","20S0KPNOIM","12341234123412341234", "87654321", "hello","", "from is invalid"},
	{"plivo1","20S0KPNOIM","12345678", "87654", "hello","", "to is invalid"},
	{"plivo1","20S0KPNOIM","12345678", "8765487654876548765487654", "hello","", "to is invalid"},
	{"plivo1","20S0KPNOIM","12345678", "87654321", "hellofromtheoutsidehellofromtheoutsidehellofromtheoutsidehellofromtheoutsidehellofromtheoutsidehellofromtheoutsidehellofromtheoutside","", "text is invalid"},

}



func TestInboundSms(t *testing.T) {

	for _, inboundtest := range inboundtests {

		rec := httptest.NewRecorder()

		data:=url.Values{}
		data.Add("from", inboundtest.from)
		data.Add("to", inboundtest.to)
		data.Add("text", inboundtest.text)

		req, _ := http.NewRequest("POST", "/inbound/sms", strings.NewReader(data.Encode()))

		req.SetBasicAuth(inboundtest.username, inboundtest.auth_id)
		req.Header.Add("Content-Type","application/x-www-form-urlencoded")


		env := Env{db: &mockDB{}, client:&mockClient{}}
		http.HandlerFunc(env.BasicAuth(env.InboundSms)).ServeHTTP(rec, req)


		var output jsonOutput
		if err := json.Unmarshal(rec.Body.Bytes(), &output); err != nil {
			t.Errorf("\n.. Unable to parse reponse, recieved=  %v", rec.Body.String())
		}
		if inboundtest.outmessage != output.Message || inboundtest.outerror !=output.Error {
			t.Errorf("\n... request(from=%v,to=%v,text=%v)  \n...expected message = %v\n...obtained message = %v   \n.... expected error = %v\n....obtained error  = %v  ",
				inboundtest.from,inboundtest.to,inboundtest.text,inboundtest.outmessage, output.Message, inboundtest.outerror, output.Error)
		}
	}
}


func TestOutboundSms(t *testing.T) {

	for _, outboundtest := range outboundtests {

		rec := httptest.NewRecorder()

		data:=url.Values{}
		data.Add("from", outboundtest.from)
		data.Add("to", outboundtest.to)
		data.Add("text", outboundtest.text)

		req, _ := http.NewRequest("POST", "/outbound/sms", strings.NewReader(data.Encode()))

		req.SetBasicAuth(outboundtest.username, outboundtest.auth_id)
		req.Header.Add("Content-Type","application/x-www-form-urlencoded")


		env := Env{db: &mockDB{}, client:&mockClient{}}
		http.HandlerFunc(env.BasicAuth(env.OutboundSms)).ServeHTTP(rec, req)


		var output jsonOutput
		if err := json.Unmarshal(rec.Body.Bytes(), &output); err != nil {
			t.Errorf("\n.. Unable to parse reponse, recieved=  %v", rec.Body.String())
		}
		if outboundtest.outmessage != output.Message || outboundtest.outerror !=output.Error {
			t.Errorf("\n... request(from=%v,to=%v,text=%v)  \n...expected message = %v\n...obtained message = %v   \n.... expected error = %v\n....obtained error  = %v  ",
				outboundtest.from,outboundtest.to,outboundtest.text,outboundtest.outmessage, output.Message, outboundtest.outerror, output.Error)
		}
	}
}



func TestBasicAuth(t *testing.T) {

	rec := httptest.NewRecorder()
	request, err := http.NewRequest("POST", "/inbound/sms", nil)
	request.SetBasicAuth("plivotest", "JKHGkjgd")


	env := Env{db: &mockDB{}, client:&mockClient{}}
	http.HandlerFunc(env.BasicAuth(env.InboundSms)).ServeHTTP(rec, request)

	if err != nil {
		t.Error(err)
	}

	if rec.Code != 403 {
		t.Errorf("Unauthorised access  expected: %d", rec.Code)
	}

}