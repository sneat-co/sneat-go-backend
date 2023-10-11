# sneat-go

[![Go CI](https://github.com/sneat-co/sneat-go-backend/actions/workflows/ci.yml/badge.svg)](https://github.com/sneat-co/sneat-go-backend/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sneat-co/sneat-go-backend)](https://goreportcard.com/report/github.com/sneat-co/sneat-go-backend)

Go lang backend for sneat apps:

- https://sneat.app/
- https://dailyscrums.app/ - free open source tool to run your stand-up meetings

## 3-d party dependencies

- AWS - to send emails

## Running GAE app

To run Google AppEngine app locally execute:

```shell
> cd sneatgaeapp
> go run main.go 
```

### Running with emulators

If you want Firebase SDKs to use Firebase emulators you would need to set next local variables:

```shell
export GCLOUD_PROJECT="demo-sneat"
export FIREBASE_AUTH_EMULATOR_HOST="localhost:9099"
export FIRESTORE_EMULATOR_HOST="localhost:8080"
```

Set path to Google Application Credentials files. For example:

```shell
export GOOGLE_APPLICATION_CREDENTIALS="~/projects/sneat/private_keys/sneat-54237a268b5a.json"
```

### Running with Google Firebase emulators:

<div style="text-decoration: line-through;">
At the time of writing Firebase Go Admin SDK does not support `Authentication` emulator. It is planned to be supported
in the future.

See [Admin SDK availability](https://firebase.google.com/docs/emulator-suite/install_and_configure#admin_sdk_availability)
to check current status.
</div>

Firebase Admin Go SDK is supposed
to [support Authentication emulator](https://github.com/firebase/firebase-admin-go/issues/409).
It has been merged with PR # [419](https://github.com/firebase/firebase-admin-go/pull/419)
on 21st April 2021 with
commit # [27ac52](https://github.com/firebase/firebase-admin-go/commit/27ac52fcc217798733768f26c2eb58cab54f5039).

So to use Firebase Authentication & Firestore emulators run as:

```shell
cd firestore
firebase emulators:start --only auth,firestore --project sneat-team
```

To run with real authentication but only emulated firestore execute:

```shell
cd firestore
firebase emulators:start --only firestore --project sneat-team
```

#### Note**:

It looks like a bug that Admin SDK is not able to use "demo-*" project ID.
We have to use a real project ID and provide an environment variable `GOOGLE_APPLICATION_CREDENTIALS`.
That's wrong and should not be used for scaled development.

#### Google emulators documentation

Read more about connecting an app to:

- [Authentication Emulator](https://firebase.google.com/docs/emulator-suite/connect_auth)
    - [Admin SKD](https://firebase.google.com/docs/emulator-suite/connect_auth#admin_sdks)
- [Cloud Firestore Emulator](https://firebase.google.com/docs/emulator-suite/connect_firestore)
    - [Admin SKD](https://firebase.google.com/docs/emulator-suite/connect_firestore#admin_sdks)

## Testing chatbots locally

There is a dedicated section regards how to test Telegram bots locally in [src/bots](src/sneatgae/sneatgaeapp/bots).