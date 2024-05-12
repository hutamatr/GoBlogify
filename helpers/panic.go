package helpers

func PanicError(err error, msg string) {
	if err != nil {
		LogError("%v : %s", err.Error(), msg)
		panic(err)
	}
}
