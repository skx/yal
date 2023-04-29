;;; trig.lisp - Our trigonometry functions.

;; These are implemented in golang, and called via reflection, so
;; the only thing we do here is setup suitable aliases and define
;; the help information for them.


(alias acos math.Acos)
(help "acos" "Acos returns the arccosine, in radians, of n.")

(alias asin math.Asin)
(help "asin" "Asin returns the arcsine, in radians, of n.")

(alias atan math.Atan)
(help "atan" "Atan returns the arctangent, in radians, of n.")


(alias cos math.Cos)
(help "cos" "Cos returns the cosine of the radian argument.")
(alias cosh math.Cosh)
(help "cosh" "Cosh returns the hyperbolic cosine of n.")


(alias sin math.Sin)
(help "sin" "Sin returns the sine of the radian argument.")
(alias sinh math.Sinh)
(help "sinh" "Sinh returns the hyperbolic sine of n.")

(alias tan math.Tan)
(help "tan" "Tan returns the tangent of the radian argument.")
(alias tanh math.Tanh)
(help "tanh" "Tanh returns the hyperbolic tangent of n.")
