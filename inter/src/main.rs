use inter::Program;
use std::{
    env,
    process::{Command, Stdio}, io,
};

fn main() {
    let ast = if env::args().any(|e| e == "--pipe") {
        io::read_to_string(io::stdin()).expect("unable to read from stdin")
    } else {
        let acorn = Command::new("acorn")
            .arg("--ecma2024")
            .stdin(Stdio::inherit())
            .output()
            .expect("failed to exec acorn");

        String::from_utf8(acorn.stdout).unwrap()
    };

    let parser_output: serde_json::Value =
        serde_json::from_str(&ast).expect("unable to deser parser json");

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
