mod chunk;
mod debug;
mod value;
mod vm;

use std::env;
use std::io;
use std::process::exit;
use std::fs;

use crate::chunk::{Chunk, OpCode};
use crate::debug::disassemble_chunk;
use crate::value::Value;
use crate::vm::VM;

fn repl() {
    let mut line = String::new();
    let stdin = io::stdin();
    loop {
        print!("> ");

        if let Ok(xs) = stdin.read_line(&mut line) {
            interpret(xs);
        } else {
            println!("");
            break;
        }
    }
}

fn run_file(path: &String) {
    let source = fs::read_to_string(path).expect(fmt!("Could not read file \"{}\".", ));
    let result: InterpretResult = interpret(source);

    if result == InterpretResult::CompileError { exit(65); }
    if result == InterpretResult::RuntimeError { exit(70); }
}

fn main() {
    let mut vm: VM = VM {
        chunk: Chunk::new(),
        ip: 0,
        stack: Vec::new(),
    };
    vm.init();

    let args: Vec<String> = env::args().collect();
    if args.len() == 1 {
        repl();
    } else if args.len() == 2 {
        run_file(args[1]);
    } else {
        eprintln!("Usage: clox [path]");
        exit(64);
    }
}
