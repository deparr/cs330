#lang plait

; I'm not sure how to properly format lisp-like code
; to be nice and readable. And racket-langserver doesn't seem
; to do helpful formatting, sorry.

; sum value of change
(define (sum-coins pennies nickels dimes quarters)
  (+ pennies (+ (* nickels 5) (+ (* dimes 10) (* quarters 25)))))

(define (square num)
  (* num num))

; surface area of cylinder
(define pi 3.141592653589)
(define (area-cylinder radius height)
  (+
   ; ends
   (* 2 (* pi (square radius)))
   ; sides
   (* height (* 2 (* pi radius)))
   )
  )

; inside : 2 * pi * inner-rad * length
; outside : 2 * pi * outer-rad * length
; ends : pi * (outer-rad^2 - inner-rad^2)
(define (area-pipe inner-rad thickness length)
  (+
   ; ends
   (* 2 (* pi (- (square (+ inner-rad thickness)) (square inner-rad))))
   (+
    ; inside
    (* length (* 2 (* pi (+ inner-rad thickness))))

    ; outside
    (* length (* 2 (* pi inner-rad)))
    )))

; do recursively
(define (tax gross-pay)
  (* gross-pay
     (cond))
  )




(test (sum-coins 5 1 2 3) 105)
(test (sum-coins 100 100 100 100) 4100)

(test (< (- (area-cylinder 1 10) 69.115) 0.01) #t)

(area-pipe 10 1 10)

