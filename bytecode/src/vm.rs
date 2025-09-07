use crate::chunk::{Chunk, OpCode};
use crate::value::{Value, print_value};

pub struct VM {
    pub chunk: Chunk,
    pub ip: usize,  // Unlike the book, we'll store the instruction pointer relatively, as an offset
}

pub enum InterpretResult {
    Ok,
    // CompileError,
    RuntimeError,
}

impl VM {
    pub fn init() {
    }

    pub fn interpret(&mut self, chunk: Chunk) -> InterpretResult {
        self.chunk = chunk;
        self.ip = 0;
        unsafe {
            return self.run();
        }
    }

    unsafe fn run(&mut self) -> InterpretResult {
        loop {
            let instruction: u8 = self.ip as u8;
            self.ip += 1;
            match instruction {
                x if x == OpCode::Constant as u8 => {
                    let constant: &Value = &self.chunk.constants[self.ip];
                    self.ip += 1;
                    print_value(&constant);
                    println!("");
                    break;
                },
                x if x == OpCode::Return as u8 => {
                    return InterpretResult::Ok;
                }
                _ => {
                    return InterpretResult::RuntimeError;
                }
            }
        }
        return InterpretResult::Ok;
    }
}
