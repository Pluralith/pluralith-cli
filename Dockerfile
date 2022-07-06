FROM alpine:3.15

# Set environment variables
ENV PLURALITH_CI=true

# Install dependencies
RUN apk upgrade --no-cache && apk --no-cache add bash jq curl wget gcompat libgcc libstdc++ npm

# Set shell to bash to use bash array
SHELL ["/bin/bash", "-ec"]

# Download and install Pluralith CLI
COPY ./scripts/ci ./scripts
RUN ./scripts/download.sh
RUN ./scripts/install-pluralith.sh

# Make terraform installation script executable in finished image
RUN chmod +x ./scripts/install-terraform.sh

# Install Compost for pull request commenting
RUN npm install -g @infracost/compost

ENTRYPOINT [ "pluralith" ]

