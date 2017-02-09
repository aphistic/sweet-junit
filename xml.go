package junit

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"
)

// A nice JUnit XML file format document can be found here:
// http://llg.cubic.org/docs/junit/

type testSuites struct {
	XMLName xml.Name     `xml:"testsuites"`
	Suites  []*testSuite `xml:"testsuite"`
}

func newTestSuites() *testSuites {
	return &testSuites{
		Suites: make([]*testSuite, 0),
	}
}

func (s *testSuites) AddSuite(suite *testSuite) {
	s.Suites = append(s.Suites, suite)
}

func (s *testSuites) GetSuite(name string) *testSuite {
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
	s.AddSuite(newSuite)

	return newSuite
}

type testSuite struct {
	Name      string  `xml:"name,attr"`
	ID        int     `xml:"id,attr"`
	Package   string  `xml:"package,attr"`
	Timestamp string  `xml:"timestamp,attr"`
	Hostname  string  `xml:"hostname,attr"`
	Tests     int     `xml:"tests,attr"`
	Failures  int     `xml:"failures,attr"`
	Errors    int     `xml:"errors,attr"`
	Time      float64 `xml:"time,attr"`

	Properties *suiteProperties `xml:"properties,omitempty"`
	TestCases  []*testCase      `xml:"testcase"`
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
