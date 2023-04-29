;;; golang.lisp - Implement compatibility functions for code moved to golang

(alias getenv os.Getenv)
(help "getenv" "Return the contents of the specified environmental variable")

(alias setenv os.Setenv)
(help "setenv" "Set the value of the specified environmental variable")
