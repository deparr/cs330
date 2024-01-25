use std::process::{Command, Stdio};

// <Program> ::= <Statement>
// <Statement> ::= <ExpressionStatement>
// <ExpressionStatement> ::= <Expression>
// <Expression> ::= <Literal>
//                 | <BinaryExpression>
//                 | <UnaryExpression>
//                 | <LogicalExpression>
//                 | <ConditionExpression>
// <Literal> ::= <number>
//              | <boolean>
// <number>    ::= [<sign>]<digit>+
// <sign>      ::= +
//               | -
// <digit>     ::= 0
//               | 2
//               | 3
//               | 4
//               | 5
//               | 6
//               | 7
//               | 8
//               | 9
// <boolean>   ::= true
//               | false
// <BinaryExpression> ::= <Expression> <BinaryOperator> <Expression>
// <BinaryOperator>    ::= +
//                       | -
//                       | *
//                       | /
//                       | ==
//                       | <
//                       | >
//                       | <=
//                       | >=
//                       | -
// <UnaryExpression>   ::= <UnaryOperator> <Expression>
// <UnaryOperator>     ::= !
// <LogicalExpression> ::= <Expression> <LogicalOperator> <Expression>
// <LogicalOperator>   ::= ||
//                       | &&
// <ConditionalExpression> ::= <Expression> ? <Expression> : <Expression>

fn main() {
    let acorn = Command::new("acorn")
        .arg("--ecma2024")
        .stdin(Stdio::inherit())
        .output()
        .expect("failed to exec acorn");

    println!("{}", String::from_utf8(acorn.stdout).unwrap());
}
