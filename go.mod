module github.com/sneat-co/sneat-go-backend

go 1.21.5

// TODO: Get ret rid of: github.com/dal-go/dalgo2datastore

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

require (
	github.com/julienschmidt/httprouter v1.3.0
	github.com/sneat-co/sneat-core-modules v0.12.4
	github.com/sneat-co/sneat-go-core v0.23.0
	github.com/sneat-co/sneat-go-firebase v0.4.21
	github.com/sneat-co/sneat-go-modules v0.3.1
	github.com/stretchr/testify v1.8.4
	github.com/strongo/log v0.3.0
)

require (
	cloud.google.com/go v0.111.0 // indirect
	cloud.google.com/go/compute v1.23.3 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/firestore v1.14.0 // indirect
	cloud.google.com/go/iam v1.1.5 // indirect
	cloud.google.com/go/longrunning v0.5.4 // indirect
	cloud.google.com/go/storage v1.36.0 // indirect
	firebase.google.com/go/v4 v4.13.0 // indirect
	github.com/MicahParks/keyfunc v1.9.0 // indirect
	github.com/alexsergivan/transliterator v1.0.0 // indirect
	github.com/bots-go-framework/bots-fw-store v0.4.0 // indirect
	github.com/crediterra/money v0.2.1 // indirect
	github.com/dal-go/dalgo v0.12.0 // indirect
	github.com/dal-go/dalgo2firestore v0.1.43 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.0 // indirect
	github.com/gosimple/slug v1.13.1 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.11.0 // indirect
	github.com/strongo/decimal v0.0.1 // indirect
	github.com/strongo/random v0.0.1 // indirect
	github.com/strongo/slice v0.1.4 // indirect
	github.com/strongo/strongoapp v0.17.0 // indirect
	github.com/strongo/validation v0.0.6 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.46.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.46.1 // indirect
	go.opentelemetry.io/otel v1.21.0 // indirect
	go.opentelemetry.io/otel/metric v1.21.0 // indirect
	go.opentelemetry.io/otel/trace v1.21.0 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/oauth2 v0.15.0 // indirect
	golang.org/x/sync v0.6.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	golang.org/x/xerrors v0.0.0-20231012003039-104605ab7028 // indirect
	google.golang.org/api v0.155.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/appengine/v2 v2.0.5 // indirect
	google.golang.org/genproto v0.0.0-20240102182953-50ed04b92917 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240102182953-50ed04b92917 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240102182953-50ed04b92917 // indirect
	google.golang.org/grpc v1.60.1 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
