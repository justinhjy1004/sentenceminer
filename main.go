package main

import (
	"os"

	"github.com/gocarina/gocsv"
	"github.com/justinhjy1004/sentenceminer/builder"
	"github.com/justinhjy1004/sentenceminer/sampler"
	"github.com/justinhjy1004/sentenceminer/tts"
)

func main() {

	numSample := 3
	maxWords := 10

	sample := sampler.SampleGermanSentence(numSample, maxWords)

	tts.GenerateGermanSpeechAudio("audio", sample)

	cards := builder.GenerateCards(sample)

	file, _ := os.Create("cards.csv")
	gocsv.MarshalFile(&cards, file)

}
