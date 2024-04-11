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

(define (take-n l n)
  (type-case (Listof 'a) l
    [empty
     empty]
	[(cons v rls)
	 (if (> n 0)
	 (cons v (take-n rls (- n 1)))
	 empty)]))

;; (build-infinite-list f) → (Listof 'a)
;;	f : (Number -> 'a)
(define (build-infinite-list f)
  (map f (seq 0)))

;; nats : (Listof Number)
(define nats
  (build-infinite-list (lambda (n) n)))

(define (all pred l)
  (type-case (Listof 'a) l
    [empty
     #t]
    [(cons v rls)
     (and (pred v) (all pred rls))]))

;; (prime? n) → Boolean
;;	n : Number
#;
(define (prime? n)
  (all (lambda (d) (> (remainder n d) 0))
       (take-while (lambda (d) (< (* d d) n))
                   (rest (rest nats)))))
(define (prime? n)
  (foldl (lambda (d a) (and a (> (remainder n d) 0))) #t
         (take-while (lambda (d) (< (* d d) n))
                     (rest (rest nats)))))

;; primes : (Listof Number)
(define primes
  (filter prime? (rest (rest nats))))

(list-ref primes 1000)

;; (prime?/fast n) → Boolean
;;	n : Number
(define (prime?/fast n)
  (foldl (lambda (d a) (and a (> (remainder n d) 0))) #t
         (take-n primes/fast n)))
#;
(define (prime?/fast n)
  (foldl (lambda (d a) (and a (> (remainder n d) 0))) #t
         (take-while (lambda (d) (< (* 2 d) n))
                     primes/fast)))

;; primes/fast : (Listof Number)
(define primes/fast
  (filter prime?/fast (rest (rest nats))))

(list-ref primes/fast 0)

