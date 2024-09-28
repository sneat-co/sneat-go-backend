module github.com/sneat-co/sneat-go-backend

go 1.22.3

//replace github.com/sneat-co/debtstracker-translations => ../debtstracker-translations

//replace github.com/sneat-co/sneat-go-core => ../sneat-go-core

//replace github.com/sneat-co/sneat-go-firebase => ../sneat-go-firebase

//replace github.com/sneat-co/sneat-core-modules => ../sneat-core-modules
//replace github.com/sneat-co/sneat-go-modules => ../sneat-go-modules

//replace github.com/bots-go-framework/bots-fw => ../../bots-go-framework/bots-fw

//replace github.com/bots-go-framework/bots-fw-telegram => ../../bots-go-framework/bots-fw-telegram

//
//replace github.com/bots-go-framework/bots-fw-store => ../../bots-go-framework/bots-fw-store
//

//replace github.com/bots-go-framework/bots-fw-telegram-models => ../../bots-go-framework/bots-fw-telegram-models
//replace github.com/bots-go-framework/bots-host-gae => ../../bots-go-framework/bots-host-gae
//replace github.com/bots-go-framework/dalgo4botsfw => ../../bots-go-framework/dalgo4botsfw

//replace github.com/dal-go/dalgo => ../../dal-go/dalgo

//replace github.com/dal-go/dalgo2firestore => ../../dal-go/dalgo2firestore

//replace github.com/strongo/app => ../../strongo/app
//replace github.com/strongo/i18n => ../../strongo/i18n
//replace github.com/strongo/strongoapp => ../../strongo/strongoapp

require (
	github.com/bots-go-framework/bots-api-telegram v0.7.1
	github.com/bots-go-framework/bots-fw v0.40.2
	github.com/bots-go-framework/bots-fw-store v0.8.2
	github.com/bots-go-framework/bots-fw-telegram v0.13.5
	github.com/bots-go-framework/bots-fw-telegram-models v0.3.6
	github.com/bots-go-framework/bots-fw-telegram-webapp v0.3.1
	github.com/bots-go-framework/bots-go-core v0.0.3
	github.com/crediterra/go-interest v0.0.0-20180510115340-54da66993b85
	github.com/crediterra/money v0.3.0
	github.com/dal-go/dalgo v0.14.0
	github.com/dal-go/mocks4dalgo v0.1.25
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/golang/mock v1.6.0
	github.com/gosimple/slug v1.14.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/olebedev/when v1.0.0
	github.com/pquerna/ffjson v0.0.0-20190930134022-aa0246cd15f7
	github.com/sanity-io/litter v1.5.5
	github.com/sendgrid/sendgrid-go v3.16.0+incompatible
	github.com/shiyanhui/hero v0.0.2
	github.com/sneat-co/debtstracker-translations v0.2.2
	github.com/sneat-co/sneat-go-core v0.36.0
	github.com/stretchr/testify v1.9.0
	github.com/strongo/decimal v0.1.1
	github.com/strongo/delaying v0.0.1
	github.com/strongo/facebook v1.8.1
	github.com/strongo/gamp v0.0.1
	github.com/strongo/gotwilio v0.0.0-20160123000810-f024bbefe80f
	github.com/strongo/i18n v0.6.1
	github.com/strongo/logus v0.2.0
	github.com/strongo/random v0.0.1
	github.com/strongo/slice v0.3.0
	github.com/strongo/strongoapp v0.25.2
	github.com/strongo/validation v0.0.7
	github.com/yaa110/go-persian-calendar v1.2.1
	golang.org/x/crypto v0.27.0
	golang.org/x/net v0.29.0
)

require (
	github.com/AlekSi/pointer v1.2.0 // indirect
	github.com/alexsergivan/transliterator v1.0.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	github.com/sendgrid/rest v2.6.9+incompatible // indirect
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
