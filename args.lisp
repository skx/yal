#!/home/skx/go/bin/yal
;;
;;  Usage:
;;
;;     ./args.lisp 2 34 4
;;

(print "I received %s command-line arguments." (length os.args))
(print "Args: %s" os.args)
