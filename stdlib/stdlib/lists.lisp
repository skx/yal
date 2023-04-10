;;; lists.lisp - Some list-utility functions

;; Some of these functions were adapted from Rob Pike's lisp
;;
;;    https://github.com/robpike/lisp
;;
;; which in turn were derived from code in "LISP 1.5 Programmer's Manual"
;; by McCarthy, Abrahams, Edwards, Hart, and Levin, from MIT in 1962.
;;



(set! find (fn* (item lst)
                "Return the offsets of any occurence of the item in the given list, nil on failure.

See-also: intersection, member, occurrences, union."
                     (let* (len (length lst)
                                res (list      ))
                        (repeat len (lambda (n)
                                      (if (eq (nth lst (- n 1)) item)
                                          (set! res (cons (- n 1) res) true))))
                        (if (= 0 (length res)) nil res ))))


(set! flatten (fn* (L)
                   "Converts a list of nested lists to a single list, flattening it."
                   (if (nil? L)
                       nil
                     (if (! (list? (first L)))
                         (cons (first L) (flatten (rest L)))
                       (append (flatten (first L)) (flatten (rest L)))))))


(set! intersection (fn* (x y)
                        "Return the values common to the two specified lists.

See-also: find, member, occurrences, union."
                        (cond
		         (nil? x) nil
		         (member (car x) y) (cons (car x) (intersection (cdr x) y))
		         true (intersection (cdr x) y))))

(set! member (fn* (item lst)
                  "Return true if the specified item is found within the given list.

See-also: find, intersection, occurrences, union."

                  (cond
		   (nil? lst) false
		   (eq item (car lst)) true
		   true (member item (cdr lst)))))


;;
;; NOTE: This could be implemented as follows:
;;
;;   (set! occurrences (fn* (item lst) (length (find item lst))))
;;
(set! occurrences (fn* (item lst)
                 "Count the number of times the given item is found in the specified list.

See-also: find, intersection, member, union"
                 (if lst
                     (+
                      (if (eq item (car lst)) 1 0)
                      (occurrences item (cdr lst)))
                     0)))

(set! union (fn* (x y)
      "Return the union of the two specified lists.

See-also: find, intersection, member, occurrences."
      (cond
       (nil? x) y
       (member (car x) y) (union (cdr x) y)
       true (cons (car x) (union (cdr x) y))
       )))
