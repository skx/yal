;;; dynamic.lisp - Execute code by name, via introspection.

;;
;; I'm not sure whether to be pleased with this or not.
;;
;; Given the (string) name of a function to be called, and some
;; arguments .. call it.
;;
;; (env) returns a lists of hashes, so we can find the function with
;; a given name via `filter`.  Assuming only one response then we're
;; able to find it by name, and execute it.
;;
(define call-by-name
  (lambda (name:string &args)
    (let ((out nil)  ; out is the result of the filter
          (nm  nil)  ; nm is the name of the result == name
          (fn  nil)) ; fn is the function of the result

      ;; find the entry in the list with the right name
      (set! out (filter (env) (lambda (x) (eq (get x :name) name))))

      ;; there should be only one matching entry
      (if (= (length out) 1)
          (begin
           (set! nm (get (car out) :name))   ;; nm == name
           (set! fn (get (car out) :value))  ;; fn is the function to call
           (if fn (fn args)))))))            ;; if we got it, invoke it


;; Print a string
(call-by-name "print" "Hello, world!")

;; Get an environmental variable
(print (call-by-name "getenv" "HOME"))

;; Call print with no arguments
(call-by-name "print")

;; The other way to do dynamic calls
(eval "(print (+ 34 43))")
