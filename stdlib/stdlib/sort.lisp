;;; sort.lisp - Implementation of quick-sort with a user-defined comparison.

(set! sort-by (fn* (cmp:function l:list)
                   "sort-by is a generic quick-sort implementation, which makes use of a user-defined comparison method.

The function specified will be called with two arguments, and should return true if the first is less than the second.

See-also: sort"
                 (if (nil? l)
                     nil
                     (let* (cur (car l))
                       (append (sort-by cmp (filter (cdr l) (lambda (n) (cmp n cur))))
                               (append (cons (car l) null)
                                        (sort-by cmp (filter (cdr l) (lambda  (n) (! (cmp n cur)))))))))))
