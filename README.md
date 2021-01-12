# inputServer

:warning: **ALPHA** this is a personal project, you may use it but it is not a production release

This is a simple implementation of [robotgo](https://github.com/go-vgo/robotgo) combined with an http server

I needed to send keystrokes to my Raspberry Pi from my PC, but could find what I was looking for so I decided to make it myself.

## Features

For keys reference please see [this](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md) document from [robotgo](https://github.com/go-vgo/robotgo)

- /tap?key=\<mainkey>&args=\<arg1>,\<arg2>,...
Emulates keystrokes. "key" is the primary key that should be pressed. "args" is an optional comma seperetated list of keys to be pressed simultaneous
  `/tap?key=r&args="shift,ctrl"`
- /print?value=\<phrase> Types out the given value `/print?value=test`

## Installation

:warning: If you don't know how to do this you probably shouldn't use it as this project is in alpha

1. Pull this repo
2. Build for your needed platform
3. Enjoy ;D

## Configuration

The default port is `25139`

You can optionally add a config.yaml file either alongside the binary or in data folder alongside the binary:

`port: ":25139"` Specify the port you like, has to be preceede by a ":"

`verbosity: 2` Specify what should be printed to console and logged: 0 - Only Errors, 1 - Errors and Warnings, 2 - Everything

`log-file: "./log.txt"` Specify path to a .txt-file for logging, not required! If left empty or unset no log will be generated