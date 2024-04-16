#lang plait
(require "json.rkt")

(define (good-error who message detail)
  (error who (string-append message (string-append " " (string-append detail " (unreachable by grader input)")))))

(define-type BinOp
  (opAdd)
  (opSub)
  (opMul)
  (opDiv)
  (opEq)
  (opLt))

(define-type LogOp
  (opAnd)
  (opOr))

(define-type Expr
  (numE [n : Number])
  (booE [b : Boolean])
  (refE [x : Symbol])
  (binE [op : BinOp]
        [lhs : Expr]
        [rhs : Expr])
  (logE [op : LogOp]
        [lhs : Expr]
        [rhs : Expr])
  (unaE [arg : Expr])
  (_ifE [tst : Expr]
        [thn : Expr]
        [els : Expr])
  (funE [params : (Listof Symbol)]
        [body : Expr])
  (appE [fun : Expr]
        [args : (Listof Expr)])
  (setE [lvalue : LVExpr]
        [arg : Expr])
  (objE [fields : (Listof (Symbol * Expr))])
  (memE [obj : Expr]
        [x : Symbol])
  (letE [bindings : (Listof (Symbol * Expr))]
        [body : Expr])
  (begE [first : Expr]
        [last : Expr])
  (whileE [test : Expr]
          [body : Expr])
  (forE [binds : (Listof (Symbol * Expr))]
        [test : Expr]
        [updt : Expr]
        [body : Expr]))

(define-type LVExpr
  (varLVE [x : Symbol])
  (memLVE [obj : Expr]
          [x : Symbol]))

(define (has-type? json type)
  (equal? (JSON->string (JSON-extract json ".type")) type))

(parse-binary-op : (String -> BinOp))
(define (parse-binary-op op)
  (cond
    [(equal? op "+")
     (opAdd)]
    [(equal? op "-")
     (opSub)]
    [(equal? op "*")
     (opMul)]
    [(equal? op "/")
     (opDiv)]
    [(equal? op "==")
     (opEq)]
    [(equal? op "<")
     (opLt)]
    [else
     (good-error 'parse-binary-op "unknown binary operator" op)]))

(parse-logical-op : (String -> LogOp))
(define (parse-logical-op op)
  (cond
    [(equal? op "&&")
     (opAnd)]
    [(equal? op "||")
     (opOr)]
    [else
     (good-error 'parse-logical-op "unknown logical operator" op)]))

(define (parse-expr json)
  (cond
    [(has-type? json "Literal")
     (let ([value (JSON-extract json ".value")])
       (cond
         [(JSON-number? value)
          (numE (JSON->number value))]
         [(JSON-boolean? value)
          (booE (JSON->boolean value))]
         [else
          (good-error 'parse-expr "unknown literal type" (to-string value))]))]
    [(has-type? json "Identifier")
     (refE (string->symbol (JSON->string (JSON-extract json ".name"))))]
    [(has-type? json "BinaryExpression")
     (binE (parse-binary-op (JSON->string (JSON-extract json ".operator")))
           (parse-expr (JSON-extract json ".left"))
           (parse-expr (JSON-extract json ".right")))]
    [(has-type? json "LogicalExpression")
     (logE (parse-logical-op (JSON->string (JSON-extract json ".operator")))
           (parse-expr (JSON-extract json ".left"))
           (parse-expr (JSON-extract json ".right")))]
    [(has-type? json "ConditionalExpression")
     (_ifE (parse-expr (JSON-extract json ".test"))
           (parse-expr (JSON-extract json ".consequent"))
           (parse-expr (JSON-extract json ".alternate")))]
    [(has-type? json "UnaryExpression")
     (if (equal? (JSON->string (JSON-extract json ".operator")) "!")
       (unaE (parse-expr (JSON-extract json ".argument")))
       (good-error 'parse-expr "bady unary operator" (JSON->string (JSON-extract json ".operator"))))]
    [(has-type? json "FunctionExpression")
     (funE (map
            (λ (param)
              (if (has-type? param "Identifier")
                (string->symbol (JSON->string (JSON-extract param ".name")))
                (good-error 'parse-expr "not an identifier" (JSON->string (JSON-extract param ".type")))))
            (JSON->list (JSON-extract json ".params")))
           (let ([body (JSON-extract json ".body")])
             (if (has-type? body "BlockStatement")
               (parse-block (JSON->list (JSON-extract body ".body")) #false)
               (good-error 'parse-expr "not a block statement" (JSON->string (JSON-extract body ".type"))))))]
    [(has-type? json "CallExpression")
     (appE (parse-expr (JSON-extract json ".callee"))
           (map parse-expr (JSON->list (JSON-extract json ".arguments"))))]
    [(has-type? json "AssignmentExpression")
     (setE (parse-lvalue-expr (JSON-extract json ".left"))
           (parse-expr (JSON-extract json ".right")))]
    [(has-type? json "ObjectExpression")
     (objE (map
            (λ (property)
              (if (has-type? property "Property")
                (pair (string->symbol (JSON->string (JSON-extract property ".key.name")))
                      (parse-expr (JSON-extract property ".value")))
                (good-error 'parse-expr "not a property"  (JSON->string (JSON-extract property ".type")))))
            (JSON->list (JSON-extract json ".properties"))))]
    [(has-type? json "MemberExpression")
     (memE (parse-expr (JSON-extract json ".object"))
           (string->symbol (JSON->string (JSON-extract json ".property.name"))))]
    [else
     (good-error 'parse-expr "unknown expression type"  (JSON->string (JSON-extract json ".type")))]))

(parse-lvalue-expr : (JSON -> LVExpr))
(define (parse-lvalue-expr json)
  (cond
    [(has-type? json "Identifier")
     (varLVE (string->symbol (JSON->string (JSON-extract json ".name"))))]
    [(has-type? json "MemberExpression")
     (memLVE (parse-expr (JSON-extract json ".object"))
             (string->symbol (JSON->string (JSON-extract json ".property.name"))))]))

(parse-block : ((Listof JSON) Boolean -> Expr))
(define (parse-block jsons top-level?)
  (type-case (Listof JSON) jsons
    [empty
     (error 'parse-block "empty block (unreachable by grader input)")]
    [(cons json jsons)
     (if (empty? jsons)
       ; this is the final statement in the sequence
       (cond
         [(and (has-type? json "ExpressionStatement")
               top-level?)
          (parse-expr (JSON-extract json ".expression"))]
         [(and (has-type? json "ReturnStatement")
               (not top-level?))
          (parse-expr (JSON-extract json ".argument"))]
         [else
          (good-error 'parse-block "incorrect final statement type" (JSON->string (JSON-extract json ".type")))])
       (cond
         [(has-type? json "VariableDeclaration")
          (letE (map
                 (λ (declaration)
                   (if (equal? (JSON->string (JSON-extract declaration ".type"))
                               "VariableDeclarator")
                     (pair (string->symbol (JSON->string (JSON-extract declaration ".id.name")))
                           (parse-expr (JSON-extract declaration ".init")))
                     (good-error 'parse-block "incorrect declarator type" (JSON->string (JSON-extract declaration ".type")))))
                 (JSON->list (JSON-extract json ".declarations")))
                (parse-block jsons top-level?))]
         [(has-type? json "ExpressionStatement")
          (begE (parse-expr (JSON-extract json ".expression"))
                (parse-block jsons top-level?))]

         ;; wrapping these in a (begE) feels like cheating...
         [(has-type? json "WhileStatement")
          (begE (whileE (parse-expr (JSON-extract json ".test"))
                        (let ([body (JSON-extract json ".body")])
                            (if (has-type? body "BlockStatement")
                                (parse-block (JSON->list (JSON-extract body ".body")) #true)
                                ((good-error 'parse-expr "not a block statement" (JSON->string (JSON-extract body ".type")))))))
                (parse-block jsons top-level?))]

         [(has-type? json "ForStatement")
            (begE (forE
                    (map
                        (λ (declaration)
                            (if (equal? (JSON->string (JSON-extract declaration ".type")) "VariableDeclarator")
                                (pair (string->symbol (JSON->string (JSON-extract declaration ".id.name")))
                                      (parse-expr (JSON-extract declaration ".init")))
                                (good-error 'parse-block "incorrect declarator type" (JSON->string (JSON-extract declaration ".type")))))
                        (JSON->list (JSON-extract json ".init.declarations")))
                    (parse-expr (JSON-extract json ".test"))
                    (parse-expr (JSON-extract json ".update"))
                    (let ([body (JSON-extract json ".body")])
                        (if (has-type? body "BlockStatement")
                            (parse-block (JSON->list (JSON-extract body ".body")) #true)
                            ((good-error 'parse-expr "not a block statement" (JSON->string (JSON-extract body ".type")))))))
                  (parse-block jsons top-level?))]

         [else
          (good-error 'parse-block "incorrect body statement type" (JSON->string (JSON-extract json ".type")))]))]))

(parse : (JSON -> Expr))
(define (parse json)
  (if (and (equal? (JSON->string (JSON-extract json ".type")) "Program")
           (equal? (JSON->string (JSON-extract json ".sourceType")) "script"))
    (parse-block (JSON->list (JSON-extract json ".body")) #true)
    (error 'parse "bad top-level program (unreachable by grader input)")))

(define-type-alias Variable Symbol)
(define-type-alias Address Number)

(define-type-alias Environment (Hashof Variable Address))

(define env-empty (hash empty))
(define (env-extend y v env)
  (hash-set env y v))
(define (env-lookup env x)
  (type-case (Optionof Address) (hash-ref env x)
    [(some addr)
     addr]
    [(none)
     (error 'env-lookup "unbound identifier")]))

(define-type Value
  (numV [n : Number])
  (booV [b : Boolean])
  (funV [xs : (Listof Symbol)]
        [body : Expr]
        [env : Environment])
  (objV [assocs : Environment])
  (voiV))

(define (unit x)
  (λ (sto vk ek) (vk x sto)))

(define (>>= c f)
  (λ (sto vk ek) (c sto (λ (x sto) ((f x) sto vk ek)) ek)))

(define (>> c₀ c₁) (>>= c₀ (λ (_) c₁)))

(define (try c f)
  (λ (sto vk ek) (c sto vk (λ (x sto) ((f x) sto vk ek)))))

(define (raise v)
  (λ (sto vk ek) (ek v sto)))

(define (bind1 f c) (>>= c f))
(define (bind2 f c₀ c₁) (>>= c₀ (λ (x₀) (>>= c₁ (λ (x₁) (f x₀ x₁))))))

(define (fmap1 f c) (bind1 (λ (x) (unit (f x))) c))
(define (fmap2 f c₀ c₁) (bind2 (λ (x y) (unit (f x y))) c₀ c₁))

(define (foldlm f u xs)
  (type-case (Listof 'a) xs
    [empty
     (unit u)]
    [(cons x rest-xs)
     (>>= (f x u) (λ (u) (foldlm f u rest-xs)))]))

(define alloc
  (λ (sto vk ek) (vk (length (hash-keys sto)) sto)))

(define store-empty (hash empty))

(define (store-lookup addr)
  (λ (sto vk ek)
    (vk (type-case (Optionof Value) (hash-ref sto addr)
          [(some v)
           v]
          [(none)
           (error 'store-lookup "segmentation fault")])
        sto)))

(define (store-update addr v)
  (λ (sto vk ek) (vk (voiV) (hash-set sto addr v))))

(define (num×num→num f)
  (λ (v₀ v₁)
    (type-case Value v₀
      [(numV n₀)
       (type-case Value v₁
         [(numV n₁)
          (fmap1 numV (f n₀ n₁))]
         [else
          (raise "not a number (banana)")])]
      [else
       (raise "not a number (banana)")])))

(define (num×num→bool f)
  (λ (v₀ v₁)
    (type-case Value v₀
      [(numV n₀)
       (type-case Value v₁
         [(numV n₁)
          (fmap1 booV (f n₀ n₁))]
         [else
          (raise "not a number (banana)")])]
      [else
       (raise "not a number (banana)")])))

(define (lift2 f)
  (λ (x y) (unit (f x y))))

(define (interp-binary-op op)
  (type-case BinOp op
    [(opAdd)
     (num×num→num (lift2 +))]
    [(opSub)
     (num×num→num (lift2 -))]
    [(opMul)
     (num×num→num (lift2 *))]
    [(opDiv)
     (num×num→num (λ (n d) (if (zero? d) (raise "divide by zero (banana)") (unit (floor (/ n d))))))]
    [(opEq)
     (num×num→bool (lift2 =))]
    [(opLt)
     (num×num→bool (lift2 <))]))

(define (boolean v)
  (type-case Value v
    [(booV b)
     (unit b)]
    [else
     (raise "not a boolean (banana)")]))

(define (function v)
  (type-case Value v
    [(funV xs body env)
     (unit (values xs body env))]
    [else
     (raise "not a function (banana)")]))

(define (object v)
  (type-case Value v
    [(objV assocs)
     (unit assocs)]
    [else
     (raise "not an object (banana)")]))

(define (linterp lexp env)
  (type-case LVExpr lexp
    [(varLVE x)
     (let ([addr (env-lookup env x)])
       (unit (λ (v) (store-update addr v))))]
    [(memLVE obj x)
     (bind1
      (λ (v)
        (type-case Value v
          [(objV assocs)
           (let ([addr (env-lookup assocs x)])
             (unit (λ (v) (store-update addr v))))]
          [else
           (raise "expected an object (banana)")]))
      (interp obj env))]))

(define (interp-bindings bindings env accum-env)
  (foldlm
   (λ (binding accum-env)
     (local [(define-values (x exp) binding)]
       (bind2
        (λ (v addr)
          (>> (store-update addr v)
              (unit (env-extend x addr accum-env))))
        (interp exp env) alloc)))
   accum-env
   bindings))

(define (interp exp env)
  (type-case Expr exp
    [(numE n)
     (unit (numV n))]
    [(booE b)
     (unit (booV b))]
    [(refE x)
     (store-lookup (env-lookup env x))]
    [(binE op lhs rhs)
     (bind2
      (interp-binary-op op)
      (interp lhs env)
      (interp rhs env))]
    [(logE op lhs rhs)
     (bind1
      (λ (v)
        (bind1
         (λ (b)
           (if b
             (unit (booV b))
             (fmap1 booV (bind1 boolean (interp rhs env)))))
         (boolean v)))
      (interp lhs env))]
    [(unaE arg)
     (bind1
      (λ (v)
        (fmap1
         (λ (b) (booV (not b)))
         (boolean v)))
      (interp arg env))]
    [(_ifE tst thn els)
     (bind1
      (λ (v)
        (bind1
         (λ (b) (interp (if b thn els) env))
         (boolean v)))
      (interp tst env))]
    [(funE xs body)
     (unit (funV xs body env))]
    [(appE fun args)
     (bind1
      (λ (fun-v)
        (bind1
         (λ (xs×body×env)
           (local [(define-values (xs body clo-env) xs×body×env)]
             (bind1
              (λ (env) (interp body env))
              (interp-bindings (map2 pair xs args) env clo-env))))
         (function fun-v)))
      (interp fun env))]
    [(setE lexp arg)
     (bind1
      (λ (set!) (bind1 set! (interp arg env)))
      (linterp lexp env))]
    [(objE fields)
     (fmap1
      objV
      (interp-bindings fields env env-empty))]
    [(memE obj x)
     (bind1
      (λ (v)
        (bind1
         (λ (assocs)
           (store-lookup (env-lookup assocs x)))
         (object v)))
      (interp obj env))]
    [(letE bindings body)
     (bind1
      (λ (env) (interp body env))
      (interp-bindings bindings env env))]
    [(begE fst snd)
     (>> (interp fst env)
         (interp snd env))]
    [(whileE test body)
     (bind1
       (λ (t)
         (type-case Value t
            [(booV v)
             (if v
               (>> (interp body env)
                   (interp exp env))
               (unit (voiV)))]
            [else
              (error 'interp "while loop test must eval to bool")]))
       (interp test env))]
    [(forE binds test updt body)
     (bind1
       (λ (env)
         (letrec
           ([eval-for (λ (_)
                        (bind1
                          (λ (t)
                            (type-case Value t
                              [(booV v)
                                (if v
                                  (>> (>> (interp body env)
                                          (interp updt env))
                                      (eval-for _))
                                  (unit (voiV)))]
                                [else
                                  (error 'interp "for loop test must eval to bool")]))
                          (interp test env)))])
             (eval-for #f)))
       (interp-bindings binds env env))
    ]))

#;
(module+ main
  (display (parse (read-JSON))))

(module+ main
  (define (tagged tag x)
    (string-append
     "("
     (string-append
      (symbol->string tag)
      (string-append
       " "
       (string-append
        x
        ")")))))
  
  ((interp (parse (read-JSON)) env-empty)
   store-empty
   (λ (v sto)
     (display
      (tagged
       'value
       (local [(define (value-to-string v)
                 (type-case Value v
                   [(numV n)
                    (tagged 'number (to-string n))]
                   [(booV b)
                    (tagged 'boolean (if b "true" "false"))]
                   [(funV xs body env)
                    "(function)"]
                   [(objV assocs)
                    (string-append
                     "(object"
                     (string-append
                      (foldr
                       (λ (k s)
                         (string-append
                          " ["
                          (string-append
                           (symbol->string k)
                           (string-append
                            " "
                            (string-append
                             (value-to-string (some-v (hash-ref sto (some-v (hash-ref assocs k)))))
                             (string-append
                              "]"
                              s))))))
                       ""
                       (hash-keys assocs))
                      ")"))]
                   [(voiV)
                    "(void)"]))]
         (value-to-string v)))))
   (λ (e sto) (display (tagged 'error (to-string e))))))

