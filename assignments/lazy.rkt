#lang plait #:lazy
;; David Parrott - dmparr22

(define (f10 l)
  (begin
  (display (list-ref l 0))
  (display (list-ref l 1))
  (display (list-ref l 2))
  (display (list-ref l 3))
  (display (list-ref l 4))
  (display (list-ref l 5))
  (display (list-ref l 6))
  (display (list-ref l 7))
  (display (list-ref l 8))
  (display (list-ref l 9))))


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


;; (build-infinite-list f) → (Listof 'a)
;;	f : (Number -> 'a)
(define (build-infinite-list f)
  (cons (f 0) (build-infinite-list f)))

;; nats : (Listof Number)
(define nats
  (build-infinite-list add1))

(f10 nats)

;; (prime? n) → Boolean
;;	n : Number
(define (prime? n)
  ....)

;; primes : (Listof Number)
(define primes 
  ....)

;; (prime?/fast n) → Boolean
;;	n : Number
(define (prime?/fast n)
  ....)

;; primes/fast : (Listof Number)
(define primes/fast 
  ....)


;;
