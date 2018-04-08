package main

import (
	"net/http"
	"testing"

	"github.com/max0ne/twitter_thing/back/middleware"

	"github.com/stretchr/testify/suite"
)

type TweetTestSuite struct {
	RouteTestSuite
}

func TestTweet(t *testing.T) {
	suite.Run(t, new(TweetTestSuite))
}

func (suite *TweetTestSuite) Test() {
	var token string
	addToken := func(req *http.Request) {
		if token != "" {
			req.Header.Add(middleware.TokenHeader, token)
		}
	}
	storeToken := func(resp *http.Response) {
		token = resp.Header.Get(middleware.TokenHeader)
	}

	signInToBigV := TestCase{
		desc:   "sign in to big v",
		method: "POST",
		path:   "/user/login",
		form: map[string]string{
			"uname":    "big_v",
			"password": "big_vpass",
		},
		expCode:      200,
		postTestCase: storeToken,
	}

	signInToUser := TestCase{
		desc:   "sign in to user",
		method: "POST",
		path:   "/user/login",
		form: map[string]string{
			"uname":    "u1",
			"password": "u1pass",
		},
		expCode:      200,
		postTestCase: storeToken,
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
				"uname":    "u1",
				"password": "u1pass",
			},
			expCode:      200,
			postTestCase: storeToken,
		},

		TestCase{
			desc:   "signup a big v",
			method: "POST",
			path:   "/user/signup",
			form: map[string]string{
				"uname":    "big_v",
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

		signInToBigV,

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
			desc:    "should see his own tweet",
			method:  "GET",
			path:    "/tweet/feed",
			expCode: 200,
			expBodyMapArr: []map[string]string{
				map[string]string{
					"tid":     "1",
					"content": "tweet1",
					"uname":   "big_v",
				},
			},
		},

		signInToUser,

		TestCase{
			desc:    "u1's feed should contain the tweet",
			method:  "GET",
			path:    "/tweet/feed",
			expCode: 200,
			expBodyMapArr: []map[string]string{
				map[string]string{
					"tid":     "1",
					"content": "tweet1",
					"uname":   "big_v",
				},
			},
		},

		signInToBigV,

		TestCase{
			desc:    "delete this tweet",
			method:  "POST",
			path:    "/tweet/del/1",
			expCode: 200,
		},

		TestCase{
			desc:          "no longer see tweet after delete",
			method:        "GET",
			path:          "/tweet/feed",
			expCode:       200,
			expBodyMapArr: []map[string]string{},
		},

		signInToUser,

		TestCase{
			desc:          "no longer see tweet after delete",
			method:        "GET",
			path:          "/tweet/feed",
			expCode:       200,
			expBodyMapArr: []map[string]string{},
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
