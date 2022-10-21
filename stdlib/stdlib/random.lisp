;;; random.lisp - Random things.


;; Choose a random character from a string, or a-z if unspecified
(set! random:char (fn* (&x)
                       "Return a random character by default from the set a-z.

If an optional string is provided it will be used as a list of characters to choose from."
                       (let* (chars (split "abcdefghijklmnopqrstuvwxyz" ""))
                         (if (list? x)
                             (set! chars (split (car x) "")))
                         (random:item chars))))

;; random list item
(set! random:item (fn* (lst:list)
                       "Return a random element from the specified list."
                       (nth lst (random (length lst)))))
