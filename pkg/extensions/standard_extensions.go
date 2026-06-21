package extensions

import (
	"github.com/sneat-co/assetus/backend/assetusext"
	"github.com/sneat-co/calendarius/backend/calendariusext"
	"github.com/sneat-co/contactus/backend/contactusext"
	"github.com/sneat-co/listus/backend/listusext"
	sneatcoremodules "github.com/sneat-co/sneat-core-modules"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/deedus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/retrospectus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/scrumus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/sportus"
	"github.com/sneat-co/sneat-go-core/extension"
)

func standardExtensions() []extension.Config {
	return []extension.Config{
		assetusext.Extension(),
		contactusext.Extension(),
		deedus.Extension(),
		calendariusext.Extension(),
		listusext.Extension(),
		retrospectus.Extension(),
		scrumus.Extension(),
		sportus.Extension(),
	}
}

func Extensions() []extension.Config {
	return append(
		sneatcoremodules.CoreExtensions(),
		standardExtensions()...,
	)
}
