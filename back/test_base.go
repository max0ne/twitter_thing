package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/max0ne/twitter_thing/back/util"
	"github.com/stretchr/testify/suite"
)

type TestCase struct {
	method        string
	path          string
	form          map[string]string
	expCode       int
	expBodyMap    map[string]string
	expBodyMapArr []map[string]string
	preTestCase   func(req *http.Request)
	postTestCase  func(resp *http.Response)
}

type RouteTestSuite struct {
	suite.Suite
	ts *httptest.Server
}

func (suite *RouteTestSuite) SetupTest() {
	router := NewServer().router
	ts := httptest.NewServer(router)
	suite.ts = ts
}

func (suite *RouteTestSuite) runTestCase(tc TestCase) {
	if tc.method == "GET" {
		suite.testGET(tc)
	} else if tc.method == "POST" {
		suite.testPOST(tc)
	} else {
		suite.Fail("unsuppoorted test case method" + tc.method)
	}
}

func (suite *RouteTestSuite) testGET(tc TestCase) {
	req, err := http.NewRequest("GET", suite.ts.URL+tc.path, nil)
	suite.Require().NoError(err)

	if tc.preTestCase != nil {
		tc.preTestCase(req)
	}

	resp, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	suite.assertResponse(tc, resp)

	if tc.postTestCase != nil {
		tc.postTestCase(resp)
	}
}

func (suite *RouteTestSuite) testPOST(tc TestCase) {
	req, err := http.NewRequest("POST", suite.ts.URL+tc.path, nil)
	suite.Require().NoError(err)

	data := url.Values{}
	for formkey, formval := range tc.form {
		data.Set(formkey, formval)
	}
	req.URL.RawQuery = data.Encode()

	if tc.preTestCase != nil {
		tc.preTestCase(req)
	}

	resp, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	suite.assertResponse(tc, resp)

	if tc.postTestCase != nil {
		tc.postTestCase(resp)
	}
}

func (suite *RouteTestSuite) assertResponse(tc TestCase, resp *http.Response) {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	suite.Require().NoError(err)
	bodyString := string(bodyBytes)

	suite.Require().Equal(tc.expCode, resp.StatusCode, fmt.Sprintf("%s: %s %s", tc.method, tc.path, bodyString))

	if tc.expBodyMap != nil {
		bodyAsMap := map[string]string{}
		suite.Require().NoError(json.Unmarshal(bodyBytes, &bodyAsMap))
		for expKey, expVal := range tc.expBodyMap {
			suite.Require().Equal(expVal, bodyAsMap[expKey], fmt.Sprintf("%s: %s %s", tc.method, tc.path, bodyString))
		}
	}

	if tc.expBodyMapArr != nil {
		errorMsg := fmt.Sprintf("%s: %s %s %s", tc.method, tc.path, bodyString, util.JSONMarshel(tc.expBodyMapArr))

		bodyAsMapArr := []map[string]string{}
		suite.Require().NoError(json.Unmarshal(bodyBytes, &bodyAsMapArr), errorMsg)
		suite.Require().Equal(len(tc.expBodyMapArr), len(bodyAsMapArr), errorMsg)
		for idx, aMap := range bodyAsMapArr {
			for expKey, expVal := range tc.expBodyMapArr[idx] {
				suite.Require().Equal(expVal, aMap[expKey], errorMsg)
			}
		}
	}
}
