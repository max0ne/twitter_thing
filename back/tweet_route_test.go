package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TweetTestSuite struct {
	RouteTestSuite
}

func TestTweet(t *testing.T) {
	suite.Run(t, new(TweetTestSuite))
}

func (suite *TweetTestSuite) Test() {
	var loginCookies []*http.Cookie
	addLoginCookies := func(req *http.Request) {
		for _, cok := range loginCookies {
			req.AddCookie(cok)
		}
	}

	testcases := []TestCase{

		TestCase{
			desc:    "feed is nothing when not logged in",
			method:  "GET",
			path:    "/tweet/feed",
			expCode: 401,
		},

		TestCase{
			desc:   "signup a user",
			method: "POST",
			path:   "/user/signup",
			form: map[string]string{
				"username": "u1",
				"password": "u1pass",
			},
			expCode: 200,
			postTestCase: func(resp *http.Response) {
				loginCookies = resp.Cookies()
			},
		},

		TestCase{
			desc:   "signup a big v",
			method: "POST",
			path:   "/user/signup",
			form: map[string]string{
				"username": "big_v",
				"password": "big_vpass",
			},
			expCode: 200,
		},

		TestCase{
			desc:    "user follow big v",
			method:  "POST",
			path:    "/user/follow/big_v",
			expCode: 200,
		},

		TestCase{
			desc:          "user's feed is empty",
			method:        "GET",
			path:          "/tweet/feed",
			expCode:       200,
			expBodyMapArr: []map[string]string{},
		},

		TestCase{
			desc:          "big v's feed is empty",
			method:        "GET",
			path:          "/tweet/user/big_v",
			expCode:       200,
			expBodyMapArr: []map[string]string{},
		},

		TestCase{
			desc:   "sign in to big v",
			method: "POST",
			path:   "/user/login",
			form: map[string]string{
				"username": "big_v",
				"password": "big_vpass",
			},
			expCode: 200,
			postTestCase: func(resp *http.Response) {
				loginCookies = resp.Cookies()
			},
		},

		TestCase{
			desc:   "big v post a tweet",
			method: "POST",
			path:   "/tweet/new",
			form: map[string]string{
				"content": "tweet1",
			},
			expCode: 200,
			expBodyMap: map[string]string{
				"tid":     "1",
				"content": "tweet1",
			},
		},

		TestCase{
			desc:    "big v can see this tweet",
			method:  "GET",
			path:    "/tweet/user/big_v",
			expCode: 200,
			expBodyMapArr: []map[string]string{
				map[string]string{
					"tid":     "1",
					"content": "tweet1",
				},
			},
		},

		TestCase{
			desc:          "his feed should still be empty",
			method:        "GET",
			path:          "/tweet/feed",
			expCode:       200,
			expBodyMapArr: []map[string]string{},
		},

		TestCase{
			desc:   "sign in to user",
			method: "POST",
			path:   "/user/login",
			form: map[string]string{
				"username": "u1",
				"password": "u1pass",
			},
			expCode: 200,
			postTestCase: func(resp *http.Response) {
				loginCookies = resp.Cookies()
			},
		},

		TestCase{
			desc:    "u1's feed should contain the tweet",
			method:  "GET",
			path:    "/tweet/feed",
			expCode: 200,
			expBodyMapArr: []map[string]string{
				map[string]string{
					"tid":     "1",
					"content": "tweet1",
				},
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
		suite.runTestCase(tc)
	}
}
