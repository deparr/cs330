
def gen(x):
    yield (yield x)
g = gen(0)

print(next(g))
print(next(g))
