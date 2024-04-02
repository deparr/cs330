# Generators
### David Parrott - dmparr22
---

## Program 1
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
0 + pause() + next(g)

Binds:
g => @generator#gen
```

Generator Environment
```
PC:
yield x + 1 

Binds:
x => 0
```

---

## Program 2
```py
while True:
    yield x
    x += 1
g = gen(0)
next(g) + next(g) + pause()
```

---

## Program 3
```py
def gen(x):
    yield (yield x)
g = gen(0)

next(g)
pause()
next(g)
pause()
```

First pause:  

Second pause:  

