package facade4brands

import (
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/brandus/dbo4brands"
)

var autoMakers = map[string]dbo4brands.Maker{
	"audi": {
		Title: "Audi",
	},
	"bmw": {
		Title: "BMW",
	},
}
