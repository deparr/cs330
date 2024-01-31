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

use std::fmt::Display;

#[derive(Debug)]
enum UnaryOp {
    // The `+` operator (unary plus)
    Plus,
    // The `-` operator (unary minus)
    Minus,
    // The `!` operator (logical not)
    Not,
    // The `~` operator (bitwise not)
    BitNot,
}

impl Display for UnaryOp {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        use UnaryOp::*;
        match self {
            Plus => write!(f, "+"),
            Minus => write!(f, "-"),
            Not => write!(f, "!"),
            BitNot => write!(f, "~"),
        }
    }
}

#[allow(dead_code)]
#[derive(Debug)]
enum BinOp {
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

impl Display for BinOp {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        use BinOp::*;
        match self {
            Add => write!(f, "+"),
            Sub => write!(f, "-"),
            Mul => write!(f, "*"),
            Div => write!(f, "/"),
            Rem => write!(f, "%"),
            And => write!(f, "&&"),
            Or => write!(f, "||"),
            BitXor => write!(f, "^"),
            BitAnd => write!(f, "&"),
            BitOr => write!(f, "|"),
            Shl => write!(f, ">>"),
            Shr => write!(f, "<<"),
            Eq => write!(f, "=="),
            Lt => write!(f, "<"),
            Le => write!(f, "<="),
            Ne => write!(f, "!="),
            Ge => write!(f, ">="),
            Gt => write!(f, ">"),
        }
    }
}

#[derive(Debug)]
struct BinaryExpr {
    op: BinOp,
    lhs: Expr,
    rhs: Expr,
}

impl BinaryExpr {
    fn new(expr: &serde_json::Value) -> Self {
        let op = match expr.get("operator") {
            Some(op) => op.as_str().unwrap(),
            None => todo!("none binop operator"),
        };

        use BinOp::*;
        let op = match op {
            "+" => Add,
            "-" => Sub,
            "*" => Mul,
            "/" => Div,
            "%" => Rem,
            "&&" => And,
            "||" => Or,
            "==" => Eq,
            "!=" => Ne,
            "<" => Lt,
            "<=" => Le,
            ">" => Gt,
            ">=" => Ge,
            _ => unimplemented!("unimplemented binop: {}", op),
        };
        let lhs = expr.get("left").unwrap();
        let rhs = expr.get("right").unwrap();

        BinaryExpr {
            op,
            lhs: Expr::new(lhs),
            rhs: Expr::new(rhs),
        }
    }
}

#[derive(Debug)]
struct UnaryExpr {
    op: UnaryOp,
    expr: Expr,
}

impl UnaryExpr {
    fn new(expr: &serde_json::Value) -> Self {
        let op = match expr.get("operator") {
            Some(op) => op.as_str().unwrap(),
            None => todo!("none binop operator"),
        };

        use UnaryOp::*;
        let op = match op {
            "+" => Plus,
            "-" => Minus,
            "!" => Not,
            "~" => BitNot,
            _ => unimplemented!("unimplemented unaryop: {}", op),
        };

        let child_expr = expr.get("argument").unwrap();

        UnaryExpr {
            op,
            expr: Expr::new(child_expr),
        }
    }
}

#[derive(Debug)]
struct CondExpr {
    test: Expr,
    cons: Expr,
    altr: Expr,
}

impl CondExpr {
    fn new(expr: &serde_json::Value) -> Self {
        let test = expr.get("test").unwrap();
        let cons = expr.get("consequent").unwrap();
        let altr = expr.get("alternate").unwrap();

        CondExpr {
            test: Expr::new(test),
            cons: Expr::new(cons),
            altr: Expr::new(altr),
        }
    }
}

#[derive(Debug)]
enum Expr {
    Binary(Box<BinaryExpr>),
    Unary(Box<UnaryExpr>),
    Conditional(Box<CondExpr>),
    Literal(String),
}

impl Expr {
    fn new(expr: &serde_json::Value) -> Self {
        let expr_type = match expr.get("type") {
            Some(t) => t.as_str().unwrap(),
            None => unreachable!("none expr type should be unreachable"),
        };

        match expr_type {
            "BinaryExpression" | "LogicalExpression" => {
                Expr::Binary(Box::new(BinaryExpr::new(expr)))
            }
            "UnaryExpression" => Expr::Unary(Box::new(UnaryExpr::new(expr))),
            "ConditionalExpression" => Expr::Conditional(Box::new(CondExpr::new(expr))),
            "Literal" => Expr::Literal(String::from(expr.get("raw").unwrap().as_str().unwrap())),
            "Expression" => todo!("`expression` type match arm"),
            _ => todo!("reached end of expr type match arms: {}", expr_type),
        }
    }
}

impl Display for Expr {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        use Expr::*;
        match self {
            Binary(expr) => {
                use BinOp::*;
                let op_type = match expr.op {
                    Add | Sub | Mul | Div | Rem | BitOr | BitAnd | BitXor | Shl | Shr => {
                        "arithmetic"
                    }
                    Or | And => "logical",
                    Eq | Ne | Lt | Le | Gt | Ge => "relational",
                };
                write!(f, "({} {} {} {})", op_type, expr.op, expr.lhs, expr.rhs)
            }
            Unary(expr) => {
                write!(f, "(unary {} {})", expr.op, expr.expr)
            }
            Conditional(expr) => {
                write!(f, "(conditional {} {} {})", expr.test, expr.cons, expr.altr)
            }
            Literal(raw) => {
                // TODO this needs to change later
                if raw.bytes().all(|c| c.is_ascii_digit()) {
                    write!(f, "(number {})", raw)
                } else {
                    write!(f, "(boolean {})", raw)
                }
            }
        }
    }
}

#[derive(Debug)]
pub struct Program {
    prog: Expr,
}

impl Program {
    pub fn new(expr: &serde_json::Value) -> Self {
        Program {
            prog: Expr::new(expr),
        }
    }
}

impl Display for Program {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", self.prog)
    }
}
