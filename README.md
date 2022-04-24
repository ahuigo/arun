# Arun
:cloud: :fire: Live reload for any command.

    $ arun go run main.go

        ___
       /   |  _______  ______
      / /| | / ___/ / / / __ \
     / ___ |/ /  / /_/ / / / /
    /_/  |_/_/   \__,_/_/ /_/ v0.1.1 // live reload for any command, with Go1.14.0

    watching .
    !exclude docs
    watching hack
    watching hooks

## Motivation
Inspired by [air](https://github.com/cosmtrek/air). 
I want to build a simple live-reload tool that support run any type of file (e.g. go, python, js) when file is modified.

Great thanks to air :). 

## Features
* Support reload any type of command 
* Colorful log output
* Customize build or binary command
* Support excluding subdirectories
* Allow watching new directories after Arun started
## todo
- [] Loop time run
- [] Support clean up

## Installation

### Go

The classic way to install
```bash
go get -u github.com/ahuigo/arun
```

## Usage

    $ arun -h 
    Usage:
        arun [options] command arguments......
    options:
        -d <item> 	Debug Item
        -h			Help
        -i 			Ignore Directory
        -v <level>	Verbose Level
        -s			Keep Silent, without log output
    example:
        arun go run main.go


## Sponsor
Most arun's code is migrated from [air](https://github.com/cosmtrek/air). If you like it, you could buy Rick Yu a beer.

<a href="https://www.buymeacoffee.com/36lcNbW" target="_blank">
    <img src="https://cdn.buymeacoffee.com/buttons/default-orange.png" alt="Buy Me A Coffee" style="height: 51px !important;width: 217px !important;" >
</a>

## License
[GNU General Public License v3.0](LICENSE)
