# Run Locally

First install  [act](https://github.com/nektos/act).

```bash
brew install act
```

And you need to have a container runtime. For example, [Docker Desktop](https://www.docker.com/products/docker-desktop/) or [Colima](https://github.com/abiosoft/colima)


```bash
brew install colima
colima start
```

## Generate GPG Key

You need a local GPG key. If you don't have one, you can generate one following the instructions [here](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-release-publish)

Next, ensure you have secrets that would have been set on the environment of the GitHub repo locally puto in a `.secrets` file. Example:

```bash
GITHUB_TOKEN="<github_access_token_goes_here"
GPG_PASSPHRASE="<gpg_password_goes_here>"
GPG_PRIVATE_KEY="<gpg_private_key_goes_here"
```

## Run specific workflow

```bash
act -W '.github/workflows/build.yml' -P macos-latest=-self-hosted
```
