;;; sentences.yal -- Generate random sentences.

;; Adapted from Norvig, Paradigms of Artificial Intelligence
;; Programming, pp. 36-43 (MIT License).
;;
;; See worked example:
;;
;;  https://github.com/norvig/paip-lisp/blob/main/docs/chapter2.md
;;


;; Sentence
(set! sentence (fn* () (list (noun-phrase) (verb-phrase))))

;; Parts
(set! noun-phrase (fn* () (list (Article) (Noun))))
(set! verb-phrase (fn* () (list (Verb) (noun-phrase))))

;; Words
(set! Article (fn* () (random:item '(the a))))
(set! Noun    (fn* () (random:item '(man ball woman table chair sofa))))
(set! Verb    (fn* () (random:item '(hit took saw liked))))

;; Show some random sentences
(repeat 5 (lambda (n) (print "%s." (join (flatten (sentence)) " "))))
