export GCLOUD_PROJECT="demo-local-sneat-app"
export FIREBASE_AUTH_EMULATOR_HOST="localhost:9099"
export FIRESTORE_EMULATOR_HOST="localhost:8080"
export SNEAT_TG_DEV_BOTS="AlextDevBot:sneat_bot"
#export GOOGLE_APPLICATION_CREDENTIALS="/Users/alexandertrakhimenok/projects/sneat/sneat-team-go/private_keys/demo-sneat.json"

script_parent_dir=$(dirname "$(dirname "$(realpath "$0")")")
go_main_path="$script_parent_dir"
go run "$go_main_path"
#
#SCRIPT=$(realpath "$0")
#SCRIPT_PATH=$(dirname "$SCRIPT")
#go run "$SCRIPT_PATH"/../src/sneatgaemain --with-bots
