;;; fuzz.lisp - Generate random expressions and evaluate them, forever.


;;
;; Cache of functions that the fuzzer can evaluate.
;;
(set! functions nil)

;;
;; Retrieve the values of functions that can be executed, using the
;; cache, or calculating as appropriate.
;;
(set! get_fns (fn* ()
                   "Get the functions we can fuzz-eval with a cache."
                   (if (nil? functions)
                       (do
                        (set! functions (fns) true)
                        functions
                        )
                     functions)))

;;
;; Actually look over our builtins, stdlib, and specials, to return the list
;; of functions that are safe for execution,
;;
;; Specifically we want to rule out functions that never terminate.
;;
(set! fns (fn* ()
               "Return specials/builtins/stdlib functions that are safe for evaluation.

We remove exit to avoid termination, forever to avoid infinite loops, and while for the same reason."
               (flatten (map (append (append (specials) (builtins)) (stdlib))
                             (lambda (x)
                               (cond
                                (eq (str x) "exit") nil
                                (eq (str x) "while") nil
                                (eq (str x) "forever") nil
                                true x))))))


;;
;; Generate a random argument.
;;
(set! arg (fn* ()
             "Generate a random argument; a number, a list, a string, etc."
             (let* (c (random:item '(1 2 3 4)))
               (cond
                (= 1 c) (random:item '("foo" 0 nil true false 1 2 3 4 5 6 7 8 9 10 '( 1 2 3 )))
                (= 2 c) (random:item '(0 1 2 3 4 5 6 7 8 9))
                (= 3 c) nil
                (= 4 c) (flatten (args))))))

;;
;; Generate a random list of arguments, of random length.
;;
(set! args (fn* ()
               "Generate a random number of arguments, and return as a list"
               (let* (c (random:item '(0 1 2 3 4 5 6 7 8 9 10))
                      r (list))
                 (repeat c (lambda (n)
                             (set! r (list r (random:item '("foo" 0 nil true false 1 2 3 4 '( 1 2 3 )))) true)
                             ))
                 (list 'list r))))



;;
;; Real core - generate a random function, with up to three arguments.
;;
;; Output it, and execute it.
;;
(set! gen (fn* ()
               "Generate a random sexp, with 1-3 arguments.
Show the expression and evaluate it."
               (let* (
                     name   (random:item (get_fns))
                     param1 (arg)
                     param2 (arg)
                     param3 (arg)
                     )
                 (let* (c (random:item '(1 2 3)))
                   (cond
                    (= 1 c)
                      (do
                       (print (list name param1))
                       (eval (list name param1)))
                    (= 2 c)
                      (do
                       (print (list name param1 param2))
                       (eval (list name param1 param2)))
                    (= 3 c)
                      (do
                       (print (list name param1 param2 param3))
                       (eval (list name param1 param2 param3))))))))



;;
;; Invoke a single fuzz iteration, catching any/all errors which are produced.
;;
(set! invoke (fn* ()
                  "Invoke a generated function once, catching any errors."
                  (try
                   (gen)
                   (catch e
                     (print "\t%s" e)))))


;;
;; Now we fuzz, forever ..
;;
(forever (invoke))
