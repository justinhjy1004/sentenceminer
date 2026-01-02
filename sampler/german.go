package sampler

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"os"
	"slices"
	"strings"

	"github.com/gocarina/gocsv"
)

type Sentence struct {
	ID       int    `csv:"0"`
	Language string `csv:"1"`
	Text     string `csv:"2"`
}

var sentenceFile string = "./sampler/deu_sentences.tsv"

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

func LoadGermanSentenceFile(maxWords int) []*Sentence {

	file, _ := os.Open(sentenceFile)

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

	if maxWords != -1 {

		sentences = slices.DeleteFunc(sentences, func(s *Sentence) bool {
			return len(strings.Fields(s.Text)) > maxWords
		})
	}

	return sentences

}

func SampleGermanSentence(k int, maxWords int) []*Sentence {

	sentences := LoadGermanSentenceFile(maxWords)

	sampledSentences := SampleWithoutReplacement(sentences, k)

	return sampledSentences

}
