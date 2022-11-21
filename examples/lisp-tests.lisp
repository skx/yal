;;; tests.lisp - Some simple tests of our lisp-written primitives

;;; About
;;
;; This file contains a bunch of simple test-cases which demonstrate
;; that our lisp-implemented functions work as expected.
;;
;; The file will attempt to output the results in a TAP format, such
;; that it can be processed via automated tools.
;;
;; For example the "tapview" shell-script can consume our output and
;; will present something like this;
;;
;;   $ yal tests.lisp | tapview
;;   ....
;;   4 tests, 0 failures.
;;
;; When a test fails it will be shown:
;;
;;   $ yal tests.lisp | tapview
;;   not ok add:mult failed %!s(int=40) != %!s(int=10)
;;   4 tests, 1 failures.
;;
;; tapview can be found here:
;;
;;   https://gitlab.com/esr/tapview
;;
;;
;;; Note
;;
;; Of course the results can also be expected manually, the tapview is
;; just one of the many available TAP-protocol helpers.
;;
;;   $ yal tests.lisp | grep "not ok"
;;
;;
;;
;;; Details
;;
;; In terms of our implementation we use a macro to register
;; test functions.  Test functions are expected to return a list
;; of two elements - a test passes if those elements are identical,
;; and fails otherwise.
;;
;; The macro which defines a test-case will store the details in the
;; global *tests* hash:
;;
;;   key -> name of the test
;;   val -> The lambda body
;;
;; When we come to execute the tests we'll just iterate over the key/val
;; pairs appropriately.
;;



;;
;; A hash of all known test-cases.
;;
;; This is updated via the `deftest` macro, and iterated over by the
;; `run-tests` function.
;;
(set! *tests* {} )

;;
;; Define a new test.
;;
(defmacro! deftest (fn* (name body)
                        "Create a new test, storing details in the global *tests* hash.

If the name of the test is not unique then that will cause an error to be printed."
                        `(if (get *tests* `~name)
                             (print "not ok - name is not unique %s" `~name)
                           (set *tests* `~name (lambda () (do ~body))))
                        ))



;;
;; Test cases now follow, defined with the macro above.
;;

;;
;; Each test-case should return a list of two values:
;;
;; 1. If the two values are equal we have a pass.
;; 2. If the two values are not equal the test fails.
;;
;; If the test case returns anything other than a two-element
;; list it is also a failure, as is a non-unique test-name.
;;

;; +
(deftest add:simple (list (+ 3 4) 7))
(deftest add:mult   (list (+ 1 2 3 4) 10))

;; /
(deftest div:1 (list (/ 2  ) 0.5))  ; "/ x" == "1/x"
(deftest div:2 (list (/ 9 3) 3))
(deftest div:3 (list (/ 8 2) 4))

;; *
(deftest mul:1 (list (* 2      ) 2))  ; "* x" == "1 * x"
(deftest mul:2 (list (* 2 2    ) 4))
(deftest mul:3 (list (* 2 2 2  ) 8))
(deftest mul:4 (list (* 2 2 2 3) 24))

;; -
(deftest minus:1 (list (- 1 2   ) -1))
(deftest minus:2 (list (- 10 2  ) 8))
(deftest minus:3 (list (- 10 2 3) 5))

;; sqrt
(deftest sqrt:1 (list (sqrt 100) 10))
(deftest sqrt:2 (list (sqrt   9)  3))

;; power
(deftest pow:1 (list (# 10 2) 100))
(deftest pow:2 (list (# 2  3) 8))

;; neg
(deftest neg:1 (list (neg 100) -100))
(deftest neg:2 (list (neg -33)  33))

;; abs
(deftest abs:1 (list (abs 100) 100))
(deftest abs:2 (list (abs -33)  33))
(deftest abs:3 (list (abs   0)   0))

;; sign
(deftest sign:1 (list (sign 100)  1))
(deftest sign:2 (list (sign -33) -1))
(deftest sign:3 (list (sign   0)  1))

;; neg?
(deftest neg?:1 (list (neg? 100)   false))
(deftest neg?:2 (list (neg? -33)   true))
(deftest neg:3 (list (neg?   0.1) false))
(deftest neg:4 (list (neg?  -0.1) true))

;; pos?
(deftest pos:1 (list (pos? 100)   true))
(deftest pos:2 (list (pos? -33)   false))
(deftest pos:3 (list (pos?   0.1) true))
(deftest pos:4 (list (pos?  -0.1) false))

;; inc
(deftest inc:1 (list (inc  1)  2))
(deftest inc:2 (list (inc -1)  0))
(deftest inc:3 (list (inc 1.3) 2.3))

;; dec
(deftest dec:1 (list (dec  1)  0))
(deftest dec:2 (list (dec -1) -2))
(deftest dec:3 (list (dec 1.5) 0.5))

;; and
(deftest and:1 (list (and (list      false)) false))
(deftest and:2 (list (and (list       true)) true))
(deftest and:3 (list (and (list true  true)) true))
(deftest and:4 (list (and (list true false)) false))

;; not
(deftest not:1 (list (not    true) false))
(deftest not:2 (list (not   false) true))
(deftest not:3 (list (not "steve") false))
(deftest not:4 (list (not       3) false))
(deftest not:5 (list (not      ()) false))
(deftest not:6 (list (not     nil) true))   ; not nil -> true is expected

;; or
(deftest or:1 (list (or (list       false)) false))
(deftest or:2 (list (or (list        true)) true))
(deftest or:3 (list (or (list true   true)) true))
(deftest or:4 (list (or (list true  false)) true))
(deftest or:5 (list (or (list false false)) false))


;; numeric parsing
(deftest parse:int:1 (list 0b1111  15))
(deftest parse:int:2 (list 0xff   255))
(deftest parse:int:3 (list 332.2  332.2))

;; Upper-case a string
(deftest string:upper:ascii (list (upper "steve")   "STEVE"))
(deftest string:upper:utf   (list (upper "π!狐犬")   "π!狐犬"))
(deftest string:upper:mixed (list (upper "π-steve") "π-STEVE"))

;; Lower-case a string
(deftest string:lower:ascii (list (lower "STEVE")   "steve"))
(deftest string:lower:utf   (list (lower "π!狐犬")   "π!狐犬"))
(deftest string:lower:mixed (list (lower "π-STEVE") "π-steve"))

;; Left-pad
(deftest string:pad:left:ascii (list (pad:left "me" "x" 4)   "xxme"))
(deftest string:pad:left:utf   (list (pad:left "狐犬π" "x" 4) "x狐犬π"))
(deftest string:pad:left:mixed (list (pad:left "fπ" "x" 4)   "xxfπ"))

;; Right-pad
(deftest string:pad:right:ascii (list (pad:right "me" "x" 8)   "mexxxxxx"))
(deftest string:pad:right:utf   (list (pad:right "狐犬π" "x" 8) "狐犬πxxxxx"))
(deftest string:pad:right:mixed (list (pad:right "fπ" "x" 8)   "fπxxxxxx"))

;; Time should have two-digit length HH, MM, SS fields.
(deftest time:hms:len  (list (strlen (hms)) 8))

;; Year should be four digits, always.
(deftest year:len (list (strlen (str (date:year))) 4))

;; < test
(deftest cmp:lt:1 (list (< 1 10) true))
(deftest cmp:lt:2 (list (< -1 0) true))
(deftest cmp:lt:3 (list (< 10 0) false))

;; > test
(deftest cmp:gt:1 (list (> 1   10) false))
(deftest cmp:gt:2 (list (> 1    0) true))
(deftest cmp:gt:3 (list (> 10 -10) true))

;; <= test
(deftest cmp:lte:1 (list (<= 1 10)  true))
(deftest cmp:lte:2 (list (<= -1 0)  true))
(deftest cmp:lte:3 (list (<= 10 0)  false))
(deftest cmp:lte:4 (list (<= 10 10) true))

;; >= test
(deftest cmp:gte:1 (list (>= 1   10) false))
(deftest cmp:gte:2 (list (>= 1    0) true))
(deftest cmp:gte:3 (list (>= 10 -10) true))
(deftest cmp:gte:4 (list (>= 10  10) true))

;; eq test
(deftest cmp:eq:1 (list (eq 1            10) false))
(deftest cmp:eq:2 (list (eq 1             1) true))
(deftest cmp:eq:3 (list (eq 10          -10) false))
(deftest cmp:eq:4 (list (eq "steve" "steve") true))
(deftest cmp:eq:5 (list (eq "steve"  "kemp") false))
(deftest cmp:eq:6 (list (eq 32      "steve") false))
(deftest cmp:eq:7 (list (eq ()         nil ) false))
(deftest cmp:eq:8 (list (eq ()          () ) true))
(deftest cmp:eq:9 (list (eq nil        nil ) true))

;; = test
(deftest cmp:=:1 (list (eq 1        1) true))
(deftest cmp:=:2 (list (eq 1  (- 3 2)) true))
(deftest cmp:=:3 (list (eq 1       -1) false))
(deftest cmp:=:4 (list (eq .5 (/ 1 2)) true))

;;TODO char<
;;TODO char>
;;TODO char>=
;;TODO char<=

;; zero? test
(deftest tst:zero:1 (list (zero?  0) true))
(deftest tst:zero:2 (list (zero? 10) false))

;; one? test
(deftest tst:one:1 (list (one?  1) true))
(deftest tst:one:2 (list (one? 10) false))

;; even? test
(deftest tst:even:1 (list (even?  1) false))
(deftest tst:even:2 (list (even?  2) true))
(deftest tst:even:3 (list (even?  3) false))

;; odd? test
(deftest tst:odd:1 (list (odd? 1) true))
(deftest tst:odd:2 (list (odd? 2) false))
(deftest tst:odd:3 (list (odd? 3) true))

;; true? test
(deftest tst:true:1 (list (true? true)  true))
(deftest tst:true:2 (list (true? nil)   false))
(deftest tst:true:3 (list (true? false) false))
(deftest tst:true:4 (list (true? 32111) false))
(deftest tst:true:5 (list (true? ())    false))

;; false? test
(deftest tst:false:1 (list (false? false) true))
(deftest tst:false:2 (list (false? nil)   false))
(deftest tst:false:3 (list (false? true)  false))
(deftest tst:false:4 (list (false? 32111) false))
(deftest tst:false:5 (list (false? ())    false))

;; nil? test
(deftest tst:nil:1 (list (nil?   false) false))
(deftest tst:nil:2 (list (nil?      ()) true))
(deftest tst:nil:3 (list (nil?     nil) true))
(deftest tst:nil:4 (list (nil? "steve") false))
(deftest tst:nil:5 (list (nil? 3223232) false))

;; member test
(deftest member:1 (list (member "foo" (list "foo" "bar" "baz"))  true))
(deftest member:2 (list (member "luv" (list "foo" "bar" "baz"))  false))

;; union test
(deftest union:1 (list (union (list "foo") (list "foo" "bar" "baz")) (list "foo" "bar" "baz")))
(deftest union:2 (list (union (list "foo") (list "bar" "baz"))       (list "foo" "bar" "baz")))

;; intersection
(deftest intersection:1 (list (intersection (list "foo") (list "foo" "bar" "baz")) (list "foo")))
(deftest intersection:2 (list (intersection (list 1 2 3) (list 2 3 4 )) (list 2 3)))

;; TODO / FIXME / BUG - should intersection return nil if there are no common elements?
(deftest intersection:3 (list (intersection (list 1) (list 2 3 4 )) nil))


;; reverse
(deftest reverse:1 (list (reverse  (list "m" "e")) (list "e" "m")))
(deftest reverse:2 (list (reverse  (list "狐" "犬" "π")) (list "π" "犬" "狐")))

;; seq
(deftest seq:0 (list (seq 0) (list     0)))
(deftest seq:1 (list (seq 1) (list   0 1)))
(deftest seq:2 (list (seq 2) (list 0 1 2)))

;; nat
(deftest nat:0 (list (nat 0) (list   )))
(deftest nat:1 (list (nat 1) (list   1)))
(deftest nat:2 (list (nat 2) (list 1 2)))

;; take
(deftest take:1 (list (take 0 (list 0 1 2 3)) nil))
(deftest take:2 (list (take 1 (list 0 1 2 3)) (list 0)))
(deftest take:3 (list (take 2 (list 0 1 2 3)) (list 0 1)))

;; drop
(deftest drop:1 (list (drop 0 (list 0 1 2 3)) (list 0 1 2 3)))
(deftest drop:2 (list (drop 1 (list 0 1 2 3)) (list   1 2 3)))
(deftest drop:3 (list (drop 2 (list 0 1 2 3)) (list     2 3)))

;; butlast
(deftest butlast:1 (list (butlast (list 0 1 2 3)) (list 0 1 2)))
(deftest butlast:2 (list (butlast       (list 0))          nil))
(deftest butlast:3 (list (butlast            nil)          nil))

;; append
(deftest append:1 (list (append () "2") "2"))
(deftest append:2 (list (append (list 2) "2") (list 2 "2")))
(deftest append:3 (list (append (list 2 3) 5) (list 2 3 5)))

;; strlen
(deftest strlen:1 (list (strlen      "") 0))
(deftest strlen:2 (list (strlen "steve") 5))
(deftest strlen:3 (list (strlen  "狐犬π") 3))

;; repeated
(deftest repeated:0 (list (repeated 0 "x") nil))
(deftest repeated:1 (list (repeated 1 "x") (list "x")))
(deftest repeated:2 (list (repeated 2 "x") (list "x" "x")))
(deftest repeated:3 (list (repeated 3 "x") (list "x" "x" "x")))

;; hex
(deftest hex:1 (list (dec2hex 255) "ff"))
(deftest hex:2 (list (dec2hex  10) "a"))

;; binary - note that the shortest form will be returned
(deftest binary:1 (list (dec2bin 3) "11"))
(deftest binary:2 (list (dec2bin 4) "100"))

;; structures
(deftest struct:1 (list (do (struct person name) (type (person "me")))
                        "struct-person"))
(deftest struct:2 (list (do (struct person name) (person? (person "me")))
                        true))
(deftest struct:3 (list (do (struct person name) (person.name (person "me")))
                        "me"))


;;
;; Define a function to run all the tests, by iterating over the hash.
;;
(set! run-tests (fn* (hsh)
                     "Run all the registered tests, by iterating over the global supplied hash.

The hash will contain a key naming the test.   The value of the hash will be a function to
invoke to run the test."
                     (do
                         (print "TAP version 14")
                         (apply-hash hsh (lambda (test fun)
                                           (let* (out (fun))
                                             (if (! (list? out))
                                                 (print "not ok %s should have returned a list, instead got %v" test out)
                                               (if (! (= (count out) 2 ))
                                                   (print "not ok %s should have been a list of 2 elements, instead got %s" test out)
                                                 (let* (a (car out)
                                                          b (car (cdr out)))
                                                   (if (! (eq a b))
                                                       (print "not ok %s failed %s != %s" test a b)
                                                     (print "ok %s" test))))))))
                         (print "1..%d" (count (keys hsh))))))


;;
;; Now run the tests.
;;
(run-tests *tests*)
