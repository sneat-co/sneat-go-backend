package sneatteamgo

import (
	"github.com/sneat-co/sneat-go-backend/src/sneatgae/sneatgaeapp"
)

var start = sneatgaeapp.Start

func init() {
	if start == nil {
		panic("Start is not defined")
	}
}
