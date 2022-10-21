;;; directory.lisp - Directory-related functions


;; Handy function to invoke a callback on files
(set! directory:walk (fn* (path:string fn:function)
                          "Invoke the specified callback on every file beneath the given path."

                          (apply (directory:entries path) fn)))
