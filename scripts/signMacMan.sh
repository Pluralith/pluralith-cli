cd app
env GOOS=darwin GOARCH=amd64 go build -o ../dist/darwin_darwin_amd64/pluralith
cd ..

# Sign & notarize binary
gon gon-config.json
