# Algebraic Data Types
### David Parrott - dmparr22
---

1. **What advantages does `type-case` have over `cond` for the task of dispatching on an ADT value? What disadvantages?**

    - Advantages: `type-case` only allows you to match on variants of the ADT, where  `cond` allows matching against arbitrary predicates. This ensures you only have one valid case for each variant, which can make your code more clear. `type-case` also allows you to lift out the values (properties?) associated with each variant instead of having to use the potentially unsafe accessors.
    - Disadvantages: You have to have define an ADT to use `type-case`, which you may not want or be able to do. If you do want to match against multiple cases on a single variant `cond` *could* make your code more clear.

2. **Discuss a danger that accessors pose related to the number of variants in the datatype. For example, when is it safe to use accessors without qualification and when is it unsafe? Why?**

    - It is safe to use accessors without qualification when there is only one variant of the ADT because every instance is guaranteed to have those properties. It is unsafe whenever the type has more than one variant because accessors raise errors when the property doesn't exist on the variant. This is because plait's type system does not consider the variants as different types; both are of type `('ADT)` 

3. **How is `(none)` distinct from `NULL` and `null` in C and Java respectively? You may answer in the context of a comparison between `(Optionof 'a)` and the behavior of `NULL`/`null` generally.**

    - The main difference is that `(Optionof 'a)` is a distinct type from `'a`. If something returns an option, you have to deal with it, the type checker won't let you keep using it as an `'a`. On the other hand, `null` is still of type `'a`, so nothing is stopping you from using that value without further validation. Options being a different type also allow for compile/build time checks instead of having to do type validation at runtime.

    - Similarly, passing an option into a function makes it clear to the function that it was given a potentially empty value that needs to be validated. Whereas a param of type `'a` also could be null, but nothing about the type tells you that.
