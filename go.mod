module github.com/sneat-co/sneat-go-backend

go 1.22.3

//replace github.com/sneat-co/sneat-core-modules => ../sneat-core-modules
//
//replace go mod tidy => ../sneat-mod-debtus-go

//replace github.com/sneat-co/debtstracker-translations => ../debtstracker-translations
//replace github.com/sneat-co/sneat-go-core => ../sneat-go-core
//replace github.com/sneat-co/sneat-go-modules => ../sneat-go-modules
//replace github.com/bots-go-framework/bots-fw => ../../bots-go-framework/bots-fw
//replace github.com/bots-go-framework/bots-fw-telegram => ../../bots-go-framework/bots-fw-telegram
//replace github.com/bots-go-framework/bots-fw-store => ../../bots-go-framework/bots-fw-store
//replace github.com/bots-go-framework/bots-fw-telegram-models => ../../bots-go-framework/bots-fw-telegram-models
//replace github.com/bots-go-framework/bots-host-gae => ../../bots-go-framework/bots-host-gae
//replace github.com/bots-go-framework/dalgo4botsfw => ../../bots-go-framework/dalgo4botsfw
//replace github.com/dal-go/dalgo => ../../dal-go/dalgo
//replace github.com/dal-go/dalgo2firestore => ../../dal-go/dalgo2firestore
//replace github.com/strongo/app => ../../strongo/app
//replace github.com/strongo/i18n => ../../strongo/i18n
//replace github.com/strongo/strongoapp => ../../strongo/strongoapp

require (
	github.com/bots-go-framework/bots-api-telegram v0.7.2
	github.com/bots-go-framework/bots-fw v0.40.4
	github.com/bots-go-framework/bots-fw-store v0.8.2
	github.com/bots-go-framework/bots-fw-telegram v0.13.8
	github.com/bots-go-framework/bots-go-core v0.0.3
	github.com/crediterra/money v0.3.0
	github.com/dal-go/dalgo v0.14.1
	github.com/dal-go/mocks4dalgo v0.1.27
	github.com/golang/mock v1.6.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/olebedev/when v1.1.0
	github.com/sneat-co/debtstracker-translations v0.2.4
	github.com/sneat-co/sneat-core-modules v0.15.11
	github.com/sneat-co/sneat-go-core v0.37.4
	github.com/stretchr/testify v1.10.0
	github.com/strongo/decimal v0.1.1
	github.com/strongo/delaying v0.1.0
	github.com/strongo/i18n v0.6.1
	github.com/strongo/logus v0.2.0
	github.com/strongo/random v0.0.1
	github.com/strongo/slice v0.3.1
	github.com/strongo/strongoapp v0.25.4
	github.com/strongo/validation v0.0.7
)

require (
	github.com/AlekSi/pointer v1.2.0 // indirect
	github.com/alexsergivan/transliterator v1.0.1 // indirect
	github.com/bots-go-framework/bots-fw-telegram-models v0.3.8 // indirect
	github.com/bots-go-framework/bots-fw-telegram-webapp v0.3.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/gosimple/slug v1.14.0 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/pquerna/ffjson v0.0.0-20190930134022-aa0246cd15f7 // indirect
	github.com/strongo/facebook v1.8.1 // indirect
	github.com/strongo/gamp v0.0.1 // indirect
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/net v0.32.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
