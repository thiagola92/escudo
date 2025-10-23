package assert

var Err error

func Catch() {
	if r := recover(); r != nil {
		Err = r.(error)
	}
}
