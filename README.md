# parrot 🦜

A simple service to relay Grafana Alerts via webhooks to NTFY.

Simply set the destination address for your Grafana alert webhook to the endpoint of parrot. The notification parmameters are read from the webhook body and passed to NTFY. If the `Authorization` header is set in the Grafana webhook, it is simply passed on to NTFY.

## Setup

You can simply use the precompiled binaries from the [Releases](https://github.com/studio-b12/parrot/releases) page to set up parrot on your system.

### Config

The service is configured either via CLI flags or via environment variables. Documentation about the configuration variables can be found using the `--help` page of the service.

```
Usage: parrot [--bind-address BIND-ADDRESS] --ntfy-upstream NTFY-UPSTREAM [--log-level LOG-LEVEL]

Options:
  --bind-address BIND-ADDRESS
                         HTTP bind address [default: 0.0.0.0:8080, env: BIND_ADDRESS]
  --ntfy-upstream NTFY-UPSTREAM
                         Address of the upstream NTFY server [env: NTFY_UPSTREAM]
  --log-level LOG-LEVEL
                         Log level [default: info, env: LOG_LEVEL]
  --help, -h             display this help and exit
```

### dpkg

There are also dpkg-Packages you can use to set up parrot on a Debian-Server as systemd service. The configuration is done via the environment file in `/etc/parrot/vars.env`.

```bash
dpkg -i parrot_v1.0.0_amd64.deb
```
