;;
;; This is a simple example of the YAL interpreter
;;
;; You'll see here that we can define functions, that we have
;; primitives such as "zero?" and that we have a built-in "cond"
;; function too.
;;
;; cond here will take a list, which is processed in pairs:
;;
;;  (cond
;;    (quote
;;      EVAL1  ACTION1
;;      EVAL2  ACTION2
;;    )
;;  )
;;
;; We take each pair "EVAL1 ACTION1", or "EVAL2 ACTION2", and if the
;; result of evaluating the first part is true we run the action.
;;
;; If not we continue down the list.  Quote is used to ensure we don't
;; evaluate the list in advance.
;;

;; Is the given number divisible by 3?
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
(apply (nat 50) fizz)
