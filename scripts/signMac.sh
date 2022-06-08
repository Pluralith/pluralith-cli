# 1) Sign & notarize binary
gon gon-config.json

# 2) Replace unsigned with signed binary
cd dist/darwin
rm pluralith
unzip -o pluralith.zip
rm *.zip
