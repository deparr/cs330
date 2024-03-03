# WAT
### David Parrott - dmparr22
---
## The C main "function"
Apparently some compilers do not require `main` to be an actual function, they'll just take whatever they get and run with it.
Which means this is a valid C* program, no special compilation flags required:
```c
char main[]__attribute__
((section(".text")))=
"\xb8\x01\x00\x00\x00\xbf\x01\x00\x00\x00\x48\x8d\x35\x0a\x00\x00\x00\x48\xc7\xc2\x0d\x00\x00\x00\x0f\x05\xc3\x48\x65\x6c\x6c\x6f\x20\x57\x6f\x72\x6c\x64\x21\x0a";
```
Outputs:
```
Hello World!

```
_* using gcc-v12.2.0 and clang-v15.0.7 on glibc void-linux-x86-64. Couldn't get it working on windows with mingw. Might work on mac. Will work on cs lab machines_  

I found this out a while ago when I was trying to get the generated object code of a single function and `gcc` errored with 'undefined reference to main'.
I was curious that it only said 'undefined reference' instead of something like 'missing main function'.
I was super surprised when it compiled after adding `int main = 0;`, so I looked into it and discovered you can do ridiculous things like this.  

**dissasembly of main**
```
┌──────────────────────────────────────────────────────────────────────────────────────────┐
│B+ 0x555555555140 <main>           mov    $0x1,%eax             # do a write() syscall... │
│   0x555555555145 <main+5>         mov    $0x1,%edi             # to fd 1 (stdout)        │
│   0x55555555514a <main+10>        lea    0xa(%rip),%rsi        # move str ptr into %rsi  │
│   0x555555555151 <main+17>        mov    $0xd,%rdx             # move strlen into %rdx   │
│   0x555555555158 <main+24>        syscall                      # make the syscall        │
│  >0x55555555515a <main+26>        ret                          # return                  │
│   0x55555555515b <main+27>        rex.W                        # string starts here      │
│   0x55555555515c <main+28>        gs insb (%dx),%es:(%rdi)                               │
│   0x55555555515e <main+30>        insb   (%dx),%es:(%rdi)                                │
│   0x55555555515f <main+31>        outsl  %ds:(%rsi),(%dx)                                │
│   0x555555555160 <main+32>        and    %dl,0x6f(%rdi)                                  │
│   0x555555555163 <main+35>        jb     0x5555555551d1                                  │
│   0x555555555165 <main+37>        and    %ecx,%fs:(%rdx)                                 │
│   0x555555555165 <main+37>        and    %ecx,%fs:(%rdx)                                 │
│   0x555555555168 <main+40>        add    %al,(%rax)            # ends here               │
│                                                                                          │
└──────────────────────────────────────────────────────────────────────────────────────────┘

The string and corresponding instructions:
"\xb8\x01\x00\x00\x00\xbf\x01\x00\x00\x00\x48\x8d\x35\x0a\x00\x00\x00\x48\xc7\xc2\x0d\x00\x00\x00\x0f\x05\xc3\x48\x65\x6c\x6c\x6f\x20\x57\x6f\x72\x6c\x64\x21\x0a"
 |   mov $0x1,%eax   |   mov $0x1,%edi   |    lea 0x0a(%rip),%rsi    |       mov $0xd,%rdx       |syscall|ret| H   e   l   l   o   _   W   o   r   l   d   !   \n|
```
---

## Catching SIGSEGV


---
## C++ Public Cast
Using C++ templates it is possible to gain access to a classes private data and function members, effectively ruining encapsulation.
I originally saw this code in [this youtube video](https://www.youtube.com/watch?v=SmlLdd1Q2V8&t=488s).
I probably won't be able to explain it super well since it goes a little over my head, but here we go:

```cpp
#include <cstdio>
class C {
    int x = 42 // private member
};

int main() {
    C c = C();
    int x = c.x;
    printf("c.x = %d\n", x);
}
```
This code does not compile because `x` is clearly a private member of `C`, which main does not have access to.
However, using some template metaprogramming we can trick the compiler into saving a reference to `x` which we can access later.
This is possible because access modifers are **ignored** during explicit instantiations.

```cpp
#include <cstdio>
class C {
	int x = 42; // private member
};

// this is where &C::x will be stored
//  `M` is the type of the member we are trying to access, int* in our case
//  `Secret` is a dummy type used as a key to refer to a particular instance (and it's corresponding &C::x value)
//  `m` is not declared const, which allows assignment to be chained
template <class M, class Secret>
struct public_cast {
	static inline M m{};
};

//  The chained assignment expresison is the key here:
//      `public_cast<decltype(M), Secret>::m = M`,
//      sets public_cast::m to &C::x and then returns M
//      so it can be assigned to access::m.
template <class Secret, auto M>
struct access {
	static const inline auto m
		= public_cast<decltype(M), Secret>::m = M;
};

// explicit instantitaion of `access`
//  &C::x is accessible here, which allows us to pass it in as `M`.
template struct access<class CxSecret, &C::x>;

int main() {
	C c = C();
    // use CxSecret to 'lookup' the previously stored &C::x
	int x = c.*public_cast<int C::*, CxSecret>::m;
	printf("c.x = %d\n", x);
}

// outputs: 42
```

What happens (as far as I understand it) is that `&C::x` is accessible during the type-explicit instantiation of the `access` struct.
Then, as part of instantiating the `access` struct, a `public_cast` struct is also created that records the value of `&C::x`.
This `public_cast` struct can then be used to access `c.x` as if it were public because it does not require referencing `&C::x` by name (like we did when instantiating `access`).
I'm not 100% sure why the `Secret` and `CxSecret` type parameters are needed, but I think it has to do with refering back to the specific `public_cast` that stored the private reference.
