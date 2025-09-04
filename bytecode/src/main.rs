mod chunk;
mod debug;
mod memory;
mod value;

use crate::chunk::{Chunk, OpCode};
use crate::debug::disassemble_chunk;
use crate::value::{Value};

fn main() {
    let mut chunk: Chunk = Chunk::new();

    let constant = chunk.add_constant(Value(1.2));
    chunk.write(OpCode::Constant as u8, 123);
    chunk.write(constant as u8, 123);

    chunk.write(OpCode::Return as u8, 123);

    disassemble_chunk(&mut chunk, "test chunk");
}
