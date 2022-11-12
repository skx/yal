;;; string-substr.lisp - Fetch substrings from a string

;; String handling is good to have, and here we implement substr in the
;; naive way:
;;
;; Split the string into a list, and take parts of it using "take" and
;; "drop".
;;


(set! substr (fn* (str start &len)
                  "Return a substring of the given string, by starting index.

The length of the substring is optional."
                  (if (> start (strlen str))  ; out of bounds?
                      ""
                    (if (nil? len)  ; start at the given offset
                        (join (drop start (explode str)))
                      (join (take (car len) (drop start (explode str))))))))
