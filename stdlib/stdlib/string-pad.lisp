;;; string-pad.lisp - String padding, prefix and postfix, functions.

;;
;; This file contains functions for padding a string to a specified
;; length, using a supplied character-string to extend it.
;;


(set! pad:left (fn* (str add len)
                    "Pad the given string to a specified length, by pre-pending the given string to it.

See also: pad:right"
                    (if (>= (strlen str) len)
                        val
                      (pad:left (join (list str val)) add len))))


(set! pad:right (fn* (str add len)
                     "Pad the given string to the specified length, by repeatedly appending the given char to the value.

See also: pad:left"
                     (if (>= (strlen str) len)
                         val
                       (pad:right (join (list str add)) add len))))
