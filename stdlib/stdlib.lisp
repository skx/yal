;;
;; stdlib.lisp - Standard library executed alongside any user-supplied
;; program.
;;
;; This implements behaviour which is useful for users.
;;

;; There is a built in `type` function which returns the type of an object.
;;
;; Use this to define some simple methods to test argument-types
(define boolean?  (lambda (x) (eq (type x) "boolean")))
(define error?    (lambda (x) (eq (type x) "error")))
(define function? (lambda (x) (or (eq (type x) "procedure(lisp)")
                                  (eq (type x) "macro")
                                  (eq (type x) "procedure(golang)"))))
(define macro?    (lambda (x) (eq (type x) "macro")))
(define list?     (lambda (x) (eq (type x) "list")))
(define number?   (lambda (x) (eq (type x) "number")))
(define string?   (lambda (x) (eq (type x) "string")))
(define symbol?   (lambda (x) (eq (type x) "symbol")))

;; We've defined "<" in golang, we can now implement the missing
;; functions in terms of that:
;;
;; >
;; <=
;; >=
;;
(define >  (lambda (a b) (< b a)))
(define >= (lambda (a b) (! (< a b))))
(define <= (lambda (a b) (! (> a b))))


;; Traditionally we use `car` and `cdr` for accessing the first and rest
;; elements of a list.  For readability it might be nice to vary that
(define first (lambda (x:list) (car x)))
(define rest  (lambda (x:list) (cdr x)))

;; inc/dec are useful primitives to have
(define inc  (lambda (n:number) (+ n 1)))
(define dec  (lambda (n:number) (- n 1)))

;; Not is useful
(define !     (lambda (x) (if x #f #t)))

;; Some simple tests of numbers
(define zero? (lambda (n) (= n 0)))
(define one?  (lambda (n) (= n 1)))
(define even? (lambda (n) (zero? (% n 2))))
(define odd?  (lambda (n) (! (even? n))))


;; Square root
(define sqrt (lambda (x:number) (# x 0.5)))


;; A useful helper to apply a given function to each element of a list.
(define apply (lambda (lst:list fun)
  (if (nil? lst)
      ()
      (begin
         (fun (car lst))
         (apply (cdr lst) fun)))))


;; Return the length of the given list
(define length (lambda (arg:list)
   (if (nil? arg) 0
      (inc (length (cdr arg))))))

;; Find the Nth item of a list
(define nth (lambda (lst:list i:number)
  (if (= 0 i)
    (car lst)
      (nth (cdr lst) (- i 1)))))

;; More mathematical functions relating to negative numbers.
(define neg  (lambda (n:number) (- 0 n)))
(define neg? (lambda (n:number) (< n 0)))
(define pos? (lambda (n:number) (> n 0)))
(define abs  (lambda (n:number) (if (neg? n) (neg n) n)))
(define sign (lambda (n:number) (if (neg? n) (neg 1) 1)))


;; Create ranges of numbers in a list
(define range (lambda (start:number end:number step:number)
  (if (< start end)
     (cons start (range (+ start step) end step))
        ())))

;; Create sequences from 0/1 to N
(define seq (lambda (n:number) (range 0 n 1)))
(define nat (lambda (n:number) (range 1 n 1)))


;; Remove items from a list where the predicate function is not T
(define filter (lambda (xs:list f:function)
  (if (nil? xs)
     ()
     (if (f (car xs))
        (cons (car xs)(filter (cdr xs) f))
           (filter (cdr xs) f)))))

;; Replace a list with the contents of evaluating the given function on
;; every item of the list
(define map (lambda (xs:list f:function)
  (if (nil? xs)
     ()
       (cons (f (car xs)) (map (cdr xs) f)))))


;; reduce function
(define reduce (lambda (xs f acc)
  (if (nil? xs)
    acc
      (reduce (cdr xs) f (f acc (car xs))))))

;; Now define min/max using reduce
(define min (lambda (xs:list)
  (if (nil? xs)
    ()
      (reduce xs
        (lambda (a b)
           (if (< a b) a b))
              (car xs)))))

(define max (lambda (xs:list)
  (if (nil? xs)
    ()
      (reduce xs
        (lambda (a b)
           (if (< a b) b a))
              (car xs)))))


; O(n^2) behavior with linked lists
(define append (lambda (xs:list el)
  (if (nil? xs)
    (list el)
      (cons (car xs) (append (cdr xs) el)))))


(define reverse (lambda (x:list)
  (if (nil? x)
    x
      (append (reverse (cdr x)) (car x)))))
