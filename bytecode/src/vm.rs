use std::convert::TryInto;

use crate::chunk::{Chunk, OpCode};
use crate::value::{Value, print_value};
use crate::debug::disassemble_instruction;

const DEBUG_TRACE_EXECUTION: bool = true;

pub struct VM {
    pub chunk: Chunk,
    pub ip: usize,  // Unlike the book, we'll store the instruction pointer
                    // relatively, as an offset, mostly because I don't know
                    // how to do it as an actual pointer in Rust and get it to
                    // actually compile.
    pub stack: Vec<Value>,
}

pub enum InterpretResult {
    Ok,
    // CompileError,
    RuntimeError,
}

impl VM {
    pub fn init(&self) {
    }

    pub fn interpret(&mut self, chunk: Chunk) -> InterpretResult {
        self.chunk = chunk;
        self.ip = 0;  // Point to the start of the chunk.code list
        return self.run();
    }

    fn read_byte(&mut self) -> u8 {
        let result = self.chunk.code[self.ip];
        self.ip +=1;
        result
    }

    fn read_constant(&mut self) -> Value {
        let index = self.read_byte();
        self.chunk.constants[index as usize]
    }

    fn run(&mut self) -> InterpretResult {
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
                },
                Ok(OpCode::Negate) => {
                    if let Some(value) = self.stack.pop() {
                        self.stack.push(-value);
                    }
                },
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
