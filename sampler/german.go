package sampler

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"os"

	"github.com/gocarina/gocsv"
)

type Sentence struct {
	ID       int    `csv:"0"`
	Language string `csv:"1"`
	Text     string `csv:"2"`
}

func SampleWithoutReplacement[T any](input []T, k int) []T {

	n := len(input)

	if k > n {
		log.Fatalf("Input slice is only %d long but expecting output of length %d!", n, k)
	}

	set := make(map[int]struct{})
	result := make([]T, 0, k)

	for len(result) < k {
		r := rand.IntN(n)
		if _, exists := set[r]; !exists {
			set[r] = struct{}{}
			result = append(result, input[r])
		}
	}

	return result

}

func SampleGermanSentence(k int) []*Sentence {

	file, _ := os.Open("./sampler/deu_sentences.tsv")

	var sentences []*Sentence

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = '\t'
		r.LazyQuotes = true
		return r
	})

	err := gocsv.UnmarshalWithoutHeaders(file, &sentences)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	sampledSentences := SampleWithoutReplacement(sentences, k)

	return sampledSentences

}
