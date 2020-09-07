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

Gork currently reads stdin (this will change)

    ./gork < test/test.gork
