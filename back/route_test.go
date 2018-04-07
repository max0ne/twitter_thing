package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	method       string
	path         string
	form         map[string]string
	expCode      int
	expBody      map[string]string
	preTestCase  func(req *http.Request)
	postTestCase func(resp *http.Response)
}

func (tc TestCase) test(t *testing.T, ts *httptest.Server) {
	if tc.method == "GET" {
		tc.testGET(t, ts)
	} else if tc.method == "POST" {
		tc.testPOST(t, ts)
	} else {
		assert.Fail(t, "unsuppoorted test case method"+tc.method)
	}
}

func (tc TestCase) testGET(t *testing.T, ts *httptest.Server) {
	req, err := http.NewRequest("GET", ts.URL+tc.path, nil)
	assert.NoError(t, err)

	if tc.preTestCase != nil {
		tc.preTestCase(req)
	}

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, tc.expCode, resp.StatusCode)

	if tc.expBody != nil {
		body := map[string]string{}

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		bodyString := string(bodyBytes)

		assert.NoError(t, json.Unmarshal(bodyBytes, &body))
		for expKey, expVal := range tc.expBody {
			assert.Equal(t, expVal, body[expKey], fmt.Sprintf("%s: %s %s", tc.method, tc.path, bodyString))
		}
	}

	if tc.postTestCase != nil {
		tc.postTestCase(resp)
	}
}

func (tc TestCase) testPOST(t *testing.T, ts *httptest.Server) {
	req, err := http.NewRequest("POST", ts.URL+tc.path, nil)
	assert.NoError(t, err)

	data := url.Values{}
	for formkey, formval := range tc.form {
		data.Set(formkey, formval)
	}
	req.URL.RawQuery = data.Encode()

	if tc.preTestCase != nil {
		tc.preTestCase(req)
	}

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, tc.expCode, resp.StatusCode)

	if tc.expBody != nil {
		body := map[string]string{}
		assert.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		for expKey, expVal := range tc.expBody {
			assert.Equal(t, expVal, body[expKey], fmt.Sprintf("%s: %s", tc.method, tc.path))
		}
	}

	if tc.postTestCase != nil {
		tc.postTestCase(resp)
	}
}

func TestUserRegister(t *testing.T) {
	router := NewServer().router
	ts := httptest.NewServer(router)
	var loginCookies []*http.Cookie
	addLoginCookies := func(req *http.Request) {
		for _, cok := range loginCookies {
			req.AddCookie(cok)
		}
	}

	testcases := []TestCase{
		TestCase{
			method: "POST",
			path:   "/user/login",
			form: map[string]string{
				"username": "u1",
				"password": "u1pass",
			},
			expCode: 401,
			expBody: nil,
		},

		TestCase{
			method: "POST",
			path:   "/user/login",
			form: map[string]string{
				"username": "u1",
			},
			expCode: 400,
			expBody: nil,
		},

		TestCase{
			method: "POST",
			path:   "/user/signup",
			form: map[string]string{
				"username": "u1",
				"password": "u1pass",
			},
			expCode: 200,
			expBody: nil,
			postTestCase: func(resp *http.Response) {
				loginCookies = resp.Cookies()
			},
		},

		TestCase{
			method:  "GET",
			path:    "/user/get/u1",
			expCode: 200,
			expBody: map[string]string{
				"username": "u1",
			},
		},
	}

	for _, tc := range testcases {
		oldPreTestCase := tc.preTestCase
		tc.preTestCase = func(req *http.Request) {
			if oldPreTestCase != nil {
				oldPreTestCase(req)
			}
			addLoginCookies(req)
		}
		tc.test(t, ts)
	}
}
