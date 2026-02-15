package extensions

import (
	sneatcoremodules "github.com/sneat-co/sneat-core-modules"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/assetus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarium"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/deedus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/listus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/retrospectus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/scrumus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/sportus"
	"github.com/sneat-co/sneat-go-core/extension"
)

func standardExtensions() []extension.Config {
	return []extension.Config{
		assetus.Extension(),
		deedus.Extension(),
		calendarium.Extension(),
		listus.Extension(),
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
