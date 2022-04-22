# irsa-tokengen

Generate an `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, and `AWS_SESSION_TOKEN` from the IRSA `AWS_WEB_IDENTITY_TOKEN_FILE` and `AWS_ROLE_ARN`. Some applications that do not support IRSA  might find this handy, like Terraform's Go Getter.

## Build

Download the repo and make a tiny binary for your pod:

```bash
GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o irsa-tokengen main.go
```

## Install

Place the build binary somewhere in `$PATH`.

## Usage

Requires the role-arn to assume and the JWT file generally placed into pods automatically with IRSA.

```bash
if [[ -s $AWS_WEB_IDENTITY_TOKEN_FILE ]]; then
    export $(irsa-tokengen)
fi
```

**Caveats**

The token is only good for 1 hour, even if the IAM Role has a longer duration.