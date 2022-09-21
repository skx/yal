;;; mtest.lisp - Simple tests of our macro system.

;;
;; Handy reference
;;
;;  https://lisp-journey.gitlab.io/blog/common-lisp-macros-by-example-tutorial/
;;
;;  https://lispcookbook.github.io/cl-cookbook/macros.html


;; To implement a macro-system there are things that are required,
;; the groundwork, such as a decent set of quote/unquote primitives.
;;
;; Simple tests of those here, from the MAL text
;;
(define lst (quote (b c)))

;; `(a lst d) -> (a lst d)
(if (! (eq (str `(a lst d)) "(a lst d)"))
    (print "Looks like our quote is broken"))

;; (quasiquote (a (unquote lst) d)) -> (a (b c) d)
(if (! (eq (str (quasiquote (a (unquote lst) d))) "(a (b c) d)"))
    (print "Looks like our quasiquote/unquote is broken"))

;; (quasiquote (a (splice-unquote lst) d)) -> (a (b c) d)
(if (! (eq (str (quasiquote (a (splice-unquote lst) d))) "(a b c d)"))
    (print "Looks like our quasiquote/splice-unquote is broken"))



;;
;; Here is our first macro, given a variable-name show both the
;; name and the current value.
;;
(define debug (macro (x) `(print "Variable '%s' has value: %s" '~x ~x)))

;;
;; Here's a similar example, which asserts a condition is true.
;;
;; The working is similar to the above, we get given a condition and
;; we both evaluate it, and show it literally (in the case where things
;; failed).
;;
(define assert (macro (exp)
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
;;    (begin
;;      (set! v1 e)
;;      (set! v2 e)
;;    )
;:
;; Something like this should work:
;;
;; NOTE:  This has a short-coming, that the "e" parameter is executed
;;        or evaluated twice.
;;
;;        We'll refine to fix this
;;
(define set2! (macro (v1 v2 e)
                     `(begin
                       (set! ~v1 ~e)
                       (set! ~v2 ~e))))


;;
;; You can see this in the following code:
;;
;;   (set2! a c (begin (print "EXECUTED TWICE!") (+ 32 23)))

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
(define set2! (macro (v1 v2 e)
                     (let ((tmp (gensym)))
                       `(begin (let ((~tmp ~e))
                           (set! ~v1 ~tmp)
                           (set! ~v2 ~tmp))))))


;;
;; The third/final attempt uses a temporary variable to store the new
;; value, so that the evaluation of the argument only occurs once.
;;
;; The difference here is we use the three-argument form of the (set!..)
;; form, to update the global/parent scope.
;;
(define set2! (macro (v1 v2 e)
                     (let ((tmp (gensym)))
                       `(begin (let ((~tmp ~e))
                           (set! ~v1 ~tmp t)
                           (set! ~v2 ~tmp t))))))

;;
;; Lets test it out.
;;
;; Define three variables A, B, & C
;;
(define a 1)
(define b 2)
(define c 3)

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
;; NOTE This expression is only evaluated once.
;;
(set2! a c (begin (print "ONLY ONE TIME EXECUTED!!") (+ 32 23)))

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
;; Define a simple "unless" macro, with a mandatory case and
;; an optional one.
;;
(define unless (macro (pred a &b) `(if (! ~pred) ~a ~b)))

;;
;; Use that to operate a series of expressions.
;;
(unless false (list
               (print "unless-test one")
               (print "unless-test two")
               (print "unless-test three")))

(unless true (print "FAIL") (print "OK - (unless ..) is good"))
(unless false (print "OK - (unless ..) is good.") (print "FAIL"))

(define truthy (macro () true))
(print (type truthy))
(if (macro? truthy)
    (print "macro? works"))
