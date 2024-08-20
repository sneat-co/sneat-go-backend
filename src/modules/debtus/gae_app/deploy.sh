#!/usr/bin/env bash

rm -rf ./ionic-app/assets/
rm -rf ./ionic-app/build/
cp -r ../../ionic-apps/public/platforms/browser/www/ ./ionic-app

while true; do
    read -p "Where do you want to deploy? (dev|prod): " app
    case $app in
        dev )
        	sed -i '' 's/^application: *[[:alpha:]]*-[[:alnum:]]*/application: debtusbot-dev1/' app.yaml
        	break;;
        prod )
        	sed -i '' 's/^application: *[[:alpha:]]*-[[:alnum:]]*/application: debtusbot-io/' app.yaml
        	break;;
        * ) echo "Please answer 'dev' or 'prod'.";;
    esac
done

echo "You selected: $app"

echo "Starting tests..."
# TODO: Add testing for ../../../../github.com/strongo/bots-framework/...
goapp test ../...
if [ $? -ne 0 ]; then
    echo "Tests failed"
	sed -i '' 's/^application: *[[:alpha:]]*-[[:alnum:]]*/application: debtusbot-local/' app.yaml
	exit 1
fi

goapp deploy
sed -i '' 's/^application: *[[:alpha:]]*-[[:alnum:]]*/application: debtusbot-local/' app.yaml

