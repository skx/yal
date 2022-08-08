;; This file contains some simple examples of working with hashes.



;; Create a hash, with some details
(set! person { :name "Steve"
      :age 46
      :location "Helsinki"
      })


;; This function is used as a callback by apply-hash.
(define hash-element (lambda (key val)
   (print "KEY:%s VAL:%s" key val)))

;; The `apply-hash` function will trigger a callback for each key and value
;; within a hash.
;;
;; It is similar to the `apply` function which will apply a function to every
;; element of a lisp.
(apply-hash person hash-element)


;; Here we see a type-restriction, the following function can only be
;; invoked with a hash-argument.
(define blah (lambda (h:hash) (print "Got argument of type %s" (type h))))

;; Call it
(blah person)

;; Use get/set to update the hash properties
(print "Original name: %s" (get person :name))
(set person :name "Bobby")
(print "Updated  name: %s" (get person :name))

;; The (env) function returns a list of hashes, one for each value in
;; the environment.
;;
;; Here we filter the output to find any functions that match the
;; regular expression /int/
(set! out (filter (env) (lambda (x) (match "int" (get x :name)))))

;; Show the results
(print "Values in the environment matching the regexp /int/\n%s" out)

;;
(apply out (lambda (x) (begin
                        (print "Function in env. matching regexp /int/")
                        (print "\t%s" (get x :name)))))