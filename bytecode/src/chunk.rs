use std::ptr;

use crate::memory::{free_array, grow_array, grow_capacity};
use crate::value::{Value, ValueArray};

#[repr(u8)]
pub enum OpCode {
    Constant = 0,
    Return,
}

pub struct DynamicArray<T> {
    pub count: isize,
    pub capacity: isize,
    pub values: *mut T,
}

impl<T> DynamicArray<T> {
    pub fn init(array: &mut DynamicArray<T>) {
        array.count = 0;
        array.capacity = 0;
        array.values = ptr::null::<T>() as *mut T;
    }

    pub fn free(array: &mut DynamicArray<T>) {
        free_array::<T>(array.values, array.capacity);
        Self::init(array);
    }

    pub fn write(array: &mut DynamicArray<T>, value: T) {
        if array.capacity < array.count + 1 {
            let old_capacity = array.capacity;
            array.capacity = grow_capacity(old_capacity);
            array.values = grow_array::<T>(array.values, old_capacity, array.capacity);
        }

        unsafe {
            *array.values.add(array.count as usize) = value;
        }
        array.count += 1;
    }
}

pub struct Chunk {
    pub chunk: DynamicArray<u8>,
    pub lines: DynamicArray<i32>,
    pub constants: ValueArray,
}

impl Chunk {
    pub fn init(chunk: &mut Chunk) {
        DynamicArray::<u8>::init(&mut chunk.chunk);
        DynamicArray::<i32>::init(&mut chunk.lines);
        ValueArray::init(&mut chunk.constants);
    }

    pub fn free(chunk: &mut Chunk) {
        DynamicArray::<u8>::free(&mut chunk.chunk);
        DynamicArray::<i32>::free(&mut chunk.lines);
        ValueArray::free(&mut chunk.constants);
    }

    pub fn add_constant(chunk: &mut Chunk, value: Value) -> isize {
        ValueArray::write(&mut chunk.constants, value);
        return chunk.constants.count - 1;
    }

    pub fn write(chunk: &mut Chunk, byte: u8, line: i32) {
        DynamicArray::<u8>::write(&mut chunk.chunk, byte);
        DynamicArray::<i32>::write(&mut chunk.lines, line);
    }
}
