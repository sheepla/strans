# strans

A command line translator with GNU Readline like interactive mode (`--repl`) inspired by [translate-shell](https://github.com/soimort/translate-shell)

## Usage

Non-interactive mode usage: 

Specify the translation source language, translation destination language, engine name and instance name in the options and specify the text in arguments execute.


```
NAME:
   strans - a command line SimplyTranslate client with GNU Readline like interactive mode

USAGE:
   strans [-e|--engine ENGINE][-i|--instance INSTANCE] [-s|--source SOURCE_LANG] -t|--target TARGET_LANG TEXT...
   strans [-r|--repl]

VERSION:
   unknown

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --source value, -s value, --from value  Source language to translate [$STRANS_SOURCE_LANG]
   --target value, -t value, --to value    Target language to translate [$STRANS_TARGET_LANG]
   --engine value, -e value                Name of translate engine (google, libre, deepl)
   --instance value, -i value              Instance host name of SimplyTranslate [$STRANS_INSTANCE]
   --repl, -r                              Start bash-like REPL mode (default: false)
   --help, -h                              show help (default: false)
   --version, -v                           print the version (default: false)
Required flag "to" not set
```

Interactive mode usage: Execute the command with the `-r`, `--repl` flag. 
Enter your text and it will be translated instantly. Exit with `Ctrl-D`

```sh
[you@your-computer]$ strans --repl
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

## Installation

```sh
go install github.com/sheepla/strans@latest
```

## Thanks

- [SimplyTranslate](https://simple-web.org/projects/simplytranslate.html)
- [translate-shell](https://github.com/soimort/translate-shell)


