;;; adder.lisp - Demonstrate creating an adder with closures.


;; Generate a function that adds two numbers.
;;
;; The first number is a constant, which is set when
;; the function is called.
;;
(set! make-adder (fn* (n)
   (fn* (m) (+ n m))))


;; Now create two adders.
;;
;; Here we see how the first value, N, is set.
;;
;; The second value in the definition above, M, is
;; set when the generated function is called.
(set! addFive (make-adder 5))
(set! addTen  (make-adder 10))


;;
;; Finally we can invoke our generated functions,
;; which work as you'd expect:
;;
(print "(+ 10 5) => %d" (addFive 10))
(print "(+ 10 3) => %d" (addTen   3))
