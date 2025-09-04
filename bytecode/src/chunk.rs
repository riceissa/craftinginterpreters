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

impl<T> Drop for DynamicArray<T> {
    fn drop(&mut self) {
        free_array::<T>(self.values, self.capacity);
    }
}

impl<T> DynamicArray<T> {
    pub fn new() -> Self {
        DynamicArray::<T> {
            count: 0,
            capacity: 0,
            values: ptr::null::<T>() as *mut T,
        }
    }

    pub fn write(&mut self, value: T) {
        if self.capacity < self.count + 1 {
            let old_capacity = self.capacity;
            self.capacity = grow_capacity(old_capacity);
            self.values = grow_array::<T>(self.values, old_capacity, self.capacity);
        }

        unsafe {
            *self.values.add(self.count as usize) = value;
        }
        self.count += 1;
    }
}

pub struct Chunk {
    pub code: DynamicArray<u8>,
    pub lines: DynamicArray<i32>,
    pub constants: ValueArray,
}

impl Chunk {
    pub fn new() -> Self {
        Chunk {
            code: DynamicArray::<u8>::new(),
            lines: DynamicArray::<i32>::new(),
            constants: ValueArray::new(),
        }
    }

    pub fn add_constant(&mut self, value: Value) -> isize {
        self.constants.write(value);
        return self.constants.count - 1;
    }

    pub fn write(&mut self, byte: u8, line: i32) {
        self.code.write(byte);
        self.lines.write(line);
    }
}
