;;; golang.lisp - Implement compatibility functions for code moved to golang

(alias getenv os.Getenv)
(help "getenv" "Return the contents of the specified environmental variable")

(alias setenv os.Setenv)
(help "setenv" "Set the value of the specified environmental variable")

(alias random rand.Intn)
(help "random" "Return a number between zero and one less than the value specified.

See also: random:char random:item
Example: (random 100) ; A number between 0 and 99")
