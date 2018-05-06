package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/max0ne/twitter_thing/back/middleware"
	"github.com/stretchr/testify/suite"
)

type GetNewUserTestSuite struct {
	RouteTestSuite
}

func TestGetNewUserTest(t *testing.T) {
	suite.Run(t, new(GetNewUserTestSuite))
}

func signUpForAnotherBunchOfUsers(unameStart, unameEnd int) []TestCase {
	tcs := []TestCase{}
	for idx := unameStart; idx < unameEnd; idx++ {
		tcs = append(tcs,
			TestCase{
				desc:   "sign up 1 user",
				method: "POST",
				path:   "/user/signup",
				form: map[string]string{
					"uname":    fmt.Sprintf("u%d", idx),
					"password": fmt.Sprintf("u%dpass", idx),
				},
				expCode: 200,
			},
		)
	}
	return tcs
}

func (suite *GetNewUserTestSuite) Test() {
	var token string
	addToken := func(req *http.Request) {
		if token != "" {
			req.Header.Add(middleware.TokenHeader, token)
		}
	}
	storeToken := func(resp *http.Response, bodyString string) {
		token = resp.Header.Get(middleware.TokenHeader)
	}

	testcases := []TestCase{

		TestCase{
			desc:          "empty before any user sign up",
			method:        "GET",
			path:          "/user/new",
			expCode:       200,
			expBodyMapArr: []map[string]string{},
		},

		TestCase{
			desc:   "sign up 1 user",
			method: "POST",
			path:   "/user/signup",
			form: map[string]string{
				"uname":    "u1",
				"password": "u1pass",
			},
			expCode:    200,
			expBodyMap: nil,
		},

		TestCase{
			desc:    "has this 1 user",
			method:  "GET",
			path:    "/user/new",
			expCode: 200,
			expBodyMapArr: []map[string]string{
				map[string]string{
					"uname": "u1",
				},
			},
		},

		TestCase{
			desc:   "sign up another user",
			method: "POST",
			path:   "/user/signup",
			form: map[string]string{
				"uname":    "u2",
				"password": "u2pass",
			},
			expCode:    200,
			expBodyMap: nil,
		},

		TestCase{
			desc:    "has these 2 user",
			method:  "GET",
			path:    "/user/new",
			expCode: 200,
			expBodyMapArr: []map[string]string{
				map[string]string{
					"uname": "u1",
				},
				map[string]string{
					"uname": "u2",
				},
			},
		},
	}

	testcases = append(testcases, signUpForAnotherBunchOfUsers(3, 12)...)

	testcases = append(testcases, []TestCase{
		TestCase{
			desc:    "has the last 10 users, excluding the very first one",
			method:  "GET",
			path:    "/user/new",
			expCode: 200,
			expBodyMapArr: []map[string]string{
				map[string]string{
					"uname": "u2",
				},
				map[string]string{
					"uname": "u3",
				},
				map[string]string{
					"uname": "u4",
				},
				map[string]string{
					"uname": "u5",
				},
				map[string]string{
					"uname": "u6",
				},
				map[string]string{
					"uname": "u7",
				},
				map[string]string{
					"uname": "u8",
				},
				map[string]string{
					"uname": "u9",
				},
				map[string]string{
					"uname": "u10",
				},
				map[string]string{
					"uname": "u11",
				},
			},
		},
		TestCase{
			desc:   "sign in to user 2",
			method: "POST",
			path:   "/user/login",
			form: map[string]string{
				"uname":    "u2",
				"password": "u2pass",
			},
			expCode:      200,
			postTestCase: storeToken,
		},

		TestCase{
			desc:    "new users does not contain current user",
			method:  "GET",
			path:    "/user/new",
			expCode: 200,
			expBodyMapArr: []map[string]string{
				map[string]string{
					"uname": "u3",
				},
				map[string]string{
					"uname": "u4",
				},
				map[string]string{
					"uname": "u5",
				},
				map[string]string{
					"uname": "u6",
				},
				map[string]string{
					"uname": "u7",
				},
				map[string]string{
					"uname": "u8",
				},
				map[string]string{
					"uname": "u9",
				},
				map[string]string{
					"uname": "u10",
				},
				map[string]string{
					"uname": "u11",
				},
			},
		},
	}...)

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
