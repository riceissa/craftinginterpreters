use std::ptr;

use crate::memory::{grow_capacity, grow_array, free_array};

#[repr(u8)]
pub enum OpCode {
    OpReturn = 0,
}

pub struct Chunk {
    pub count: isize,
    pub capacity: isize,
    pub code: *mut u8,
}

pub fn init_chunk(chunk: &mut Chunk) {
    chunk.count = 0;
    chunk.capacity = 0;
    chunk.code = ptr::null::<u8>() as *mut u8;
}

pub fn free_chunk(chunk: &mut Chunk) {
    free_array::<u8>(chunk.code, chunk.capacity);
    init_chunk(chunk);
}

pub fn write_chunk(chunk: &mut Chunk, byte: u8) {
    if chunk.capacity < chunk.count + 1 {
        let old_capacity = chunk.capacity;
        chunk.capacity = grow_capacity(old_capacity);
        chunk.code = grow_array::<u8>(chunk.code, old_capacity, chunk.capacity);
    }

    unsafe {
        *chunk.code.add(chunk.count as usize) = byte;
    }
    chunk.count += 1;
}
