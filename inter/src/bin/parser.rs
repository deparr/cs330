use inter::Program;
use std::{
    env, io,
    process::{Command, Stdio},
};

fn main() {
    let ast = if env::args().any(|e| e == "--exec") {
        let acorn = Command::new("acorn")
            .arg("--ecma2024")
            .stdin(Stdio::inherit())
            .output()
            .expect("exec acorn");

        String::from_utf8(acorn.stdout).unwrap()
    } else {
        io::read_to_string(io::stdin()).expect("reading from stdin")
    };

    let parser_output: serde_json::Value = serde_json::from_str(&ast).expect("parsing json input");

    // yeah this sucks
    let expr = parser_output
        .get("body")
        .unwrap()
        .get(0)
        .unwrap()
        .get("expression")
        .unwrap();

    let program = Program::new(expr);
    println!("{}", program);
}
