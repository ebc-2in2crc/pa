[English](README.md) | [日本語](README_ja.md)

# pa

[![Version](https://img.shields.io/github/release/ebc-2in2crc/pa.svg?label=version)](https://img.shields.io/github/release/ebc-2in2crc/pa.svg?label=version)

pa ([pʌ]) is a CLI tool for Pixela (pixe.la).

## Description

pa ([pʌ]) is a CLI tool for Pixela (pixe.la).

pa provides a unified command to manage Pixela's services.

pa takes advantage of the shell's completion capabilities to enter commands quickly, easily and reliably.

## Usage

### User API

```
$ pa user create \
    --username=yourname \
    --token=thisissecret \
    --agree-terms-of-service \
    --not-minor
```

or

```
$ export PA_USERNAME=yourname
$ export PA_TOKEN=thisissecret
$ pa user create --agree-terms-of-service --not-minor
```

User API sub commands.

- create
- delete
- update

### Graph API

```
$ pa graph create \
    --id=your-graph-id \
    --name=your-graph-name \
    --type=int \
    --unit=count \
    --color=ichou

$ pa graph get | jq
{
  "graphs": [
    {
      "id": "your-graph-id",
      "name": "your-graph-name",
      "unit": "count",
      "type": "int",
      "color": "ichou",
      "timezone": "",
      "purgeCacheURLs": null,
      "selfSufficient": "none",
      "isSecret": false,
      "publishOptionalData": false
    }
  ]
}
```

Graph API sub commands.

- create
- delete
- detail
- get
- list
- pixels
- stats
- stopwatch
- svg
- update

### Pixel API

```
$ pa pixel create --graph-id=your-graph-id --date 20200101 --quantity=1
```

Pixel API sub commands.

- create
- decrement
- delete
- get
- increment
- update

### Channel API

```
$ pa channel create \
    --id=your-channel-id \
    --name=your-channel-name \
    --type=slack \
    --slack-username=your-slack-user \
    --slack-channel-name=your-slack-channel \
    --slack-url=your-slack-url
```

Channel API sub commands.

- create
- delete
- get
- update

### Notification API

```
$ pa notification create \
    --id=your-notification-id \
    --name=your-notification-name \
    --channel-id=your-channel-id \
    --condition=">" \
    --target=quantity \
    --threshold=5 \
    --graph-id=your-graph-id
```

Notification API sub commands.

- create
- delete
- get
- update

### Webhook

```
$ pa webhook create --graph-id=your-graph-id --type=increment
```

Webhook API sub commands.

- create
- delete
- get
- invoke

### User name for Pixela, and token for Pixela

Specify the Pixela username with the `--username` flag and Specify the Pixela token with the `--token` flag.

```
$ pa user create \
    --username=yourname \
    --token=thisissecret \
    --agree-terms-of-service \
    --not-minor
```

You can also specify the Pixela username and Pixela token with the environment variables.

Specify the Pixela username with the `PA_USERNAME` environment variable and Specify the Pixela token with the `PA_TOKEN` environment variable.

```
$ export PA_USERNAME=yourname
$ export PA_TOKEN=thisissecret
$ pa user create --agree-terms-of-service --not-minor
```

### Generating shell completions

You can generate zsh, bash, fish and PowerShell completions and use it.

```
$ pa completion <SHELL> > /path/to/completion
```

### Help

Global help.

```
$ pa --help
The Pixela Command Line Interface is a unified tool to manage your Pixela services

Usage:
  pa [command]

Available Commands:
  channel      Channel
  completion   Generate shell completion
  graph        Graph
  help         Help about any command
  notification Notification
  pixel        Pixel
  user         User
  webhook      Webhook

Flags:
  -h, --help              help for pa
  -t, --token string      Pixela user token
  -u, --username string   Pixela user name
  -v, --version           version for pa

Use "pa [command] --help" for more information about a command.
```

Commands help available.

```
$ pa user --help
User

Usage:
  pa user [flags]
  pa user [command]

Available Commands:
  create      Create a new Pixela user
  delete      Delete a Pixela user
  update      Updates user token

Flags:
  -h, --help   help for user

Global Flags:
  -t, --token string      Pixela user token
  -u, --username string   Pixela user name

Use "pa user [command] --help" for more information about a command.
```

Sub commands help available.

```
$ pa user create --help
Create a new Pixela user

Usage:
  pa user create [flags]

Flags:
  -a, --agree-terms-of-service   Agree to the terms of service
  -h, --help                     help for create
  -m, --not-minor                You are not a minor or if you are a minor and you have the parental consent of using this service
  -c, --thanks-code string       Like a registration code obtained when you register for Patreon support

Global Flags:
  -t, --token string      Pixela user token
  -u, --username string   Pixela user name
```

## Installation

### Developer

```
$ go get -u github.com/ebc-2in2crc/pa/...
```

### User

Download from the following url.

- https://github.com/ebc-2in2crc/pa/releases

Or, you can use Homebrew (Only macOS).

```
$ brew tap ebc-2in2crc/tap
$ brew install pa
```

## References

[Pixela API Document](https://docs.pixe.la/)

## Contribution

1. Fork this repository
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Add some feature')
4. Rebase your local changes against the master branch
5. Run test suite with the go test ./... command and confirm that it passes
6. Run gofmt -s
7. Create new Pull Request

## License

[MIT](https://github.com/ebc-2in2crc/wareki/blob/master/LICENSE)

## Author

[ebc-2in2crc](https://github.com/ebc-2in2crc)
