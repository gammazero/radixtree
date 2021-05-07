package radixtree

import (
	"testing"
)

// Test splitting a path keys into segemnts, /some/file => /some, /file
func TestPathNext(t *testing.T) {
	cases := []struct {
		key     string
		parts   []string
		indexes []int
	}{
		{"", []string{""}, []int{-1}},
		{"/", []string{""}, []int{-1}},
		{"//", []string{""}, []int{-1}},
		{"/a/b/c", []string{"a", "b", "c"}, []int{3, 5, -1}},
		{"some_file", []string{"some_file"}, []int{-1}},
		{"/home/aripley", []string{"home", "aripley"}, []int{6, -1}},
		{"home/ntaylor", []string{"home", "ntaylor"}, []int{5, -1}},
		{"/home//jwaits/", []string{"home", "jwaits"}, []int{7, -1}},
		{"home/dverlaine/", []string{"home", "dverlaine"}, []int{5, -1}},
	}

	for _, c := range cases {
		var partNum int
		for part, next := pathNext(c.key, 0); ; part, next = pathNext(c.key, next) {
			if part != c.parts[partNum] {
				t.Errorf("expected part %d of key %q to be %q, got '%s'", partNum, c.key, c.parts[partNum], part)
			}
			if next != c.indexes[partNum] {
				t.Errorf("at split %d, expected next index of key %q to be %d, got %d", partNum, c.key, c.indexes[partNum], next)
			}
			partNum++
			if next == -1 {
				break
			}
		}
		if partNum != len(c.parts) {
			t.Errorf("expected %q to have %d parts, got %d", c.key, len(c.parts), partNum)
		}
	}
}

func TestPathNextMultichar(t *testing.T) {
	PathSeparator = "--"
	defer func() {
		PathSeparator = "/"
	}()

	cases := []struct {
		key     string
		parts   []string
		indexes []int
	}{
		{"", []string{""}, []int{-1}},
		{"--", []string{""}, []int{-1}},
		{"----", []string{""}, []int{-1}},
		{"--a--b--c", []string{"a", "b", "c"}, []int{5, 8, -1}},
		{"some_file", []string{"some_file"}, []int{-1}},
		{"--home--aripley", []string{"home", "aripley"}, []int{8, -1}},
		{"home--ntaylor", []string{"home", "ntaylor"}, []int{6, -1}},
		{"--home----jwaits--", []string{"home", "jwaits"}, []int{10, -1}},
		{"home--dverlaine--", []string{"home", "dverlaine"}, []int{6, -1}},
	}

	for _, c := range cases {
		var partNum int
		for part, next := pathNext(c.key, 0); ; part, next = pathNext(c.key, next) {
			if part != c.parts[partNum] {
				t.Errorf("expected part %d of key %q to be %q, got '%s'", partNum, c.key, c.parts[partNum], part)
			}
			if next != c.indexes[partNum] {
				t.Errorf("at split %d, expected next index of key %q to be %d, got %d", partNum, c.key, c.indexes[partNum], next)
			}
			partNum++
			if next == -1 {
				break
			}
		}
		if partNum != len(c.parts) {
			t.Errorf("expected %q to have %d parts, got %d", c.key, len(c.parts), partNum)
		}
	}
}

func TestPathNextBeginEnd(t *testing.T) {
	cases := []struct {
		path  string
		start int
		part  string
		next  int
	}{
		{"", 0, "", -1},
		{"", 100, "", -1},

		{" /", 0, " ", -1},
		{" /", 1, "", -1},

		{"/", 0, "", -1},
		{"/", 1, "", -1},
		{"/", 100, "", -1},
		{"/", -100, "", -1},

		{"//", 0, "", -1},
		{"//", 1, "", -1},
		{"//", 2, "", -1},

		{"///", 0, "", -1},
		{"///", 1, "", -1},
		{"///", 2, "", -1},
		{"///", 3, "", -1},
	}

	for _, c := range cases {
		part, next := pathNext(c.path, c.start)
		if part != c.part {
			t.Errorf("expected part %q starting at %d in path %q, got %q", c.part, c.start, c.path, part)
		}
		if next != c.next {
			t.Errorf("expected next %d starting at %d in path %q, got %d", c.next, c.start, c.path, next)
		}
	}
}
