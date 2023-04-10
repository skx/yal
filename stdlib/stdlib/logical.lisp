;;; logical.lisp - Logical functions.


;; Not is useful
(set! not (fn* (x)
             "Return the inverse of the given boolean value."
             (if x #f #t)))

(alias ! not)


;; This is a bit sneaky.  NOTE there is no short-circuiting here.
;;
;; Given a list use `filter` to return those items which are "true".
;;
;: If the length of the input list, and the length of the filtered list
;; are the same then EVERY element was true so our AND result is true.
(set! and (fn* (xs:list)
               "Return true if every item in the specified list is true.

NOTE: This is not a macro, so all arguments are evaluated."
               (let* (res nil)
                 (set! res (filter xs (lambda (x) (if x true false))))
                 (if (= (length res) (length xs))
                     true
                   false))))

(alias && and)

;; This is also a bit sneaky.  NOTE there is no short-circuiting here.
;;
;; Given a list use `filter` to return those items which are "true".
;;
;; If the output list has at least one element that was true then the
;; OR result is true.
(set! or (fn* (xs:list)
              "Return true if any value in the specified list contains a true value.

NOTE: This is not a macro, so all arguments are evaluated."
              (let* (res nil)
                (set! res (filter xs (lambda (x) (if x true false))))
                (if (> (length res) 0)
                    true
                  false))))

(alias || or)


;; every is useful and almost a logical operation
(set! every (fn* (xs:list fun:function)
                 "Return true if applying every element of the list through the specified function resulted in a true result.

NOTE: This is not a macro, so all arguments are evaluated."
                 (let* (res (map xs fun))
                   (if (and res)
                       true
                     false))))
