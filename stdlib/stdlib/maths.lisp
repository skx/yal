;;; maths.lisp - Some simple maths-related primitives

;; inc/dec are useful primitives to have
(set! inc (fn* (n:number)
               "inc will add one to the supplied value, and return the result."
               (+ n 1)))

(set! dec (fn* (n:number)
               "dec will subtract one from the supplied value, and return the result."
               (- n 1)))

;; PI
(set! pi (fn* ()
              "Return the value of PI, calculated via arctan, as per https://en.m.wikibooks.org/wiki/Trigonometry/Calculating_Pi"
              (* 4 (+ (* 6 (atan (/ 1 8))) (* 2 (atan (/ 1 57))) (atan (/ 1 239))))
              ))

;; Square root
(set! sqrt (fn* (x:number)
                "Calculate the square root of the given value."
                (# x 0.5)))

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
