module github.com/sneat-co/sneat-go-backend

go 1.23.0

//replace github.com/sneat-co/sneat-core-modules => ../sneat-core-modules
//replace github.com/sneat-co/sneat-go-core => ../sneat-go-core

//replace github.com/sneat-co/debtstracker-translations => ../debtstracker-translations
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
	github.com/crediterra/money v0.3.0
	github.com/dal-go/dalgo v0.18.1
	github.com/dal-go/mocks4dalgo v0.2.3
	github.com/julienschmidt/httprouter v1.3.0
	github.com/sneat-co/sneat-core-modules v0.24.28
	github.com/sneat-co/sneat-go-core v0.47.6
	github.com/stretchr/testify v1.10.0
	github.com/strongo/decimal v0.1.1
	github.com/strongo/delaying v0.1.0
	github.com/strongo/logus v0.2.1
	github.com/strongo/random v0.0.1
	github.com/strongo/slice v0.3.1
	github.com/strongo/strongoapp v0.26.6
	github.com/strongo/validation v0.0.7
	go.uber.org/mock v0.5.0
)

require (
	github.com/alexsergivan/transliterator v1.0.1 // indirect
	github.com/bots-go-framework/bots-fw-store v0.10.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gosimple/slug v1.15.0 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/net v0.37.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
