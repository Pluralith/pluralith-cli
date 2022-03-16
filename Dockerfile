FROM alpine:3.15

# Set environment variables
ENV PLURALITH_CI=true

# Install dependencies
RUN apk upgrade --no-cache && apk --no-cache add jq curl wget gcompat libgcc libstdc++

# Donwload and install Pluralith CLI
COPY ./scripts/ci ./scripts
RUN ./scripts/download.sh
RUN ./scripts/install.sh

CMD tail -f /dev/null

# ENTRYPOINT [ "pluralith" ]


