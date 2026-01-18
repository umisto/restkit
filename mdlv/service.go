package mdlv

import (
	"github.com/netbill/logium"
)

type Service struct {
	skUser string
	ctxKey any

	log logium.Logger
}

func New(skUser string, ctxKey any, log logium.Logger) Service {
	return Service{
		skUser: skUser,
		ctxKey: ctxKey,
		log:    log,
	}
}
