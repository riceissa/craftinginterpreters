use std::convert::TryFrom;

use crate::value::{Value, ValueArray};

#[repr(u8)]
pub enum OpCode {
    Constant = 0,
    Add,
    Subtract,
    Multiply,
    Divide,
    Negate,
    Return,
}

impl TryFrom<u8> for OpCode {
    type Error = String;

    fn try_from(v: u8) -> Result<Self, Self::Error> {
        match v {
            x if x == OpCode::Constant as u8 => Ok(OpCode::Constant),
            x if x == OpCode::Add as u8 => Ok(OpCode::Add),
            x if x == OpCode::Subtract as u8 => Ok(OpCode::Subtract),
            x if x == OpCode::Multiply as u8 => Ok(OpCode::Multiply),
            x if x == OpCode::Divide as u8 => Ok(OpCode::Divide),
            x if x == OpCode::Negate as u8 => Ok(OpCode::Negate),
            x if x == OpCode::Return as u8 => Ok(OpCode::Return),
            _ => Err(format!("Unknown opcode {}", v)),
        }
    }
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
