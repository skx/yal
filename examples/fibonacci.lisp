;;; fibonacci.lisp - Calculate the first 25 fibonacci numbers.

;;
;; This is a sample input file for our minimal lisp interpreter.
;;
;; We use it to demonstrate and test some basic features.
;;
;; Here we use "while" from our standard library, and have defined a
;; function to turn "1" into "1st", etc, as appropriate.  This uses our
;; "match" primitive, which is implemented in golang.
;;


;; Add a suitable suffix to a number.
;;
;; e.g.  1 -> 1st
;;      11 -> 11th
;;      21 -> 21st
;;     333 -> 333rd
(set! add-numeric-suffix (fn* (n)
                              "Add a trailing suffix to make a number readable."
                              (cond
                                (match "(^|[^1]+)1$" n) (sprintf "%dst" n)
                                (match "(^|[^1]+)2$" n) (sprintf "%dnd" n)
                                (match "(^|[^1]+)3$" n) (sprintf "%drd" n)
                                true  (sprintf "%dth" n)
                                )))

;; Fibonacci function
(set! fibonacci (fn* (n)
                     "Calculate the Nth fibonacci number."
                     (if (<= n 1)
                         n
                       (+ (fibonacci (- n 1)) (fibonacci (- n 2))))))


;; Now call our function in a loop, twenty times.
(let* (n 1)
  (while (<= n 25)
    (print "%s fibonacci number is %d" (add-numeric-suffix n) (fibonacci n))
    (set! n (+ n 1) true)))
