#lang racket
(require racket/generator)

(define x 0)
(define (a b)
  (b)
  (set! x 3))

(define c
  (generator ()
    (a (lambda () (yield 1)))
    (yield x)))

(c)
(c)
