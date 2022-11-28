
* [Language Server Provider](#language-server-provider)
  * [Our LSP Features](#our-lsp-features)
* [Configuration](#configuration)
  * [Emacs](#configuration-emacs)
  * [NeoVim](#configuration-neovim)
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


## Our LSP Features

We only support the bare minimum LSP features:

* Provide completion for the names of all standard-library functions.
* Shows information on standard-library functions, on hover.




# Configuration

To use our LSP implementation you'll need to configure your editor, IDE, or environment appropriately.  Configuration will vary depending on what you're using.

Typically configuration will involve at least:

* Specifying the type of files that hould use LSP (i.e. a filename suffixes).
* Specifying the name/arguments to use for the LSP server.

For our implementation you'll need to launch "`yal -lsp`" to startup the LSP-process.


## Configuration: Emacs

For GNU Emacs the following file should provide all the help you need:

* [_misc/yal.el](_misc/yal.el)


## Configuration: neovim

Create the file `~/.config/nvim/init.lua` with contents as below:

* [_misc/init.lua](_misc/init.lua)



# Screenshots

Here we see what completion might look like:

![Completion](_misc/complete.png?raw=true "Completion")

Here's our help-text being displayed on-hover:

![Help](_misc/help.png?raw=true "Help")



# Notes

As stated above we only support hover-text, and completion, from the
standard library.  Supporting the users' own code is harder because that
would involve evaluating it - and that might cause side-effects.

It should be noted that our completion-support is very naive - it literally
returns the names of __all__ available methods, and relies upon the editor
to narrow down the selection - that seems to work though.



# See Also

* [README.md](README.md)
  * More details of the project.
* [PRIMITIVES.md](PRIMITIVES.md)
  * The list of built-in functions, whether implemented in Golang or YAL.
* [INTRODUCTION.md](INTRODUCTION.md)
  * Getting started setting variables, defining functions, etc.
