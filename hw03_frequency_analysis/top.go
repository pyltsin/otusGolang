package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"regexp"
	"sort"
)

type wordCounter struct {
	word  string
	count int
}

var reg = regexp.MustCompile("[[:space:]]")

func Top10(text string) []string {
	if len(text) == 0 {
		return []string{}
	}
	wordsBag := createWordsBag(text)
	allWords := createWordCounters(wordsBag)
	sortCounters(allWords)
	return getTop10(allWords)
}

func createWordsBag(text string) map[string]wordCounter {
	var countWords = make(map[string]wordCounter)

	for _, word := range reg.Split(text, -1) {
		if len(word) == 0 {
			continue
		}

		counter, ok := countWords[word]
		if !ok {
			counter = wordCounter{
				word:  word,
				count: 0,
			}
		}
		counter.count++
		countWords[word] = counter
	}
	return countWords
}

func createWordCounters(wordsBag map[string]wordCounter) []wordCounter {
	var allWords = make([]wordCounter, 0, len(wordsBag))
	for _, v := range wordsBag {
		allWords = append(allWords, v)
	}
	return allWords
}

func sortCounters(allWords []wordCounter) {
	sort.SliceStable(allWords, func(i, j int) bool {
		return allWords[i].count > allWords[j].count
	})
}

func getTop10(allWords []wordCounter) []string {
	result := make([]string, 0, 10)
	upperCase := min(10, len(allWords))
	for _, counter := range allWords[0:upperCase] {
		result = append(result, counter.word)
	}
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
