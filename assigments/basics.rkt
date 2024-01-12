#lang plait

(print-only-errors #t)

; common definitions
(define (pow2 num)
  (* num num))
(define pi 3.141592653589)

(define (sum-coins pennies nickels dimes quarters)
  (let ([nickels (* nickels 5)]
        [dimes (* dimes 10)]
        [quarters (* quarters 25)])
    (+ pennies (+ nickels (+ dimes quarters)))))

(define (area-cylinder radius height)
  (+
   ; ends
   (* 2 (* pi (pow2 radius)))
   ; sides
   (* height (* 2 (* pi radius)))
   )
  )

; inside : 2 * pi * inner-rad * length
; outside : 2 * pi * outer-rad * length
; ends : pi * (outer-rad^2 - inner-rad^2)
(define (area-pipe inner-rad thickness length)
  (let ([outer-rad (+ inner-rad thickness)])
    (+
     ; ends
     (* 2 (* pi (- (pow2 outer-rad) (pow2 inner-rad))))
     (+
      ; inside
      (* length (* 2 (* pi outer-rad)))
      ; outside
      (* length (* 2 (* pi inner-rad)))))))

(define (tax gross-pay)
  (cond
    [(> gross-pay 480)
     (+ (* (- gross-pay 480) 0.28) (* 240 0.15 ))]
    [(> gross-pay 240)
     (* (- gross-pay 240) 0.15)]
    [else 0]))

(define (gross-pay hours wage)
  (let ([gross (* hours wage)])
    (- gross (tax gross))))

(define-type QuadraticClassification
  (degenerate)
  (zeroSolutions)
  (oneSolution)
  (twoSolutions))

(define (what-kind a b c)
  (let ([discriminant (- (pow2 b) (* 4 (* a c)))])
    (cond
      [(zero? a) (degenerate)]
      [(> discriminant 0) (twoSolutions)]
      [(zero? discriminant) (oneSolution)]
      [(< discriminant 0) (zeroSolutions)])))

(define-type Time
  (hms [hours : Number] [minutes : Number] [seconds : Number]))

(define (hms-to-seconds t)
  (+ (* 3600 (hms-hours t))
     (+ (* 60 (hms-minutes t))
        (hms-seconds t))))

(define (time-diff t1 t2)
  (- (hms-to-seconds t2)
     (hms-to-seconds t1)))

(define-type Position
  (position [x : Number] [y : Number]))

(define-type Shape
  (circle [center : Position]
          [radius : Number])
  (square [upper-left : Position]
          [side-length : Number])
  (rectangle [upper-left : Position]
             [width : Number]
             [height : Number]))

(define (area shape)
  (type-case Shape shape
    [(circle _ r)
     (* pi (pow2 r))]

    [(square _ l)
     (pow2 l)]

    [(rectangle _ w h)
     (* w h)]))

(define (translate-shape shape delta)
  (type-case Shape shape
    [(circle c r)
     (circle (position (+ (position-x c) delta) (position-y c)) r)]

    [(square c l)
     (square (position (+ (position-x c) delta) (position-y c)) l)]

    [(rectangle c w h)
     (rectangle (position (+ (position-x c) delta) (position-y c)) h w)]))


; this is so bad, sorry.
(define (in-shape shape p)
  (let ([p-x (position-x p)]
        [p-y (position-y p)])
    (type-case Shape shape
      [(circle c r) 
	   (let ([c-x (position-x c)]
			 [c-y (position-y c)])
		 ; (x - c_x)^2 + (y - c_y)^2 < r^2
		 (< (+ (pow2 (- p-x c-x))
			   (pow2 (- p-y c-y)))
			(pow2 r)))]

      [(square c l)
       (let ([left-x (position-x c)]
             [top-y (position-y c)]
             [right-x (+ (position-x c) l)]
             [bot-y (+ (position-y c) l)])

         (and (>= p-x left-x)
              (<= p-x right-x)
              (>= p-y top-y)
              (<= p-y bot-y)))]

      [(rectangle c w h)
       (let ([left-x (position-x c)]
             [top-y (position-y c)]
             [right-x (+ (position-x c) w)]
             [bot-y (+ (position-y c) h)])

         (and (>= p-x left-x)
              (<= p-x right-x)
              (>= p-y top-y)
              (<= p-y bot-y)))])))

(test (sum-coins 5 1 2 3) 105)
(test (sum-coins 100 100 100 100) 4100)

(test (< (- (area-cylinder 1 10) 69.115) 0.01) #t)

(test (area-pipe 10 1 10) 1451.4158)

(test (tax 600) 69.6)
(test (gross-pay 6 100) 530.4)

(test (what-kind 0 1 1) (degenerate))
(test (what-kind 1 10 1) (twoSolutions))
(test (what-kind 2 4 2) (oneSolution))
(test (what-kind 4 2 4) (zeroSolutions))

(test (time-diff (hms 23 50 0) (hms 23 52 0)) 120)
(test (time-diff (hms 23 52 0) (hms 23 50 0)) -120)

(test (area (circle (position 0 0 ) 1)) pi)
(test (area (circle (position 0 0 ) 10)) (* pi 100))
(test (area (square (position 0 0 ) 6)) 36)
(test (area (rectangle (position 0 0 ) 6 10)) 60)

(test (translate-shape (circle (position 0 0) 10) 5) (circle (position 5 0) 10))
(test (translate-shape (circle (position 0 0) 10) -5) (circle (position -5 0) 10))
(test (translate-shape (circle (position 0 10 ) 10) 5) (circle (position 5 10) 10))

(test (in-shape (square (position 0 0 ) 6) (position 2 4)) #t)
(test (in-shape (square (position 0 0 ) 6) (position 7 4)) #f)
(test (in-shape (square (position 0 0 ) 6) (position 4 7)) #f)
(test (in-shape (rectangle (position 0 0 ) 6 10) (position 4 7)) #t)
(test (in-shape (rectangle (position 0 0 ) 6 10) (position 7 4)) #f)
(test (in-shape (rectangle (position 0 0 ) 6 10) (position 7 4)) #f)
(test (in-shape (circle (position 0 0 ) 6) (position 4 4)) #t)
(test (in-shape (circle (position 0 0 ) 6) (position 6.1 0)) #f)

