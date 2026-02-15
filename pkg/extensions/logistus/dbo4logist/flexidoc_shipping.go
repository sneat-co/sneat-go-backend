package dbo4logist

import (
	"github.com/sneat-co/sneat-go-core/flexidoc"
)

// ShippingDoc is a document definition for shipping
var ShippingDoc = flexidoc.DocumentDefinitionBase{
	Fields: []*flexidoc.Field{
		{
			ID: "number",
			Titles: map[string]string{
				"en": "Number",
			},
		},
		{
			ID: "number",
			Titles: map[string]string{
				"en": "Number",
			},
		},
	},
}
