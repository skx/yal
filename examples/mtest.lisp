;;; mtest.lisp - Simple tests of our macro system.

;;
;; Handy reference
;;
;;  https://lisp-journey.gitlab.io/blog/common-lisp-macros-by-example-tutorial/
;;
;;  https://lispcookbook.github.io/cl-cookbook/macros.html
;;


;; Define a simple list for testing-purposes.
(set! lst (quote (b c)))

;;
;; Here is our first macro, given a variable-name show both the
;; name and the current value.
;;
(defmacro! debug (fn* (x) `(print "Variable '%s' has value %s" '~x ~x)))
(debug lst)

;;
;; Here's a similar example, which asserts a condition is true.
;;
;; The working is similar to the above, we get given a condition and
;; we both evaluate it, and show it literally (in the case where things
;; failed).
;;
(defmacro! assert (fn* (exp)
                      `(if ~exp
                         ()
                           (print "Assertion failed: %s" `~exp))))



;; Suppose you want a version of setq that sets two variables to the
;; same value. So if you write:
;;
;;   (set2! x y (+ z 3))
;;
;; When z=8 then both x and y are set to 11.
;;
;; When you (the Lisp system) see:
;;
;;    (set2! v1 v2 e)
;;
;; We want to treat it as:
;:
;;    (do
;;      (set! v1 e)
;;      (set! v2 e)
;;    )
;:
;; Something like this should work:
;;
;; NOTE:  This has a short-coming, that the "e" parameter is executed
;;        or evaluated twice.
;;
;;        We'll refine to fix this.
;;
(defmacro! set2! (fn* (v1 v2 e)
                     `(do
                       (set! ~v1 ~e)
                       (set! ~v2 ~e))))


;;
;; You can see this in the following code:
;;
;;   (set2! a c (do (print "EXECUTED TWICE!") (+ 32 23)))

;;
;; The second attempt would use a temporary variable to store the new
;; value, so that the evaluation of the argument only occurs once.
;;
;; This looks like it should work:
;;
;; NOTE: This does not work.
;;
;;       The "(set!..)" calls operate in a new scope.  So they can't modify
;;       the global environment.
;;
(defmacro! set2! (fn* (v1 v2 e)
                     (let* (tmp (gensym))
                       `(do (let* (~tmp ~e)
                           (set! ~v1 ~tmp)
                           (set! ~v2 ~tmp))))))


;;
;; The third/final attempt uses a temporary variable to store the new
;; value, so that the evaluation of the argument only occurs once.
;;
;; The difference here is we use the three-argument form of the (set!..)
;; form, to update the global/parent scope.
;;
(defmacro! set2! (fn* (v1 v2 e)
                     (let* (tmp (gensym))
                       `(do (let* (~tmp ~e)
                           (set! ~v1 ~tmp true)
                           (set! ~v2 ~tmp true))))))

;;
;; Lets test it out.
;;
;; Define three variables A, B, & C
;;
(set! a 1)
(set! b 2)
(set! c 3)

;;
;; Confirm they have expected values
;;
(assert (= a 1))
(assert (= b 2))
(assert (= c 3))

;;
;; Update A + B, leaving C alone.
;;
(set2! a b 33)

;;
;; Confirm the values are changed, as expected.
;;
(assert (= a 33))
(assert (= b 33))
(assert (= c 3))


;; Confirm it works with an expression too.
;;
;; NOTE This expression is only evaluated once, which is what we wanted.
;;
(set2! a c (do (print "ONLY EXECUTED ONCE!") (+ 32 23)))

;;
;; So the values will be changed, again.
;;
(assert (= a 55))
(assert (= b 33))
(assert (= c 55))


;;
;; That's a very simple macro.
;;
;; Lets add some more simple ones.
;;



;;
;; if2 is a simple macro which allows you to run two actions if an
;; (if ..) test succeeds.
;;
;; This means you can write:
;;
;;   (if2 true (print "1") (print "2"))
;;
;; Instead of having to add (do:
;;
;;   (if true (do (print "1") (print "2")))
;;
;; The downside here is that you don't get a negative branch, but running
;; two things is very common - see for example the "(while)" and "(repeat)"
;; macros in our standard library.
;;
;; See also "(when) in the standard-library, which allows a list of operations
;; when a condition is true rather than two, and only two.
;;
(defmacro! if2 (fn* (pred one two)
  `(if ~pred (do ~one ~two))))


;;
;; Increment the given variable by one.
;;
(defmacro! incr (fn* (x) `(set! ~x (+ ~x 1))))

;;
;; Show macro expansion
;;
(print "The (incr a) macro expands to %s" (macroexpand (incr a)))

;;
;; Use the if2 macro to run two increment options
;;
(set! a 32)
(if2 true (incr a) (incr a))
(assert (= a 34))


;;
;; Finally we'll ensure that our type-checking understands what a macro is,
;; and that it is different from a (user) function or a builtin function.
;;


;; Type of a macro is "macro"
(defmacro! truthy (fn* () true))
(print "The type of a macro is (type truthy):%s" (type truthy))

;; The macro? predicate will recognize one too.
(if (macro? truthy)
    (print "(macro? truthy) -> true"))
