# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.137.0/containers/go/.devcontainer/base.Dockerfile
ARG VARIANT="1"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

# [Optional] Uncomment this section to install additional OS packages.
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends unzip

# [Optional] Uncomment the next line to use go get to install anything else you need
# RUN go get -x <your-dependency-or-tool>
RUN go get -tags 'mysql' -u \
    github.com/golang-migrate/migrate/cmd/migrate

# Install ngrok
COPY ./ngrok /bin/ngrok

# Install Terraform
RUN curl -s -o /var/tmp/terraform.zip \
    -L "https://releases.hashicorp.com/terraform/0.12.26/terraform_0.12.26_linux_386.zip"
RUN unzip /var/tmp/terraform.zip -d /usr/bin
