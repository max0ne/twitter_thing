package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/max0ne/twitter_thing/back/config"
	"github.com/max0ne/twitter_thing/back/db"
)

func getTestCases(cnt int) chan string {
	channel := make(chan string)
	go func() {
		for ii := 0; ii < cnt; ii++ {
			channel <- fmt.Sprintf("%d", ii)
		}
		close(channel)
	}()
	return channel
}

type TestSetGetSuite struct {
	suite.Suite
	dbServer *db.Server
}

func TestSetGet(t *testing.T) {
	suite.Run(t, new(TestSetGetSuite))
}

func (suite *TestSetGetSuite) SetupTest() {
	dbServer, err := newDB()
	suite.Require().NoError(err)
	suite.Require().NoError(dbServer.Start())
	suite.dbServer = dbServer
}

func (suite *TestSetGetSuite) TearDownTest() {
	fmt.Println("teardown")
}

func (suite *TestSetGetSuite) BeforeTest(suiteName, testName string) {
	fmt.Println(suiteName, testName, "start")
}

func (suite *TestSetGetSuite) AfterTest(suiteName, testName string) {
	fmt.Println(suiteName, testName, "done")
}

func (suite *TestSetGetSuite) TestGetSetSerial() {
	client, err := db.NewClient(config.Config{
		DBAddr: "localhost",
		DBPort: suite.dbServer.Port(),
	})
	suite.Require().NoError(err)
	t1 := client.NewTable("t1")
	for cs := range getTestCases(1000) {
		t1.Put(cs, fmt.Sprintf("%s_val", cs))
	}
	for cs := range getTestCases(1000) {
		got, err := t1.Get(cs)
		suite.Require().NoError(err)
		suite.Require().Equal(fmt.Sprintf("%s_val", cs), got)
	}
}

func (suite *TestSetGetSuite) TestGetParallel() {
	client, err := db.NewClient(config.Config{
		DBAddr: "localhost",
		DBPort: suite.dbServer.Port(),
	})
	suite.Require().NoError(err)
	t1 := client.NewTable("t1")
	putTestCaseChan := getTestCases(1000)
	putTableChan := make(chan bool)
	for ii := 0; ii < 100; ii++ {
		go func(ii int) {
			for ii := 0; ii < 10; ii++ {
				tc := <-putTestCaseChan
				t1.Put(tc, fmt.Sprintf("%s_val", tc))
			}
			putTableChan <- true
		}(ii)
	}
	for idx := 0; idx < 100; idx++ {
		<-putTableChan
	}

	getTestCaseChan := getTestCases(1000)
	getTableChan := make(chan bool)
	for ii := 0; ii < 1000; ii++ {
		go func(ii int) {
			tc := <-getTestCaseChan
			got, err := t1.Get(tc)
			suite.Require().NoError(err)
			suite.Require().Equal(fmt.Sprintf("%s_val", tc), got)
			getTableChan <- true
		}(ii)
	}
	for idx := 0; idx < 1000; idx++ {
		<-getTableChan
	}
}

func (suite *TestSetGetSuite) TestSetDelParallel() {
	client, err := db.NewClient(config.Config{
		DBAddr: "localhost",
		DBPort: suite.dbServer.Port(),
	})
	suite.Require().NoError(err)
	t1 := client.NewTable("t1")
	putTestCaseChan := getTestCases(1000)
	putTableChan := make(chan string)
	delTableChan := make(chan bool)
	for ii := 0; ii < 100; ii++ {
		go func() {
			for ii := 0; ii < 10; ii++ {
				tc := <-putTestCaseChan
				suite.Require().NoError(t1.Put(tc, fmt.Sprintf("%s_val", tc)))
				putTableChan <- tc
			}
		}()
	}
	for ii := 0; ii < 100; ii++ {
		go func() {
			for ii := 0; ii < 10; ii++ {
				tc := <-putTableChan
				suite.Require().NoError(t1.Del(tc))
				delTableChan <- true
			}
		}()
	}

	for ii := 0; ii < 1000; ii++ {
		<-delTableChan
	}

	fmt.Println("ha")

	getTestCaseChan := getTestCases(1000)
	getTableChan := make(chan bool)
	for ii := 0; ii < 1000; ii++ {
		go func() {
			tc := <-getTestCaseChan
			got, err := t1.Get(tc)
			suite.Require().NoError(err)
			suite.Require().Equal("", got)
			getTableChan <- true
		}()
	}
	for idx := 0; idx < 1000; idx++ {
		<-getTableChan
	}
}
