---
meta:
  - name: description
    content: "lamb | Installation documentation"
---
# Installation

## asdf

We have an [asdf](https://asdf-vm.com/#/) plugin [here](https://github.com/danielpickens/asdf-lamb). You can install with:

```
asdf plugin-add lamb
asdf list-all lamb
asdf install lamb <latest version>
asdf local lamb <latest version>
```

## Binary

Install the binary from our [releases](https://github.com/danielpickens/lamb/releases) page.

## Homebrew Tap

```
brew installdanielpicks/tap/lamb
```

## Scoop (Windows)
Note: This is not maintained by danielpickens, but should stay up to date with future releases.

```
scoop install lamb
```

# Verify Artifacts
danielpickens signs the lamb docker image and the checksums file with [cosign](https://github.com/sigstore/cosign). Our public key is available at https://artifacts/danielpickens.com/cosign.pub

You can verify the checksums file from the [releases](https://github.com/danielpicks/lamb/releases) page with the following command:

```
cosign verify-blob checksums.txt --signature=checksums.txt.sig  --key https://artifacts/danielpickens.com/cosign.pub
```

Verifying docker images is even easier:

```
cosign verify us-docker.pkg.dev/danielpickens/oss/lamb:v5 --key https://artifacts/danielpickens.com/cosign.pub
```

