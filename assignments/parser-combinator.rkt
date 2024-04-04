#lang plait #:untyped

(module+ test
  (print-only-errors #t))

(define-type-alias (Parser 'a 'b)
  ((Listof Char) ('a (Listof Char) (-> 'b) -> 'b) (-> 'b) -> 'b))

(parse : ((Parser 'a 'a) String -> 'a))
(define (parse p s)
  (p (string->list s)
     (λ (x cs fk)
       (if (empty? cs) x (fk)))
     (λ () (error 'parse "couldn't parse"))))

(define-syntax-rule
  (delayp p)
  (λ (cs sk fk) (p cs sk fk)))

(unitp : ('a -> (Parser 'a 'b)))
(define (unitp x)
  (λ (cs sk fk) (sk x cs fk)))

(module+ test
  (test (parse (unitp 42) "") 42)) 

(failp : (Parser 'a 'b))
(define failp
  (λ (cs sk fk) (fk)))

(module+ test
  (test/exn (parse failp "") ""))

(charp : (Parser Char 'b))
(define charp
  (λ (cs sk fk)
    (type-case (Listof Char) cs
      [empty
       (fk)]
      [(cons c cs)
       (sk c cs fk)])))

(module+ test
  (test (parse charp "c") #\c))

(bind1 : (('a₀ -> (Parser 'c 'b)) (Parser 'a₀ 'b) -> (Parser 'c 'b)))
(define (bind1 f p₀)
  (λ (cs sk fk)
    (p₀ cs
        (λ (x₀ cs fk)
          ((f x₀) cs sk fk))
        fk)))

(module+ test
  (test (parse (bind1 (λ (x) (unitp (+ x 1))) (unitp 42)) "") 43))

(bind2 : (('a₀ 'a₁ -> (Parser 'c 'b)) (Parser 'a₀ 'b) (Parser 'a₁ 'b) -> (Parser 'c 'b)))
(define (bind2 f p₀ p₁)
  (bind1 (λ (x₀) (bind1 (λ (x₁) (f x₀ x₁)) p₁)) p₀))

(bind3 : (('a₀ 'a₁ 'a₂ -> (Parser 'c 'b)) (Parser 'a₀ 'b) (Parser 'a₁ 'b) (Parser 'a₂ 'b) -> (Parser 'c 'b)))
(define (bind3 f p₀ p₁ p₂)
  (bind2 (λ (x₀ x₁) (bind1 (λ (x₂) (f x₀ x₁ x₂)) p₂)) p₀ p₁))

(bind4 : (('a₀ 'a₁ 'a₂ 'a₄ -> (Parser 'c 'b)) (Parser 'a₀ 'b) (Parser 'a₁ 'b) (Parser 'a₂ 'b) (Parser 'a₃ 'b) -> (Parser 'c 'b)))
(define (bind4 f p₀ p₁ p₂ p₃)
  (bind3 (λ (x₀ x₁ x₂) (bind1 (λ (x₃) (f x₀ x₁ x₂ x₃)) p₃)) p₀ p₁ p₂))

(define (fmap1 f p₀)
  (bind1 (λ (x₀) (unitp (f x₀))) p₀))

(define (fmap2 f p₀ p₁)
  (bind2 (λ (x₀ x₁) (unitp (f x₀ x₁))) p₀ p₁))

(define (fmap3 f p₀ p₁ p₂)
  (bind3 (λ (x₀ x₁ x₂) (unitp (f x₀ x₁ x₂))) p₀ p₁ p₂))

(define (fmap4 f p₀ p₁ p₂ p₃)
  (bind4 (λ (x₀ x₁ x₂ x₃) (unitp (f x₀ x₁ x₂ x₃))) p₀ p₁ p₂ p₃))

(altp : ((Listof (Parser 'a 'b)) -> (Parser 'a 'b)))
(define (altp ps)
  (λ (cs sk fk)
    ((foldr (λ (p fk) (λ () (p cs sk fk))) fk ps))))

(module+ test
  (test (parse (altp (list failp
                           failp
                           (unitp 10)))
               "")
        10))

(starp : ((Parser 'a 'b) -> (Parser (Listof 'a) 'b)))
(define (starp p)
  (altp (list (fmap2 cons p (delayp (starp p)))
            (unitp empty))))

(module+ test
  (test (parse (starp charp) "abc") (list #\a #\b #\c))
  (test (parse (bind1
                (λ (cs) (unitp (list->string cs)))
                (starp charp))
               "abc")
        "abc"))

(?p : ((Parser 'a 'b) -> (Parser (Optionof 'a) 'b)))
(define (?p p)
  (altp (list (fmap1 some p)
              (unitp (none)))))

(module+ test
  (test (parse (?p charp) "") (none))
  (test (parse (?p charp) "c") (some #\c)))

(begp : ((Parser 'a 'c) (Parser 'b 'c) -> (Parser 'b 'c)))
(define (begp p₀ p₁)
  (bind1 (λ (_) p₁) p₀))

(define-syntax beginp
  (syntax-rules ()
    [(_ p)
     p]
    [(_ p ps ...)
     (begp p (beginp ps ...))]))

(module+ test
  (test (parse (begp charp charp) "ab") #\b)
  (test (parse (beginp (list charp charp charp)) "abcd") #\d)) 

(?cp : ((Char -> Boolean) -> (Parser Char 'b)))
(define (?cp ?)
  (bind1
   (λ (c)
     (if (? c)
       (unitp c)
       failp))
   charp))

(cp : (Char -> (Parser Char 'b)))
(define (cp c₀)
  (?cp (λ (c) (equal? c₀ c))))

(foldp : (('a 'b -> 'b) 'b (Parser 'b 'c) -> (Parser 'b 'c)))
(define (foldp f u p)
  (altp (list (bind1 (λ (x) (foldp f (f x u) p)) p)
              (unitp u))))

