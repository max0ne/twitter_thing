package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UserFollowTestSuite struct {
	RouteTestSuite
}

func TestUserFollow(t *testing.T) {
	suite.Run(t, new(UserFollowTestSuite))
}

func (suite *UserFollowTestSuite) TestUserFollow() {
	var loginCookies []*http.Cookie
	addLoginCookies := func(req *http.Request) {
		for _, cok := range loginCookies {
			req.AddCookie(cok)
		}
	}

	testcases := []TestCase{

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
			method:        "GET",
			path:          "/user/following/u1",
			expCode:       200,
			expBodyMapArr: []map[string]string{},
		},

		TestCase{
			method:  "POST",
			path:    "/user/follow/big_v",
			expCode: 200,
		},

		TestCase{
			method:        "GET",
			path:          "/user/following/u1",
			expCode:       200,
			expBodyMapArr: []map[string]string{},
		},

		TestCase{
			method: "POST",
			path:   "/user/signup",
			form: map[string]string{
				"username": "big_v",
				"password": "big_vpass",
			},
			expCode:    200,
			expBodyMap: nil,
		},

		TestCase{
			method:  "GET",
			path:    "/user/following/u1",
			expCode: 200,
			expBodyMapArr: []map[string]string{
				map[string]string{
					"username": "big_v",
				},
			},
		},

		TestCase{
			method:  "GET",
			path:    "/user/follower/big_v",
			expCode: 200,
			expBodyMapArr: []map[string]string{
				map[string]string{
					"username": "u1",
				},
			},
		},

		TestCase{
			method:  "POST",
			path:    "/user/unfollow/small_v",
			expCode: 200,
		},

		TestCase{
			method:  "POST",
			path:    "/user/unfollow/big_v",
			expCode: 200,
		},

		TestCase{
			method:        "GET",
			path:          "/user/following/u1",
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
			addLoginCookies(req)
		}
		suite.runTestCase(tc)
	}
}
