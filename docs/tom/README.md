# tom

`tom` is the cli interface for the `TOM` service.

## Configuration file

The default location for the configuration file is `~/.tom/tom.conf`. An
alternate location can be provided via `-c|-config` commandline flag.

```
api: https://127.0.0.1:443/
logdir: ~/.tom/log
json.output.processor: jq --unbuffered --sort-keys
ca.file: ~/.tom/ca.pem
authentication {
    credential.path: ~/.tom/root/
    username: root
    identity.library: system
    use.fingerprint.as.username: no
}
```

## First start

During initial setup, started with `enforcement: false` parameter, the
service will only allow registering keymaterial with pre-existing accounts,
ie. `root`.

Create a special `root.conf` in `~/.tom/root/` and run a command using this
configuration:

```
tom -config ~/.tom/root/root.conf namespace list
```

This will create the following three files and register the key with the
account:

```
~/.tom/root/machinekey.epk -- the encrypted Ed25519 privatekey
~/.tom/root/machinekey.pub -- the Ed25519 publickey
~/.tom/root/passphrase     -- the unshifted base information for unlocking
                              the private key
```

## Installating a user identity library and user enrolment

Work in progress. Sorry.
