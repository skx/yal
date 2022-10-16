;;; mal.lisp - Compatability with MAL, implemented in lisp.

;; This is essentially prepended to any program the user tries to run,
;; and implements functions that are expected by any MAL implementation.
;;
;; More general functions can be found in stdlib.lisp.
;;


;; Traditionally we use `car` and `cdr` for accessing the first and rest
;; elements of a list.  For readability it might be nice to vary that
(set! first (fn* (x:list)
                 "Return the first element of the specified list.  This is an alias for 'car'."
                 (car x)))

(set! rest (fn* (x:list)
                 "Return all elements of the specified list, except the first.  This is an alias for 'cdr'."
                 (cdr x)))

;; Some simple tests of numbers
(set! zero? (fn* (n)
                 "Return true if the number supplied as the first argument to this function is equal to zero."
                 (= n 0)))

(set! one? (fn* (n)
                "Return true if the number supplied as the argument to this function is equal to one."
                (= n 1)))

(set! even? (fn* (n)
                 "Return true if the number supplied as the argument to this function is even."
                 (zero? (% n 2))))

(set! odd?  (fn* (n)
                 "Return true if the number supplied as the argument to this function is odd."
                 (! (even? n))))

;; is the given argument "true", or "false"?
(def! true?  (fn* (arg)
                  "Return true if the argument supplied to this function is true."
                  (if (eq #t arg) true false)))

(def! false? (fn* (arg)
                  "Return true if the argument supplied to this function is false."
                  (if (eq #f arg) true false)))


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
                     "when is a macro which runs the specified body, providing the specified predicate is true.    It is similar to an if-statement, however there is no provision for an 'else' clause, and the body specified may contain more than once expression to be evaluated."
                     `(if ~pred (do ~@rest))))

;;
;; If the specified predicate is true, then run the body.
;;
;; NOTE: This recurses, so it will eventually explode the stack.
;;
(defmacro! while (fn* (condition &body)
                      "while is a macro which repeatedly runs the specified body, while the condition returns a true-result"
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
                 "Return the result of calling the specified function on every element in the given list"
                 (if (nil? lst)
                     ()
                     (do
                      (fun (car lst))
                      (apply (cdr lst) fun)))))


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

;; Find the Nth item of a list
(set! nth (fn* (lst:list i:number)
               "Return the Nth item of the specified list.

Note that offset starts from 0, rather than 1, for the first item."
               (if (> i (length lst))
                   (error "Out of bounds on list-length")
                 (if (= 0 i)
                     (car lst)
                   (nth (cdr lst) (- i 1))))))


(set! map (fn* (xs:list f:function)
               "Return a list with the contents of evaluating the given function on every item of the supplied list."
               (if (nil? xs)
                   ()
                 (cons (f (car xs)) (map (cdr xs) f)))))


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



;;
;; Read a file
;;
(def! load-file (fn* (filename)
                     "Load and execute the contents of the supplied filename."
                     (eval (join (list "(do " (slurp filename) "\nnil)")))))
