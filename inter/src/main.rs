use std::io;
use std::process::{Command, Stdio};

fn main() {
    // let acorn = Command::new("acorn")
    //     .arg("--ecma2024")
    //     .stdin(Stdio::inherit())
    //     .output()
    //     .expect("failed to exec acorn");

    let parser_output =
        io::read_to_string(io::stdin()).expect("unable to read parser ouput on stdin");
    let parser_output: serde_json::Value =
        serde_json::from_str(&parser_output).expect("unable to deser parser json");
    println!("{:?}", parser_output);
    let body = parser_output.get("body");
    println!("{:?}", body);
}
