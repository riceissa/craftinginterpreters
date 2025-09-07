mod chunk;
mod debug;
mod value;
mod vm;

use crate::chunk::{Chunk, OpCode};
use crate::debug::disassemble_chunk;
use crate::value::Value;
use crate::vm::{VM};

fn main() {
    let mut vm: VM = VM{chunk: Chunk::new(), ip: 0, stack: Vec::new()};
    vm.init();

    let mut chunk: Chunk = Chunk::new();

    let constant: isize = chunk.add_constant(Value(1.2));
    chunk.write(OpCode::Constant as u8, 123);
    chunk.write(constant as u8, 123);
    chunk.write(OpCode::Negate as u8, 123);

    chunk.write(OpCode::Return as u8, 123);

    disassemble_chunk(&mut chunk, "test chunk");

    vm.interpret(chunk);
}
