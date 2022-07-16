;;
;; This is a simple FizzBuzz example, which we can execute.
;;
;; You'll see here that we can define functions, that we have
;; primitives such as "zero?" and that we have a built-in "cond"
;; function too.
;;
;; cond here will take a list, which is processed in pairs:
;;
;;  (cond
;;    (quote
;;      TEST1  ACTION1
;;      TEST2  ACTION2
;;    )
;;  )
;;
;; For each pair (e.g. `TEST1 ACTION1`) we run the first statement, and if
;; the result is `true` we evaluate the action, and stop.
;;
;; When the test returns nil/false/similar then we continue running until
;; we do get success.  That means it is important to end with something that
;; will always succeed.
;;
;; `(quote) is used to ensure we don't evaluate the list in advance of the
;; statement.
;;

;; Is the given number divisible by 3?
;;
;; Note that we add ":number" to the end of the argument, which means
;; a fatal error will be raised if we invoke this function with a non-number,
;; for example:
;;
;;   (divByThree "Steve")
;;   (divByThree true)
;;
(define divByThree (lambda (n:number) (zero? (% n 3))))

;; Is the given number divisible by 5?
(define divByFive  (lambda (n:number) (zero? (% n 5))))

;; Run the fizz-buzz test for the given number, N
(define fizz (lambda (n:number)
  (cond
    (quote
      (and (divByThree n) (divByFive n))  (print "fizzbuzz")
      (divByThree n)                      (print "fizz")
      (divByFive  n)                      (print "buzz")
      #t                                  (print n)))))


;; Apply the function fizz, for each number 1-50
(apply (nat 51) fizz)
