use std::{collections::HashMap, fmt::Display};

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
struct BindExpr {
    ident: String,
    bind: Expr,
}

impl BindExpr {
    fn new(expr: &serde_json::Value) -> Self {
        let ident = String::from(
            expr.get("id")
                .unwrap()
                .get("name")
                .unwrap()
                .as_str()
                .unwrap(),
        );

        let bind = Expr::new(expr.get("init").unwrap());

        BindExpr { ident, bind }
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

    fn extract_bool(self) -> Result<bool, &'static str> {
        if let Self::Bool(b) = self {
            Ok(b == true)
        } else {
            Err("non bool value in bool operator")
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

impl Clone for Value {
    fn clone(&self) -> Self {
        use Value::*;
        match self {
            String(s) => String(s.clone()),
            Bool(b) => Bool(*b),
            Int(i) => Int(*i),
            Float(f) => Float(*f),
        }
    }
}

#[derive(Debug)]
struct Environ {
    env: HashMap<String, Value>,
}

impl Environ {
    fn empty() -> Self {
        Environ {
            env: HashMap::new(),
        }
    }

    fn extend(&self, ident: String, val: Value) -> Self {
        let mut extended = self.env.clone();
        extended.insert(ident, val);
        Environ { env: extended }
    }

    fn lookup(&self, ident: &String) -> Option<&Value> {
        self.env.get(ident)
    }
}

#[derive(Debug)]
enum Expr {
    Binary(Box<BinaryExpr>),
    Unary(Box<UnaryExpr>),
    Conditional(Box<CondExpr>),
    Bind(Vec<BindExpr>),
    Ref(String),
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
            "Identifier" => Expr::Ref(expr.get("name").unwrap().as_str().unwrap().into()),
            "ExpressionStatement" => Expr::new(expr.get("expression").unwrap()),
            "AssignmentExpression" => todo!("Assignment expressions"),
            "VariableDeclaration" => {
                let mut binds = Vec::new();
                let decs = expr.get("declarations").unwrap().as_array().unwrap();
                for dec in decs {
                    binds.push(BindExpr::new(dec))
                }
                Expr::Bind(binds)
            }
            _ => todo!("reached end of expr type match arms: {}", expr_type),
        }
    }

    fn eval(&self, env: &Environ) -> Result<Value, &str> {
        use Expr::*;
        use Value::*;
        // I really hate this but don't really have time to find a better way
        match self {
            Binary(expr) => {
                use BinOp::*;
                match expr.op {
                    // number -> number
                    Add | Sub | Mul | Div => {
                        let lhs = expr.lhs.eval(env)?;
                        let rhs = expr.rhs.eval(env)?;
                        match (lhs, rhs) {
                            (Int(l), Int(r)) => match expr.op {
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
                            },
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
                    Lt | Le | Gt | Ge | Eq | Ne => {
                        let lhs = expr.lhs.eval(env)?;
                        let rhs = expr.rhs.eval(env)?;
                        match (lhs, rhs) {
                            (Int(l), Int(r)) => Ok(Bool(match expr.op {
                                Lt => l < r,
                                Le => l <= r,
                                Gt => l > r,
                                Ge => l >= r,
                                Eq => l == r,
                                Ne => l != r,
                                _ => unreachable!("op in BinOp::Rel::Int"),
                            })),
                            (Float(l), Float(r)) => Ok(Bool(match expr.op {
                                Lt => l < r,
                                Le => l <= r,
                                Gt => l > r,
                                Ge => l >= r,
                                Eq => l == r,
                                Ne => l != r,
                                _ => unreachable!("op in BinOp::Rel::Float"),
                            })),
                            _ => eval_error("relation bin op on non number or mixed numbers"),
                        }
                    }
                    // int -> int
                    Rem | BitXor | BitOr | BitAnd | Shl | Shr => {
                        let lhs = expr.lhs.eval(env)?;
                        let rhs = expr.rhs.eval(env)?;
                        match (lhs, rhs) {
                            (Int(l), Int(r)) => Ok(Int(match expr.op {
                                Rem => l % r,
                                BitXor => l ^ r,
                                BitOr => l | r,
                                BitAnd => l & r,
                                Shl => l << r,
                                Shr => l >> r,
                                _ => unreachable!(),
                            })),
                            _ => eval_error("bit op on non ints"),
                        }
                    }
                    // boolean -> boolean
                    Or => {
                        let lhs = expr.lhs.eval(env)?;
                        let lhs = lhs.extract_bool()?;
                        if lhs {
                            Ok(Bool(true))
                        } else {
                            let rhs = expr.rhs.eval(env)?;
                            let rhs = rhs.extract_bool()?;
                            if rhs {
                                return Ok(Bool(true));
                            }
                            Ok(Bool(false))
                        }
                    }
                    And => {
                        let lhs = expr.lhs.eval(env)?;
                        let lhs = lhs.extract_bool()?;
                        if !lhs {
                            Ok(Bool(false))
                        } else {
                            let rhs = expr.rhs.eval(env)?;
                            let rhs = rhs.extract_bool()?;
                            if rhs {
                                return Ok(Bool(true));
                            }
                            Ok(Bool(false))
                        }
                    }
                }
            }
            Unary(expr) => {
                let arg = expr.expr.eval(env)?;

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
                let test = expr.test.eval(env)?;
                if let Bool(cond) = test {
                    if cond {
                        Ok(expr.cons.eval(env)?)
                    } else {
                        Ok(expr.altr.eval(env)?)
                    }
                } else {
                    eval_error("conditional non bool test")
                }
            }
            // this sucks, can you really not impl copy if you have a string???
            //   theres gotta be a way, they're small strings but still
            Literal(val) => Ok(match val {
                Int(v) => Int(*v),
                Float(v) => Float(*v),
                Bool(v) => Bool(*v),
                String(s) => String(s.clone()),
            }),
            Bind(_) => {
                todo!("eval bind exprs in Expr::eval");
            }
            Ref(ident) => match env.lookup(&ident) {
                Some(val) => Ok(val.clone()),
                None => eval_error("unbound identifier"),
            },
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
            Ref(ident) => {
                write!(f, "{}", ident)
            }
            Bind(exprs) => {
                write!(f, "(let ").expect("unable to print bind");
                exprs.iter().for_each(|expr| {
                    write!(f, "{} = {}, ", expr.ident, expr.bind).expect("unable to print bind 2");
                });
                write!(f, ")")
            }
        }
    }
}

#[derive(Debug)]
pub struct Program {
    binds: Vec<Expr>,
    statement: Expr,
}

impl Program {
    pub fn new(body: &serde_json::Value) -> Self {
        let body = body.as_array().unwrap();
        let mut i = 0;
        // last item is always expr
        while i < body.len() {
            let expr = body[i].get("type").unwrap().as_str().unwrap();
            if expr != "VariableDeclaration" {
                break;
            }
        }

        // hack
        let statement = match expressions.pop() {
            Some(expr) => expr,
            None => Expr::Literal(Value::String(String::new())),
        };

        Program {
            binds: expressions,
            statement,
        }
    }

    pub fn run(&self) -> Result<Value, &str> {
        let mut env = Environ::empty();

        // eval binds
        for bind in &self.binds {
            if let Expr::Bind(bind_exprs) = bind {
                for expr in bind_exprs {
                    let val = expr.bind.eval(&env)?;
                    env = env.extend(expr.ident.clone(), val);
                }
            } else {
                unreachable!("non bind expr in Prog.binds");
            }
        }

        // eval final statement
        self.statement.eval(&env)
    }
}

impl Display for Program {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", self.statement)
    }
}
