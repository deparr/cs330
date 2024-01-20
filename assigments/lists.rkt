#lang plait
(print-only-errors #t)

(define (find-temperate temps)
  (type-case (Listof 'a) temps
    [empty
      (none)]
    [(cons t restts)
     (if (and (<= t 95) (>= t 5))
       (some t)
       (find-temperate restts))]))

(define (check-temperate temps)
  (check-temps temps 5 95))

(define (check-temps temps low high)
  (type-case (Listof Number) temps
    [empty #t]
    [(cons t restts)
     (and
      (<= t high)
      (>= t low)
      (check-temps restts low high))]))

(define (number-from-digits digits mag)
  (type-case (Listof Number) digits
    [empty 0]
    [(cons d restds)
     (+ (* d mag)
        (number-from-digits restds (* 10 mag)))]))

(define (convert digits)
  (number-from-digits digits 1))

(define (average nums sum count)
  (type-case (Listof Number) nums
    [empty (/ sum count)]
    [(cons n restns)
     (average restns (+ sum n) (add1 count))]))

(define (average-price prices)
  (if (cons? prices)
      (average prices 0 0)
      (error 'badaverage "bad div")))

(define (convertFC fahrenheits)
  (map (lambda (temp)
         (* (- temp 32) 5/9))
       fahrenheits))

(define (eliminate-exp ua lop)
  (filter (lambda (val) (<= val ua)) lop))

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

(define (sum-ages pedigree)
  (type-case Pedigree pedigree
    [(person _n b _e f m)
     (+ (+ (- 2023 b) (sum-ages f)) (sum-ages m))]
    [(unknown) 0]))

(define (average-age pedigree)
  (let ([count (count-persons pedigree)])
    (if (zero? count)
        (error 'badaverage "bad div")
        (floor (/ (sum-ages pedigree) count)))))

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
(test (check-temps '() 0 100) #t)

(test (convert '(1 2 3)) 321)
(test (convert '(5 4 3 2 1)) 12345)
(test (convert '(5 0 0 0 1)) 10005)

(test (average-price '(5 5 5 5 5 5)) 5)
(test (average-price '(1 2 3 4 5 6 7)) 4)
(test/exn (average-price '()) "bad div")

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
      25)
