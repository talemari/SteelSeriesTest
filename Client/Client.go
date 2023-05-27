package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type TestVideo struct {
	InputVideoPath  string `json:"InputVideoPath"`
	Duration        string `json:"Duration"`
	Resolution      string `json:"Resolution"`
	OutputVideoPath string `json:"OutputVideoPath"`
}

type TextEffect struct {
	TextString string `json:"TextString"`
	X          string `json:"X"`
	Y          string `json:"Y"`
	FontSize   string `json:"FontSize"`
	FontColor  string `json:"FontColor"`
	StartTime  string `json:"StartTime"`
	EndTime    string `json:"EndTime"`
}

type TextRequest struct {
	Video TestVideo  `json:"Video"`
	Text  TextEffect `json:"Text"`
}

func main() {
	if len(os.Args) > 1 {
		fmt.Printf("Found argument : %s\n", os.Args[1])
	}
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
	defer resp.Body.Close()

	fmt.Printf("Response : %s\n", responseBody)
}
