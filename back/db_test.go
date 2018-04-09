package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

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
}

func TestSetGet(t *testing.T) {
	suite.Run(t, new(TestSetGetSuite))
}

func (suite *TestSetGetSuite) TestGetSetSerial() {
	db := db.NewStore()
	t1 := db.NewTable("t1")
	for cs := range getTestCases(1000) {
		t1.Put(cs, fmt.Sprintf("%s_val", cs))
	}
	for cs := range getTestCases(1000) {
		suite.Require().Equal(fmt.Sprintf("%s_val", cs), t1.Get(cs))
	}
}

func (suite *TestSetGetSuite) TestGetParallel() {
	db := db.NewStore()
	t1 := db.NewTable("t1")
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
			suite.Require().Equal(fmt.Sprintf("%s_val", tc), t1.Get(tc))
			getTableChan <- true
		}(ii)
	}
	for idx := 0; idx < 1000; idx++ {
		<-getTableChan
	}
}
