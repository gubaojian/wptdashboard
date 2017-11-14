package wptdashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
)

// This handler is responsible for all pages that display test results.
// It fetches the latest TestRun for each browser then renders the HTML
// page with the TestRuns encoded as JSON. The Polymer app picks those up
// and loads the summary files based on each entity's TestRun.ResultsURL.
//
// The browsers initially displayed to the user are defined in browsers.json.
// The JSON property "initially_loaded" is what controls this.
func testHandler(w http.ResponseWriter, r *http.Request) {
	runSHA, err := ParseSHAParam(r)
	if err != nil {
		http.Error(w, "Invalid query params", http.StatusBadRequest)
		return
	}

	var testRunSources []string
	var browserNames []string
	browserNames, err = GetBrowserNames()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	const sourceURL = `/api/run?browser=%s&sha=%s`
	for _, browserName := range browserNames {
		testRunSources = append(testRunSources, fmt.Sprintf(sourceURL, browserName, runSHA))
	}

	testRunSourcesBytes, err := json.Marshal(testRunSources)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		TestRunSources string
		SHA            string
	}{
		string(testRunSourcesBytes),
		runSHA,
	}

	if err := templates.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ParseSHAParam parses and validates the 'sha' param for the request.
// It returns "latest" by default (and in error cases).
func ParseSHAParam(r *http.Request) (runSHA string, err error) {
	// Get the SHA for the run being loaded (the first part of the path.)
	runSHA = "latest"
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return runSHA, err
	}

	runParam := params.Get("sha")
	regex := regexp.MustCompile("[0-9a-fA-F]{10}")
	if regex.MatchString(runParam) {
		runSHA = runParam
	}
	return runSHA, err
}
