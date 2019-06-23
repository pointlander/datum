// Copyright 2019 The Datum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iris

import (
	"testing"
)

func TestIris(t *testing.T) {
	datum, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	if value := len(datum.Fisher); value != 150 {
		t.Fatalf("length of fisher should be 150 but is %d", value)
	}
	if value := len(datum.Bezdek); value != 150 {
		t.Fatalf("length of bezdek should be 150 but is %d", value)
	}
}
