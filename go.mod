module github.com/sneat-co/sneat-go-backend

go 1.22.2

//replace github.com/sneat-co/sneat-go-core => ../sneat-go-core
//replace github.com/sneat-co/sneat-core-modules => ../sneat-core-modules
//replace github.com/sneat-co/sneat-go-modules => ../sneat-go-modules
//replace github.com/bots-go-framework/bots-fw => ../../bots-go-framework/bots-fw
//replace github.com/bots-go-framework/bots-fw-store => ../../bots-go-framework/bots-fw-store
//replace github.com/bots-go-framework/bots-fw-telegram => ../../bots-go-framework/bots-fw-telegram
//replace github.com/bots-go-framework/bots-fw-telegram-models => ../../bots-go-framework/bots-fw-telegram-models
//replace github.com/bots-go-framework/bots-host-gae => ../../bots-go-framework/bots-host-gae
//replace github.com/bots-go-framework/dalgo4botsfw => ../../bots-go-framework/dalgo4botsfw
//replace github.com/dal-go/dalgo => ../../dal-go/dalgo
//replace github.com/strongo/app => ../../strongo/app
//replace github.com/strongo/i18n => ../../strongo/i18n
//replace github.com/strongo/strongoapp => ../../strongo/strongoapp

require (
	github.com/bots-go-framework/bots-api-telegram v0.4.3
	github.com/bots-go-framework/bots-fw v0.25.2
	github.com/bots-go-framework/bots-fw-store v0.4.0
	github.com/bots-go-framework/bots-fw-telegram v0.8.5
	github.com/bots-go-framework/bots-fw-telegram-models v0.1.4
	github.com/bots-go-framework/bots-go-core v0.0.2
	github.com/bots-go-framework/bots-host-gae v0.5.7
	github.com/captaincodeman/datastore-mapper v0.0.0-20170320145307-cb380a4c4d13
	github.com/crediterra/go-interest v0.0.0-20180510115340-54da66993b85
	github.com/crediterra/money v0.2.1
	github.com/dal-go/dalgo v0.12.1
	github.com/dal-go/mocks4dalgo v0.1.19
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/golang/mock v1.6.0
	github.com/gorilla/sessions v1.2.2
	github.com/gosimple/slug v1.14.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/olebedev/when v1.0.0
	github.com/pquerna/ffjson v0.0.0-20190930134022-aa0246cd15f7
	github.com/sanity-io/litter v1.5.5
	github.com/sendgrid/sendgrid-go v3.14.0+incompatible
	github.com/shiyanhui/hero v0.0.2
	github.com/sneat-co/debtstracker-translations v0.0.21
	github.com/sneat-co/sneat-go-core v0.23.4
	github.com/sneat-co/sneat-go-firebase v0.4.39
	github.com/stretchr/testify v1.9.0
	github.com/strongo/app-host-gae v0.1.26
	github.com/strongo/decimal v0.0.1
	github.com/strongo/delaying v0.0.1
	github.com/strongo/facebook v1.8.1
	github.com/strongo/gamp v0.0.1
	github.com/strongo/gotwilio v0.0.0-20160123000810-f024bbefe80f
	github.com/strongo/i18n v0.0.4
	github.com/strongo/log v0.3.0
	github.com/strongo/random v0.0.1
	github.com/strongo/slice v0.1.4
	github.com/strongo/slices v0.0.0-20231201223919-29a6c669158a
	github.com/strongo/strongoapp v0.18.2
	github.com/strongo/validation v0.0.6
	github.com/yaa110/go-persian-calendar v1.2.0
	golang.org/x/crypto v0.22.0
	golang.org/x/net v0.24.0
	google.golang.org/appengine/v2 v2.0.6
)

require (
	cloud.google.com/go v0.112.2 // indirect
	cloud.google.com/go/compute/metadata v0.3.0 // indirect
	cloud.google.com/go/firestore v1.15.0 // indirect
	cloud.google.com/go/iam v1.1.7 // indirect
	cloud.google.com/go/longrunning v0.5.6 // indirect
	cloud.google.com/go/storage v1.40.0 // indirect
	firebase.google.com/go/v4 v4.14.0 // indirect
	github.com/AlekSi/pointer v1.2.0 // indirect
	github.com/MicahParks/keyfunc v1.9.0 // indirect
	github.com/alexsergivan/transliterator v1.0.0 // indirect
	github.com/captaincodeman/datastore-locker v0.0.0-20170308203336-4eddc467484e // indirect
	github.com/dal-go/dalgo2firestore v0.1.59 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.3 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/sendgrid/rest v2.6.9+incompatible // indirect
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.50.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.50.0 // indirect
	go.opentelemetry.io/otel v1.25.0 // indirect
	go.opentelemetry.io/otel/metric v1.25.0 // indirect
	go.opentelemetry.io/otel/trace v1.25.0 // indirect
	golang.org/x/oauth2 v0.19.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	google.golang.org/api v0.173.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto v0.0.0-20240415180920-8c6c420018be // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240415180920-8c6c420018be // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240415180920-8c6c420018be // indirect
	google.golang.org/grpc v1.63.2 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
