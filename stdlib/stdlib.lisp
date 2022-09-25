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
(define function? (lambda (x) (or
                                (list
                                   (eq (type x) "procedure(lisp)")
                                   (eq (type x) "procedure(golang)")))))
(define hash?     (lambda (x) (eq (type x) "hash")))
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

;;
;; This is a bit sneaky.  NOTE there is no short-circuiting here.
;;
;; Given a list use `filter` to return those items which are "true".
;;
;: If the length of the input list, and the length of the filtered list
;; are the same then EVERY element was true so our AND result is true.
;;
(define and (lambda (xs:list)
  (let ((res nil))
    (set! res (filter xs (lambda (x) (if x true false))))
    (if (= (length res) (length xs))
        true
      false))))

;;
;; This is also a bit sneaky.  NOTE there is no short-circuiting here.
;;
;; Given a list use `filter` to return those items which are "true".
;;
;; If the output list has at least one element that was true then the
;; OR result is true.
;;
(define or (lambda (xs:list)
  (let ((res nil))
    (set! res (filter xs (lambda (x) (if x true false))))
    (if (> (length res) 0)
        true
      false))))


;; Traditionally we use `car` and `cdr` for accessing the first and rest
;; elements of a list.  For readability it might be nice to vary that
(define first (lambda (x:list) (car x)))
(define rest  (lambda (x:list) (cdr x)))

;; inc/dec are useful primitives to have
(define inc  (lambda (n:number) (+ n 1)))
(define dec  (lambda (n:number) (- n 1)))

;; We could also define the incr/decr operations as macros.
(define incr (macro (x) `(set! ~x (+ ~x 1))))
(define decr (macro (x) `(set! ~x (- ~x 1))))

;; Not is useful
(define !     (lambda (x) (if x #f #t)))

;; Some simple tests of numbers
(define zero? (lambda (n) (= n 0)))
(define one?  (lambda (n) (= n 1)))
(define even? (lambda (n) (zero? (% n 2))))
(define odd?  (lambda (n) (! (even? n))))


;; Square root
(define sqrt (lambda (x:number) (# x 0.5)))



;;
;; if2 is a simple macro which allows you to run two actions if an
;; (if ..) test succeeds.
;;
;; This means you can write:
;;
;;   (if2 true (print "1") (print "2"))
;;
;; Instead of having to use (begin), like so:
;;
;;   (if true (begin (print "1") (print "2")))
;;
;; The downside here is that you don't get a negative branch, but running
;; two things is very common - see for example the "(while)" and "(repeat)"
;; macros later in this file.
;;
(define if2 (macro (pred one two)
  `(if ~pred (begin ~one ~two))))


;;
;; Run an arbitrary series of statements, if the given condition is true.
;;
;; This is the more general/useful version of the "if2" macro, given above.
;;
;; Sample usage:
;;
;;  (when (= 1 1) (print "OK") (print "Still OK") (print "final statement"))
;;
(define when (macro (pred &rest) `(if ~pred (begin ~@rest))))

;;
;; Part of our while-implementation.
;; If the specified predicate is true, then run the body.
;;
;; NOTE: This recurses, so it will eventually explode the stack.
;;
;; NOTE: We use "if2" not "if".
;;
(define while-fun (lambda (predicate body)
  (if2 (predicate)
    (body)
    (while-fun predicate body))))

;;
;; Now a macro to use the while-fun body as part of a while-function
;;
;; NOTE: We use "if2" not "if".
;;
(define while (macro (expression body)
                     (list 'while-fun
                           (list 'lambda '() expression)
                           (list 'lambda '() body))))


;;
;; cond is a useful thing to have.
;;
(define cond (macro (&xs)
  (if (> (length xs) 0)
      (list 'if (first xs)
            (if (> (length xs) 1)
                (nth xs 1)
              (error "An odd number of forms to (cond..)"))
            (cons 'cond (rest (rest xs)))))))

;; Setup a simple function to run a loop N times
;;
;; NOTE: We use "if2" not "if".
;;
(define repeat (lambda (n body)
  (if2 (> n 0)
     (body n)
     (repeat (- n 1) body))))

;; A useful helper to apply a given function to each element of a list.
(define apply (lambda (lst:list fun:function)
  (if (nil? lst)
      ()
      (begin
         (fun (car lst))
         (apply (cdr lst) fun)))))

;; A helper to apply a function to each key/value pair of a hash
(define apply-hash (lambda (hs:hash fun:function)
  (let ((lst (keys hs)))
    (apply lst (lambda (x) (fun x (get hs x)))))))

;; Return the length of the given string or list.
(define length (lambda (arg)
  (if (list? arg)
    (begin
      (if (nil? arg) 0
        (inc (length (cdr arg)))))
    0
    )))


;; Alias to (length)
(define count (lambda (arg) (length arg)))


;; Count the length of a string
(define strlen (lambda (str:string) (length (split str "" ))))


;; Find the Nth item of a list
(define nth (lambda (lst:list i:number)
  (if (> i (length lst))
    (error "Out of bounds on list-length")
    (if (= 0 i)
      (car lst)
        (nth (cdr lst) (- i 1))))))

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
    ()
      (append (reverse (cdr x)) (car x)))))

;;
;; This is either gross or cool.
;;
;; Define a hash which has literal characters and their upper-case, and
;; lower-cased versions
;;
(define upper-table {
  a "A"
  b "B"
  c "C"
  d "D"
  e "E"
  f "F"
  g "G"
  h "H"
  i "I"
  j "J"
  k "K"
  l "L"
  m "M"
  n "N"
  o "O"
  p "P"
  q "Q"
  r "R"
  s "S"
  t "T"
  u "U"
  v "V"
  w "W"
  x "X"
  y "Y"
  z "Z"
  } )

(define lower-table {
  A "a"
  B "b"
  C "c"
  D "d"
  E "e"
  F "f"
  G "g"
  H "h"
  I "i"
  J "j"
  K "k"
  L "l"
  M "m"
  N "n"
  O "o"
  P "p"
  Q "q"
  R "r"
  S "s"
  T "t"
  U "u"
  V "v"
  W "w"
  X "x"
  Y "y"
  Z "z"
  } )


;; Translate the elements of the string using the specified hash
(define translate (lambda (x:string hsh:hash)
  (let ((chrs (split x "")))
    (join (map chrs (lambda (x)
                  (if (get hsh x)
                      (get hsh x)
                    x)))))))

;; Convert the given string to upper-case, via the lookup table.
(define upper (lambda (x:string)
                (translate x upper-table)))

;; Convert the given string to upper-case, via the lookup table.
(define lower (lambda (x:string)
                (translate x lower-table)))


;; This is required for our quote/quasiquote/unquote/splice-unquote handling
;;
;; Testing is hard, but
;;
;; (define lst (quote (b c)))                      ; b c
;; (print (quasiquote (a lst d)))                  ; (a lst d)
;; (print (quasiquote (a (unquote lst) d)))        ; (a (b c) d)
;; (print (quasiquote (a (splice-unquote lst) d))) ; (a b c d)
;;
(define concat (lambda (seq1 seq2)
  (if (nil? seq1)
      seq2
      (cons (car seq1) (concat (cdr seq1) seq2)))))
