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
                     "Returns true if the argument specified is a hash."
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
(set! >  (fn* (a b)
              "Return true if the first value is greater than the second."
              (< b a)))

(set! >= (fn* (a b)
              "Return true if the first value is greater than, or equal to the second."
              (! (< a b))))
(set! <= (fn* (a b)
              "Return true if the first value is less than, or equal to, the second."
              (! (> a b))))

;; We have a built in function "date" to return the current date
;; as a list (DD MM YYYY).  We also ahve a builtin function (time)
;; to return the time as a list (HH MM SS).
;;
;; Here we create some helper functions for retrieving the various
;; parts of the date/time, as well as some aliases for ease of typing.
(set! date:year (fn* ()
                "Return the current year, as an integer."
                (nth (date) 3)))
(set! year date:year)

(set! date:month (fn* ()
                 "Return the number of the current month, as an integer."
                 (nth (date) 2)))
(set! month date:month)

(set! date:day (fn* ()
               "Return the day of the current month, as an integer."
               (nth (date) 1)))
(set! day date:day)

(set! date:weekday (fn* ()
                   "Return a string containing the current day of the week."
                   (nth (date) 0)))
(set! weekday date:weekday)

(set! time:hour (fn* ()
                "Return the current hour, as an integer."
                (nth (time) 0)))
(set! hour time:hour)

(set! time:minute (fn* ()
                  "Return the current minute, as an integer."
                  (nth (time) 1)))
(set! minute time:minute)

(set! time:second (fn* ()
                  "Return the current seconds, as an integer."
                  (nth (time) 2)))
(set! second time:second)

(set! time:hms (fn* ()
               "Return the current time, formatted as 'HH:MM:SS', as a string."
               (sprintf "%s:%s:%s" (hour) (minute) (second))))
(set! hms time:hms)


;;
;; This is a bit sneaky.  NOTE there is no short-circuiting here.
;;
;; Given a list use `filter` to return those items which are "true".
;;
;: If the length of the input list, and the length of the filtered list
;; are the same then EVERY element was true so our AND result is true.
;;
(set! and (fn* (xs:list)
               "Return true if every item in the specified list is true."
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
              "Return true if any value in the specified list contains a true value."
              (let* (res nil)
                (set! res (filter xs (lambda (x) (if x true false))))
                (if (> (length res) 0)
                    true
                  false))))


;; every is useful and almost a logical operation
(set! every (fn* (xs:list fun:function)
                 "Return true if applying every element of the list through the specified function resulted in a true result"
                 (let* (res (map xs fun))
                   (if (and res)
                       true
                     false))))


;; Useful for creating a list of numbers
(set! repeated (fn* (n:number x)
                  "Return a list of length n whose elements are all x"
                  (when (pos? n)
                    (cons x (repeated (dec n) x)))))

;; inc/dec are useful primitives to have
(set! inc (fn* (n:number)
               "inc will add one to the supplied value, and return the result."
               (+ n 1)))

(set! dec (fn* (n:number)
               "dec will subtract one from the supplied value, and return the result."
               (- n 1)))

;; We could also define the incr/decr operations as macros.
(defmacro! incr (fn* (x)
                     "incr is a macro which will return the given value, incremented by one."
                     `(set! ~x (+ ~x 1))))

(defmacro! decr (fn* (x)
                     "decr is a macro which will return the given value, decremented by one."
                     `(set! ~x (- ~x 1))))

;; Not is useful
(set! ! (fn* (x)
             "Return the inverse/NOT of the given boolean value"
             (if x #f #t)))

;; Square root
(set! sqrt (fn* (x:number)
                "Calculate the square root of the given value."
                (# x 0.5)))

;; Return the last element of a list
(set! last (fn* (lst:list)
                "last returns the last item in the specified list, it is the opposite of cdr."
                (let* (c (cdr lst))
                  (if (! (nil? c))
                      (last c)
                    (car lst)))))

;; Setup a simple function to run a loop N times
;;
(set! repeat (fn* (n body)
                  "Execute the supplied body of code N times."
                  (if (> n 0)
                      (do
                          (body n)
                          (repeat (- n 1) body)))))

;; A helper to apply a function to each key/value pair of a hash
(set! apply-hash (fn* (hs:hash fun:function)
                      "Apply given function to every key in the specified hash"
                      (let* (lst (keys hs))
                        (apply lst (lambda (x) (fun x (get hs x)))))))


;; Count the length of a string
(set! strlen (fn* (str:string)
                  "Calculate and return the length of the supplied string."
                  (length (split str "" ))))


;; More mathematical functions relating to negative numbers.
(set! neg  (fn* (n:number)
                "Negate the supplied number, and return it."
                (- 0 n)))

(set! neg? (fn* (n:number)
                "Return true if the supplied number is negative."
                (< n 0)))

(set! pos? (fn* (n:number)
                "Return true if the supplied number is positive."
                (> n 0)))

(set! abs  (fn* (n:number)
                "Return the absolute value of the supplied number."
                (if (neg? n) (neg n) n)))

(set! sign (fn* (n:number)
                "Return 1 if the specified number is positive, and -1 if it is negative."
                (if (neg? n) (neg 1) 1)))


;; Create ranges of numbers in a list
(set! range (fn* (start:number end:number step:number)
                 "Create a list of numbers between the start and end bounds, incrementing by the given offset each time."
                 (if (< start end)
                     (cons start (range (+ start step) end step))
                   ())))

;; Create sequences from 0/1 to N
(set! seq (fn* (n:number)
               "Create, and return, list of number ranging from 0-N"
               (range 0 n 1)))
(set! nat (fn* (n:number)
               "Create, and return, a list of numbers ranging from 1 to N."
               (range 1 n 1)))


;; Remove items from a list where the predicate function is not T
(set! filter (fn* (xs:list f:function)
                  "Remove any items from the specified list, if the result of calling the provided function on that item is not true."
                  (if (nil? xs)
                      ()
                      (if (f (car xs))
                          (cons (car xs)(filter (cdr xs) f))
                          (filter (cdr xs) f)))))




;; reduce function
(set! reduce (fn* (xs f acc)
                  "Reduce"
                  (if (nil? xs)
                      acc
                      (reduce (cdr xs) f (f acc (car xs))))))

;; Now define min/max using reduce
(set! min (fn* (xs:list)
               "Return the smallest integer from the list of numbers supplied."
               (if (nil? xs)
                   ()
                 (reduce xs
                         (lambda (a b)
                           (if (< a b) a b))
                         (car xs)))))

(set! max (fn* (xs:list)
               "Return the maximum integer from the list of numbers supplied."
               (if (nil? xs)
                   ()
                 (reduce xs
                         (lambda (a b)
                           (if (< a b) b a))
                         (car xs)))))


; O(n^2) behavior with linked lists
(set! append (fn* (xs:list el)
                  "Append the given element to the specified list"
                  (if (nil? xs)
                      (list el)
                      (cons (car xs) (append (cdr xs) el)))))


(set! reverse (fn* (x:list)
                   "Return a list containing all values in the supplied list, in reverse order."
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
                     "Translate each character in the given string, via the means of the supplied lookup-table.  This is used by 'upper' and 'lower'."
                     (let* (chrs (split x ""))
                       (join (map chrs (lambda (x)
                                         (if (get hsh x)
                                             (get hsh x)
                                           x)))))))

;; Convert the given string to upper-case, via the lookup table.
(set! upper (fn* (x:string)
                 "Convert each character from the supplied string to upper-case, and return that string."
                 (translate x upper-table)))

;; Convert the given string to upper-case, via the lookup table.
(set! lower (fn* (x:string)
                 "Convert each character from the supplied string to lower-case, and return that string."
                (translate x lower-table)))


;; Get the first N items from a list.
(set! take (fn* (n l)
                "Return the first N items from the specified list."
                (cond (zero? n) nil
                      (nil? l) nil
                      true (cons (car l) (take (- n 1) (cdr l))))))

;; Remove the first N items from a list.
(set! drop (fn* (n l)
                "Remove the first N items from the specified list."
                (cond (zero? n) l
                      (nil? l) nil
                      true (drop (- n 1) (cdr l)))))

;; Return everything but the last element.
(set! butlast (fn* (l)
                   "Return everything but the last element from the specified list."

                   (take (dec (length l)) l)))

;; Wrappers for our file functions
(set! file:stat:size (fn* (path)
                          "Return the size of the given file, return -1 on error."
                          (let* (info (file:stat path))
                            (cond
                             (nil? info) -1
                             true (nth info 1)))))

(set! file:stat:uid (fn* (path)
                          "Return the UID of the given file owner, return '' on error."
                          (let* (info (file:stat path))
                            (cond
                             (nil? info) ""
                             true (nth info 2)))))


(set! file:stat:gid (fn* (path)
                          "Return the GID of the given file owner, return '' on error."
                          (let* (info (file:stat path))
                            (cond
                             (nil? info) ""
                             true (nth info 3)))))

(set! file:stat:mode (fn* (path)
                          "Return the mode of the given file, return '' on error."
                          (let* (info (file:stat path))
                            (cond
                             (nil? info) ""
                             true (nth info 4)))))

; Slurp used to be a primitive for reading file contents
(set! slurp file:read)
