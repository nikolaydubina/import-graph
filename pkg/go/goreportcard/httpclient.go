package goreportcard

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// GoReportCardHTTPClient is unofficial interface to fetch data from goreportcard.com
type GoReportCardHTTPClient struct {
	HTTPClient *http.Client
	BaseURL    string
}

type GradeEnum string

var (
	GradeAP GradeEnum = "A+"
	GradeA  GradeEnum = "A"
	GradeB  GradeEnum = "B"
	GradeC  GradeEnum = "C"
	GradeD  GradeEnum = "D"
	GradeE  GradeEnum = "E"
	GradeF  GradeEnum = "F"
)

type Report struct {
	Average   float64   `json:"average"`
	Grade     GradeEnum `json:"grade"`
	NumFiles  uint      `json:"files"`
	NumIssues uint      `json:"issues"`
}

// redirectResponse can be returned for 200 response meaning browser should redirect to it
type redirectResponse struct {
	RedirectPath string `json:"redirect"` // e.g. "/report/github.com/go-playground/validator"
}

func (c *GoReportCardHTTPClient) GetReport(modName string) (*Report, error) {
	path := fmt.Sprintf("/report/%s", modName)

	// try fetch redirect path
	if resp, err := c.HTTPClient.Get(fmt.Sprintf("https://%s/checks?repo=%s", c.BaseURL, modName)); err == nil {
		defer func() { resp.Body.Close() }()
		var buf bytes.Buffer
		buf.ReadFrom(resp.Body)

		var redirect redirectResponse
		if err := json.Unmarshal(buf.Bytes(), &redirect); err == nil {
			path = redirect.RedirectPath
		}
	}

	resp, err := c.HTTPClient.Get(fmt.Sprintf("https://%s%s", c.BaseURL, path))
	if err != nil {
		return nil, fmt.Errorf("can not make GET: %w", err)
	}
	defer func() { resp.Body.Close() }()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	return extractResponse(buf.String())
}

// extractResponse parses HTML, finds variable with value, unescapes and unmarshals it
func extractResponse(htmlResp string) (*Report, error) {
	idxStart := strings.Index(htmlResp, `var response =  "`)
	if idxStart == -1 {
		return nil, errors.New("not found response object in html")
	}

	offEnd := strings.Index(htmlResp[idxStart:], "\" ;\n")
	if offEnd == -1 {
		return nil, errors.New("not found response end object in html")
	}

	respStrEscaped := htmlResp[idxStart+len(`var response =  "`) : idxStart+offEnd]
	respStr, err := strconv.Unquote(`"` + respStrEscaped + `"`)
	if err != nil {
		return nil, fmt.Errorf("can not unescape: %w", err)
	}

	var rep Report
	if err := json.Unmarshal([]byte(respStr), &rep); err != nil {
		return nil, fmt.Errorf("can not unmarshal: %w", err)
	}

	return &rep, nil
}
