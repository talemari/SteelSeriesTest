package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
	fmt.Println("Server started.")
	server()
}

func checkForErrors(data *TextRequest) error {
	textStartTime, err := strconv.ParseFloat(strings.Split(data.Text.StartTime, "s")[0], 32)
	if err != nil {
		return err
	}
	if textStartTime < 0 {
		return fmt.Errorf("the text start time cannot be less than 0")
	}
	textEndTime, err := strconv.ParseFloat(strings.Split(data.Text.EndTime, "s")[0], 32)
	if err != nil {
		return err
	}
	videoDuration, err := strconv.ParseFloat(strings.Split(data.Video.Duration, "s")[0], 32)
	if err != nil {
		return err
	}
	if textEndTime > videoDuration {
		return fmt.Errorf("the text end time cannot be more than the duration of the video (%f)", videoDuration)
	}
	resolution := strings.Split(strings.ReplaceAll(data.Video.Resolution, " ", ""), "x")
	textPositionX, err := strconv.Atoi(data.Text.X)
	if err != nil {
		return err
	}
	resolutionX, err := strconv.Atoi(resolution[0])
	if err != nil {
		return err
	}
	if textPositionX < 0 || textPositionX > resolutionX {
		return fmt.Errorf("text position (%d) cannot be less than 0 or more than the video resolution (%d)", textPositionX, resolutionX)
	}
	textPositionY, err := strconv.Atoi(data.Text.Y)
	if err != nil {
		return err
	}
	resolutionY, err := strconv.Atoi(resolution[1])
	if err != nil {
		return err
	}
	if textPositionY < 0 || textPositionY > resolutionY {
		return fmt.Errorf("text position (%d) cannot be less than 0 or more than the video resolution (%d)", textPositionY, resolutionY)
	}
	return nil
}

func generateFfmpegCommandLine(data *TextRequest) (string, error) {
	err := checkForErrors(data)
	if err != nil {
		return err.Error(), err
	}
	return `ffmpeg -i ` + data.Video.InputVideoPath + ` -vf drawtext="enable='between(t,` +
		strings.Split(data.Text.StartTime, "s")[0] + `,` + strings.Split(data.Text.EndTime, "s")[0] +
		`)': text='` + data.Text.TextString + `':fontcolor=` + data.Text.FontColor + `:fontsize=` + data.Text.FontSize +
		`: x=` + data.Text.X + `:y=` + data.Text.Y + `" ` + data.Video.OutputVideoPath + `"`, nil
}

func server() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		textRequest := &TextRequest{}
		err := json.NewDecoder(r.Body).Decode(textRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		readableJson, err := json.MarshalIndent(textRequest, "", "    ")
		if err != nil {
			fmt.Println("Failed to make the data human readable.")
			readableJson = []byte("I am a disgrace to my clan.")
		}
		fmt.Println("Received request:", string(readableJson))
		ffmpegCommand, err := generateFfmpegCommandLine(textRequest)
		if err != nil {
			fmt.Printf("Failed to parse request. Error: %s", err)
		} else {
			fmt.Printf("Final command : %s\n", ffmpegCommand)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(ffmpegCommand))
	})

	if err := http.ListenAndServe(":8080", nil); err != http.ErrServerClosed {
		panic(err)
	}
}
