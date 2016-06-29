package  main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"io/ioutil"
)

//func fakeApp(msg string) *httptest.Server{
//
//	return httptest.NewServer()
//
//}

func TestInboundSms(t *testing.T) {
	t.Log("hello form testing")
}

func get(t *testing.T, s *httptest.Server, path string) string {
	resp, err := http.Get(s.URL + path)
	if err != nil {
		t.Error(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	return string(body)
}