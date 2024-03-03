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
I was super surprised when it compiled after adding `int main = 0;`, so I looked into it and found out you can do ridiculous things like this.  

**dissasembly of main**
```
┌──────────────────────────────────────────────────────────────────────────────────────────┐
│B+ 0x555555555140 <main>           mov    $0x1,%eax             # do a write syscall...   │
│   0x555555555145 <main+5>         mov    $0x1,%edi             # to fd 1, stdout         │
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
"\xb8\x01\x00\x00\x00\xbf\x01\x00\x00\x00\x48\x8d\x35\x0a\x00\x00\x00\x48\xc7\xc2\x0d\x00\x00\x00\x0f\x05\xc3\x48\x65\x6c\x6c\x6f\x20\x57\x6f\x72\x6c\x64\x21\x0a"
 |   mov $0x1,%eax   |   mov $0x1,%edi   |    lea 0x0a(%rip),%rsi    |       mov $0xd,%rdx       |syscall|ret| H   e   l   l   o   _   W   o   r   l   d   !   \n|
```
---

