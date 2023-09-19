export GCLOUD_PROJECT="demo-local-sneat-app"
export FIREBASE_AUTH_EMULATOR_HOST="localhost:9099"
export FIRESTORE_EMULATOR_HOST="localhost:8080"
export SNEAT_TG_DEV_BOTS="AlextDevBot:listus"
#export GOOGLE_APPLICATION_CREDENTIALS="/Users/alexandertrakhimenok/projects/sneat/sneat-team-go/private_keys/demo-sneat.json"

SCRIPT=$(realpath "$0")
SCRIPT_PATH=$(dirname "$SCRIPT")
go run "$SCRIPT_PATH"/../src/sneatgae/sneatgaemain
