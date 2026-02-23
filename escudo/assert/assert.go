// Go error handling is annoying and they don't intend to improve
// https://go.dev/blog/error-syntax
//
// - I don't want to write 3 extra lines for each line of code (annoying to read)
// - I don't want to put function call and nil check in same line (ugly as hell)
// - I can't handle the error at that function (only the user can decide)
//
// So my approach is to panic when the error is detected,
// catch the error during a defered call
// and overwrite the previous return value.
package assert

import (
	"errors"
	"os"
)

func NoErr(err error) {
	if err != nil {
		panic(err)
	}
}

func NoErrOn2(a any, err error) {
	if err != nil {
		panic(err)
	}
}

func Closed(err error) {
	if err != nil && !errors.Is(err, os.ErrClosed) {
		panic(err)
	}
}

func FileExist(err error) {
	if err != nil && !errors.Is(err, os.ErrExist) {
		panic(err)
	}
}

func FileNotExist(err error) {
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}
}
