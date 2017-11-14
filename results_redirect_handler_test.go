package wptdashboard

import (
	"testing"
)

type Case struct {
	testRun  TestRun
	testFile string
	expected string
}

const sha = "abcdef0123"
const resultsURLBase = "https://storage.googleapis.com/wptd/" + sha + "/"
const platform = "chrome-63.0-linux"
const resultsURL = resultsURLBase + "/" + platform + "-summary.json.gz"

func TestGetResultsURL_EmptyFile(t *testing.T) {
	checkResult(
		t,
		Case{
			TestRun{
				ResultsURL: resultsURL,
				Revision:   sha,
			},
			"",
			resultsURL,
		})
}

func TestGetResultsURL_TestFile(t *testing.T) {
	file := "css/vendor-imports/mozilla/mozilla-central-reftests/flexbox/flexbox-root-node-001b.html"
	checkResult(
		t,
		Case{
			TestRun{
				ResultsURL: resultsURL,
				Revision:   sha,
			},
			file,
			resultsURLBase + platform + "/" + file,
		})
}

func TestGetResultsURL_TrailingSlash(t *testing.T) {
	checkResult(
		t,
		Case{
			TestRun{
				ResultsURL: resultsURL,
				Revision:   sha,
			},
			"/",
			resultsURL,
		})
}

func checkResult(t *testing.T, c Case) {
	got := getResultsURL(c.testRun, c.testFile)
	if got != c.expected {
		t.Errorf("\nGot:\n%q\nExpected:\n%q", got, c.expected)
	}
}
