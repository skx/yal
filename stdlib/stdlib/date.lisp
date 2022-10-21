;;; date.lisp - Date-related functions.

;; We have a built in function "date" to return the current date
;; as a list (XXXX DD MM YYYY).
;;
;; Here we create some helper functions for retrieving the various
;; parts of the date, as well as some aliases for ease of typing.
(set! date:day (fn* ()
               "Return the day of the current month, as an integer."
               (nth (date) 1)))

(set! date:month (fn* ()
                 "Return the number of the current month, as an integer."
                 (nth (date) 2)))


(set! date:weekday (fn* ()
                   "Return a string containing the current day of the week."
                   (nth (date) 0)))

(set! date:year (fn* ()
                "Return the current year, as an integer."
                (nth (date) 3)))


;;
;; define legacy aliases
;;
(alias day     date:day)
(alias month   date:month)
(alias weekday date:weekday)
(alias year    date:year)
