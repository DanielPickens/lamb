#!/bin/bash

if [[ -z "${VERSION}" ]]; then
  VERSION=latest
fi

curl -s https://api.github.com/danielpickens/lamb/releases/${VERSION} \
| grep "browser_download_url.*linux_amd64.tar.gz" \
| cut -d '"' -f 4 \
| wget -qi -

tarball="$(find . -name "*linux_amd64.tar.gz")"
tar -xzf $tarball

sudo mv lamb /bin

location="$(which lamb)"
echo "lamb binary location: $location"

version="$(lamb version)"
echo "lamb binary version: $version"
