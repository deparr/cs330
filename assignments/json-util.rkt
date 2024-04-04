#lang plait #:untyped
(require "parser-combinator.rkt")

(module+ test
  (print-only-errors #t))

(define digitp
  (bind1
   (λ (c)
     (letrec ([loop (λ (i ds)
                      (type-case (Listof Char) ds
                        [empty
                         failp]
                        [(cons d ds)
                         (if (equal? c d)
                           (unitp i)
                           (loop (add1 i) ds))]))])
       (loop 0 (string->list "0123456789"))))
   charp))

(numberp : (Parser Number 'b))
(define numberp
  (bind1 (λ (d) (foldp (λ (d u) (+ (* u 10) d)) d digitp))
         digitp))

(module+ test
  (test (parse numberp "123") 123))

(define stringp
  (bind1
   (λ (cs) (unitp (list->string cs)))
   (beginp (cp #\")
         (letrec ([loop (λ ()
                          (bind1
                           (λ (c)
                             (cond
                               [(equal? c #\")
                                (unitp empty)]
                               [(equal? c #\\)
                                (bind1
                                 (λ (c)
                                   (cond
                                     [(equal? c #\")
                                      (bind1 (λ (cs) (unitp (cons c cs))) (loop))]
                                     [(equal? c #\n)
                                      (bind1 (λ (cs) (unitp (cons #\newline cs))) (loop))]
                                     [(equal? c #\n)
                                      (bind1 (λ (cs) (unitp (cons #\return cs))) (loop))]
                                     [(equal? c #\\)
                                      (bind1 (λ (cs) (unitp (cons #\\ cs))) (loop))]
                                     [else
                                      failp]))
                                 charp)]
                               [else
                                (bind1 (λ (cs) (unitp (cons c cs))) (loop))]))
                           charp))])
           (loop)))))

(module+ test
  (test (parse stringp "\"abc\"") "abc")
  (test (parse stringp "\"a\\nc\"") "a\nc"))

(literalp : (String -> (Parser Void 'b)))
(define (literalp s)
  (foldr (λ (c p) (beginp (cp c) p)) (unitp (void)) (string->list s)))

(module+ test
  (test (parse (literalp "hello") "hello") (void)))

(define (char-whitespace? c)
  (member c (list #\space #\tab #\return #\newline)))

(sp : ((Parser 'a 'b) -> (Parser 'a 'b)))
(define (sp p)
  (beginp (starp (?cp char-whitespace?)) p))

(module+ test
  (test (parse (sp charp) "  c") #\c))

(define (delimitedp leftp rightp commap p)
  (beginp leftp
          (altp (list (beginp (sp rightp)
                              (unitp empty))
                      (bind1
                       (λ (x)
                         (bind1
                          (λ (xs) (unitp (cons x xs)))
                          (letrec ([loop (λ ()
                                           (altp (list (beginp (sp rightp)
                                                               (unitp empty))
                                                       (bind1
                                                        (λ (x) (bind1 (λ (xs) (unitp (cons x xs))) (loop)))
                                                        (beginp (sp commap)
                                                                (sp p))))))])
                            (loop))))
                       (sp p))))))

(define (listp p)
  (delimitedp (literalp "(")
              (literalp ")")
              (literalp ",")
              p))

(module+ test
  (test (parse (listp charp) "(a,b,c)") (list #\a #\b #\c))
  (test (parse (listp charp) "( a , b , c )") (list #\a #\b #\c)))

