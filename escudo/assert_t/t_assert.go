// Copy of the package assert but edited to be used in tests.
package assert_t

import (
	"errors"
	"os"
	"testing"
)

func NoErr(t *testing.T, err error) {
	if err != nil {
		t.Error(err.Error())
	}
}

func Closed(t *testing.T, err error) {
	if err != nil && !errors.Is(err, os.ErrClosed) {
		t.Error(err.Error())
	}
}

func FileExist(t *testing.T, err error) {
	if err != nil && !errors.Is(err, os.ErrExist) {
		t.Error(err.Error())
	}
}

func FileNotExist(t *testing.T, err error) {
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Error(err.Error())
	}
}
