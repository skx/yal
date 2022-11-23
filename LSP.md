
* [Language Server Provider](#language-server-provider)
  * [LSP Features](#lsp-features)
  * [Configuration](#configuration)
    * [Emacs](#configuration-emacs)
  * [Screenshots](#screenshots)
  * [Notes](#notes)
  * [See Also](#see-also)


# Language Server Provider

LSP is a new standard which makes it easier for editors to support
advanced features in a portable way.

In brief a "language server" provides the ability to show help,
provide completions, etc, and a users' editor just needs to speak the
appropriate protocol to get all those features - without implementing them
directly.



## LSP Features

We support a minimal LSP implementation that:

* Provides completion for the names of all standard-library functions.
* Shows information on standard-library functions, on hover.




## Configuration

The specific configuration will depend upon which editor/environment you're using.
You'll want to configure things such that "`yal -lsp`" is launched to provide the LSP-support though.


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
