# strans

A command line translator with GNU Readline like interactive mode (`--repl`) inspired by [translate-shell](https://github.com/soimort/translate-shell)

## Usage

This tool supports both interactive and non-interactive usage.

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
   --engine value, -e value                Name of translate engine [$STRANS_ENGINE]
   --instance value, -i value              Instance host name of SimplyTranslate [$STRANS_INSTANCE]
   --repl, -r                              Start interactive mode (default: false)
   --help, -h                              show help (default: false)
   --version, -v                           print the version (default: false)
```


## Non-interactive mode

If you specify text as a non-option argument, that text will be translated. Multiple arguments are allowed, and arguments are joined by spaces.

```
strans [OPTIONS] TEXT...

# e.g.
strans -s en -t ja "Hello, World" # => "こんにちは世界"
```

A non-option argument of `-` will read text from standard input and translate it.

```
strans [OPTIONS] -

# e.g.
echo "Hello, World" | strans -s en -t ja - # => "こんにちは世界"
strans -s en -t ja < README.md # => The contents of the README.md will translated.
```

## Interactive mode

You can use GNU Readline-like interactive mode.

Execute the command with the `-r`, `--repl` flag. 
Enter your text and it will be translated instantly.

Empty inputs are ignored and no translation is performed.

Typing `Ctrl-D` exits interactive mode and returns you to the shell you were running from.

```sh
[you@your-computer]$ strans --repl -t ja
REPL mode. Type Ctrl-D to exit.
>
> Hello, World
こんにちは世界
>
> The quick brown fox jumps over the lazy dog.
素早い茶色のキツネが怠け者の犬を飛び越えます。
>
>
[you@your-computer]$ 
```

## Options

Specify the source language name (e.g. `en`, `ja`, etc.) for `--source` option, and specify the target language name for `--target` option.

Various options can specify default values not only from command line arguments, 
but also by setting environment variables.

```sh
STRANS_TARGET_LANG="ja" strans "Hello, World" # => "こんにちは世界"
export STRANS_TARGET_LANG="ja"
strans "Hello, World" # => こんにちは世界
```

## Installation

```sh
go install github.com/sheepla/strans@latest
```

## Thanks

- [SimplyTranslate](https://simple-web.org/projects/simplytranslate.html)
- [translate-shell](https://github.com/soimort/translate-shell)


