# cleura-client-go

Cleura API Client and a CLI application

## Installation

```terminal
go install github.com/aztekas/cleura-client-go/cmd/cleura@latest
```

> [!TIP]
> Check latest available version on the release page

## Usage and functionality

<cmd exec="cleura -h" src="."></cmd>

## Credentials

Cleura CLI requires a token to query Cleura API. The easiest way to get token is to generate it providing `username` and `password` to the `cleura token get` command (or use corresponding environment variables).

<cmd exec="cleura token get -h" src="."></cmd>

For convenience, you can, first,  generate a cleura configuration file with `generate-template` command:

<cmd exec="cleura config generate-template -h" src="."></cmd>

and then issue `cleura token get -u <username> -p <password> --update-config` command. Token will then be written to the configuration file in **open text**. Following `cleura` CLI commands will first try to use configuration file for receiving `username` and `token` values. Use the same command if token is revoked or outdated.

Commands that require `username` and `token` values would also attempt to read `CLEURA_API_USERNAME` and `CLEURA_API_TOKEN` environmental variables.
