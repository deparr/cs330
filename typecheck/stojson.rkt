#lang racket/base
(require json)

#|

reads s-expression on standard input
produces JSON encoding of it on standard output

examples:

5
=>
{
   "type": "number",
   "value": 5
}

#true
=>
{
   "type": "boolean",
   "value": true
}

"hello"
=>
{
   "type": "string",
   "value": "hello"
}

hello
=>
{
   "type": "symbol",
   "value": "hello"
}

(42 #false hello "goodbye" ())
=>
{
   "type": "list",
   "value": [{
               "type": "number",  
               "value": 42
             },
             {
               "type": "boolean",  
               "value": false
             },
             {
               "type": "symbol",  
               "value": "hello"
             },
             {
               "type": "string",  
               "value": "goodbye"
             },
             {
               "type": "list",  
               "value": []
             }]
}

|#

(write-json
 (let loop ([x (read)])
   (cond
     [(number? x)
      (hasheq 'type "number"
              'value x)]
     [(boolean? x)
      (hasheq 'type "boolean"
              'value x)]
     [(string? x)
      (hasheq 'type "string"
              'value x)]
     [(symbol? x)
      (hasheq 'type "symbol"
              'value (symbol->string x))]
     [(list? x)
      (hasheq 'type "list"
              'value (map loop x))]
     [else
      (error 's-exp->json "unhandled s-exp ~s" x)])))

