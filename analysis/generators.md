# Generators
### David Parrott - dmparr22
---
## Part 1: Language Comparison
Python3:
```py
x = 0

def a(b):
    global x
    b()
    x = 3

def c():
    a(lambda: (yield 1))
    yield x

gen = c()
print(next(gen))
print(next(gen))
```

Racket:
```scheme
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
```

These are not equivalent because python generators keep track of their own state.

I'm pretty sure we said that Racket's are the ones that do this is class, but I think I'm confused because it makes much more sense to me for python's to do this.
Calling `b` in the python code does not yield `1`, instead it creates a separate generator object independent of the one bound to `g`.
If the call to `b` were to be replaced with `next(b())`, *it would yield `1`*, but not to the top level (`print(next(gen))`), only to `a`.

In contrast, `yield` in Racket is a (seemingly) normal function that is able to walk the callstack or otherwise know that it is being called in the context of a generator.
This allows the call to `b` to yield `1` to the generator `c` at the top level. 

---

## Part 2: Implementation State
### Program 1
```py
def gen(x):
    yield x
    yield x + 1
g = gen(0)
next(g) + pause() + next(g)
```

Top-level Environment:
```
PC:
0 + • + next(g)

Binds:
g => @generator#gen
```

Generator Environment:
```
PC:
yield x
yield x + 1 <-- 

Binds:
x => 0
```

---

### Program 2
```py
def gen(x):
    while True:
        yield x
        x += 1
g = gen(0)
next(g) + next(g) + pause()
```

Top-level Environment:
```
PC:
0 + 1 + •

Binds:
g => @generator#gen
```

Generator Environment:
```
PC:
while True:
    yield x <-- 
    x += 1

Binds:
x => 2
```

---

### Program 3
```py
def gen(x):
    yield (yield x)
g = gen(0)

next(g)
pause()
next(g)
pause()
```

#### First pause:  
Top-level Environment:
```
PC:
0
•
next(g)
pause()


Binds:
g => @generator#gen
```

Generator Environment:
```
PC:
yield None <-- 

Binds:
x => 0
```

#### Second pause:  
Top-level Environment:
```
PC:
0
0
None
•

Binds:
g => @generator#gen
```

Generator Environment:
```
PC:
<nil>

Binds:
```
