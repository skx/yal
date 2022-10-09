;;; mal.lisp - Compatability with MAL, implemented in lisp.

;; This is essentially prepended to any program the user tries to run,
;; and implements functions that are expected by any MAL implementation.
;;
;; More general functions can be found in stdlib.lisp.
;;


;; Traditionally we use `car` and `cdr` for accessing the first and rest
;; elements of a list.  For readability it might be nice to vary that
(set! first (fn* (x:list) (car x)))
(set! rest  (fn* (x:list) (cdr x)))

;; Some simple tests of numbers
(set! zero? (fn* (n) (= n 0)))
(set! one?  (fn* (n) (= n 1)))
(set! even? (fn* (n) (zero? (% n 2))))
(set! odd?  (fn* (n) (! (even? n))))

;; is the given argument "true", or "false"?
(def! true?  (fn* (arg) (if (eq #t arg) true false)))
(def! false? (fn* (arg) (if (eq #f arg) true false)))


;; Run an arbitrary series of statements, if the given condition is true.
;;
;; This is the more general/useful version of the "if2" macro, which
;; we demonstrate in mtest.lisp.
;;
;; Sample usage:
;;
;;  (when (= 1 1) (print "OK") (print "Still OK") (print "final statement"))
;;
(defmacro! when (fn* (pred &rest) `(if ~pred (do ~@rest))))

;;
;; If the specified predicate is true, then run the body.
;;
;; NOTE: This recurses, so it will eventually explode the stack.
;;
(defmacro! while (fn* (condition &body)
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
  (if (> (length xs) 0)
      (list 'if (first xs)
            (if (> (length xs) 1)
                (nth xs 1)
              (error "An odd number of forms to (cond..)"))
            (cons 'cond (rest (rest xs)))))))

;; A useful helper to apply a given function to each element of a list.
(set! apply (fn* (lst:list fun:function)
  (if (nil? lst)
      ()
      (do
         (fun (car lst))
         (apply (cdr lst) fun)))))


;; Return the length of the given list.
(set! length (fn* (arg)
  (if (list? arg)
    (do
      (if (nil? arg) 0
        (inc (length (cdr arg)))))
    0
    )))


;; Alias to (length)
(set! count (fn* (arg) (length arg)))


;; Find the Nth item of a list
(set! nth (fn* (lst:list i:number)
  (if (> i (length lst))
    (error "Out of bounds on list-length")
    (if (= 0 i)
      (car lst)
        (nth (cdr lst) (- i 1))))))


;; Replace a list with the contents of evaluating the given function on
;; every item of the list
(set! map (fn* (xs:list f:function)
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
  (if (nil? seq1)
      seq2
      (cons (car seq1) (concat (cdr seq1) seq2)))))



;;
;; Read a file
;;
(def! load-file (fn* (f)
                     (eval (join (list "(do " (slurp f) "\nnil)")))))
