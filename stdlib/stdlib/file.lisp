;;; file.lisp - File-related primitives


;; Wrappers for accessing results of (file:stat)
(set! file:stat:size (fn* (path)
                          "Return the size of the given file, return -1 on error."
                          (let* (info (file:stat path))
                            (cond
                             (nil? info) -1
                             true (nth info 1)))))

(set! file:stat:uid (fn* (path)
                          "Return the UID of the given file owner, return '' on error."
                          (let* (info (file:stat path))
                            (cond
                             (nil? info) ""
                             true (nth info 2)))))


(set! file:stat:gid (fn* (path)
                          "Return the GID of the given file owner, return '' on error."
                          (let* (info (file:stat path))
                            (cond
                             (nil? info) ""
                             true (nth info 3)))))

(set! file:stat:mode (fn* (path)
                          "Return the mode of the given file, return '' on error."
                          (let* (info (file:stat path))
                            (cond
                             (nil? info) ""
                             true (nth info 4)))))

(set! file:which (fn* (binary)
                 "Return the complete path to the specified binary, found via the users' PATH setting.

If the binary does not exist in a directory located upon the PATH nil will be returned.

NOTE: This is a non-portable function!

      1.  It assumes that the environmental variable PATH exists.
      2.  It assumes $PATH can be split by ':'
      3.  It assumes '/' works as a directory separator.
"
                 (let* (path (split (getenv "PATH") ":")
                             res (filter path (lambda (dir) (exists? (join (list dir "/" binary))))))
                   (if res
                       (join (list (car res) "/" binary))))))


;; Define a legacy alias
(alias slurp file:read)

;; Read a file, and execute the contents
(def! load-file (fn* (filename)
                     "Load and execute the contents of the supplied filename."
                     (eval (join (list "(do " (slurp filename) "\nnil)")))))


;; Similar to load-file, but with automatic suffix addition and error-testing.
(def! require (fn* (name)
                   "Load and execute the given file, adding a .yal suffix.

To load and execute foo.yal, returning nil if that doesn't exist:

Example: (require 'foo)"
                   (let* (fname (sprintf "%s%s" name ".yal"))
                     (if (file:stat fname)
                         (load-file fname)
                       nil))))
