[English](README.md) | [日本語](README_ja.md)

# pa

[![Build Status](https://travis-ci.com/ebc-2in2crc/pa.svg?branch=master)](https://travis-ci.com/ebc-2in2crc/pa)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/ebc-2in2crc/pa)](https://goreportcard.com/report/github.com/ebc-2in2crc/pa)
[![Version](https://img.shields.io/github/release/ebc-2in2crc/pa.svg?label=version)](https://img.shields.io/github/release/ebc-2in2crc/pa.svg?label=version)

pa は [Pixela](https://pixe.la/) の CLI ツールです。

pa は「ぱ」と読みます。

![demo](https://user-images.githubusercontent.com/1122547/87226990-a959bb00-c3d2-11ea-96ab-f9b4015a7ed2.gif)

## Description

pa は [Pixela](https://pixe.la/) の CLI ツールです。

pa は Pixela のサービスを管理するための統一されたコマンドを提供します。

pa はシェルの補完機能を利用してコマンドを素早く、簡単にそして確実に入力できます。

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

### User Profile API

```
$ pa profile update \
    --display-name=display-name \
    --gravatar-icon-email=gravatar-icon-email \
    --title=title \
    --timezone=Asia/Tokyo \
    --about-url=about-URL \
    --contribute-urls=url \
    --pinned-graph-id=pinned-graph-id
```

User Profile API sub commands.

- update
- get

### Graph API

```
$ pa graph create \
    --id=your-graph-id \
    --name=your-graph-name \
    --type=int \
    --unit=count \
    --color=ichou

$ pa graph get-all | jq
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
- get-all
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

### Webhook

```
$ pa webhook create --graph-id=your-graph-id --type=increment
```

Webhook API sub commands.

- create
- delete
- get
- invoke

### Pixela のユーザー名とトークン

Pixela のユーザー名は `--username` フラグで指定して Pixela のトークンは `--token` フラグで指定します。

```
$ pa user create \
    --username=yourname \
    --token=thisissecret \
    --agree-terms-of-service \
    --not-minor
```

Pixela のユーザー名とトークンは環境変数で指定することもできます。

Pixela のユーザー名は `PA_USERNAME` 環境変数で指定して Pixela のトークンは `PA_TOKEN` フラグで指定します。

```
$ export PA_USERNAME=yourname
$ export PA_TOKEN=thisissecret
$ pa user create --agree-terms-of-service --not-minor
```

Pixela のユーザー名とトークンは設定ファイルで指定することもできます。

```
$ cat ~/.pa
username = "yourname"
token = "thisissecret"
$ pa user create --agree-terms-of-service --not-minor
```

pa は次の優先順位でユーザー名とトークンを使用します。
それぞれの項目はその下の項目よりも優先されます。

- フラグ
- 環境変数
- カレントディレクトリの設定ファイル
- ホームディレクトリの設定ファイル

### シェルの補完スクリプトの生成

Zsh, Bash, Fish, PowerShell の補完スクリプトを生成して利用できます。

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
  completion  Generate shell completion
  graph       Graph
  help        Help about any command
  pixel       Pixel
  profile     Profile
  user        User
  webhook     Webhook

Flags:
      --config string     config file (default is $HOME/.pa)
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
      --config string     config file (default is $HOME/.pa)
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
      --config string     config file (default is $HOME/.pa)
  -t, --token string      Pixela user token
  -u, --username string   Pixela user name
```

## Installation

### Developer

```
$ go get -u github.com/ebc-2in2crc/pa/...
```

### User

[https://github.com/ebc-2in2crc/pa/releases](https://github.com/ebc-2in2crc/pa/releases) からダウンロードします。

Homebrew を使うこともできます (Mac のみ)

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
