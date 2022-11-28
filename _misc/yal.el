;;; yal.el -- Sample configuration for using Emacs LSP-mode with YAL


;; Create a keyboard-map for use within YAL files
(defvar yal-mode-map
  (let ((map (make-sparse-keymap)))
    (define-key map (kbd "C-c TAB") 'completion-at-point)
    map))

;; Define a hook which will run when entering YAL mode.
(add-hook 'yal-mode-hook 'lsp-deferred)

;; Now create a trivial "yal-mode"
(define-derived-mode yal-mode
  lisp-mode "YAL"
  "Major mode for working with yet another lisp, YAL.")

;; yal-mode will be invoked for *.yal files
(add-to-list 'auto-mode-alist '("\\.yal\\'" . yal-mode))

;; Load the library
(require 'lsp-mode)

;; Register an LSP helper
(lsp-register-client
 (make-lsp-client :new-connection (lsp-stdio-connection '("yal" "-lsp"))
                  :major-modes '(yal-mode)
                  :priority -1
                  :server-id 'yal-ls))

;; Not sure what this does, but it seems to be necessary
(add-to-list 'lsp-language-id-configuration '(yal-mode . "yal"))
