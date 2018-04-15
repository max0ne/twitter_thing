package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"

	"github.com/max0ne/twitter_thing/back/config"
	"github.com/max0ne/twitter_thing/back/db"

	"github.com/max0ne/twitter_thing/back/util"
	"github.com/stretchr/testify/suite"
)

// TestCase an http based test case
type TestCase struct {
	desc          string
	method        string
	path          string
	form          map[string]string
	expCode       int
	expBodyMap    map[string]string
	expBodyMapArr []map[string]string
	preTestCase   func(req *http.Request)
	postTestCase  func(resp *http.Response, bodyString string)
}

// RouteTestSuite test http route base suite
type RouteTestSuite struct {
	suite.Suite
	ts       *httptest.Server
	dbServer *db.Server
}

var incrementingDBPort = 4000

func newDB() (*db.Server, error) {
	incrementingDBPort++
	return db.NewServer(config.Config{
		Role:   "db",
		DBAddr: "localhost",
		DBPort: strconv.FormatInt(int64(incrementingDBPort), 10),
	})
}

// SetupTest - -
func (suite *RouteTestSuite) SetupTest() {
	dbServer, err := newDB()
	suite.Require().NoError(err)
	suite.dbServer = dbServer
	suite.Require().NoError(dbServer.Start())

	suite.ts = httptest.NewServer(NewServer(config.Config{
		Role:   "api",
		DBAddr: "localhost",
		DBPort: dbServer.Port(),
	}).router)
}

func (suite *RouteTestSuite) runTestCase(tc TestCase) {
	if tc.desc != "" {
		fmt.Println("")
		fmt.Println(">-", tc.desc)
	}
	suite.Require().NotNil(suite.ts)
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

	bodyString := suite.assertResponse(tc, resp)

	if tc.postTestCase != nil {
		tc.postTestCase(resp, bodyString)
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

	bodyString := suite.assertResponse(tc, resp)

	if tc.postTestCase != nil {
		tc.postTestCase(resp, bodyString)
	}
}

// returns body string
func (suite *RouteTestSuite) assertResponse(tc TestCase, resp *http.Response) string {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	suite.Require().NoError(err)
	bodyString := string(bodyBytes)

	errorMsg := fmt.Sprintf("%s: %s %s", tc.method, tc.path, bodyString)
	suite.Require().Equal(tc.expCode, resp.StatusCode, errorMsg)

	if tc.expBodyMap != nil {
		errorMsg := fmt.Sprintf("%s: %s %s %s", tc.method, tc.path, bodyString, util.JSONMarshel(tc.expBodyMap))
		bodyAsMap := map[string]string{}
		suite.Require().NoError(json.Unmarshal(bodyBytes, &bodyAsMap), errorMsg)
		for expKey, expVal := range tc.expBodyMap {
			suite.Require().Equal(expVal, bodyAsMap[expKey], errorMsg)
		}
	}

	if tc.expBodyMapArr != nil {
		errorMsg := fmt.Sprintf("%s: %s \nexpected: %s\ngot: %s", tc.method, tc.path, util.JSONMarshel(tc.expBodyMapArr), bodyString)

		bodyAsMapArr := []map[string]string{}
		suite.Require().NoError(json.Unmarshal(bodyBytes, &bodyAsMapArr), errorMsg)
		suite.Require().Equal(len(tc.expBodyMapArr), len(bodyAsMapArr), errorMsg)
		for idx, aMap := range bodyAsMapArr {
			for expKey, expVal := range tc.expBodyMapArr[idx] {
				suite.Require().Equal(expVal, aMap[expKey], errorMsg)
			}
		}
	}

	return bodyString
}
