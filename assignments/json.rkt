#lang plait #:untyped
(require "parser-combinator.rkt"
         "json-util.rkt")
(print-only-errors #t)

(define-type JSON
  (JSON-null)
  (JSON-number [number  : Number])
  (JSON-string [string  : String])
  (JSON-array  [array   : (Listof JSON)])
  (JSON-object [entries : (Hashof String JSON)]))

(define json-numberp
  (bind1 (λ (n) (unitp (JSON-number n)))
         numberp))
(define json-stringp
  (bind1 (λ (n) (unitp (JSON-string n)))
         stringp))
(define json-nullp
  (bind1 (λ (_) (unitp (JSON-null)))
         (literalp "null")))

(JSONp : (Parser JSON 'b))
(define JSONp
  (altp (list (bind1 (λ (l) (unitp (JSON-array l)))
                     (delimitedp (literalp "[")
                                 (literalp "]")
                                 (literalp ",")
                                 (altp (list (literalp "null")
                                             numberp
                                             stringp))))
              (delimitedp (literalp "{") (literalp "}") (literalp ",") (unitp 'balls)))))

(test (parse JSONp #<<JSON
[]
JSON
             )
      (JSON-array (list)))


(test (parse JSONp #<<JSON
[null]
JSON
             )
      (JSON-array (list (JSON-null))))

(test (parse JSONp #<<JSON
[null,1]
JSON
             )
      (JSON-array (list (JSON-null) (JSON-number 1))))


(test (parse JSONp #<<JSON
[1]
JSON
             )
      (JSON-array (list (JSON-number 1))))

(test (parse JSONp #<<JSON
[1,2,"hello"]
JSON
             )
      (JSON-array (list (JSON-number 1)
                        (JSON-number 2)
                        (JSON-string "hello"))))

(test (parse JSONp #<<JSON
{}
JSON
             )
      (JSON-object (hash (list))))

(test (parse JSONp #<<JSON
{
}
JSON
             )
      (JSON-object (hash (list))))

(test (parse JSONp #<<JSON
{"a":null}
JSON
             )
      (JSON-object (hash (list (values "a" (JSON-null))))))

(test (parse JSONp #<<JSON
{ "a" : null }
JSON
             )
      (JSON-object (hash (list (values "a" (JSON-null))))))

(test (parse JSONp #<<JSON
{ "a" : null
, "b" : null
}
JSON
             )
      (JSON-object (hash (list (values "a" (JSON-null))
                               (values "b" (JSON-null))))))

(test (parse JSONp #<<JSON
{ "a" : null
, "b" : null
, "b" : 42
}
JSON
             )
      (JSON-object (hash (list (values "a" (JSON-null))
                               (values "b" (JSON-number 42))))))

(test (parse JSONp #<<JSON
{ "a" : null
, "b" : null
, "b" : 42
, "c" : [ 1, 2, null ]
}
JSON
             )
      (JSON-object (hash (list (values "a" (JSON-null))
                               (values "b" (JSON-number 42))
                               (values "c" (JSON-array (list (JSON-number 1)
                                                             (JSON-number 2)
                                                             (JSON-null))))))))

