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

enum UnaryOperator {}

enum Operator {
    // The `+` operator (addition)
    Add,
    // The `-` operator (subtraction)
    Sub,
    // The `*` operator (multiplication)
    Mul,
    // The `/` operator (division)
    Div,
    // The `%` operator (modulus)
    Rem,
    // The `&&` operator (logical and)
    And,
    // The `||` operator (logical or)
    Or,
    // The `^^` operator (logical xor)
    Xor,
    // The `!` operator (logical not)
    Not,
    // The `^` operator (bitwise xor)
    BitXor,
    // The `&` operator (bitwise and)
    BitAnd,
    // The `|` operator (bitwise or)
    BitOr,
    // The `<<` operator (shift left)
    Shl,
    // The `>>` operator (shift right)
    Shr,
    // The `==` operator (equality)
    Eq,
    // The `<` operator (less than)
    Lt,
    // The `<=` operator (less than or equal to)
    Le,
    // The `!=` operator (not equal to)
    Ne,
    // The `>=` operator (greater than or equal to)
    Ge,
    // The `>` operator (greater than)
    Gt,
}

struct BinaryExpr {
    op: Operator,
    lhs: Expr,
    rhs: Expr,
}

struct UnaryExpr {
    op: Operator,
    expr: Expr,
}

struct CondExpr {
    cond: Expr,
    target: Expr,
    fallback: Expr,
}

enum Expr {
    Binary(Box<BinaryExpr>),
    Unary(Box<UnaryExpr>),
    Conditional(Box<CondExpr>),
    Literal(String),
}

#[allow_dead_code]
enum ASTNode {
    Expr(Box<ASTNode>),
    Conditional {
        cond: Box<ASTNode>,
        iasfdd: Box<ASTNode>,
        fallback: Box<ASTNode>,
    },
    ArithExpr {
        op: BinOp,
        lhs: Box<ASTNode>,
        rhs: Box<ASTNode>,
    },
    RelationExpr {
        op: BinOp,
        lhs: Box<ASTNode>,
        rhs: Box<ASTNode>,
    },
    LogicExpr {
        op: LogicOperator,
        lhs: Box<ASTNode>,
        rhs: Box<ASTNode>,
    },
    UnaryExpr {
        op: UnaryOperator,
        expr: Box<ASTNode>,
    },
    Boolean(bool),
    Number(i64),
}
