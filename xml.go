package junit

import (
	"encoding/xml"
	"fmt"
	"os"
	"sync"
	"time"
)

// A nice JUnit XML file format document can be found here:
// http://llg.cubic.org/docs/junit/

type testSuites struct {
	suiteLock sync.Mutex

	XMLName xml.Name     `xml:"testsuites"`
	Suites  []*testSuite `xml:"testsuite"`
}

func newTestSuites() *testSuites {
	return &testSuites{
		Suites: make([]*testSuite, 0),
	}
}

func (s *testSuites) GetSuite(name string) *testSuite {
	s.suiteLock.Lock()
	defer s.suiteLock.Unlock()

	for _, suite := range s.Suites {
		if suite.Name == name {
			return suite
		}
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}

	newSuite := &testSuite{
		Name:      name,
		ID:        len(s.Suites),
		Package:   name,
		Timestamp: time.Now().Format("2006-01-02T15:04:05"),
		Hostname:  hostname,
	}
	s.Suites = append(s.Suites, newSuite)

	return newSuite
}

type testSuite struct {
	// These must be at the top of the struct due to how 64 bit values are
	// handled by the atomic package in 32 bit Go builds
	Tests    int64   `xml:"tests,attr"`
	Failures int64   `xml:"failures,attr"`
	Errors   int64   `xml:"errors,attr"`

	Name      string `xml:"name,attr"`
	ID        int    `xml:"id,attr"`
	Package   string `xml:"package,attr"`
	Timestamp string `xml:"timestamp,attr"`
	Hostname  string `xml:"hostname,attr"`

	Time     float64 `xml:"time,attr"`

	Properties *suiteProperties `xml:"properties,omitempty"`

	testCaseLock sync.Mutex
	TestCases    []*testCase `xml:"testcase"`
}

func (s *testSuite) AddProperty(name, value string) {
	if s.Properties == nil {
		s.Properties = &suiteProperties{
			Properties: make([]*suiteProperty, 1),
		}
	}
	s.Properties.Properties = append(s.Properties.Properties, &suiteProperty{
		Name:  name,
		Value: value,
	})
}
func (s *testSuite) AddTestCase(newCase *testCase) {
	s.testCaseLock.Lock()
	defer s.testCaseLock.Unlock()

	if s.TestCases == nil {
		s.TestCases = make([]*testCase, 1)
	}
	s.TestCases = append(s.TestCases, newCase)
}

type suiteProperties struct {
	Properties []*suiteProperty `xml:"property"`
}

type suiteProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type testCase struct {
	ClassName string  `xml:"classname,attr"`
	Name      string  `xml:"name,attr"`
	Time      float64 `xml:"time,attr"`

	Failure *caseFailure `xml:"failure"`
}

func (c *testCase) SetFailure(file string, line int, message string) {
	c.Failure = &caseFailure{
		Message: fmt.Sprintf("%s:%d", file, line),
		Type:    "assertion",
		Text:    message,
	}
}

type caseFailure struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
	Text    string `xml:",chardata"`
}
