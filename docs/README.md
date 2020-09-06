# GORK

Gork is an experimental text adventure game system.

# BUILD

To build you first need to make and install goyacc and nex and put them in the bin directory.

To build the binary:

    ./build

To see some debug info as you build (e.g. lexer conflicts):

    ./build -v

To clean the intermediary files (for checking into Git etc):

    ./build --clean

# USAGE

Gork currently reads stdin (this will change)

    ./gork test/gork.test
