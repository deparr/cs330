#lang plait

(module access racket/base
  (require json
           racket/match)

  (struct json (content)
    #:property prop:equal+hash
    (list (λ (a b eql?)
            (eql? (json-content a)
                  (json-content b)))
          (λ (a hc)
            (hc (json-content a)))
          (λ (a hc)
            (hc (json-content a))))
    #:property prop:custom-print-quotable
    'never
    #:property prop:custom-write
    (λ (s out depth)
      (display "#<json:" out)
      (write-json (json-content s) out)
      (display ">" out)))

  (define (read-JSON) (json (read-json)))
  (define JSON? jsexpr?)

  (define (JSON-null? js) (eq? (json-content js) 'null))
  (define JSON-null (json 'null))
  (define (JSON-boolean? js) (boolean? (json-content js)))
  (define (JSON->boolean js)
    (let ([v (json-content js)])
      (if (boolean? v) v (error 'json->boolean "expected a JSON boolean"))))
  (define boolean->JSON json)
  (define (JSON-number? js) (number? (json-content js)))
  (define (JSON->number js)
    (let ([v (json-content js)])
      (if (number? v) v (error 'json->number "expected a JSON number"))))
  (define number->JSON json)
  (define(JSON-string? js) (string? (json-content js)))
  (define (JSON->string js)
    (let ([v (json-content js)])
      (if (string? v) v (error 'json->string "expeceted a JSON string"))))
  (define string->JSON json)
  (define (JSON-array? js) (list? (json-content js)))
  (define (JSON->list js)
    (let ([v (json-content js)])
      (if (list? v) (map json v) (error 'json->list "expected a JSON array"))))
  (define (list->JSON jss) (json (map json-content jss)))
  (define (JSON-object? js) (hash? (json-content js)))
  (define (JSON->hash js)
    (let ([v (json-content js)])
      (if (hash? v)
          (for/hasheq ([(k v) (in-hash v)])
            (values k (json v)))
          (error 'json->hash "expected a JSON object"))))
  (define (hash->JSON hjs)
    (json
     (for/hasheq ([(k js) (in-hash hjs)])
       (values k (json-content js)))))

  (define (JSON-extract js path)
    (match path
      [""
       js]
      [(regexp #px"^\\.([^\\.\\[]+)(.*)$" (list _ (app string->symbol name) path))
       (JSON-extract (hash-ref (JSON->hash js) name) path)]
      [(regexp #px"^\\[(\\d+)\\](.*)$" (list _ (app string->number idx) path))
       (JSON-extract (list-ref (JSON->list js) idx) path)]))

  (provide (all-defined-out)))

(require (opaque-type-in (submod "." access) [JSON JSON?])
         (rename-in (typed-in (submod "." access)
                              [read-JSON : (-> JSON)]
                              [JSON-null? : (JSON -> Boolean)]
                              [JSON-null : JSON]
                              [JSON-boolean? : (JSON -> Boolean)]
                              [JSON->boolean : (JSON -> Boolean)]
                              [boolean->JSON : (Boolean -> JSON)]
                              [JSON-number? : (JSON -> Boolean)]
                              [JSON->number : (JSON -> Number)]
                              [number->JSON : (Number -> JSON)]
                              [JSON-string? : (JSON -> Boolean)]
                              [JSON->string : (JSON -> String)]
                              [string->JSON : (String -> JSON)]
                              [JSON-array? : (JSON -> Boolean)]
                              [JSON->list : (JSON -> (Listof JSON))]
                              [list->JSON : ((Listof JSON) -> JSON)]
                              [JSON-object? : (JSON -> Boolean)]
                              [JSON->hash : (JSON -> (Hashof Symbol JSON))]
                              [hash->JSON : ((Hashof Symbol JSON) -> JSON)]
                              [JSON-extract : (JSON String -> JSON)])
                    [read-JSON *read-JSON]
                    [JSON-null? *JSON-null?]
                    [JSON-null *JSON-null]
                    [JSON-boolean? *JSON-boolean?]
                    [JSON->boolean *JSON->boolean]
                    [boolean->JSON *boolean->JSON]
                    [JSON-number? *JSON-number?]
                    [JSON->number *JSON->number]
                    [number->JSON *number->JSON]
                    [JSON-string? *JSON-string?]
                    [JSON->string *JSON->string]
                    [string->JSON *string->JSON]
                    [JSON-array? *JSON-array?]
                    [JSON->list *JSON->list]
                    [list->JSON *list->JSON]
                    [JSON-object? *JSON-object?]
                    [JSON->hash *JSON->hash]
                    [hash->JSON *hash->JSON]
                    [JSON-extract *JSON-extract]))

(define read-JSON *read-JSON)
(define JSON-null? *JSON-null?)
(define JSON-null *JSON-null)
(define JSON-boolean? *JSON-boolean?)
(define JSON->boolean *JSON->boolean)
(define boolean->JSON *boolean->JSON)
(define JSON-number? *JSON-number?)
(define JSON->number *JSON->number)
(define number->JSON *number->JSON)
(define JSON-string? *JSON-string?)
(define JSON->string *JSON->string)
(define string->JSON *string->JSON)
(define JSON-array? *JSON-array?)
(define JSON->list *JSON->list)
(define list->JSON *list->JSON)
(define JSON-object? *JSON-object?)
(define JSON->hash *JSON->hash)
(define hash->JSON *hash->JSON)
(define JSON-extract *JSON-extract)

