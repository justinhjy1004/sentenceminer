package main

import (
	"context"
	"fmt"

	"github.com/justinhjy1004/sentenceminer/sampler"
	"github.com/justinhjy1004/sentenceminer/translator"
)

func main() {
	sample := sampler.SampleGermanSentence(10)

	translatorService, _ := translator.NewTranslationService("localhost:50051")

	defer translatorService.Close()

	for _, s := range sample {
		fmt.Println(s.Text)
		t, _ := translatorService.TranslateText(context.Background(), s.Text)
		fmt.Println(t)
	}
}
