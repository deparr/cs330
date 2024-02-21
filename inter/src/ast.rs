use std::{collections::HashMap, fmt::Display, rc::Rc};

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
struct BinaryExpr<'a, 'b> {
    op: BinOp,
    lhs: Expr<'a, 'b>,
    rhs: Expr<'a, 'b>,
}

impl BinaryExpr<'_, '_> {
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
struct UnaryExpr<'a, 'b> {
    op: UnaryOp,
    expr: Expr<'a, 'b>
}

impl UnaryExpr<'_, '_> {
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
struct CondExpr<'a, 'b> {
    test: Expr<'a, 'b>,
    cons: Expr<'a, 'b>,
    altr: Expr<'a, 'b>,
}

impl CondExpr<'_> {
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
struct BindExpr<'a, 'b> {
    binds: Vec<(String, Expr<'a, 'b>)>,
    body: Expr<'a, 'b>,
}

// TODO: don't particularly like this
impl BindExpr<'_, '_> {
    fn new(expr: &serde_json::Value, rest: &[serde_json::Value]) -> Self {
        let mut binds = Vec::new();

        let declarations = expr.get("declarations").unwrap().as_array().unwrap();
        for dec in declarations {
            let ident = dec
                .get("id")
                .unwrap()
                .get("name")
                .unwrap()
                .as_str()
                .unwrap();
            let init = dec.get("init").unwrap();
            let init = Expr::new(init);

            binds.push((ident.to_owned(), init));
        }

        let body = match rest.first() {
            Some(expr) => match expr.get("type").unwrap().as_str().unwrap() {
                "VariableDeclaration" => Expr::Bind(Box::new(BindExpr::new(expr, &rest[1..]))),
                _ => Expr::new(expr),
            },
            None => unimplemented!("Err: BindExpr with no body expr"),
        };

        BindExpr { binds, body }
    }
}

#[derive(Debug)]
struct FnExpr<'a, 'b> {
    arg: String,
    body: Expr<'a, 'b>,
}
impl FnExpr<'_, '_> {
    fn new(expr: &serde_json::Value) -> FnExpr {
        let arg = expr
            .get("params")
            .unwrap()
            .as_array()
            .unwrap()
            .first()
            .unwrap()
            .as_str()
            .unwrap();

        let body = Expr::new(expr.get("body").unwrap());

        FnExpr {
            arg: arg.into(),
            body,
        }
    }
}

#[derive(Debug)]
pub struct FnValue<'a, 'b> {
    env: Environ<'b>,
    arg: String,
    body: Rc<&'a Expr<'a>>,
}

#[derive(Debug)]
pub enum Value<'a, 'b> {
    String(String),
    Bool(bool),
    Int(i64),
    Float(f64),
    Fn(FnValue<'a, 'b>),
    Unit,
}

fn eval_error(msg: &str) -> Result<Value, &str> {
    Err(msg)
}

impl Value<'_, '_> {
    fn from_json(expr: &serde_json::Value) -> Self {
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

impl Display for Value<'_, '_> {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        use Value::*;
        match self {
            Int(v) => write!(f, "(value (number {}))", v),
            Float(v) => write!(f, "(value (number {}))", v),
            Bool(v) => write!(f, "(value (boolean {}))", v),
            String(v) => write!(f, "(value (string {}))", v),
            Fn(_) => write!(f, "(value (function))"),
            Unit => write!(f, "(value ())"),
        }
    }
}

// TODO or this
impl Clone for Value<'_, '_> {
    fn clone(&self) -> Self {
        use Value::*;
        match self {
            String(s) => String(s.clone()),
            Bool(b) => Bool(*b),
            Int(i) => Int(*i),
            Float(f) => Float(*f),
            Unit => Unit,
            Fn(_) => unimplemented!("Literal(Value::Func)::clone()"),
        }
    }
}

#[derive(Debug, Clone)]
struct Environ<'a, 'b> {
    env: HashMap<String, Value<'a, 'b>>,
}

impl Environ<'_, '_> {
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
enum Expr<'a, 'b> {
    Binary(Box<BinaryExpr<'a>>),
    Unary(Box<UnaryExpr<'a>>),
    Conditional(Box<CondExpr<'a>>),
    Fn(Box<FnExpr<'a>>),
    Bind(Box<BindExpr<'a>>),
    Ref(String),
    Literal(Value<'a, 'b>),
}

impl Expr<'_, '_> {
    fn new(expr: &serde_json::Value) -> Self {
        let expr_type = match expr.get("type") {
            Some(t) => t.as_str().unwrap(),
            None => unimplemented!("none expr type should be unreachable"),
        };

        match expr_type {
            "BinaryExpression" | "LogicalExpression" => {
                Expr::Binary(Box::new(BinaryExpr::new(expr)))
            }
            "BlockStatement" => {
                match Expr::from_body(expr.get("body").unwrap().as_array().unwrap()) {
                    Some(expr) => expr,
                    None => Expr::Literal(Value::Unit),
                }
            }
            "FunctionExpression" => Expr::Fn(Box::new(FnExpr::new(expr))),
            "UnaryExpression" => Expr::Unary(Box::new(UnaryExpr::new(expr))),
            "ConditionalExpression" => Expr::Conditional(Box::new(CondExpr::new(expr))),
            "Literal" => Expr::Literal(Value::from_json(expr)),
            "Identifier" => Expr::Ref(expr.get("name").unwrap().as_str().unwrap().into()),
            "ExpressionStatement" => Expr::new(expr.get("expression").unwrap()),
            "ReturnStatement" => Expr::new(expr.get("argument").unwrap()),
            "AssignmentExpression" => todo!("Expr::new AssignmentExpression"),
            _ => todo!("reached end of expr type match arms: {}", expr_type),
        }
    }

    fn from_body(body: &[serde_json::Value]) -> Option<Self> {
        if body.len() < 1 {
            return None;
        }
        let first = &body[0];
        let statement = match first.get("type").unwrap().as_str().unwrap() {
            "VariableDeclaration" => Expr::Bind(Box::new(BindExpr::new(first, &body[1..]))),
            "ExpressionStatement" | "ReturnStatement" => Expr::new(first),
            _ => todo!("unknown expression type in program::new"),
        };

        Some(statement)
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
                                _ => unimplemented!("op in BinOp::Int"),
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
                                _ => unimplemented!("op in BinOp::Int"),
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
                                _ => unimplemented!("op in BinOp::Rel::Int"),
                            })),
                            (Float(l), Float(r)) => Ok(Bool(match expr.op {
                                Lt => l < r,
                                Le => l <= r,
                                Gt => l > r,
                                Ge => l >= r,
                                Eq => l == r,
                                Ne => l != r,
                                _ => unimplemented!("op in BinOp::Rel::Float"),
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
                                _ => unimplemented!(),
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
            Expr::Fn(expr) => {
                Ok(Value::Fn(FnValue { arg: expr.arg.clone(), body: Rc::new(&expr.body), env: env.clone()  } ))
            }
            Bind(expr) => {
                let mut new_env = env.clone();
                for (ident, bind_expr) in &expr.binds {
                    let bound_val = bind_expr.eval(&new_env)?;
                    new_env = new_env.extend(ident.clone(), bound_val);
                }

                expr.body.eval(&new_env)
            }
            Ref(ident) => match env.lookup(&ident) {
                Some(val) => Ok(val.clone()),
                None => eval_error("unbound identifier"),
            },
            Literal(val) => Ok(match val {
                Int(v) => Int(*v),
                Float(v) => Float(*v),
                Bool(v) => Bool(*v),
                String(s) => String(s.clone()),
                Unit => Unit,
                Value::Fn(_) => unimplemented!("Expr::Eval::Literal(Value::Func)"),
            }),
        }
    }
}

impl Display for Expr<'_, '_> {
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
            Fn(expr) => {
                write!(f, "(fn ({}) {})", expr.arg, expr.body)
            }
            Literal(val) => {
                write!(f, "{}", val)
            }
            Ref(ident) => {
                write!(f, "{}", ident)
            }
            Bind(expr) => {
                write!(f, "(let ").expect("unable to print bind");
                expr.binds.iter().for_each(|expr| {
                    write!(f, "{} = {}, ", expr.0, expr.1).expect("unable to print bind 2");
                });
                write!(f, ")")
            }
        }
    }
}

#[derive(Debug)]
pub struct Program<'a, 'b> {
    statement: Expr<'a, 'b>,
}

impl Program<'_, '_> {
    pub fn new(body: &serde_json::Value) -> Self {
        let body = body.as_array().unwrap();

        let statement = match Expr::from_body(body) {
            Some(statement) => statement,
            None => Expr::Literal(Value::Unit),
        };

        Program { statement }
    }

    pub fn run(&self) -> Result<Value, &str> {
        let env = Environ::empty();
        self.statement.eval(&env)
    }
}

impl Display for Program<'_, '_> {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", self.statement)
    }
}
