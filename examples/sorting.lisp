;;; sorting.lisp - Demonstrate generating random lists, and sorting them.

;;
;; This example demonstrates creating lists of random (integer)
;; numbers, and then sorting them.
;;
;; We have three sorting methods to test:
;;
;;  1.  insert-sort, implemented in lisp.
;;
;;  2.  quick-sort, implemented in lisp.
;;
;;  3.  sort, implemented in golang
;;
;;
;; See-also the "sorting.lisp" file in our standard-library, which is
;; an unfolded version of this quicksort - with the use of a user-defined
;; comparison function.
;;

(set! random:list (fn* (n max)
                       "Return a list of random numbers, of the length n, ranging from -max to +max."
                       (map (nat n) (lambda (n)
                                      (let* (sign 1)
                                        ; optionally this might be negative
                                        (if (= 0 (random 2))
                                            (set! sign -1))
                                        (* sign (random max)))))))



;;
;; insertion-sort
;;

(set! insert (fn* (item lst)
                  "Insert the specified item into the given list, in the correct (sorted) order."
                  (if (nil? lst)
                      (cons item lst)
                    (if (> item (car lst))
                        (cons (car lst) (insert item (cdr lst)))
                      (cons item lst)))))

(set! insertsort (fn* (lst)
                      "An insert-sort implementation.  For each item in the given list, use insert to place it into the list in the correct order."
                      (if (nil? lst)
                          nil
                        (insert (car lst) (insertsort (cdr lst))))))



;;
;; quick-sort
;;

(set! append3 (fn* (a b c)
                   "Like append, but with three items, not two."
                   (append a (append b c))))

(set! list>= (fn* (m list)
                  "Return all items of the given list which are greater than, or equal to, the specified item."
                  (filter list (lambda (n) (! (< n m))))))


(set! list< (fn* (m list)
                 "Return all items of the given list which are less than the specified item."
                 (filter list (lambda (n) (< n m)))))

(set! qsort (fn* (l)
                 "A recursive quick-sort implementation."
                 (if (nil? l)
                     nil
                   (append3 (qsort (list< (car l) (cdr l)))
                            (cons (car l) null)
                            (qsort (list>= (car l) (cdr l)))))))




;;
;; Now we have defined functions to generate a list of random integers,
;; and we also have our two sorting methods, so we can test things :)
;;
;; Spoiler qsort is the faster lisp sort, but the native golang sort
;; is significantly faster - due to lack of recursion & etc.
;;
(let* (
       count 512                     ; how many items to work with
       lst (random:list count 4096)  ; create a random list of integers
                                     ; between -4096 and +4096.

       bis (ms)                      ; before-insert-sort take a timestamp
       isrt (insertsort lst)         ;   run insert-sort
       ais (ms)                      ; after insert-sort take a timestamp

       bqs (ms)                      ; before-quick-sort take a timestamp
       qsrt (qsort lst)              ;   run the quick-sort
       aqs (ms)                      ; after-quick-sort take a timestamp

       bgs (ms)                      ; before-go-sort take a timestamp
       gsrt (sort lst)               ;   run the go sort
       ags (ms)                      ; after-go-sort take a timestamp
       )
  (print "insert sort took %d ms " (- ais bis))
  (print "quick sort took %d ms " (- aqs bqs))
  (print "go sort took %d ms " (- ags bgs))

  ; a simple sanity-check of the results
  (print "Testing results, to ensure each sort produced identical results.")
  (print "offset,insert-sort,quick-sort,go-sort")
  (apply (seq count)
         (lambda (x)
           (let* (a (nth isrt x)
                  b (nth qsrt x)
                  c (nth gsrt x))

                  (if (! (eq a b))
                      (print "List element %d differs!" x))
                  (if (! (eq b c))
                      (print "List element %d differs!" x))
                  (print "%d,%d,%d,%d" x a b c))))
  (print "All done!")
  )
