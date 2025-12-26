package tts

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

func GermanTextToSpeech(directory string, id string, text string) {

	err := os.MkdirAll(directory, 0755)
	if err != nil {
		log.Fatalf("Error creating directory: %v\n", err)
		return
	}

	cmd := exec.Command("piper",
		"--model", "./tts/de-thorsten-high.onnx",
		"--length_scale", "1.2",
		"--output-file", fmt.Sprintf("%s/%s.wav", directory, id),
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
