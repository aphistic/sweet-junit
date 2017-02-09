package junit

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"math"
	"os"
	"runtime"
	"time"

	"github.com/aphistic/sweet"
)

var (
	junitOutput string
)

func init() {
	flag.StringVar(&junitOutput, "junit.output", "", "path to write a junit XML file after tests complete")
}

func roundTime(duration time.Duration) float64 {
	seconds := float64(duration) / float64(time.Second)
	ms := math.Floor(seconds * 10000)
	return ms / 10000
}

type JUnitPlugin struct {
	suites *testSuites
}

func NewPlugin() *JUnitPlugin {
	return &JUnitPlugin{
		suites: newTestSuites(),
	}
}

func (p *JUnitPlugin) Name() string {
	return "JUnit Output"
}

func (p *JUnitPlugin) Starting() {

}
func (p *JUnitPlugin) SuiteStarting(suite string) {
	s := p.suites.GetSuite(suite)
	s.AddProperty("go.version", runtime.Version())
}
func (p *JUnitPlugin) TestStarting(suite, test string) {

}
func (p *JUnitPlugin) TestPassed(suite, test string, stats *sweet.TestPassedStats) {
	s := p.suites.GetSuite(suite)
	s.Tests++
	s.AddTestCase(&testCase{
		Name:      test,
		ClassName: suite,
		Time:      roundTime(stats.Time),
	})
}
func (p *JUnitPlugin) TestFailed(suite, test string, stats *sweet.TestFailedStats) {
	s := p.suites.GetSuite(suite)
	s.Tests++
	s.Failures++

	tc := &testCase{
		Name:      test,
		ClassName: suite,
		Time:      roundTime(stats.Time),
	}
	tc.SetFailure(stats.File, stats.Line, stats.Message)
	s.AddTestCase(tc)
}
func (p *JUnitPlugin) SuiteFinished(suite string, stats *sweet.SuiteFinishedStats) {
	s := p.suites.GetSuite(suite)
	s.Time = roundTime(stats.Time)
}
func (p *JUnitPlugin) Finished() {
	if junitOutput == "" {
		return
	}

	of, err := os.OpenFile(junitOutput, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Unable to open file for JUnit output: %s\n", err)
		os.Exit(1)
	}
	defer of.Close()

	x, err := p.generateXML()
	if err != nil {
		fmt.Printf("Error generating xml: %s", err)
		return
	}

	of.WriteString(xml.Header)
	of.WriteString(x)
	of.WriteString("\n")
}

func (p *JUnitPlugin) generateXML() (string, error) {
	buffer := &bytes.Buffer{}
	enc := xml.NewEncoder(buffer)
	enc.Indent("", "    ")

	err := enc.Encode(p.suites)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
