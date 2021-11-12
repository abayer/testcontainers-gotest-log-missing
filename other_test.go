package repro_test

import (
	"testing"

	repro "github.com/abayer/testcontainers-gotest-log-missing"
)

func TestAddNumbersAgain(t *testing.T) {
	t.Log("just going straight to adding")
	result := repro.AddNumbers(1, 5)

	if result != 3 {
		t.Errorf("%d should have been 3", result)
	}
}
