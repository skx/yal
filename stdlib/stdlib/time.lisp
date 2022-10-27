;;; time.lisp - Time related functions


;; We have a built in function "time" to return the current time
;; as a list (HH MM SS).
;;
;; Here we create some helper functions for retrieving the various
;; parts of the time, as well as some aliases for ease of typing.

(set! time:hour (fn* ()
                     "Return the current hour, as an integer."
                     (nth (time) 0)))

(set! time:minute (fn* ()
                       "Return the current minute, as an integer."
                       (nth (time) 1)))

(set! time:second (fn* ()
                       "Return the current seconds, as an integer."
                       (nth (time) 2)))

;; define legacy aliases
(alias hour time:hour
       minute time:minute
       second time:second)

(set! zero-pad-single-number (fn* (num:number)
                                  "Return the given number, padded to two digits, as a string.
i.e. Add a '0' prefix to the specified number, for values less than ten.

This is designed to pad the hours, minutes, and seconds in (hms)."
                                  (if (< num 10)
                                      (sprintf "0%d" num)
                                    (str num))))

(set! time:hms (fn* ()
               "Return the current time as a string, formatted as 'HH:MM:SS'."
               (sprintf "%s:%s:%s"
                        (zero-pad-single-number (hour))
                        (zero-pad-single-number (minute))
                        (zero-pad-single-number (second)))))
(alias hms time:hms)
