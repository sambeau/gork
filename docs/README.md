# Gork

Experimental text adventure game system.

## Status

Currently gork is a glorified pretty printer. It can parse simple game scripts and print them out. It doesn't actually run them yet.

(this will change)

## Build

To install the tools and build the binary:

    ./build

To see some debug info as you build (e.g. parser conflicts):

    ./build -v

To clean the intermediary files (for checking into Git etc):

    ./build --clean

## Usage

Gork can be used to compile and run game files. Source files can be run directly or can be compiled to a binary format.

    Gork! - Experimental text adventure interpreter

    Usage:
      gork <file>
      gork -r <sourcefile>
      gork -c <sourcefile> [-o <gamefile>]
      gork -h | --help | --version

    Options:
      -h --help     Show this screen.
      --version     Show version.
      -c --compile  Compile source file.
      -o --outfile  Where to save compiled game.
      -r --run      Run source file directly.
