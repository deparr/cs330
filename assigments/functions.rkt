#lang plait
(print-only-errors #t)

(define (check-temperate temps)
  (check-temps temps 95 5))

(define (check-temps temps high low)
  (foldl (lambda (v acc)
           (and acc (<= v high) (>= v low)))
         #t
         temps))

(define (convert digits) ....)

(define (average-price nums) ....)

(define (convertFC fs)
  (map (lambda (temp)
         (* (- temp 32) 5/9))
       fs))

(define (eliminate-exp ua lop)
  (filter (lambda (v) (<= v ua)) lop))

(define (compose-func after before)
  (lambda (val)
    (after (before val))))

(define (flatten loloa) ....)

(define (flatten-foldr) ....)

(define (bucket lon) ....)

(define-type Pedigree
  (person [name : String]
          [birth-year : Number]
          [eye-color : Symbol]
          [father : Pedigree]
          [mother : Pedigree])
  (unknown))

(define (tree-map f pedigree) ....)

(define (add-last-name pedigree last-name) ....)
