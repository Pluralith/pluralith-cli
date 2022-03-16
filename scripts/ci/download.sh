# Download latest Pluralith release
curl -s "https://api.pluralith.com/v1/dist/download/cli?arch=amd64&os=linux" \
| jq '.data.url' \
| xargs wget -O "pluralith"