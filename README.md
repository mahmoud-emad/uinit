# uinit

**u**nix **init** — a small process manager written in Go, inspired by
[zinit](https://github.com/threefoldtech/zinit).

uinit runs as a daemon that loads a list of processes from a YAML file and exposes
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
# processes.yaml
processes:
  - name: python-server
    cmd: python3 -m http.server
  - name: ping
    cmd: ping 8.8.8.8
```

Start the daemon:

```sh
./uinit init processes.yaml
```

In another terminal, list the processes:

```sh
./uinit list
```

```
PROCESS        STATUS      PID      UPTIME   COMMAND
─────────────────────────────────────────────────────────────────────
python-server  ● running   45609    7s       python3 -m http.server
ping           ● running   45610    7s       ping 8.8.8.8
sleep-five     ○ exited    45611    5s       sleep 5
```

Inspect one of them, or read what it has written:

```sh
./uinit inspect ping
./uinit logs ping
```

Each process writes its stdout and stderr to `/tmp/uinit/logs/<name>.log`.
