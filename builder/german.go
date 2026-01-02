package builder

import (
	"context"
	"errors"
	"log"
	"math/rand/v2"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/justinhjy1004/sentenceminer/sampler"
	"github.com/justinhjy1004/sentenceminer/translator"
)

var translatorPort string = "localhost:50051"

var excludedWords = []string{"Tom", "Tobias"}

func Intersect[T comparable](s1, s2 []T) []T {
	set := make(map[T]struct{})
	var result []T

	for _, v := range s1 {
		set[v] = struct{}{}
	}

	for _, v := range s2 {
		if _, found := set[v]; found {
			result = append(result, v)
		}
	}
	return result
}

type Card struct {
	OriginalText string
	Translation  string
	AudioFile    string
	MaskedText   string
	AnswerText   string
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

func RemovePunctuations(input string) string {

	// Added the curly apostrophe ’ to the allowed list
	re := regexp.MustCompile(`[^\p{L}\p{N}\s'’]`)

	output := re.ReplaceAllString(input, "")

	return output

}

func ContainsExcludedWords(maskedWords []string, excluded []string) bool {
	skipSet := make(map[string]struct{})

	for _, word := range excluded {
		skipSet[word] = struct{}{}
	}

	for _, m := range maskedWords {
		_, containsExcluded := skipSet[m]

		if containsExcluded {
			return true
		}
	}

	return false
}

func MaskWords(input string, excluded []string) (string, string, error) {

	words := strings.Fields(RemovePunctuations(input))

	numWords := len(words)

	var numMask int

	switch {
	case numWords < 5:
		numMask = 1
	case numWords < 9:
		numMask = 2
	default:
		numMask = 3
	}

	maskTerms := SampleWithoutReplacement(strings.Fields(RemovePunctuations(input)), numMask)

	for i := 0; ContainsExcludedWords(maskTerms, excluded); i++ {
		maskTerms = SampleWithoutReplacement(strings.Fields(RemovePunctuations(input)), numMask)

		if i == 10000 {
			return "", "", errors.New("Loooooooped too much!")
		}
	}

	maskedText := input
	answerText := input

	for _, term := range maskTerms {

		pattern := `\b` + regexp.QuoteMeta(term) + `\b`
		re := regexp.MustCompile(pattern)

		// Replace only the first occurrence
		maskedText = re.ReplaceAllStringFunc(maskedText, func(s string) string {
			return "____"
		})

		answerText = re.ReplaceAllStringFunc(answerText, func(s string) string {
			return "<b>" + s + "</b>"
		})
	}

	return maskedText, answerText, nil

}

func NumRepetitions(input string) int {

	words := strings.Fields(input)

	numWords := len(words)

	switch {
	case numWords < 5:
		return 1
	case numWords < 9:
		return 2
	default:
		return 3
	}

}

func RepeatElements[T any](input []T, counts []int) ([]T, error) {

	if len(input) != len(counts) {
		return input, errors.New("Counts must be the same length as input!")
	}

	var output []T

	for i, val := range input {
		output = slices.Concat(output, slices.Repeat([]T{val}, counts[i]))
	}

	return output, nil

}

func GenerateCards(sample []*sampler.Sentence) []Card {

	numReps := make([]int, len(sample))

	translatorService, _ := translator.NewTranslationService(translatorPort)

	defer translatorService.Close()

	cards := make([]Card, len(sample))

	for i, s := range sample {

		translation, _ := translatorService.TranslateText(context.Background(), s.Text)

		rep := NumRepetitions(s.Text)

		card := Card{
			OriginalText: s.Text,
			Translation:  translation,
			AudioFile:    "[de-es-" + strconv.Itoa(s.ID) + ".wav]",
		}

		cards[i] = card
		numReps[i] = rep
	}

	cards, _ = RepeatElements(cards, numReps)

	for i, card := range cards {

		maskedText, answerText, err := MaskWords(card.OriginalText, excludedWords)

		if err != nil {
			log.Fatal("Tom is a problem!")
		}

		cards[i].AnswerText = answerText
		cards[i].MaskedText = maskedText
	}

	return cards
}
