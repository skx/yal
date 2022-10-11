;;; stdlib.lisp - Standard library as implemented in lisp.
;; This is essentially prepended to any program the user tries to run,
;; and implements behaviour which is useful for users.
;;
;; For compatability with MAL mal.lisp is also considered part of our
;; standard-library.
;;



;; There is a built in `type` function which returns the type of an object.
;;
;; Use this to define some simple methods to test argument-types
(set! boolean?  (fn* (x)
                     "Returns true if the argument specified is a boolean value."
                     (eq (type x) "boolean")))

(set! error?    (fn* (x)
                     "Returns true if the argument specified is an error-value."
                     (eq (type x) "error")))

(set! function? (fn* (x) "Returns true if the argument specified is a function, either a built-in function, or a user-written one."
                     (or
                      (list
                       (eq (type x) "procedure(lisp)")
                       (eq (type x) "procedure(golang)")))))

(set! hash?     (fn* (x)
                     "Returns true if the argument specified is a hash"
                     (eq (type x) "hash")))

(set! macro?    (fn* (x)
                     "Returns true if the argument specified is a macro."
                     (eq (type x) "macro")))

(set! list?     (fn* (x)
                     "Returns true if the argument specified is a list."
                     (eq (type x) "list")))

(set! number?   (fn* (x)
                     "Returns true if the argument specified is a number."
                     (eq (type x) "number")))

(set! string?   (fn* (x)
                     "Returns true if the argument specified is a string."
                     (eq (type x) "string")))

(set! symbol?   (fn* (x)
                     "Returns true if the argument specified is a symbol."
                     (eq (type x) "symbol")))

;; We've defined "<" in golang, we can now implement the missing
;; functions in terms of that:
;;
;; >
;; <=
;; >=
;;
(set! >  (fn* (a b) (< b a)))
(set! >= (fn* (a b) (! (< a b))))
(set! <= (fn* (a b) (! (> a b))))

;; We have a built in function "date" to return the current date
;; as a list (DD MM YYYY).  We also ahve a builtin function (time)
;; to return the time as a list (HH MM SS).
;;
;; create some helper functions for retrieving the various parts of
;; the date/time.
(set! year (fn* () (nth (date) 3)))
(set! month (fn* () (nth (date) 2)))
(set! day (fn* () (nth (date) 1)))
(set! weekday (fn* () (nth (date) 0)))

(set! hour (fn* () (nth (time) 0)))
(set! minute (fn* () (nth (time) 1)))
(set! second (fn* () (nth (time) 2)))
(set! hms (fn* () (sprintf "%s:%s:%s" (hour) (minute) (second))))

;;
;; This is a bit sneaky.  NOTE there is no short-circuiting here.
;;
;; Given a list use `filter` to return those items which are "true".
;;
;: If the length of the input list, and the length of the filtered list
;; are the same then EVERY element was true so our AND result is true.
;;
(set! and (fn* (xs:list)
  (let* (res nil)
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
(set! or (fn* (xs:list)
  (let* (res nil)
    (set! res (filter xs (lambda (x) (if x true false))))
    (if (> (length res) 0)
        true
      false))))


;; inc/dec are useful primitives to have
(set! inc  (fn* (n:number) (+ n 1)))
(set! dec  (fn* (n:number) (- n 1)))

;; We could also define the incr/decr operations as macros.
(defmacro! incr (fn* (x) `(set! ~x (+ ~x 1))))
(defmacro! decr (fn* (x) `(set! ~x (- ~x 1))))

;; Not is useful
(set! !     (fn* (x) (if x #f #t)))

;; Square root
(set! sqrt (fn* (x:number) (# x 0.5)))

;; Return the last element of a list
(set! last (fn* (lst:list)
  (let* (c (cdr lst))
    (if (! (nil? c))
      (last c)
      (car lst)))))

;; Setup a simple function to run a loop N times
;;
(set! repeat (fn* (n body)
  (if (> n 0)
      (do
          (body n)
          (repeat (- n 1) body)))))

;; A helper to apply a function to each key/value pair of a hash
(set! apply-hash (fn* (hs:hash fun:function)
  (let* (lst (keys hs))
    (apply lst (lambda (x) (fun x (get hs x)))))))


;; Count the length of a string
(set! strlen (fn* (str:string) (length (split str "" ))))


;; More mathematical functions relating to negative numbers.
(set! neg  (fn* (n:number) (- 0 n)))
(set! neg? (fn* (n:number) (< n 0)))
(set! pos? (fn* (n:number) (> n 0)))
(set! abs  (fn* (n:number) (if (neg? n) (neg n) n)))
(set! sign (fn* (n:number) (if (neg? n) (neg 1) 1)))


;; Create ranges of numbers in a list
(set! range (fn* (start:number end:number step:number)
  (if (< start end)
     (cons start (range (+ start step) end step))
        ())))

;; Create sequences from 0/1 to N
(set! seq (fn* (n:number) (range 0 n 1)))
(set! nat (fn* (n:number) (range 1 n 1)))


;; Remove items from a list where the predicate function is not T
(set! filter (fn* (xs:list f:function)
  (if (nil? xs)
     ()
     (if (f (car xs))
        (cons (car xs)(filter (cdr xs) f))
           (filter (cdr xs) f)))))




;; reduce function
(set! reduce (fn* (xs f acc)
  (if (nil? xs)
    acc
      (reduce (cdr xs) f (f acc (car xs))))))

;; Now define min/max using reduce
(set! min (fn* (xs:list)
  (if (nil? xs)
    ()
      (reduce xs
        (lambda (a b)
           (if (< a b) a b))
              (car xs)))))

(set! max (fn* (xs:list)
  (if (nil? xs)
    ()
      (reduce xs
        (lambda (a b)
           (if (< a b) b a))
              (car xs)))))


; O(n^2) behavior with linked lists
(set! append (fn* (xs:list el)
  (if (nil? xs)
    (list el)
      (cons (car xs) (append (cdr xs) el)))))


(set! reverse (fn* (x:list)
  (if (nil? x)
    ()
      (append (reverse (cdr x)) (car x)))))

;;
;; This is either gross or cool.
;;
;; Define a hash which has literal characters and their upper-case, and
;; lower-cased versions
;;
(set! upper-table {
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

(set! lower-table {
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
(set! translate (fn* (x:string hsh:hash)
  (let* (chrs (split x ""))
    (join (map chrs (lambda (x)
                  (if (get hsh x)
                      (get hsh x)
                    x)))))))

;; Convert the given string to upper-case, via the lookup table.
(set! upper (fn* (x:string)
                (translate x upper-table)))

;; Convert the given string to upper-case, via the lookup table.
(set! lower (fn* (x:string)
                (translate x lower-table)))
