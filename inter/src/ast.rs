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
pub enum Value {
    String(String),
    Bool(bool),
    Int(i64),
    Float(f64),
    // error?
    // Error(String),
}

fn eval_error(msg: &str) -> Result<Value, &str> {
    Err(msg)
}

impl Value {
    fn new(expr: &serde_json::Value) -> Self {
        let value = expr.get("value").unwrap();
        if value.is_i64() {
            Self::Int(value.as_i64().unwrap())
        } else if value.is_f64() {
            Self::Float(value.as_f64().unwrap())
        } else if value.is_string() {
            Self::String(String::from(value.as_str().unwrap()))
        } else if value.is_boolean() {
            Self::Bool(value.as_bool().unwrap())
        } else {
            unimplemented!("tried to create inter::ast::Value out of unsupported type")
        }
    }
}

impl Display for Value {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        use Value::*;
        match self {
            Int(v) => write!(f, "(value (number {}))", v),
            Float(v) => write!(f, "(value (number {}))", v),
            Bool(v) => write!(f, "(value (boolean {}))", v),
            String(v) => write!(f, "(value (string {}))", v),
        }
    }
}

#[derive(Debug)]
enum Expr {
    Binary(Box<BinaryExpr>),
    Unary(Box<UnaryExpr>),
    Conditional(Box<CondExpr>),
    Literal(Value),
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
            "Literal" => Expr::Literal(Value::new(expr)),
            "Expression" => todo!("`expression` type match arm"),
            _ => todo!("reached end of expr type match arms: {}", expr_type),
        }
    }

    fn eval(&self) -> Result<Value, &str> {
        use Expr::*;
        use Value::*;
        match self {
            Binary(expr) => {
                let lhs = expr.lhs.eval()?;
                let rhs = expr.rhs.eval()?;

                use BinOp::*;
                // this is so ass
                match expr.op {
                    // number -> number
                    Add | Sub | Mul | Div => {
                        match (lhs, rhs) {
                            (Int(l), Int(r)) => {
                                // not sure if i can capture the op match somehow
                                // this is ass
                                match expr.op {
                                    Add => Ok(Int(l + r)),
                                    Sub => Ok(Int(l - r)),
                                    Mul => Ok(Int(l * r)),
                                    Div => {
                                        if r == 0 {
                                            eval_error("divide by zero")
                                        } else {
                                            Ok(Int(l / r))
                                        }
                                    }
                                    _ => unreachable!("op in BinOp::Int"),
                                }
                            }
                            (Float(l), Float(r)) => match expr.op {
                                Add => Ok(Float(l + r)),
                                Sub => Ok(Float(l - r)),
                                Mul => Ok(Float(l * r)),
                                Div => {
                                    if r == 0. {
                                        eval_error("divide by zero")
                                    } else {
                                        Ok(Float(l / r))
                                    }
                                }
                                _ => unreachable!("op in BinOp::Int"),
                            },
                            _ => eval_error("non numbers in numerical binop"),
                        }
                    }

                    // number -> boolean
                    Lt | Le | Gt | Ge => {
                        match (lhs, rhs) {
                            (Int(l), Int(r)) => {
                                // not sure if i can capture the op match somehow
                                // this is ass
                                Ok(Bool(match expr.op {
                                    Lt => l < r,
                                    Le => l <= r,
                                    Gt => l > r,
                                    Ge => l >= r,
                                    _ => unreachable!("op in BinOp::Rel::Int"),
                                }))
                            }
                            (Float(l), Float(r)) => Ok(Bool(match expr.op {
                                Lt => l < r,
                                Le => l <= r,
                                Gt => l > r,
                                Ge => l >= r,
                                _ => unreachable!("op in BinOp::Rel::Float"),
                            })),
                            _ => eval_error("relation bin op on non number"),
                        }
                    }
                    // int -> int
                    Rem | BitXor | BitOr | BitAnd | Shl | Shr => match (lhs, rhs) {
                        (Int(l), Int(r)) => Ok(Int(match expr.op {
                            Rem => l % r,
                            BitXor => l ^ r,
                            BitOr => l | r,
                            BitAnd => l & r,
                            Shl => l << r,
                            Shr => r >> r,
                            _ => unreachable!(),
                        })),
                        _ => eval_error("bit op on non ints"),
                    },
                    // number/boolean -> boolean
                    Eq | Ne => match (lhs, rhs) {
                        (Int(l), Int(r)) => Ok(Bool(match expr.op {
                            Eq => l == r,
                            Ne => l != r,
                            _ => unreachable!(),
                        })),
                        (Float(l), Float(r)) => Ok(Bool(match expr.op {
                            Eq => l == r,
                            Ne => l != r,
                            _ => unreachable!(),
                        })),
                        (String(l), String(r)) => Ok(Bool(match expr.op {
                            Eq => l == r,
                            Ne => l != r,
                            _ => unreachable!(),
                        })),
                        (Bool(l), Bool(r)) => Ok(Bool(match expr.op {
                            Eq => l == r,
                            Ne => l != r,
                            _ => unreachable!(),
                        })),
                        _ => eval_error("bin op equality on separate types"),
                    },
                    // boolean -> boolean
                    And | Or => match (lhs, rhs) {
                        (Bool(l), Bool(r)) => Ok(Bool(match expr.op {
                            And => l && r,
                            Or => l || r,
                            _ => unreachable!(),
                        })),
                        _ => eval_error("bin op logical on non bools"),
                    },
                }
            }
            Unary(expr) => {
                let arg = expr.expr.eval()?;

                use UnaryOp::*;
                match expr.op {
                    Plus => match arg {
                        Int(v) => Ok(Int(if v < 0 { v * -1 } else { v })),
                        Float(v) => Ok(Float(if v < 0. { v * -1. } else { v })),
                        _ => eval_error("unary plus non number"),
                    },
                    Minus => match arg {
                        Int(v) => Ok(Int(v * -1)),
                        Float(v) => Ok(Float(v * -1.)),
                        _ => eval_error("unary minus non number"),
                    },
                    Not => match arg {
                        Bool(v) => Ok(Bool(!v)),
                        _ => eval_error("unary not non bool"),
                    },
                    BitNot => match arg {
                        Int(v) => Ok(Int(!v)),
                        _ => eval_error("unary bit not non int"),
                    },
                }
            }
            Conditional(expr) => {
                let test = expr.test.eval()?;
                if let Bool(cond) = test {
                    if cond {
                        Ok(expr.cons.eval()?)
                    } else {
                        Ok(expr.altr.eval()?)
                    }
                } else {
                    eval_error("conditional non bool test")
                }
            }
            // this sucks, can you really not impl copy if you have a string???
            //   theres gotta be a way
            Literal(val) => Ok(match val {
                Int(v) => Int(*v),
                Float(v) => Float(*v),
                Bool(v) => Bool(*v),
                String(s) => String(s.clone()),
            }),
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
            Literal(val) => {
                write!(f, "{}", val)
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

    pub fn run(&self) -> Result<Value, &str> {
        self.prog.eval()
    }
}

impl Display for Program {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", self.prog)
    }
}
