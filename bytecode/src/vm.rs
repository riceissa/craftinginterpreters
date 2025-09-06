use crate::chunk::{Chunk, OpCode};
use crate::value::{Value, print_value};

pub struct VM {
    chunk: Chunk,
    ip: *mut u8,
}

enum InterpretResult {
    Ok,
    CompileError,
    RuntimeError,
}

impl Drop for VM {
    fn drop(&mut self) {
    }
}

impl VM {
    pub fn init() {
    }

    pub fn interpret(&mut self, chunk: Chunk) -> InterpretResult {
        self.chunk = chunk;
        self.ip = self.chunk.code.values;
        unsafe {
            return self.run();
        }
    }

    unsafe fn run(&self) -> InterpretResult {
        loop {
            let instruction: u8 = *self.ip;
            *self.ip += 1;
            match instruction {
                x if x == OpCode::Constant as u8 => {
                    let index = *self.ip as usize;
                    let constant: Value = *self.chunk.constants.values.add(index);
                    *self.ip += 1;
                    print_value(&constant);
                    println!("");
                    // break;
                    return InterpretResult::Ok;
                },
                x if x == OpCode::Return as u8 => {
                    return InterpretResult::Ok;
                }
                _ => {
                    return InterpretResult::RuntimeError;
                }
            }
        }
    }
}
