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
func BenchmarkWordsBytesGet(b *testing.B) {
	benchmarkBytesGet(wordsPath, b)
}

func BenchmarkWordsBytesPut(b *testing.B) {
	benchmarkBytesPut(wordsPath, b)
}

func BenchmarkWordsBytesWalk(b *testing.B) {
	benchmarkBytesWalk(wordsPath, b)
}

func BenchmarkWordsBytesWalkPath(b *testing.B) {
	benchmarkBytesWalkPath(wordsPath, b)
}

// ----- Web2a -----
func BenchmarkWeb2aBytesGet(b *testing.B) {
	benchmarkBytesGet(web2aPath, b)
}

func BenchmarkWeb2aBytesPut(b *testing.B) {
	benchmarkBytesPut(web2aPath, b)
}

func BenchmarkWeb2aBytesWalk(b *testing.B) {
	benchmarkBytesWalk(web2aPath, b)
}

func BenchmarkWeb2aBytesWalkPath(b *testing.B) {
	benchmarkBytesWalkPath(web2aPath, b)
}

func BenchmarkWeb2aPathsPut(b *testing.B) {
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

func benchmarkMapToCompareWithPut(filePath string, b *testing.B) {
	words, err := loadWords(filePath)
	if err != nil {
		b.Skip(err.Error())
	}

	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		m := map[string]string{}
		for _, w := range words {
			m[w] = w
		}
		if len(m) != len(words) {
			panic("wrong size map")
		}
	}
}

func benchmarkBytesPut(filePath string, b *testing.B) {
	words, err := loadWords(filePath)
	if err != nil {
		b.Skip(err.Error())
	}
	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		tree := new(Bytes)
		for _, w := range words {
			tree.Put(w, w)
		}
	}
}

func benchmarkBytesGet(filePath string, b *testing.B) {
	words, err := loadWords(filePath)
	if err != nil {
		b.Skip(err.Error())
	}
	tree := new(Bytes)
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

func benchmarkBytesWalk(filePath string, b *testing.B) {
	words, err := loadWords(filePath)
	if err != nil {
		b.Skip(err.Error())
	}
	tree := new(Bytes)
	for _, w := range words {
		tree.Put(w, w)
	}
	b.ResetTimer()
	b.ReportAllocs()
	var count int
	for n := 0; n < b.N; n++ {
		count = 0
		tree.Walk("", func(k string, value interface{}) bool {
			count++
			return false
		})
	}
	if count != len(words) {
		panic("wrong count")
	}
}

func benchmarkBytesWalkPath(filePath string, b *testing.B) {
	words, err := loadWords(filePath)
	if err != nil {
		b.Skip(err.Error())
	}
	tree := new(Bytes)
	for _, w := range words {
		tree.Put(w, w)
	}
	b.ResetTimer()
	b.ReportAllocs()
	var count int
	for n := 0; n < b.N; n++ {
		count = 0
		for _, w := range words {
			tree.WalkPath(w, func(key string, value interface{}) bool {
				count++
				return false
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
		tree.Walk("", func(k string, value interface{}) bool {
			count++
			return false
		})
		if count != len(words) {
			panic("wrong count")
		}
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
			tree.WalkPath(w, func(key string, value interface{}) bool {
				count++
				return false
			})
		}
		if count < len(words) {
			panic("wrong count")
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
