package main

import (
	"bytes"
	"testing"
)

func TestCountWords(t *testing.T) {
	buffer := bytes.NewBufferString("word1 word2 word3 word4\n")

	exp := 4
	actual := count(buffer, false)

	if exp != actual {
		t.Errorf("Expected word count %d, but got %d instead", exp, actual)
	}
}

func TestCountLines(t *testing.T) {
	buffer := bytes.NewBufferString("hello world\nhow are you?")

	expected := 2
	actual := count(buffer, true)

	if expected != actual {
		t.Errorf("Expected word count %d, but got %d instead", expected, actual)
	}
}
