package goreportcard

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractResponse(t *testing.T) {
	t.Run("basic success", func(t *testing.T) {
		htmlBytes, err := ioutil.ReadFile("testdata/success_response.html")
		require.NoError(t, err)

		respBytes, err := ioutil.ReadFile("testdata/response.json")
		require.NoError(t, err)

		var expResp Report
		require.NoError(t, json.Unmarshal(respBytes, &expResp))

		resp, err := extractResponse(string(htmlBytes))
		assert.NoError(t, err)
		assert.Equal(t, &expResp, resp)
	})
}
