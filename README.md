# uinit

**u**nix **init** — a small process supervisor written in Go, inspired by
[zinit](https://github.com/threefoldtech/zinit).

uinit runs as a daemon that loads a list of services from a YAML file and exposes
them over a Unix socket. A CLI client talks to the daemon to inspect them.

This is a learning project for exploring Linux processes, signals, and Go systems
programming.

## Install

```sh
git clone https://github.com/mahmoud-emad/uinit.git
cd uinit
make build
```

Requires Go 1.26+. The binary is written to `./uinit`.

## Usage

Write a config file:

```yaml
# services.yaml
services:
  - name: python-server
    cmd: python3 -m http.server
  - name: ping
    cmd: ping 8.8.8.8
```

Start the daemon:

```sh
./uinit init services.yaml
```

In another terminal, list the services:

```sh
./uinit list
```

```
SERVICE        STATUS      PID      UPTIME   COMMAND
─────────────────────────────────────────────────────────────────────
python-server  ● running   45609    7s       python3 -m http.server
ping           ● running   45610    7s       ping 8.8.8.8
sleep-five     ○ exited    45611    5s       sleep 5
```
