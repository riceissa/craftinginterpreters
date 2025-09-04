mod chunk;
mod debug;
mod memory;
mod value;

use crate::chunk::{Chunk, OpCode, DynamicArray};
use crate::debug::disassemble_chunk;
use crate::value::{ValueArray, Value};

fn main() {
    let mut chunk: Chunk = Chunk {
        chunk: DynamicArray::<u8>{
            count: 0,
            capacity: 0,
            values: 0 as *mut u8,
        },
        lines: DynamicArray::<i32>{
            count: 0,
            capacity: 0,
            values: 0 as *mut i32,
        },
        constants: ValueArray{
            count: 0,
            capacity: 0,
            values: 0 as *mut Value,
        },
    };
    Chunk::init(&mut chunk);

    let constant = Chunk::add_constant(&mut chunk, Value(1.2));
    Chunk::write(&mut chunk, OpCode::Constant as u8, 123);
    Chunk::write(&mut chunk, constant as u8, 123);

    Chunk::write(&mut chunk, OpCode::Return as u8, 123);

    disassemble_chunk(&mut chunk, "test chunk");
    Chunk::free(&mut chunk);
}
