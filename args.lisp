#!/usr/bin/env yal
;;
;;  Usage:
;;
;;     ./args.lisp 2 34 4
;;

;; Show the count.
(print "I received %d command-line arguments." (length os.args))

;; Show the actual arguments
(print "Args: %s" os.args)

;; And followup with the username
(print "The current user is %s, running on %s (arch:%s)"
       (getenv "USER")
       (os)
       (arch)
       )
