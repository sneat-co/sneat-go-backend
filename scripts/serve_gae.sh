export GCLOUD_PROJECT="demo-local-sneat-app"
export FIREBASE_AUTH_EMULATOR_HOST="localhost:9099"
export FIRESTORE_EMULATOR_HOST="localhost:8080"
export SNEAT_TG_DEV_BOTS="AlextDevBot:listus"

script_parent_dir=$(dirname "$(dirname "$(realpath "$0")")")
go_main_path="$script_parent_dir"/src/sneatgae/sneatgaemain
go run "$go_main_path"
