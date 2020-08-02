// Copyright 2019 The Datum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bible

import (
	"bytes"
	"compress/bzip2"
	"errors"
	"io/ioutil"
	"regexp"
	"strings"
)

const (
	// NumberOfVerses is the number of verses in the bible
	NumberOfVerses = 31102
)

var (
	// PatternBook marks the start of a book
	PatternBook = regexp.MustCompile(`\r\n\r\n\r\n\r\n[A-Za-z]+([ \t]+[A-Za-z:]+)*\r\n\r\n`)
	// PatternVerse is a verse
	PatternVerse = regexp.MustCompile(`\d+[:]\d+[A-Za-z:.,?!;"' ()\t\r\n]+`)
	// PatternSentence is a sentence
	PatternSentence = regexp.MustCompile(`[.,?!;]`)
	// PatternWord is for splitting into words
	PatternWord = regexp.MustCompile(`[ \t\r\n]+`)
	// WordCutSet is the trim cut set for a word
	WordCutSet = ".,?!:;\"' ()\t\r\n"
)

// Bible is a bible
type Bible []Testament

// Testament is a bible testament
type Testament struct {
	Name  string
	Books []Book
}

// Book is a book of the bible
type Book struct {
	Name   string
	Verses []Verse
}

// Verse is a bible verse
type Verse struct {
	Testament string
	Book      string
	Number    string
	Verse     string
	Sentences []string
	Words     []string
}

// LoadBible returns the bible
func Load() (Bible, error) {
	reader := bzip2.NewReader(bytes.NewReader(AssetBible))
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	bible := string(data)
	beginning := strings.Index(bible, "*** START OF THIS PROJECT GUTENBERG EBOOK THE KING JAMES BIBLE ***")
	ending := strings.Index(bible, "End of the Project Gutenberg EBook of The King James Bible")
	bible = bible[beginning:ending]
	testaments := make([]Testament, 2)
	testaments[0].Name = "The Old Testament of the King James Version of the Bible"
	testaments[1].Name = "The New Testament of the King James Bible"
	countVerses := 0

	a := strings.Index(bible, testaments[0].Name)
	b := strings.Index(bible, testaments[1].Name)
	parse := func(t *Testament, testament string) {
		books := PatternBook.FindAllStringIndex(testament, -1)
		for i, book := range books {
			b := Book{
				Name: strings.TrimSpace(testament[book[0]:book[1]]),
			}
			end := len(testament)
			if i+1 < len(books) {
				end = books[i+1][0]
			}
			content := testament[book[1]:end]
			lines := PatternVerse.FindAllStringIndex(content, -1)
			for _, line := range lines {
				l := strings.TrimSpace(strings.ReplaceAll(content[line[0]:line[1]], "\r\n", " "))
				a := strings.Index(l, " ")
				verse := Verse{
					Testament: t.Name,
					Book:      b.Name,
					Number:    strings.TrimSpace(l[:a]),
					Verse:     strings.TrimSpace(l[a:]),
				}

				verseSentences := PatternSentence.Split(verse.Verse, -1)
				for _, sentence := range verseSentences {
					sentence = strings.Trim(sentence, WordCutSet)
					if len(sentence) == 0 {
						continue
					}
					verse.Sentences = append(verse.Sentences, sentence)
				}

				verseWords := PatternWord.Split(verse.Verse, -1)
				for _, word := range verseWords {
					word = strings.Trim(word, WordCutSet)
					if len(word) == 0 {
						continue
					}
					verse.Words = append(verse.Words, word)
				}

				countVerses++
				b.Verses = append(b.Verses, verse)
			}
			t.Books = append(t.Books, b)
		}
	}
	parse(&testaments[0], bible[a:b])
	parse(&testaments[1], bible[b:])

	if countVerses != NumberOfVerses {
		return nil, errors.New("wrong number of verses")
	}

	return testaments, nil
}

// GetVerses gets all of the verses in the bible
func (b Bible) GetVerses() (verses []Verse) {
	for _, testament := range b {
		for _, book := range testament.Books {
			verses = append(verses, book.Verses...)
		}
	}
	return verses
}

// GetSentences gets all of the sentences in the bible
func (b Bible) GetSentences() (sentences []string) {
	for _, testament := range b {
		for _, book := range testament.Books {
			for _, verse := range book.Verses {
				sentences = append(sentences, verse.Sentences...)
			}
		}
	}
	return sentences
}

// GetWords gets all of the words in the bible
func (b Bible) GetWords() (words []string) {
	seen := make(map[string]bool, 8)
	for _, testament := range b {
		for _, book := range testament.Books {
			for _, verse := range book.Verses {
				for _, word := range verse.Words {
					if seen[word] {
						continue
					}
					seen[word] = true
					words = append(words, word)
				}
			}
		}
	}
	return words
}
