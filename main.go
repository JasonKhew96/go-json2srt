package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/asticode/go-astisub"
)

// json subtitle body
type JsonSubBody struct {
	From     float64 `json:"from"`     // start time (seconds)
	To       float64 `json:"to"`       // end time (seconds)
	Location int     `json:"location"` // subtitle location
	Content  string  `json:"content"`  // subtitle content
}

// json subtitle
type JsonSub struct {
	FontSize        float64       `json:"font_size"`        // font size
	FontColor       string        `json:"font_color"`       // font color in hex
	BackgroundAlpha float64       `json:"background_alpha"` // font transparency
	BackgroundColor string        `json:"background_color"` // font background color
	Stroke          string        `json:"Stroke"`           // font stroke
	Body            []JsonSubBody `json:"body"`             // json subtile body
}

// json2srt converts json subtitle into srt subtitle
func json2srt(inputPath, outputPath string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	newSub := astisub.NewSubtitles()

	index := 1

	subtitle := &JsonSub{}
	err = json.Unmarshal(data, subtitle)
	if err != nil {
		return err
	}

	var items []*astisub.Item

	body := subtitle.Body
	for _, value := range body {
		from := value.From
		content := value.Content
		to := value.To

		line := []astisub.Line{
			{
				Items: []astisub.LineItem{
					{
						Text: content,
					},
				},
			},
		}

		item := &astisub.Item{
			Index:   index,
			EndAt:   time.Duration(to * float64(time.Second)),
			Lines:   line,
			StartAt: time.Duration(from * float64(time.Second)),
		}
		items = append(items, item)

		index++
	}

	newSub.Items = items

	f, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	return newSub.WriteToSRT(f)
}

// any2json converts srt subtitle into json subtitle
func any2json(inputPath, outputPath string) error {
	sub, err := astisub.OpenFile(inputPath)
	if err != nil {
		return err
	}
	jsonSubtitle := &JsonSub{
		FontSize:        0.4,
		FontColor:       "#FFFFFF",
		BackgroundAlpha: 0.5,
		BackgroundColor: "#9C27B0",
		Stroke:          "none",
	}
	for _, item := range sub.Items {
		var lines []string
		from := item.StartAt.Seconds()
		to := item.EndAt.Seconds()
		for _, line := range item.Lines {
			for _, lineItem := range line.Items {
				lines = append(lines, lineItem.Text)
			}
		}
		newLine := strings.Join(lines, "\n")
		jsonSubtitle.Body = append(jsonSubtitle.Body, JsonSubBody{
			From:     from,
			To:       to,
			Location: 2,
			Content:  newLine,
		})
	}

	data, err := json.Marshal(jsonSubtitle)
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, data, 0666)
}

// processJsonSub process json subtitle
func processJsonSub(inputPath, outputPath string) error {
	outputExt := strings.ToLower(filepath.Ext(outputPath))
	switch outputExt {
	case ".srt":
		return json2srt(inputPath, outputPath)
	default:
		return fmt.Errorf("\"%s\" file type is not supported", outputExt)
	}
}

// processSrtSub process srt subtitle
func processSrtSub(inputPath, outputPath string) error {
	outputExt := strings.ToLower(filepath.Ext(outputPath))
	switch outputExt {
	case ".json":
		return any2json(inputPath, outputPath)
	default:
		return fmt.Errorf("\"%s\" file type is not supported", outputExt)
	}
}

func main() {
	arguments := os.Args[1:]
	argsLen := len(arguments)
	if argsLen != 2 {
		fmt.Printf("Expected 2 arguments, but only found %d\nExample: convert input.json output.srt\n", argsLen)
		return
	}
	inputPath := arguments[0]
	outputPath := arguments[1]

	var err error

	inputExt := strings.ToLower(filepath.Ext(inputPath))
	switch inputExt {
	case ".json":
		err = processJsonSub(inputPath, outputPath)
	case ".srt", ".ass":
		err = processSrtSub(inputPath, outputPath)
	default:
		fmt.Printf("\"%s\" file type is not supported\n", inputExt)
		return
	}
	if err != nil {
		fmt.Println(err)
		return
	}
}
