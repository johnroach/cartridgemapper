# Endeca Cartridge Mapper

[![Build Status](https://travis-ci.org/JohnRoach/cartridgemapper.svg?branch=master)](https://travis-ci.org/JohnRoach/cartridgemapper)

This cool tool scans a given Endeca Application an generates documentation regarding the cartridges in HTML or JSON form.

The documentation includes information such as:
- Name of cartridge
- ID of cartridge
- Description of cartridge
- Endeca rules that use cartridge
- Sites that use cartridge
- Pages that use said cartridge


Sample run:

```
$ ./cartridgemapper
Endeca cartridge mapper maps the Endeca application cartridge usage.
This is very useful in understanding how cartridges are used and which cartridges
are available. One would use this tool to point to a given Endeca application and
get an ouput of a certain format.

For example:
    cartridgemapp mapEndecaApp /full/path/to/endeca/App/lication
    cartridgemapp mapEndecaApp /full/path/to/endeca/App/lication --output json

Usage:
  cartridgemapper [command]

Available Commands:
  help         Help about any command
  mapEndecaApp mapEndecaApp maps the Endeca cartridges used in an Endeca Application
  version      Print the version number of cartridgemapper

Flags:
      --config string   config file (default is $HOME/.cartridgemapper.yaml)
      --debug           adding debug to the logging
      --disable-color   disable color for logging output
  -h, --help            help for cartridgemapper

Use "cartridgemapper [command] --help" for more information about a command.
```
