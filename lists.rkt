#lang plait

(define (find-temperate temps)
  (cond
    [(empty? temps) (none)]
    [(cons? temps) (let ([val (first temps)])
                     (if
                      (and (<= val 95) (>= val 5))
                      (some val)
                      (find-temperate (rest temps))))]))

(define (check-temperate temps)
  (check-temps temps 5 95))

(define (check-temps temps low high)
  (letrec ([check-bounds
            (lambda (vals)
              (cond
                [(empty? vals) #t]
                [(cons? vals) (let ([val (first vals)])
                                (if
                                 (or (> val high) (< val low))
                                 #f
                                 (check-bounds (rest vals))))]))])
    (check-bounds temps)))

(define (convert digits)
  (cond
    [(empty? digits) 1]
    [(cons? digits) (* 10 (convert (rest digits)))]))


(define (average-price prices)
  (letrec ([average
            (lambda (nums sum count)
              (cond
                [(empty? nums) (/ sum count)]
                [(cons? nums) (average (rest nums) (+ sum (first nums)) (add1 count))]))])
    (if (cons? prices)
        (average prices 0 0)
        0)))

; dont think we're allowed to use map here
(define (convertFC fahrenheits)
  (map (lambda (temp)
         (* (- temp 32) 5/9))
       fahrenheits))

(define (eliminate-exp ua lop)
  (letrec ([filter (lambda (vals proc)
                     (cond
                       [(empty? vals) (list)]
                       [(cons? vals) (if (proc (first vals))
                                         (cons (first vals) (filter vals proc))
                                         (filter vals proc)
                                         )]))])
    (filter lop (lambda (val) (> val ua)))))

(test (find-temperate '(100 100 100 100 100 100 100 60 100 100 100)) (some 60))
(test (find-temperate '(100 100 100 100 100 100 100 100 100 100 95)) (some 95))
(test (find-temperate '(100 100 100 100 100 100 100 100 100 100 100)) (none))

(test (check-temperate '(10 20 30 40 50 60 70 80 90)) #t)
(test (check-temperate '(10 20 30 40 50 60 70 80 90 0)) #f)
(test (check-temperate '(4 4 4 4 4 96 96 96 96)) #f)

(test (check-temps '(0 10 20 30 40 50 60 70 80 90 100) 0 100) #t)
(test (check-temps '(-10) 0 100) #f)
(test (check-temps '(200) 0 100) #f)

(test (average-price '(5 5 5 5 5 5)) 5)
(test (average-price '(1 2 3 4 5 6 7)) 4)
(test (average-price '()) 0)

(test (convertFC '(32 212 -40)) '(0 100 -40))
