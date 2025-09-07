use crate::value::{Value, ValueArray};

#[repr(u8)]
pub enum OpCode {
    Constant = 0,
    Return,
}

pub struct Chunk {
    pub code: Vec<u8>,
    pub lines: Vec<i32>,
    pub constants: ValueArray,
}

impl Chunk {
    pub fn new() -> Self {
        Chunk {
            code: Vec::new(),
            lines: Vec::new(),
            constants: ValueArray::new(),
        }
    }

    pub fn add_constant(&mut self, value: Value) -> isize {
        self.constants.push(value);
        return (self.constants.len() as isize) - 1;
    }

    pub fn write(&mut self, byte: u8, line: i32) {
        self.code.push(byte);
        self.lines.push(line);
    }
}
