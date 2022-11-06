;;; lists.lisp - Some list-utility functions

;; These were adapted from Rob Pike's lisp
;;
;;    https://github.com/robpike/lisp
;;
;; which in turn were derived from code in "LISP 1.5 Programmer's Manual"
;; by McCarthy, Abrahams, Edwards, Hart, and Levin, from MIT in 1962.
;;


(set! member (fn* (item list)
                  "Return true if the specified item is found within the given list.

See-also: intersection, union."

                  (cond
		   (nil? list) false
		   (eq item (car list)) true
		   true (member item (cdr list)))))

(set! union (fn* (x y)
      "Return the union of the two specified lists.

See-also: intersection, member"
      (cond
       (nil? x) y
       (member (car x) y) (union (cdr x) y)
       true (cons (car x) (union (cdr x) y))
       )))

(set! intersection (fn* (x y)
                        "Return the values common to both specified lists

See-also: member, union."
                        (cond
		         (nil? x) nil
		         (member (car x) y) (cons (car x) (intersection (cdr x) y))
		         true (intersection (cdr x) y))))
