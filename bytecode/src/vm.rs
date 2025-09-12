use std::convert::TryInto;

use crate::chunk::{Chunk, OpCode};
use crate::debug::disassemble_instruction;
use crate::value::{print_value, Value};
use crate::compiler::compile;

const DEBUG_TRACE_EXECUTION: bool = true;

pub struct VM {
    pub chunk: Chunk,

    // Unlike the book, we'll store the instruction pointer relatively, as an
    // offset, mostly because I don't know how to do it as an actual pointer in
    // Rust and get it to actually compile.
    pub ip: usize,

    pub stack: Vec<Value>,
}

pub enum InterpretResult {
    Ok,
    CompileError,
    RuntimeError,
}

impl VM {
    pub fn init(&self) {}

    pub fn interpret(&mut self, source: &str) -> InterpretResult {
        compile(source);
        InterpretResult::Ok
    }

    fn read_byte(&mut self) -> u8 {
        let result = self.chunk.code[self.ip];
        self.ip += 1;
        result
    }

    fn read_constant(&mut self) -> Value {
        let index = self.read_byte();
        self.chunk.constants[index as usize]
    }

    fn run(&mut self) -> InterpretResult {
        macro_rules! binary_op {
            ($op:tt) => {
                {
                    let b = self.stack.pop().expect("Stack underflow: failed to pop 'b'");
                    let a = self.stack.pop().expect("Stack underflow: failed to pop 'a'");
                    self.stack.push(a $op b);
                }
            };
        }

        loop {
            if DEBUG_TRACE_EXECUTION {
                print!("          ");
                for slot in &self.stack {
                    print!("[ ");
                    print_value(&slot);
                    print!(" ]");
                }
                println!("");
                disassemble_instruction(&mut self.chunk, self.ip as i32);
            }
            let instruction: u8 = self.read_byte();
            match instruction.try_into() {
                Ok(OpCode::Constant) => {
                    let constant: Value = self.read_constant();
                    self.stack.push(constant);
                }
                Ok(OpCode::Add) => binary_op!(+),
                Ok(OpCode::Subtract) => binary_op!(-),
                Ok(OpCode::Multiply) => binary_op!(*),
                Ok(OpCode::Divide) => binary_op!(/),
                Ok(OpCode::Negate) => {
                    if let Some(value) = self.stack.pop() {
                        self.stack.push(-value);
                    }
                }
                Ok(OpCode::Return) => {
                    if let Some(value) = self.stack.pop() {
                        print_value(&value);
                        println!("");
                    }
                    return InterpretResult::Ok;
                }
                Err(_) => {
                    return InterpretResult::RuntimeError;
                }
            }
        }
    }
}
