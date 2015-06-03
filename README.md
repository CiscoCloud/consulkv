# consulkv
Command line interface to the [Consul K/V HTTP API](https://consul.io/docs/agent/http/kv.html)

## Installation
You can download a released `consulkv` artifact from [the consulkv release page][Releases] on Github. If you wish to compile from source, you will need to have buildtools and [Go][] installed:

```shell
$ git clone https://github.com/CiscoCloud/consulkv.git
$ cd consulkv
$ make
```

## Basic Usage

```
usage: consulkv [--version] [--help] <command> [<args>]

Available commands are:
    delete    Delete a path
    lock      Lock a node
    read      Read a value
    unlock    Unlock a node
    write     Write a value

```

### Common arguments

| Option | Default | Description |
| ------ | ------- | ----------- |
| `--consul` | `127.0.0.1:8500` | HTTP address of the Consul Agent
| `--ssl` | `false` | Use HTTPS while talking to Consul
| `--ssl-verify` | `true` | Verify certificates when connecting via SSL. Requires `--ssl`
| `--ssl-cert` | `unset` | Path to an SSL client certificate to use to authenticate to the consul server
| `--ssl-ca-cert` | `unset` | Path to a CA certificate file, containing one or more CA certificates to use to validate the certificate sent by the consul server to us.
| `--token` | `unset` | The [Consul API token][Consul ACLs].


### delete command


#### Usage

```
consulkv delete [options] path

  Delete a given path
```

#### Arguments

| Option | Default | Description |
|--------|---------|-------------|
| `--modifyIndex`* | `unset` | Perform a [Check-and-Set delete](https://consul.io/docs/agent/http/kv.html#DELETE)
| `--recurse` | `false` | Recursively delete path

\* Returns 1 if the Check-and-Set delete fails

#### Example

```shell
$ consulkv delete --ssl --consul=consul.service.consul:8500 --recurse nodes/config/test
```

### read command

#### Basic Usage

```
consulkv read [options] path

  Read a value from a given path.
```

#### Arguments

| Option | Default | Description |
| ------ | ------- | ----------- |
| `--fields`* | `value` | Comma separated list of fields to return
| `--format` | `raw` | Output format. One of `raw`, `json` or `prettyjson`

\* Currently supported fields: `Key`, `CreateIndex`, `ModifyIndex`, `LockIndex`, `Flags`, `Value`, `Session`

#### Examples

```shell
$ consulkv read --ssl nodes/config/test
Hello world
$ consulkv read --ssl --fields=modifyindex,value --format=prettyjson nodes/config/test
CAUBUCHO-M-X5AZ:consulkv chris$ ./consulkv read --ssl --ssl-verify=false --fields=modifyindex,value --format=prettyjson nodes/config/test
{
  "ModifyIndex": 916,
  "Value": "Hello world"
}
```

### write command

#### Basic Usage

```
consulkv write [options] path value

  Write a value to a given path.
```

#### Arguments

| Option | Default | Description |
| ------ | ------- | ----------- |
| `--modifyindex` | `unset` | Perform a [Check-and-Set write](https://consul.io/docs/agent/http/kv.html#PUT)
| `--flags` | `unset` | Integer value between 0 and 2<sup>64-1</sup>
| `path` | N/A | Path of the entry
| `value` | N/A | Entry value. `@file` syntax is supported

#### Examples

```shell
$ consulkv write --ssl nodes/config/test This is a test
$ consulkv read --ssl nodes/config/test
This is a test
$ cat test.data
Test data
$ consulkv write --ssl nodes/config/test @test.data
$ consulkv read --ssl nodes/config/test
Test data
```

### lock command

The lock command creates a lock at `path`. It returns the session ID string to use by the subsequent `unlock` command

#### Basic Usage

```
consulkv lock [options] path

  Acquire a lock on a given path
```

#### Arguments

| Option | Default | Description |
| ------ | ------- | ----------- |
| `--behavior` | `release` | Lock behavior. One of `release` or `delete`
| `--ttl` | `15s` | Lock time to live
| `--lock-delay` | `5s` | Lock delay

#### Examples

```shell
$ consulkv lock --ssl --ttl=10m locks/test
e8483859-fb12-2a5c-83f6-6f6f4af90c28
```

### unlock command

The unlock command releases a lock at `path`. The session ID returned by a previous `lock` command is a required argument.

#### Basic Usage

```
consulkv unlock [options] path

  Release a lock on a given path
```

#### Arguments

| Option | Default | Description |
| ------ | ------- | ----------- |
| `-session`* | `unset` | Session ID of the lock holder

\* Required

#### Example

```shell
$ consulkv unlock --ssl --session=e8483859-fb12-2a5c-83f6-6f6f4af90c28 locks/test
```

[Consul ACLs]: http://www.consul.io/docs/internals/acl.html "Consul ACLs"
[Releases]: https://github.com/CiscoCloud/consulkv/releases "consulkv releases page"
[Go]: http://golang.org "Go the language"

