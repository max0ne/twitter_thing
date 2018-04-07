package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type RegisterTestSuite struct {
	RouteTestSuite
}

func TestUserRegister(t *testing.T) {
	suite.Run(t, new(RegisterTestSuite))
}

func (suite *RegisterTestSuite) TestUserRegister() {
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
			expCode:    401,
			expBodyMap: nil,
		},

		TestCase{
			method: "POST",
			path:   "/user/login",
			form: map[string]string{
				"username": "u1",
			},
			expCode:    400,
			expBodyMap: nil,
		},

		TestCase{
			method: "POST",
			path:   "/user/signup",
			form: map[string]string{
				"username": "u1",
				"password": "u1pass",
			},
			expCode:    200,
			expBodyMap: nil,
			postTestCase: func(resp *http.Response) {
				loginCookies = resp.Cookies()
			},
		},

		TestCase{
			method:  "GET",
			path:    "/user/get/u1",
			expCode: 200,
			expBodyMap: map[string]string{
				"username": "u1",
			},
		},

		TestCase{
			method:  "POST",
			path:    "/user/unregister",
			expCode: 200,
		},

		TestCase{
			method:  "GET",
			path:    "/user/get/u1",
			expCode: 404,
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
		suite.runTestCase(tc)
	}
}
