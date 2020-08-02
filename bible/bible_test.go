// Copyright 2019 The Datum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bible

import (
	"testing"
)

func TestBible(t *testing.T) {
	datum, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	if len(datum.GetVerses()) != 31102 {
		t.Fatal("invalid number of verses")
	}
	if len(datum.GetSentences()) != 112801 {
		t.Fatal("invalid number of sentences")
	}
	if len(datum.GetWords()) != 13749 {
		t.Fatal("invalid number of words")
	}
}
