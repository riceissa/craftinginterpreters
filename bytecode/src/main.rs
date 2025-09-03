use std::ptr;

mod memory;
mod chunk;
mod debug;

use crate::chunk::{Chunk, OpCode, init_chunk, write_chunk, free_chunk};
use crate::debug::{disassemble_chunk};

fn main() {
    let mut chunk: Chunk = Chunk{count: 0, capacity: 0, code: ptr::null::<u8>() as *mut u8};
    init_chunk(&mut chunk);
    write_chunk(&mut chunk, OpCode::OpReturn as u8);
    disassemble_chunk(&mut chunk, "test chunk");
    free_chunk(&mut chunk);
}
