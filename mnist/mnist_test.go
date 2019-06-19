package mnist

import (
	"testing"
)

func TestMNIST(t *testing.T) {
	datum, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	if size := datum.Train.Width; size != Width {
		t.Fatalf("train width should be 28 but is %d", size)
	}
	if size := datum.Train.Height; size != Height {
		t.Fatalf("train height should be 28 but is %d", size)
	}
	if length := len(datum.Train.Images); length != 60000 {
		t.Fatalf("train images length should be 60000 but is %d", length)
	}
	if length := len(datum.Train.Labels); length != 60000 {
		t.Fatalf("train labels length should be 60000 but is %d", length)
	}

	if size := datum.Test.Width; size != Width {
		t.Fatalf("test width should be 28 but is %d", size)
	}
	if size := datum.Test.Height; size != Height {
		t.Fatalf("test height should be 28 but is %d", size)
	}
	if length := len(datum.Test.Images); length != 10000 {
		t.Fatalf("test images length should be 10000 but is %d", length)
	}
	if length := len(datum.Test.Labels); length != 10000 {
		t.Fatalf("test labels length should be 10000 but is %d", length)
	}
}
