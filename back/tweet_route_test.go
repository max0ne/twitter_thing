package main

import (
	"fmt"
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
	storeToken := func(resp *http.Response, bodyString string) {
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

	checkIsUser := func(uname string) TestCase {
		return TestCase{
			desc:    "check current user",
			method:  "GET",
			path:    "/user/me",
			expCode: 200,
			expBodyMap: map[string]string{
				"uname": uname,
			},
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
		checkIsUser("u1"),

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
		checkIsUser("big_v"),

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
		checkIsUser("u1"),

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

func postTweet(desc, content string, tid string) TestCase {
	expBodyMap := map[string]string{
		"content": content,
	}
	if tid != "" {
		expBodyMap["tid"] = tid
	}

	return TestCase{
		desc:   desc,
		method: "POST",
		path:   "/tweet/new",
		form: map[string]string{
			"content": content,
		},
		expCode:    200,
		expBodyMap: expBodyMap,
	}
}

func logDB() TestCase {
	return TestCase{
		desc:    "log db",
		method:  "GET",
		path:    "/db",
		expCode: 200,
		postTestCase: func(resp *http.Response, bodyString string) {
			fmt.Println(bodyString)
		},
	}
}

func (suite *TweetTestSuite) TestPostNewAfterFollow() {
	var token string
	addToken := func(req *http.Request) {
		if token != "" {
			req.Header.Add(middleware.TokenHeader, token)
		}
	}
	storeToken := func(resp *http.Response, bodyString string) {
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

	unregisterUser := TestCase{
		desc:         "unregister",
		method:       "POST",
		path:         "/user/unregister",
		expCode:      200,
		postTestCase: storeToken,
	}

	testcases := []TestCase{

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

		postTweet("user post a tweet", "t1", "1"),

		TestCase{
			desc:    "user can see this tweet",
			method:  "GET",
			path:    "/tweet/feed",
			expCode: 200,
			expBodyMapArr: []map[string]string{
				map[string]string{
					"content": "t1",
					"tid":     "1",
					"uname":   "u1",
				},
			},
		},

		postTweet("user post another tweet", "t2", "2"),

		TestCase{
			desc:    "user can see this tweet",
			method:  "GET",
			path:    "/tweet/feed",
			expCode: 200,
			expBodyMapArr: []map[string]string{
				map[string]string{
					"content": "t2",
					"tid":     "2",
					"uname":   "u1",
				},
				map[string]string{
					"content": "t1",
					"tid":     "1",
					"uname":   "u1",
				},
			},
		},

		signInToBigV,

		TestCase{
			desc:          "bigv's feed is empty",
			method:        "GET",
			path:          "/tweet/feed",
			expCode:       200,
			expBodyMapArr: []map[string]string{},
		},

		postTweet("big v post a tweet", "v tweet 1", "3"),

		TestCase{
			desc:    "bigv's contains this tweet",
			method:  "GET",
			path:    "/tweet/feed",
			expCode: 200,
			expBodyMapArr: []map[string]string{
				map[string]string{
					"content": "v tweet 1",
					"tid":     "3",
					"uname":   "big_v",
				},
			},
		},

		signInToUser,

		TestCase{
			desc:    "user can see his 2 tweets and big v's tweet",
			method:  "GET",
			path:    "/tweet/feed",
			expCode: 200,
			expBodyMapArr: []map[string]string{
				map[string]string{
					"content": "v tweet 1",
					"tid":     "3",
					"uname":   "big_v",
				},
				map[string]string{
					"content": "t2",
					"tid":     "2",
					"uname":   "u1",
				},
				map[string]string{
					"content": "t1",
					"tid":     "1",
					"uname":   "u1",
				},
			},
		},

		signInToBigV,
		unregisterUser,
		signInToUser,

		TestCase{
			desc:    "user no longer see big v's tweet after big v unregistered",
			method:  "GET",
			path:    "/tweet/feed",
			expCode: 200,
			expBodyMapArr: []map[string]string{
				map[string]string{
					"content": "t2",
					"tid":     "2",
					"uname":   "u1",
				},
				map[string]string{
					"content": "t1",
					"tid":     "1",
					"uname":   "u1",
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
			addToken(req)
		}
		suite.runTestCase(tc)
	}
}

func (suite *TweetTestSuite) TestDeleteNonExistTweet() {
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
			desc:    "delete this tweet",
			method:  "POST",
			path:    "/tweet/del/nonexisttweetid",
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
