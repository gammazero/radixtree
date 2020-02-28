package radixtree

import (
	"bufio"
	"os"
	"testing"
)

const (
	wordsPath = "/usr/share/dict/words"
	web2aPath = "/usr/share/dict/web2a"
)

//
// Benchmarks
//
func BenchmarkWordsRunesPut(b *testing.B) {
	benchmarkPut(wordsPath, b)
}

func BenchmarkWordsRunesGet(b *testing.B) {
	benchmarkGet(wordsPath, b)
}

func BenchmarkWordsRunesWalk(b *testing.B) {
	benchmarkWalk(wordsPath, b)
}

func BenchmarkWordsRunesWalkPath(b *testing.B) {
	benchmarkWalkPath(wordsPath, b)
}

func BenchmarkWeb2aRunesPut(b *testing.B) {
	benchmarkPut(web2aPath, b)
}

func BenchmarkWeb2aRunesGet(b *testing.B) {
	benchmarkGet(web2aPath, b)
}

func BenchmarkWeb2aRunesWalk(b *testing.B) {
	benchmarkWalk(web2aPath, b)
}

func BenchmarkWeb2aRunesWalkPath(b *testing.B) {
	benchmarkWalkPath(web2aPath, b)
}

func BenchmarkWeb2aPathsPut(b *testing.B) {
	PathSeparator = ' '
	benchmarkPathsPut(web2aPath, b)
}

func BenchmarkWeb2aPathsGet(b *testing.B) {
	benchmarkPathsGet(web2aPath, b)
}

func BenchmarkWeb2aPathsWalk(b *testing.B) {
	benchmarkPathsWalk(web2aPath, b)
}

func BenchmarkWeb2aPathsWalkPath(b *testing.B) {
	benchmarkPathsWalkPath(web2aPath, b)
}

func benchmarkPut(filePath string, b *testing.B) {
	words, err := loadWords(filePath)
	if err != nil {
		b.Skip(err.Error())
	}
	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		tree := new(Runes)
		for _, w := range words {
			tree.Put(w, w)
		}
	}
}

func benchmarkGet(filePath string, b *testing.B) {
	words, err := loadWords(filePath)
	if err != nil {
		b.Skip(err.Error())
	}
	tree := new(Runes)
	for _, w := range words {
		tree.Put(w, w)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		for _, w := range words {
			tree.Get(w)
		}
	}
}

func benchmarkWalk(filePath string, b *testing.B) {
	words, err := loadWords(filePath)
	if err != nil {
		b.Skip(err.Error())
	}
	tree := new(Runes)
	for _, w := range words {
		tree.Put(w, w)
	}
	b.ResetTimer()
	b.ReportAllocs()
	var count int
	for n := 0; n < b.N; n++ {
		count = 0
		tree.Walk("", func(k KeyStringer, value interface{}) error {
			count++
			return nil
		})
	}
	if count != len(words) {
		panic("wrong count")
	}
}

func benchmarkWalkPath(filePath string, b *testing.B) {
	words, err := loadWords(filePath)
	if err != nil {
		b.Skip(err.Error())
	}
	tree := new(Runes)
	for _, w := range words {
		tree.Put(w, w)
	}
	b.ResetTimer()
	b.ReportAllocs()
	var count int
	for n := 0; n < b.N; n++ {
		count = 0
		for _, w := range words {
			tree.WalkPath(w, func(key string, value interface{}) error {
				count++
				return nil
			})
		}
	}
	if count <= len(words) {
		panic("wrong count")
	}
}

func benchmarkPathsPut(filePath string, b *testing.B) {
	words, err := loadWords(filePath)
	if err != nil {
		b.Skip(err.Error())
	}
	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		tree := new(Paths)
		for _, w := range words {
			tree.Put(w, w)
		}
	}
}

func benchmarkPathsGet(filePath string, b *testing.B) {
	words, err := loadWords(filePath)
	if err != nil {
		b.Skip(err.Error())
	}
	tree := new(Paths)
	for _, w := range words {
		tree.Put(w, w)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		for _, w := range words {
			tree.Get(w)
		}
	}
}

func benchmarkPathsWalk(filePath string, b *testing.B) {
	words, err := loadWords(filePath)
	if err != nil {
		b.Skip(err.Error())
	}
	tree := new(Paths)
	for _, w := range words {
		tree.Put(w, w)
	}
	b.ResetTimer()
	b.ReportAllocs()
	var count int
	for n := 0; n < b.N; n++ {
		count = 0
		tree.Walk("", func(k KeyStringer, value interface{}) error {
			count++
			return nil
		})
	}
	if count != len(words) {
		panic("wrong count")
	}
}

func benchmarkPathsWalkPath(filePath string, b *testing.B) {
	words, err := loadWords(filePath)
	if err != nil {
		b.Skip(err.Error())
	}
	tree := new(Paths)
	for _, w := range words {
		tree.Put(w, w)
	}
	b.ResetTimer()
	b.ReportAllocs()
	var count int
	for n := 0; n < b.N; n++ {
		count = 0
		for _, w := range words {
			tree.WalkPath(w, func(key string, value interface{}) error {
				count++
				return nil
			})
		}
	}
	if count <= len(words) {
		panic("wrong count")
	}
}

func loadWords(wordsFile string) ([]string, error) {
	f, err := os.Open(wordsFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var word string
	var words []string

	// Scan through line-dilimited words.
	for scanner.Scan() {
		word = scanner.Text()
		words = append(words, word)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}
