# 1) Build unsigned binary
cd app
env GOOS=darwin GOARCH=amd64 go build -o ../dist/unsigned/pluralith_cli_darwin_amd64
env GOOS=windows GOARCH=amd64 go build -o ../dist/unsigned/pluralith_cli_windows_amd64
env GOOS=linux GOARCH=amd64 go build -o ../dist/unsigned/pluralith_cli_linux_amd64

# 2) Sign & notarize binary
cd ..
gon gon-config.json

# 3) Unzip notarized binary
cd dist/signed
unzip -o *.zip

# 4) Eliminate zips
rm *.zip

# 5) Upload to GCS
# gsutil cp ./* gs://pluralith-cli-releases