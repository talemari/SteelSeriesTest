package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	var testRequest TextRequest
	testRequest.Video = TestVideo{
		InputVideoPath:  "test_input1.mp4",
		Duration:        "60.0s",
		Resolution:      "1920 x 1080",
		OutputVideoPath: "test_output1.mp4",
	}
	testRequest.Text = TextEffect{
		TextString: "I’m sOoOo good at this game! xD",
		X:          "200",
		Y:          "100",
		FontSize:   "64",
		FontColor:  "0xFFFFFF",
		StartTime:  "23.0s",
		EndTime:    "40.0s",
	}
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(testRequest)
	if err != nil {
		fmt.Printf("Error : %s\n", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/", "application/json", b)
	if err != nil {
		fmt.Printf("Error : %s\n", err)
		return
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		return
	}
	expectedResponse := `ffmpeg -i test_input1.mp4 -vf drawtext="enable='between(t,23.0,40.0)': text='I’m sOoOo good at this game! xD':fontcolor=0xFFFFFF:fontsize=64: x=200:y=100" test_output1.mp4"`
	assert.Equal(t, expectedResponse, string(responseBody))
	defer resp.Body.Close()
}
