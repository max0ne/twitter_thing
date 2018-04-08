package main

import (
	"net/http"
	"testing"

	"github.com/max0ne/twitter_thing/back/middleware"
	"github.com/stretchr/testify/suite"
)

type UserRouteTestSuite struct {
	RouteTestSuite
}

func TestUserRoute(t *testing.T) {
	suite.Run(t, new(UserRouteTestSuite))
}

func (suite *UserRouteTestSuite) TestUserRoute() {
	var token string
	addToken := func(req *http.Request) {
		if token != "" {
			req.Header.Add(middleware.TokenHeader, token)
		}
	}
	storeToken := func(resp *http.Response) {
		token = resp.Header.Get(middleware.TokenHeader)
	}

	testcases := []TestCase{
		TestCase{
			method: "POST",
			path:   "/user/login",
			form: map[string]string{
				"uname": "u1",
				"password": "u1pass",
			},
			expCode:    401,
			expBodyMap: nil,
		},

		TestCase{
			method: "POST",
			path:   "/user/login",
			form: map[string]string{
				"uname": "u1",
			},
			expCode:    400,
			expBodyMap: nil,
		},

		TestCase{
			method: "POST",
			path:   "/user/signup",
			form: map[string]string{
				"uname": "u1",
				"password": "u1pass",
			},
			expCode:      200,
			expBodyMap:   nil,
			postTestCase: storeToken,
		},

		TestCase{
			method:  "GET",
			path:    "/user/get/u1",
			expCode: 200,
			expBodyMap: map[string]string{
				"uname": "u1",
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
			addToken(req)
		}
		suite.runTestCase(tc)
	}
}
