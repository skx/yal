;;; time.lisp - Demonstrate the date/time functions

;;
;; This is a sample input file for our minimal lisp interpreter.
;;
;; We use it to demonstrate the date and time functions.
;;
;; (date) and (time) are implemented in our golang application,
;; and each returns a list of values.  The individual fields are
;; made available by helper-functions defined in our standard-library.
;;

(print "The year is %s" (year))
(print "The date is %s/%s/%s" (day) (month) (year))
(print "The time is %s (%s seconds past the epoch)" (hms) (now))
(print "Today is a %s" (weekday))

(print "Date as a list %s" (date))
(print "Time as a list %s" (time))
