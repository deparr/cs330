#lang plait
(print-only-errors #t)

(define (find-temperate temps)
  (cond
    [(empty? temps) (none)]
    [(cons? temps) (let ([val (first temps)])
                     (if (and (<= val 95) (>= val 5))
                         (some val)
                         (find-temperate (rest temps))))]))

(define (check-temperate temps)
  (check-temps temps 5 95))

(define (check-temps temps low high)
  (letrec ([check-bounds
            (lambda (vals)
              (cond
                [(empty? vals) #t]
                [(cons? vals)
                 (let ([val (first vals)])
                   (if (or (> val high) (< val low))
                       #f
                       (check-bounds (rest vals))))]))])
    (check-bounds temps)))

(define (convert digits)
  (letrec ([walk (lambda (digits mag)
                   (cond
                     [(empty? digits) 0]
                     [(cons? digits)
                      (+ (* (first digits) mag)
                         (walk (rest digits) (* 10 mag)))]))])
    (walk digits 1)))

(define (average-price prices)
  (letrec ([average
            (lambda (nums sum count)
              (cond
                [(empty? nums) (/ sum count)]
                [(cons? nums) (average
                               (rest nums)
                               (+ sum (first nums))
                               (add1 count))]))])
    (if (cons? prices)
        (average prices 0 0)
        0)))

; dont think we're allowed to use map here
#;(define (convertFC fahrenheits)
    (map (lambda (temp)
           (* (- temp 32) 5/9))
         fahrenheits))

(define (convertFC fs)
  (cond
    [(empty? fs) empty]
    [(cons? fs)
     (cons (* (- (first fs) 32) 5/9) (convertFC (rest fs)))]))


; jank filter implementation
(define (eliminate-exp ua lop)
  (letrec ([filter (lambda (vals proc)
                     (cond
                       [(empty? vals) (list)]
                       [(cons? vals) (if (proc (first vals))
                                         (filter (rest vals) proc)
                                         (cons (first vals) (filter (rest vals) proc)))]))])
    (filter lop (lambda (val) (> val ua)))))

(define (suffixes l)
  (cond
    [(empty? l) (cons l empty)]
    [(cons? l) (cons l (suffixes (rest l)))]))

(define-type Pedigree
  (person [name : String]
          [birth-year : Number]
          [eye-color : Symbol]
          [father : Pedigree]
          [mother : Pedigree])
  (unknown))

(define (count-persons pedigree)
  (type-case Pedigree pedigree
    [(person _n _b _e f m)
     (+ 1 (+ (count-persons f) (count-persons m)))]
    [(unknown) 0]))

(define (average-age pedigree)
  (letrec ([sum-ages (lambda (pedigree)
                       (type-case Pedigree pedigree
                         [(person _n b _e f m)
                          (+ (+ (- 2023 b) (sum-ages f)) (sum-ages m))]
                         [(unknown) 0]))])
    (/ (sum-ages pedigree) (count-persons pedigree))))

(define (eye-colors pedigree)
  (type-case Pedigree pedigree
    [(person _n _b e f m)
     (append (cons e (eye-colors f)) (eye-colors m))]
    [(unknown) (list)]))

(test (find-temperate '(100 100 100 100 100 100 100 60 100 100 100)) (some 60))
(test (find-temperate '(100 100 100 100 100 100 100 100 100 100 95)) (some 95))
(test (find-temperate '(100 100 100 100 100 100 100 100 100 100 100)) (none))

(test (check-temperate '(10 20 30 40 50 60 70 80 90)) #t)
(test (check-temperate '(10 20 30 40 50 60 70 80 90 0)) #f)
(test (check-temperate '(4 4 4 4 4 96 96 96 96)) #f)

(test (check-temps '(0 10 20 30 40 50 60 70 80 90 100) 0 100) #t)
(test (check-temps '(-10) 0 100) #f)
(test (check-temps '(200) 0 100) #f)

(test (convert '(1 2 3)) 321)
(test (convert '(5 4 3 2 1)) 12345)
(test (convert '(5 0 0 0 1)) 10005)

(test (average-price '(5 5 5 5 5 5)) 5)
(test (average-price '(1 2 3 4 5 6 7)) 4)
(test (average-price '()) 0)

(test (convertFC '(32 212 -40)) '(0 100 -40))

(test (eliminate-exp 5 '(3 4 5 6 7 8)) '(3 4 5))

(test (suffixes '('a 'b 'c 'd))
      (list '('a 'b 'c 'd) '('b 'c 'd) '('c 'd) '('d) '()))

(test (count-persons
       (person "me" 0 'green
               (person "father" 0 'brown (unknown) (unknown))
               (person "mother" 0 'green (unknown) (unknown))))
      3)

(test (eye-colors
       (person "me" 0 'green
               (person "father" 0 'brown (unknown) (unknown))
               (person "mother" 0 'green (unknown) (unknown))))
      (list 'green 'brown 'green))

(test (average-age
       (person "me" 2023 'green
               (person "father" 2023 'brown (unknown) (unknown))
               (person "mother" 2023 'green (unknown) (unknown))))
      0)

(test (average-age
       (person "me" 2018 'green
               (person "father" 2013 'brown (unknown) (unknown))
               (person "mother" 2008 'green (unknown) (unknown))))
      10)

(test (average-age
       (person "me" 2013 'green
               (person "father" 1989 'brown (unknown) (unknown))
               (person "mother" 1990 'green (unknown) (unknown))))
      77/3)
