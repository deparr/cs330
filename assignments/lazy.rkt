#lang plait #:lazy
;; David Parrott - dmparr22

;; (take-while pred xs) → (Listof 'a)
;;   pred : ('a -> Boolean)
;;   xs : (Listof 'a)
(define (take-while pred xs)
  (type-case (Listof 'a) xs
    [empty
     empty]
    [(cons x rest-xs)
     (if (pred x)
         (cons x (take-while pred rest-xs))
         empty)]))

(define (seq start)
  (cons start (seq (add1 start))))

;; (build-infinite-list f) → (Listof 'a)
;;	f : (Number -> 'a)
(define (build-infinite-list f)
  (map f (seq 0)))

;; nats : (Listof Number)
(define nats
  (build-infinite-list (λ (n) n)))

;; (prime? n) → Boolean
;;	n : Number
(define (prime? n)
  (foldl (λ (d a) (and a (> (remainder n d) 0))) #t
         (take-while (λ (d) (<= (* d d) n))
                     (rest (rest nats)))))

;; primes : (Listof Number)
(define primes
  (filter prime? (rest (rest nats))))

;; (prime?/fast n) → Boolean
;;	n : Number
(define (prime?/fast n)
  (foldl (λ (d a) (and a (> (remainder n d) 0))) #t
         (take-while (λ (d) (< (* 2 d) n)) primes/fast)))

;; primes/fast : (Listof Number)
(define primes/fast
  (cons 2 (filter prime?/fast (rest (rest (rest nats))))))

