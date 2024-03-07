# Loops
### David Parrott - dmparr22
---

**Question:** What happens in the python/java/racket programs?
- Python - The `i` binding is created once and then its value is mutated every iteration. Each closure captures the same reference, so when they are all printed out at the end of the program they all have the same value (9).
- Java -  This program would have the same issue as the python one, but does not compile because Java requires lambda captures to be (effectively) immutable. This stops the surrounding scope (environment?) from modifying the closure's capture.
- Racket - A new binding of `i` is created for each iteration of the loop. Because each binding is bound to a different reference, the numbers `0-9` are printed as you would expect.

**Prompt:** Explain what happens in this python code:
```python
>>> d = [{}] * 2
>>> d
[{}, {}]

>>> d[0]["a"] = 1
>>> d[1]["a"] = 2
>>> d
[ {'a': 2}, {'a': 2} ]
```

`{}` creates a reference to a dictionary/set object. `[{}] * 2` essentially appends the list/vector to itself, so because the reference is simply copied, both `d[0]` and `d[1]` refer to the same underlying object on the heap. Which means modifying `d[0]` also modifies `d[1]`.

