;;; stack.lisp - A simple stack implemented as a list.


(defmacro! stack:push (fn* (val stck)
                           "Push the given value to the top of the stack.

This is a destructive operation.

See-also: stack:empty? stack:pop stack:size
"
   `(set! ~stck (cons ~val ~stck))))


(defmacro! stack:pop (fn* (stck)
                          "Remove and return the item from the head of the stack.  If the stack is empty nil is returned.

This is a destructive operation.

See-also: stack:empty? stack:push stack:size
"
   `(let* (val (if (list? ~stck) (car ~stck) nil)
           rem (if (list? ~stck) (cdr ~stck) nil))
           (set! ~stck rem t)  ; note global set
           val)))

(set! stack:empty? (fn* (stck)
                        "Return true if the stack is empty, false otherwise.

See-also: stack:pop stack:push stack:size"
                        (if (zero? (count stck))
                            true
                          false)))

(set! stack:size (fn* (stck)
                      "Return the number of entries present on the stack.

See-also: stack:empty? stack:pop stack:push"
                      (count stck)))
