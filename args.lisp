#!/home/skx/go/bin/yal
;;
;;  Usage:
;;
;;     ./args.lisp 2 34 4
;;

;; Show the count.
(print "I received %s command-line arguments." (length os.args))

;; Show the actual arguments
(print "Args: %s" os.args)

;; And followup with the username
(print "User: %s" (getenv "USER"))
