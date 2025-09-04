use crate::chunk::{Chunk, OpCode};
use crate::value::print_value;

pub fn disassemble_chunk(chunk: &mut Chunk, name: &str) {
    println!("== {} ==", name);

    let mut offset: i32 = 0;
    while (offset as isize) < chunk.code.count {
        offset = disassemble_instruction(chunk, offset as i32);
    }
}

fn disassemble_instruction(chunk: &mut Chunk, offset: i32) -> i32 {
    print!("{:04} ", offset);
    unsafe {
        if offset > 0
            && *chunk.lines.values.add(offset as usize)
                == *chunk.lines.values.add((offset - 1) as usize)
        {
            print!("   | ");
        } else {
            print!("{:4} ", *chunk.lines.values.add(offset as usize))
        }
        let instruction: u8 = *chunk.code.values.add(offset as usize);
        match instruction {
            x if x == OpCode::Constant as u8 => constant_instruction("OP_CONSTANT", chunk, offset),
            x if x == OpCode::Return as u8 => simple_instruction("OP_RETURN", offset),
            _ => {
                println!("Unknown opcode {}", instruction);
                offset + 1
            }
        }
    }
}

fn constant_instruction(name: &str, chunk: &Chunk, offset: i32) -> i32 {
    unsafe {
        let constant: u8 = *chunk.code.values.add((offset + 1) as usize);
        print!("{:16} {:4} '", name, constant);
        print_value(&*chunk.constants.values.add(constant as usize));
    }
    println!("'");
    return offset + 2;
}

fn simple_instruction(name: &str, offset: i32) -> i32 {
    println!("{}", name);
    offset + 1
}
