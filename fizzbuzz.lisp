;;; fizzbuzz2.lisp - A simple fizzbuzz implementation

;; Taking advantage of our (cond) primitive we can just return the
;; string to print for any given number.

(set! fizzbuzz (fn* (n)
    (print "%s"
     (cond
      (= 0 (% n 15)) "fizzbuzz"
      (= 0 (% n 3))  "fizz"
      (= 0 (% n 5))  "buzz"
      true n) )))


;; Apply the function to each number 1-50
(apply (nat 51) fizzbuzz)
