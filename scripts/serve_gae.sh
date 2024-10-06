sed -i '' 's|// replace github.com/sneat-co/sneat-go-backend => ../sneat-go-backend|replace github.com/sneat-co/sneat-go-backend => ../sneat-go-backend|g' ../../sneat-go-server/go.mod
DIR_NAME="$(dirname "$(realpath "$0")")"

bash "$DIR_NAME/../../sneat-go-server/serve_gae.sh"