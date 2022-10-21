;;; stdlib.lisp - Standard library as implemented in lisp.




;; Useful for creating a list of numbers
(set! repeated (fn* (n:number x)
                    "Return a list of length n whose elements are all x."
                  (when (pos? n)
                    (cons x (repeated (dec n) x)))))

;; Return the last element of a list
;;
;; NOTE: This could be written more simply, for example:
;;
;;   (set! last (fn* (lst:list) "Return the last element of the given list" (car (reverse lst))))
;;
(set! last (fn* (lst:list)
                "last returns the last item in the specified list, it is the inverse of (butlast) and the logical opposite of (car)."
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
                      "Call the given function to every key in the specified hash.

See-also: apply, apply-pairs"
                      (let* (lst (keys hs))
                        (apply lst (lambda (x) (fun x (get hs x)))))))


;; Count the length of a string
(set! strlen (fn* (str:string)
                  "Calculate and return the length of the supplied string."
                  (length (split str ""))))


;; Create ranges of numbers in a list
(set! range (fn* (start:number end:number step:number)
                 "Create a list of numbers between the start and end bounds, inclusive, incrementing by the given offset each time."
                 (if (<= start end)
                     (cons start (range (+ start step) end step))
                   ())))

;; Create sequences from 0/1 to N
(set! seq (fn* (n:number)
               "Create, and return, list of number ranging from 0-N, inclusive."
               (range 0 n 1)))
(set! nat (fn* (n:number)
               "Create, and return, a list of numbers ranging from 1-N, inclusive."
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
                  "This is our reduce function, which uses a list, a function, and the accumulator."
                  (if (nil? xs)
                      acc
                      (reduce (cdr xs) f (f acc (car xs))))))

; O(n^2) behavior with linked lists
(set! append (fn* (xs:list el)
                  "Append the given element to the specified list."
                  (if (nil? xs)
                      (list el)
                      (cons (car xs) (append (cdr xs) el)))))


(set! reverse (fn* (x:list)
                   "Return a list containing all values in the supplied list, in reverse order."
                   (if (nil? x)
                       ()
                     (append (reverse (cdr x)) (car x)))))


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
