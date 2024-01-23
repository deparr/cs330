#lang plait
(print-only-errors #t)

(define (check-temperate temps)
  (check-temps temps 95 5))

(define (check-temps temps high low)
  (foldl (lambda (v acc)
           (and acc (<= v high) (>= v low)))
         #t
         temps))

(define (convert digits)
  (fst (foldl (lambda (d acc)
                (let ([sum (fst acc)]
                      [mag (snd acc)])
                  (pair (+ sum (* mag d)) (* mag 10))))
              (pair 0 1)
              digits)))

(define (average-price nums)
  (let ([sum (foldr + 0 nums)]
        [count (length nums)])
    (if (zero? count)
        (error 'badaverage "bad div")
        (/ sum count))))

(define (convertFC fs)
  (map (lambda (temp)
         (* (- temp 32) 5/9))
       fs))

(define (eliminate-exp ua lop)
  (filter (lambda (v) (<= v ua)) lop))

(define (compose-func after before)
  (lambda (val)
    (after (before val))))

(define (flatten loloa)
  (type-case (Listof (Listof 'a)) loloa
    [empty (list)]
    [(cons loa r) (append loa (flatten r))]))

(define (flatten-foldr loloa)
  (foldr append (list) loloa))

(define (bucket lon)
  (foldr (lambda (n acc)
           (type-case (Listof 'a) acc
             [empty (cons (cons n empty) empty)]
             [(cons cur-bkt rst-bkts)
              (if (= n (first cur-bkt))
                  (cons (cons n cur-bkt) rst-bkts)
                  (cons (cons n empty) acc))]))
         (list) lon))

(define-type Pedigree
  (person [name : String]
          [birth-year : Number]
          [eye-color : Symbol]
          [father : Pedigree]
          [mother : Pedigree])
  (unknown))

(define (tree-map proc pedigree)
  (type-case Pedigree pedigree
    [(unknown) (unknown)]
    [(person n b e f m)
     (person (proc n) b e (tree-map proc f) (tree-map proc m))]))

(define (add-last-name pedigree last-name)
  (let ([last-name (string-append " " last-name)])
    (tree-map (lambda (n)
                (string-append n last-name))
              pedigree)))

(test (convert (list 1 2 3)) 321)

(test ((compose-func add1 (lambda (v) (* v 2))) 5) 11)
(test ((compose-func (lambda (v) (* v 2)) add1) 5) 12)

(test (flatten (list (list 1 2) (list 3 4 5) (list 6))) (list 1 2 3 4 5 6))
(test (flatten empty) empty)
(test (flatten-foldr (list (list 1 2) (list 3 4 5) (list 6))) (list 1 2 3 4 5 6))
(test (flatten-foldr empty) empty)

(test (bucket (list 1 1 2 2 2 3 1 1 1 2 3 3 ))
      (list (list 1 1) (list 2 2 2 ) (list 3) (list 1 1 1) (list 2) (list 3 3)))
(test (bucket empty) (list))

(test (add-last-name (person "child" 1997 'gray
                             (person "father" 1970 'brown (unknown) (unknown))
							 (person "mother" 1971 'blue (unknown) (unknown)))
                     "last-name")
      (person "child last-name" 1997 'gray
              (person "father last-name" 1970 'brown (unknown) (unknown))
              (person "mother last-name" 1971 'blue (unknown) (unknown))))

