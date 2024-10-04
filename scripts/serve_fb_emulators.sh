SCRIPT=$(realpath "$0")
SCRIPT_PATH=$(dirname "$SCRIPT")
FB_CONFIG="$SCRIPT_PATH"/../firebase/firebase.json
FB_DATA="$SCRIPT_PATH"/../firebase/local_data
firebase emulators:start --only auth,firestore --project demo-local-sneat-app --config "$FB_CONFIG" --import "$FB_DATA" --export-on-exit "$FB_DATA"
