package radixtree

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
)

const (
	// web2: Webster's Second International Dictionary, all 234,936 words worth.
	web2URL  = "https://raw.githubusercontent.com/openbsd/src/master/share/dict/web2"
	web2Path = "web2"
	// web2a: hyphenated terms as well as assorted noun and adverbial
	// phrasesfrom Webster's Second International Dictionary.
	web2aURL  = "https://raw.githubusercontent.com/openbsd/src/master/share/dict/web2a"
	web2aPath = "web2a"
)

func BenchmarkGet(b *testing.B) {
	err := getWords()
	if err != nil {
		b.Skip(err.Error())
	}

	b.Run("Words", func(b *testing.B) {
		benchmarkGet(b, web2Path)
	})

	b.Run("Web2a", func(b *testing.B) {
		benchmarkGet(b, web2aPath)
	})
}

func BenchmarkPut(b *testing.B) {
	b.Run("Words", func(b *testing.B) {
		benchmarkPut(b, web2Path)
	})

	b.Run("Web2a", func(b *testing.B) {
		benchmarkPut(b, web2aPath)
	})
}

func BenchmarkWalk(b *testing.B) {
	b.Run("Words", func(b *testing.B) {
		benchmarkWalk(b, web2Path)
	})

	b.Run("Web2a", func(b *testing.B) {
		benchmarkWalk(b, web2aPath)
	})
}

func BenchmarkWalkPath(b *testing.B) {
	b.Run("Words", func(b *testing.B) {
		benchmarkWalkPath(b, web2Path)
	})

	b.Run("Web2a", func(b *testing.B) {
		benchmarkWalkPath(b, web2aPath)
	})
}

func benchmarkGet(b *testing.B, filePath string) {
	words, err := loadWords(filePath)
	if err != nil {
		b.Skip(err.Error())
	}
	tree := new(Tree)
	for _, w := range words {
		tree.Put(w, w)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		for _, w := range words {
			if _, ok := tree.Get(w); !ok {
				panic("missing value")
			}
		}
	}
}

func benchmarkPut(b *testing.B, filePath string) {
	words, err := loadWords(filePath)
	if err != nil {
		b.Skip(err.Error())
	}
	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		tree := new(Tree)
		for _, w := range words {
			tree.Put(w, w)
		}
	}
}

func benchmarkWalk(b *testing.B, filePath string) {
	words, err := loadWords(filePath)
	if err != nil {
		b.Skip(err.Error())
	}
	tree := new(Tree)
	for _, w := range words {
		tree.Put(w, w)
	}
	b.ResetTimer()
	b.ReportAllocs()
	var count int
	for n := 0; n < b.N; n++ {
		count = 0
		tree.Walk("", func(k string, value any) bool {
			count++
			return false
		})
	}
	if count != len(words) {
		b.Fatalf("Walk wrong count, expected %d got %d", len(words), count)
	}
}

func benchmarkWalkPath(b *testing.B, filePath string) {
	words, err := loadWords(filePath)
	if err != nil {
		b.Skip(err.Error())
	}
	tree := new(Tree)
	for _, w := range words {
		tree.Put(w, w)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		found := false
		for _, w := range words {
			tree.WalkPath(w, func(key string, value any) bool {
				found = true
				return false
			})
		}
		if !found {
			b.Fatal("Walk did not find word")
		}
	}
}

func loadWords(wordsFile string) ([]string, error) {
	f, err := os.Open(wordsFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var words []string

	// Scan through line-dilimited words.
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

func getWords() error {
	err := downloadFile(web2URL, web2Path)
	if err != nil {
		return err
	}
	return downloadFile(web2aURL, web2aPath)
}

func downloadFile(fileURL, filePath string) error {
	_, err := os.Stat(filePath)
	if err == nil {
		return nil
	}
	rsp, err := http.Get(fileURL)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != 200 {
		return fmt.Errorf("error response getting file: %d", rsp.StatusCode)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, rsp.Body)
	return err
}
