package helpers

import (
	"github.com/rs/zerolog/log"
)

func PanicError(err error, msg string) {
	if err != nil {
		log.Error().Err(err).Msg(msg)
		panic(err)
	}
}
