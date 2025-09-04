mod chunk;
mod debug;
mod memory;
// mod value;

use crate::chunk::{Chunk, OpCode};
use crate::debug::disassemble_chunk;

fn main() {
    let mut chunk: Chunk = Chunk {
        count: 0,
        capacity: 0,
        values: 0 as *mut u8,
    };
    Chunk::init(&mut chunk);
    Chunk::write(&mut chunk, OpCode::OpReturn as u8);
    disassemble_chunk(&mut chunk, "test chunk");
    Chunk::free(&mut chunk);
}
