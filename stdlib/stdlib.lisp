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
(define function? (lambda (x) (eq (type x) "procedure")))
(define list?     (lambda (x) (eq (type x) "list")))
(define number?   (lambda (x) (eq (type x) "number")))
(define string?   (lambda (x) (eq (type x) "string")))
(define symbol?   (lambda (x) (eq (type x) "symbol")))


;; A useful helper to apply a given function to each element of a list.
(define apply (lambda (lst fun)
  (if (nil? lst)
    ()
      (begin
         (fun (car lst))
         (apply (cdr lst) fun)))))


(define !     (lambda (x) (if x #f #t)))
(define zero? (lambda (n) (if (= n 0) #t #f)))
(define even? (lambda (n) (if (zero? (% n 2)) #t #f)))
(define odd?  (lambda (n) (! (even? n))))

;; inc/dec are kinda useful
(define inc  (lambda (n) (- n 1)))
(define dec  (lambda (n) (+ n 1)))

;; We've defined "<" and ">" in golang, but not the or-equals variants.
;;
;; Add those.
(define >= (lambda (a b) (! (< a b))))
(define <= (lambda (a b) (! (> a b))))

;; More mathematical functions relating to negative numbers.
(define neg  (lambda (n) (- 0 n)))
(define neg? (lambda (n) (< n 0)))
(define abs  (lambda (n) (if (neg? n) (neg n) n)))
(define sign (lambda (n) (if (neg? n) (neg 1) 1)))
