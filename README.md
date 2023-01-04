<div align="right">

[![golangci-lint](https://github.com/sheepla/strans/actions/workflows/ci.yml/badge.svg)](https://github.com/sheepla/strans/actions/workflows/ci.yml)
[![Release](https://github.com/sheepla/strans/actions/workflows/release.yml/badge.svg)](https://github.com/sheepla/strans/actions/workflows/release.yml)

</div>

<div align="center">

# strans

</div>


<div align="center">

A command line translate tool written in Go with GNU Readline-like interactive mode (`--repl`) inspired by [translate-shell](https://github.com/soimort/translate-shell)


![Language:Go](https://img.shields.io/static/v1?label=Language&message=Go&color=blue&style=flat-square)
![License:MIT](https://img.shields.io/static/v1?label=License&message=MIT&color=blue&style=flat-square)
[![Latest Release](https://img.shields.io/github/v/release/sheepla/strans?style=flat-square)](https://github.com/sheepla/strans/releases/latest)

</div>

## Features

- **Non-interactive mode**: A mode that can be used in the same way as a general command line tool
- **Interactive mode**: GNU Readline-like line editing mode for instant translation
- **Read translated text aloud**: The option to read the translated text aloud after performing the translation

## Usage

### Non-interactive mode

If you specify text as a non-option argument, that text will be translated. Multiple arguments are allowed, and arguments are joined by spaces.

```
strans [OPTIONS] TEXT...

# e.g.
strans -s en -t ja "Hello, World" # => "こんにちは世界"
```

A non-option argument of `-` will read text from standard input and translate it.
You can also output the translated text as speech.

```
strans [OPTIONS] -

# e.g.
echo "Hello, World" | strans -s en -t ja - # => "こんにちは世界"
strans -s en -t ja < README.md # => The contents of the README.md will translated.
```

### Interactive mode

To use interactive mode, run the command with the `-r`, `--repl` flag. 
Enter your text and it will be translated instantly.

You can use GNU Readline-like line editing, scroll back (`Ctrl-N`, `Ctrl-P`) and incremental search (`Ctrl-R`) the execution history.

Empty inputs (just typing `Enter`) are ignored and no translation is performed.

Typing `Ctrl-D` exits interactive mode and returns you to the shell you were running from.

The history is kept in the file `strans_history.txt` in the OS temporary directory 
and can be recalled when executing the command again.

```sh
[you@your-computer]$ strans --repl -s ja -t en

[ja -> en]
> こんにちは世界
hello world

[ja -> en]
> bye

[you@your-computer]$ 
```

### Read text as speech (beta)

Running the command with the `--audio` flag, after executing the translation, read the translated text aloud and you can check the pronunciation.

This feature is available in both interactive and non-interactive mode.

### Options

Specify the source language name (e.g. `en`, `ja`, etc.) for `--source` option, and specify the target language name for `--target` option.

To change the instance, specify the hostname of the instance in the `--instance` option.

```
NAME:
   strans - a command line translate tool with GNU Readline like interactive mode

USAGE:
   strans [OPTIONS] TEXT...
   echo TEXT... | strans [OPTIONS] -
   strans [OPTIONS] - < FILE

VERSION:
   unknown-unknown

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --source value, -s value, --from value  Source language to translate [$STRANS_SOURCE_LANG]
   --target value, -t value, --to value    Target language to translate [$STRANS_TARGET_LANG]
   --instance value, -i value              Instance host name of SimplyTranslate [$STRANS_INSTANCE]
   --repl, -r                              Start interactive mode (default: false)
   --help, -h                              show help (default: false)
   --version, -v                           print the version (default: false)
```

Various options can specify default values not only from command line arguments, 
but also by setting environment variables.

```sh
STRANS_TARGET_LANG="ja" strans "Hello, World" # => "こんにちは世界"

export STRANS_SOURCE_LANG="en"
export STRANS_TARGET_LANG="ja"
strans "Hello, World" # => こんにちは世界
```

This tool is a program that calls Lingva Translate's public API. 
See the Lingva Translate [README.md](https://github.com/thedaviddelta/lingva-translate/blob/main/README.md) for details.

## Installation

```sh
go install github.com/sheepla/strans@latest
```

## Thanks

- [lingva-translate](https://github.com/thedaviddelta/lingva-translate)
- [translate-shell](https://github.com/soimort/translate-shell)


