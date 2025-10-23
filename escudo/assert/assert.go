// Go error handling is annoying and they don't intend to improve
// https://go.dev/blog/error-syntax
//
// - I don't want to write 3 extra lines for each line of code
// - I don't want to put function call and nil check in same line
// - I don't want to handle the error at that function
//
// So my approach is to make a little more like C,
// where it store the last error into a variable
// and the user has to check the variable.
//
// - This go against Go guidelines
// - This is not thread safe!
// - This is my package so fuck it...
package assert

import (
	"errors"
	"os"
)

func NoErr(err error) {
	Err = nil

	if err != nil {
		panic(err)
	}
}

func NoErrOn2(a any, err error) {
	Err = nil

	if err != nil {
		panic(err)
	}
}

func Closed(err error) {
	Err = nil

	if err != nil && !errors.Is(err, os.ErrClosed) {
		panic(err)
	}
}

func FileNotExist(err error) {
	Err = nil

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}
}
