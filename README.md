
# cleura-client-go

Cleura API Client and a CLI application

## Installation

```terminal
go install github.com/aztekas/cleura-client-go/cmd/cleura@latest

```

> [!TIP] Check latest available version on the release page

## Usage and functionality

```shell
$ cleura -h

NAME:
   cleura - A Cleura API CLI

USAGE:
   cleura [global options] command [command options]

VERSION:
   latest-uncommitted

COMMANDS:
   config   Command used for working with configuration file for the cleura cli
   domain   Command used to perform actions with available domains
   project  Command used to perform operations with projects
   token    Command used to perform actions with Cleura API tokens
   shoot    Command used to perform operations with shoot clusters
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --loglevel value  (default: "info")
   --help, -h        show help
   --version, -v     print the version
```

## Credentials

Cleura CLI requires a token to query Cleura API. The easiest way to get token is to generate it providing `username` and `password` to the `cleura token get` command (or use corresponding environment variables).

```shell
$ cleura token get -h

NAME:
   cleura token get - Receive token from Cleura API using username and password

USAGE:
   cleura token get [command options] [arguments...]

DESCRIPTION:
   Receive token from Cleura API using username and password

OPTIONS:
   --username value, -u value      Username for token request [$CLEURA_API_USERNAME]
   --password value, -p value      Password for token request. [$CLEURA_API_PASSWORD]
   --api-host value, --host value  Cleura API host (default: "https://rest.cleura.cloud") [$CLEURA_API_HOST]
   --update-config                 Save token to active configuration. NB: token saved in open text (default: true)
   --config-path value             Path to configuration file. $HOME/.config/cleura/config if not set
   --interactive, -i               Interactive mode. Input username and password in interactive mode (default: false)
   --two-factor, --2fa             Set this flag if two-factor authentication (sms) is enabled in your cleura profile  (default: false)
   --help, -h                      show help
```

For convenience, you can, first,  generate a cleura configuration file with `generate-template` command:

```shell
$ cleura config generate-template -h

NAME:
   cleura config generate-template - Generate configuration file template on the given path

USAGE:
   cleura config generate-template [command options] [arguments...]

DESCRIPTION:
   Generate configuration file template on the given path

OPTIONS:
   --output-file value, -o value  Path to configuration file. $HOME/.config/cleura/config if not set. NB: Overwrites existing if found
   --help, -h                     show help
```

and then issue `cleura token get -u <username> -p <password> --update-config` command. Token will then be written to the configuration file in **open text**. Following `cleura` CLI commands will first try to use configuration file for receiving `username` and `token` values. Use the same command if token is revoked or outdated.

Commands that require `username` and `token` values would also attempt to read `CLEURA_API_USERNAME` and `CLEURA_API_TOKEN` environmental variables.
