use crate::chunk::Chunk;

pub fn disassemble_chunk(chunk: &mut Chunk, name: &str) {
    println!("== {} ==", name);

    let mut offset: i32 = 0;
    while (offset as isize) < chunk.count {
        offset = disassemble_instruction(chunk, offset as i32);
    }
}

fn disassemble_instruction(chunk: &mut Chunk, offset: i32) -> i32 {
    print!("{:04} ", offset);
    unsafe {
        let instruction: u8 = *chunk.values.add(offset as usize);
        match instruction {
            _ => simple_instruction("OP_RETURN", offset),
            // _ => {
            //     println!("Unknown opcode {}", instruction);
            //     offset + 1
            // }
        }
    }
}

fn simple_instruction(name: &str, offset: i32) -> i32 {
    println!("{}", name);
    offset + 1
}
