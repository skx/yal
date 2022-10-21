;;; string-case.lisp - Convert a string to upper/lower case

;;
;; This is either gross or cool.
;;
;; Define a hash which has literal characters and their upper-case, and
;; lower-cased versions
;;
(set! upper-table {
  a "A"
  b "B"
  c "C"
  d "D"
  e "E"
  f "F"
  g "G"
  h "H"
  i "I"
  j "J"
  k "K"
  l "L"
  m "M"
  n "N"
  o "O"
  p "P"
  q "Q"
  r "R"
  s "S"
  t "T"
  u "U"
  v "V"
  w "W"
  x "X"
  y "Y"
  z "Z"
  } )

(set! lower-table {
  A "a"
  B "b"
  C "c"
  D "d"
  E "e"
  F "f"
  G "g"
  H "h"
  I "i"
  J "j"
  K "k"
  L "l"
  M "m"
  N "n"
  O "o"
  P "p"
  Q "q"
  R "r"
  S "s"
  T "t"
  U "u"
  V "v"
  W "w"
  X "x"
  Y "y"
  Z "z"
  } )


;; Translate the elements of the string using the specified hash
(set! translate (fn* (x:string hsh:hash)
                     "Translate each character in the given string, via the means of the supplied lookup-table.

This is used by both 'upper' and 'lower'."
                     (let* (chrs (split x ""))
                       (join (map chrs (lambda (x)
                                         (if (get hsh x)
                                             (get hsh x)
                                           x)))))))

;; Convert the given string to upper-case, via the lookup table.
(set! upper (fn* (x:string)
                 "Convert each character from the supplied string to upper-case, and return that string."
                 (translate x upper-table)))

;; Convert the given string to upper-case, via the lookup table.
(set! lower (fn* (x:string)
                 "Convert each character from the supplied string to lower-case, and return that string."
                (translate x lower-table)))
