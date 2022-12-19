;;; comparisons.lisp - boolean, character, and numerical comparison functions.


;; We've defined "<" in natively, in golang.
;;
;; We can define the other numerical relational comparisons in terms of that.
(set! >  (fn* (a b)
              "Return true if a is greater than b."
              (< b a)))

(set! >= (fn* (a b)
              "Return true if a is greater than, or equal to b."
              (! (< a b))))
(set! <= (fn* (a b)
              "Return true if a is less than, or equal to, b."
              (! (> a b))))


;; We've defined "char<" in natively, in golang.
;;
;; We can define the other relational comparisons in terms of that.
(set! char>  (fn* (a b)
              "Return true if a is greater than b."
              (char< b a)))

(set! char>= (fn* (a b)
              "Return true if a is greater than, or equal to b."
              (! (char< a b))))
(set! char<= (fn* (a b)
              "Return true if a is less than, or equal to, b."
              (! (char> a b))))


;; We've defined "string<" in natively, in golang.
;;
;; We can define the other relational comparisons in terms of that.
(set! string>  (fn* (a b)
              "Return true if a is greater than b."
              (string< b a)))

(set! string>= (fn* (a b)
              "Return true if a is greater than, or equal to b."
              (! (string< a b))))
(set! string<= (fn* (a b)
              "Return true if a is less than, or equal to, b."
              (! (string> a b))))


;;
;; Some simple tests of specific numbers.
;;
(set! zero? (fn* (n)
                 "Return true if the number supplied as the first argument to this function is equal to zero."
                 (= n 0)))

(set! one? (fn* (n)
                "Return true if the number supplied as the argument to this function is equal to one."
                (= n 1)))

(set! even? (fn* (n)
                 "Return true if the number supplied as the argument to this function is even."
                 (zero? (% n 2))))

(set! odd?  (fn* (n)
                 "Return true if the number supplied as the argument to this function is odd."
                 (! (even? n))))

;;
;; Some simple tests of specific boolean results
;;
(def! true?  (fn* (arg)
                  "Return true if the argument supplied to this function is true."
                  (if (eq #t arg) true false)))

(def! false? (fn* (arg)
                  "Return true if the argument supplied to this function is false."
                  (if (eq #f arg) true false)))
