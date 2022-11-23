
* [Language Server Provider](#language-server-provider)
  * [LSP Features](#lsp-features)
  * [Configuration](#configuration)
    * [Emacs](#configuration-emacs)
  * [Screenshots](#screenshots)
  * [Notes](#notes)
  * [See Also](#see-also)


# Language Server Provider

Adding features like auto complete, go to definition, or documentation
on hover for a programming language takes significant effort. Traditionally
this work had to be repeated for each development tool, as each tool
provides different APIs for implementing the same feature.

A Language Server is meant to provide the language-specific smarts and
communicate with development tools over a protocol that enables
inter-process communication.

The idea behind the Language Server Protocol (LSP) is to standardize
the protocol for how such servers and development tools communicate.
This way, a single Language Server can be re-used in multiple
development tools, which in turn can support multiple languages with
minimal effort.


## LSP Features

We support a minimal LSP implementation that:

* Provides completion for the names of all standard-library functions.
* Shows information on standard-library functions, on hover.




## Configuration

The specific configuration will depend upon which editor/environment you're using.

Typically configuration will involve specifying at least the type of files that
should use LSP (i.e. based on filename suffixes), and specifying the way to launch
the LSP-server, or communicate with a long-running one.

For our implementation you'll need to launch "`yal -lsp`" to startup the LSP-process.


### Configuration Emacs

For GNU Emacs the following file should provide all the help you need:

* [_misc/yal.el](_misc/yal.el)




## Screenshots

Here we see what completion might look like:

![Completion](_misc/complete.png?raw=true "Completion")

Here's our help-text being displayed on-hover:

![Help](_misc/help.png?raw=true "Help")



## Notes

TODO



## See Also

* [README.md](README.md)
  * More details of the project.
* [PRIMITIVES.md](PRIMITIVES.md)
  * The list of built-in functions, whether implemented in Golang or YAL.
* [INTRODUCTION.md](INTRODUCTION.md)
  * Getting started setting variables, defining functions, etc.
