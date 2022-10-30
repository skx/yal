;;; sorting.lisp - Demonstrate generating random lists, and sorting them.

;;
;; This example demonstrates creating lists of random (integer)
;; numbers, and sorting them with two different methods.
;;
;; Our random numbers range from -1024 to +1024.
;;
;; We have three sorting methods to test:
;;
;;  1.  insert-sort, implemented in lisp.
;;
;;  2.  quick-sort, implemented in lisp.
;;
;;  3.  sort, implemented in golang
;;


(set! random:number (fn* ()
                         "Return a random number from -1024 to 1024"
                         (let* (sign 1
                                max  1024)
                           ; optionally this might be negative
                           (if (= 0 (random 2))
                               (set! sign -1))
                           (* sign (random max)))))

(set! random:list (fn* (n)
                       "Return a random list of numbers, of the length n."
                       (map (seq n) random:number)))



;;
;; insertion-sort
;;

(set! insert (fn* (item lst)
                  (if (nil? lst)
                      (cons item lst)
                    (if (> item (car lst))
                        (cons (car lst) (insert item (cdr lst)))
                      (cons item lst)))))

(set! insertsort (fn* (lst)
                      (if (nil? lst)
                          nil
                        (insert (car lst) (insertsort (cdr lst))))))




;;
;; quick-sort
;;

;; append in our standard library makes "(list b)" not "b".  Changing that
;; would cause issues, so I'm duplicating here.
(set! append2 (fn* (a b)
                  (if (nil? a)
                      b
                    (cons (car a) (append2 (cdr a) b) ))))

(set! append3 (fn* (a b c)
                   (append2 a (append2 b c))))

(set! list>= (fn* (m list)
                  (filter list (lambda (n) (! (< n m))))))

(set! list< (fn* (m list)
                 (filter list (lambda (n) (< n m)))))

(set! qsort (fn* (l)
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
       lst (random:list 512)       ; create a random list - 512 elements

           bis (ms)                ; before-insert-sort take a timestamp
           isrt (insertsort lst)   ;   run insert-sort
           ais (ms)                ; after insert-sort take a timestamp

           bqs (ms)                ; before-quick-sort take a timestamp
           qsrt (qsort lst)        ;   run the quick-sort
           aqs (ms)                ; after-quick-sort take a timestamp

           bgs (ms)                ; before-go-sort take a timestamp
           gsrt (sort lst)         ;   run the go sort
           ags (ms)                ; after-go-sort take a timestamp
       )
  (print "insert sort took %d ms " (- ais bis))
  (print "quick sort took %d ms " (- aqs bqs))
  (print "go sort took %d ms " (- ags bgs))

  ; a simple sanity-check of the results
  (print "Testing results - first 20 - to make sure all sorts were equal")
  (print "insert-sort,quick-sort,go-sort")
  (apply (nat 20)
         (lambda (x)
           (let* (a (nth isrt x)
                  b (nth qsrt x)
                  c (nth gsrt x))

                  (if (! (eq a b))
                      (print "List element %d differs!" x))
                  (if (! (eq b c))
                      (print "List element %d differs!" x))
                  (print "%d,%d,%d" a b c))))
  (print "All done!")
  )
