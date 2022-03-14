cd app
env GOOS=darwin GOARCH=amd64 go build -o ../dist/unsigned/pluralith_cli_darwin_amd64
cd ..
# 1) Sign & notarize binary
gon gon-config.json

# 2) Replace unsigned with signed binary
cd dist/unsigned/pluralith_cli_darwin_amd64
# cd dist/darwin_darwin_amd64
rm pluralith
unzip -o pluralith.zip
rm *.zip
