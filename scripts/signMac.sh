
# 1) Sign & notarize binary
gon gon-config.json

# 2) Replace unsigned with signed binary
cd dist/darwin_darwin_amd64
rm pluralith
unzip -o pluralith.zip
rm *.zip
