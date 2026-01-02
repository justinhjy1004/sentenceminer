package tts

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/justinhjy1004/sentenceminer/sampler"
)

func GenerateGermanSpeechAudio(directory string, sentences []*sampler.Sentence) {

	err := os.MkdirAll(directory, 0755)
	if err != nil {
		log.Fatalf("Error creating directory: %v\n", err)
		return
	}

	for _, s := range sentences {
		GermanTextToSpeech(directory, strconv.Itoa(s.ID), s.Text)
	}

	return

}

func GermanTextToSpeech(directory string, id string, text string) {

	cmd := exec.Command("piper",
		"--model", "./tts/de-thorsten-high.onnx",
		"--length_scale", "1.2",
		"--output-file", fmt.Sprintf("%s/de-es-%s.wav", directory, id),
	)

	stdin, err := cmd.StdinPipe()

	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Start()

	if err != nil {
		log.Fatal(err)
	}

	io.WriteString(stdin, text)
	stdin.Close()

	err = cmd.Wait()

	if err != nil {
		log.Fatal(err)
	}

}
