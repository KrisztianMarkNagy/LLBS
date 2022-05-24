
# LLBS (Low Level Build System)

An amateur build system made mostly for C/CPP/ASM type of languages. It allows to build a binary-executable easily supporting local modules for code sharding, local libraries on your machine, external modules and external libraries too. The external libraries or modules would represent packages/code/... that you are not the author of, while the local libraries and modules would represent otherwise.

Potential feature to be added later on (if I do not forget and/or do not get lazy) would be a custom compilation process for external/local libraries/modules, in case a certain piece of code requires special attention.

## The reason I developed LLBS

It all starts with C and Makefile, when I was getting more and more familiar and accustomed to C and Makefile, and having fun playing around with C, I wanted the ability to shard my code-base in case I want to make a bigger project, and when I went a bit deeper into makefiles in order to see if it was possible to do what LLBS does with an existing tool I noticed that it simply did not work (or maybe I'm just stupid, meh who knows). Makefile did not have what I needed to compile C source files from many different directories and their subdirectories into one final executable.

While I know that there is CMake and other build systems one could use, even so, what I required was simplicity for a simple problem. And then I made it.

## Refactor (for a more readable code-base)

Coming soon...?


## Installation

To install LLBS, you can simply run `go install .` or you can use the Makefile to install it to `/usr/bin/`.
Read the Makefile with full confidence since it's not that big of a deal.

The four commands of the Makefile:
- install
- clean
- reinstall
- uninstall

At the moment, the Makefile's installation process only supports Linux.

