package assert

func Catch(err *error) {
	if r := recover(); r != nil {
		*err = r.(error)
	}
}
