#lang plait #:untyped
(require "parser-combinator.rkt"
         "json-util.rkt")

(parse stringp "\"balls\"")
(parse numberp "23423434")
(parse (delimitedp (literalp "[") (literalp "]") (literalp ",") numberp) "[123,12312,123,123,123,123,123,123,123,12,3]")
