;;; fizzbuzz2.lisp - Show the fizzbuzz up to 50.

;; Taking advantage of our (cond) primitive we can just return the
;; string to print for any given number.

(set! fizzbuzz (fn* (n)
                    "This function outputs the appropriate fizzbuzz-response
for the specified number.

'fizz' when the number is divisible by three, 'buzz' when divisible by five,
and 'fizzbuzz' when divisible by both."
                    (print "%s"
                           (cond
                            (= 0 (% n 15)) "fizzbuzz"
                            (= 0 (% n 3))  "fizz"
                            (= 0 (% n 5))  "buzz"
                            true (str n)) )))


;; As you can see the function above contains some help-text, or overview.
;; we can output that like so:
(print (help fizzbuzz))

;; Apply the function to each number 1-50
(apply (nat 51) fizzbuzz)
