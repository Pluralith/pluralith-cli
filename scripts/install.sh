# A little helper script to get Pluralith set up in CI
echo "Installing Pluralith"

# Get download url from latest release
url="https://api.pluralith.com/v1/dist/download/cli?os=linux&arch=amd64"
url=$(curl -s $url | jq -r '.data.url')

# Download latest release binary
curl -sL $url -o "/tmp/pluralith"

# Make binary executable
chmod +x "/tmp/pluralith"

# Move to /usr/local/bin
if [ -x "$(command -v sudo)" ]; then
  sudo mv "/tmp/pluralith" "/usr/local/bin/pluralith"
else
  mv "/tmp/pluralith" "/usr/local/bin/pluralith"
fi

echo "Pluralith successfully installed"