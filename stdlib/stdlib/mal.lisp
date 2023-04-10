;;; mal.lisp - Compatability with MAL, implemented in lisp.


;; Traditionally we use `car` and `cdr` for accessing the first and rest
;; elements of a list.  For readability it might be nice to vary that
(alias first car)
(alias rest  cdr)


;; Run an arbitrary series of statements, if the given condition is true.
;;
;; This is the more general/useful version of the "if2" macro, which
;; we demonstrate in mtest.lisp.
;;
;; Sample usage:
;;
;;  (when (= 1 1) (print "OK") (print "Still OK") (print "final statement"))
;;
(defmacro! when (fn* (pred &rest)
                     "when is a macro which runs the specified body, providing the specified predicate is true.

It is similar to an if-statement, however there is no provision for an 'else' clause, and the body specified may contain more than once expression to be evaluated."
                     `(if ~pred (do ~@rest))))

;;
;; If the specified predicate is true, then run the body.
;;
;; NOTE: This recurses, so it will eventually explode the stack.
;;
(defmacro! while (fn* (condition &body)
                      "while is a macro which repeatedly runs the specified body, while the condition returns a true-result."
                      (let* (inner-sym (gensym))
                        `(let* (~inner-sym (fn* ()
                                                (if ~condition
                                                    (do
                                                        ~@body
                                                        (~inner-sym)))))
                           (~inner-sym)))))


;;
;; cond is a useful thing to have.
;;
(defmacro! cond (fn* (&xs)
                     "cond is a macro which accepts a list of conditions and results, and returns the value of the first matching condition.  It is similar in functionality to a C case-statement."
                     (if (> (length xs) 0)
                         (list 'if (first xs)
                               (if (> (length xs) 1)
                                   (nth xs 1)
                                 (error "An odd number of forms to (cond..)"))
                               (cons 'cond (rest (rest xs)))))))

;; A useful helper to apply a given function to each element of a list.
(set! apply (fn* (lst:list fun:function)
                 "Call the specified function on every element in the given list.

See-also: apply-pairs, apply-hash"
                 (if (nil? lst)
                     ()
                     (do
                      (fun (car lst))
                      (apply (cdr lst) fun)))))

;; Apply, but walking the list in pairs.
(set! apply-pairs (fn* (lst:list fun:function)
                       "Calling the specified function with two items on the specified list.

This is similar to apply, but apply apply invokes the callback with a single list-item, and here we apply in pairs.

Note: The list-length must be even, and if not that will raise an error.

See-also: apply apply-hash
Example: (apply-pairs (list 1 2 3 4) (lambda (a b) (print \"Called with %s %s\" a b)))
"
                       (if (! (nil? lst))
                           (if (= (% (length lst) 2) 0)
                               (let* (a (car lst)
                                      b (car (cdr lst)))
                                 (fun a b)
                                 (apply-pairs (cdr (cdr lst) ) fun))
                               (error "The list passed to (apply-pairs..) should have an even length")))))

;; Return the length of the given list.
(set! length (fn* (arg)
                  "Return the length of the supplied list.  See-also strlen."
                  (if (list? arg)
                      (do
                          (if (nil? arg) 0
                            (inc (length (cdr arg)))))
                    0
                    )))

(alias count length)


(set! map (fn* (lst:list fun:function)
               "Return a list with the contents of evaluating the given function on every item of the supplied list.

See-also: map-pairs"
               (if (nil? lst)
                   ()
                 (cons (fun (car lst)) (map (cdr lst) fun)))))

(set! map-pairs (fn* (lst:list fun:function)
               "Return a list with the contents of evaluating the given function on every pair of items in the supplied list.

See-also: map"
               (if (! (nil? lst))
                   (if (= (% (length lst) 2) 0)
                       (let* (a (car lst)
                                b (car (cdr lst)))
                         (cons (fun a b) (map-pairs (cdr (cdr lst)) fun)))
                     (error "The list passed should have an even length"))
                 ())))


;; This is required for our quote/quasiquote/unquote/splice-unquote handling
;;
;; Testing is hard, but
;;
;; (define lst (quote (b c)))                      ; b c
;; (print (quasiquote (a lst d)))                  ; (a lst d)
;; (print (quasiquote (a (unquote lst) d)))        ; (a (b c) d)
;; (print (quasiquote (a (splice-unquote lst) d))) ; (a b c d)
;;
(set! concat (fn* (seq1 seq2)
                  "Join two lists"
                  (if (nil? seq1)
                      seq2
                    (cons (car seq1) (concat (cdr seq1) seq2)))))
