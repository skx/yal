;;; readme.lisp - Generate a README.md file, based on directory contents.

;;
;; This directory contains *.lisp, echo of which has a header-line prefixed
;; with three semi-colons:
;;
;;     ;;; foo.lisp - Information
;;
;; This script reads those files and outputs a simple index, of the filename
;; and the information.
;;

(set! lisp:files (fn* ()
                      "Return a list of all the lisp files in the current directory"
                      (sort (glob "*.lisp"))))


(set! lisp:info (fn* (file)
                     "Output a brief overview of the given file"
                     (let* (
                            text (file:lines file)
                            line  (nth text 0)
                            info (match "^(.*)-+[ ]+(.*)$" line))
                       (when (list? info)
                         (print "* [%s](%s)" file file)
                         (print "  * %s" (nth info 2))))))


(set! lisp:index (fn* ()
                      "Generate a README.md snippet"
                      (let* (files (lisp:files))
                        (apply files lisp:info()))))

(print "# Examples\n")
(print "This directory contains some simple lisp examples, which can be executed via `yal`.\n\n")

(lisp:index)
