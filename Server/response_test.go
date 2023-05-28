package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestResponse(jsonPath string) (string, error) {
	jsonFile, err := os.Open(jsonPath)
	if err != nil {
		return err.Error(), err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var textRequest TextRequest
	json.Unmarshal(byteValue, &textRequest)

	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(textRequest)
	if err != nil {
		return err.Error(), err
	}

	resp, err := http.Post("http://localhost:8080/", "application/json", b)
	if err != nil {
		return err.Error(), err
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err.Error(), err
	}
	return string(responseBody), nil
}

func TestResponseToCorrectQuery(t *testing.T) {
	response, err := getTestResponse("Test.json")
	if err != nil {
		t.Fatalf("Test failed with error: %s", err)
	}
	expectedResponse := `ffmpeg -i test_input1.mp4 -vf drawtext="enable='between(t,23.0,40.0)': text='I’m sOoOo good at this game! xD':fontcolor=0xFFFFFF:fontsize=64: x=200:y=100" test_output1.mp4"`
	assert.Equal(t, expectedResponse, response)
}

func TestResponseToInvalidTextEndTime(t *testing.T) {
	response, err := getTestResponse("TestInvalidEndTime.json")
	if err != nil {
		t.Fatalf("Test failed with error: %s", err)
	}
	expectedResponse := "the text end time cannot be more than the duration of the video (60.000000)"
	assert.Equal(t, expectedResponse, response)
}

func TestResponseToInvalidTextCoordinates(t *testing.T) {
	response, err := getTestResponse("TestInvalidCoordinates.json")
	if err != nil {
		t.Fatalf("Test failed with error: %s", err)
	}
	expectedResponse := "text position (9999) cannot be less than 0 or more than the video resolution (1080)"
	assert.Equal(t, expectedResponse, response)
}
