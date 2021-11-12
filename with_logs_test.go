package repro_test

import (
	repro "github.com/abayer/testcontainers-gotest-log-missing"
	"testing"
	"time"
)

func TestAddNumbers(t *testing.T) {
	t.Log("Sleeping for 15 seconds before adding numbers")
	time.Sleep(15*time.Second)
	t.Log("Done sleeping, let's add")
	result := repro.AddNumbers(1, 2)

	if result != 3 {
		t.Errorf("%d should have been 3", result)
	}
}
