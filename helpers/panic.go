package helpers

import (
	"log"

	"github.com/fatih/color"
)

func PanicError(err error) {
	errorColor := color.New(color.FgBlack).Add(color.BgRed).SprintfFunc()
	outputColor := color.New(color.FgRed).Add(color.BgBlack).SprintfFunc()
	if err != nil {
		log.Printf("%s%v", errorColor("[ERROR] An error occurred-> "), outputColor(err.Error()))
		panic(err)
	}
}
