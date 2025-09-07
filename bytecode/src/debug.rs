use std::convert::TryInto;

use crate::chunk::{Chunk, OpCode};
use crate::value::print_value;

pub fn disassemble_chunk(chunk: &mut Chunk, name: &str) {
    println!("== {} ==", name);

    let mut offset: i32 = 0;
    while (offset as isize) < (chunk.code.len() as isize) {
        offset = disassemble_instruction(chunk, offset as i32);
    }
}

fn disassemble_instruction(chunk: &mut Chunk, offset: i32) -> i32 {
    print!("{:04} ", offset);
    if offset > 0 && chunk.lines[offset as usize] == chunk.lines[(offset - 1) as usize]
    {
        print!("   | ");
    } else {
        print!("{:4} ", chunk.lines[offset as usize])
    }
    let instruction: u8 = chunk.code[offset as usize];
    match instruction.try_into() {
        Ok(OpCode::Constant) => constant_instruction("OP_CONSTANT", chunk, offset),
        Ok(OpCode::Return) => simple_instruction("OP_RETURN", offset),
        Err(err) => {
            println!("{} {}", err, instruction);
            offset + 1
        }
    }
}

fn constant_instruction(name: &str, chunk: &Chunk, offset: i32) -> i32 {
    let constant: u8 = chunk.code[(offset + 1) as usize];
    print!("{:16} {:4} '", name, constant);
    print_value(&chunk.constants[constant as usize]);
    println!("'");
    return offset + 2;
}

fn simple_instruction(name: &str, offset: i32) -> i32 {
    println!("{}", name);
    offset + 1
}
