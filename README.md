# kea-cli

This project provides a Kea DHCP Server REST API client CLI.

See more information here: https://kea.readthedocs.io/en/latest/arm/agent.html

API reference can be found here: https://kea.readthedocs.io/en/latest/api.html#api-reference

## How to build kea-cli
Build can be run the following way:

```bash
$ cd kea-cli
$ make
```

The kea-cli binary file is built in the kea-cli/bin directory

## How to run kea-cli unit tests
Tests can be run the following way:

```bash
$ cd kea-cli
$ make test
```

## How to configure kea-cli

A config file kea-cli.yaml or kea-cli.json can be provided in the following locations:

- current kea-cli directory
- $HOME/.kea-cli directory
- /etc/kea-cli directory


## How to use kea-cli

Once the kea-cli binary is built, installed in your PATH and configured to access to an exposed Kea DHCP REST API, you just need to run:

```bash
$ kea-cli -h
kea-cli is a Kea REST API client

Usage:
  kea-cli [flags]
  kea-cli [command]

Available Commands:
  completion  Generate completion script
  leases4     Perform commands on leases for ipv4

Flags:
      --api_url string                 Specify the host and port used to join the Kea REST API (default "127.0.0.1:8000")
      --ca_file string                 This provides a file path to a PEM-encoded certificate authority. This implies the use of HTTPS to connect to the Kea REST API.
      --ca_path string                 Path to a directory of PEM-encoded certificates authorities. This implies the use of HTTPS to connect to the Kea REST API.
      --cert_file string               File path to a PEM-encoded client certificate used to authenticate to the Kea REST API. This must be provided along with key-file. If one of key-file or cert-file is not provided then SSL authentication is disabled. If both cert-file and key-file are provided this implies the use of HTTPS to connect to the Kea REST API.
  -c, --config string                  Config file (default is /etc/kea-cli/kea-cli.[json|yaml])
  -h, --help                           help for kea-cli
      --http_client_timeout duration   Timeout for HTTP Client used to join the Kea REST API (default 30s)
      --key_file string                File path to a PEM-encoded client private key used to authenticate to the Kea REST API. This must be provided along with cert-file. If one of key-file or cert-file is not provided then SSL authentication is disabled. If both cert-file and key-file are provided this implies the use of HTTPS to connect to the Kea REST API.
      --log_level string               Log level: can be 'trace', 'debug', 'info', 'warning', 'error', 'fatal', 'panic'. (default "info")
      --skip_tls_verify                Controls whether a client verifies the server's certificate chain and host name. If set to true, TLS accepts any certificate presented by the server and any host name in that certificate. In this mode, TLS is susceptible to man-in-the-middle attacks. This should be used only for testing. This implies the use of HTTPS to connect to the Kea REST API.
      --ssl_enabled                    Use HTTPS to connect to the Kea REST API
  -v, --version                        version for kea-cli

Use "kea-cli [command] --help" for more information about a command.
```

You can create bash completion script with the following command:

```bash
$ sudo bash -c 'bin/kea-cli completion bash > /etc/bash_completion.d/kea-cli'
```

## How to debug kea-cli

To log in debug mode all requests and responses contents, run your command with the --log_level=debug flag:

```bash
$ kea-cli leases4 list --log_level=debug
```

## How to contribute to kea-cli

For now, only some leases4 commands have been implemented, so you're welcome to implement others.

