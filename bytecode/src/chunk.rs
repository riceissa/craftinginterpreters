use std::ptr;

use crate::memory::{free_array, grow_array, grow_capacity};

#[repr(u8)]
pub enum OpCode {
    OpReturn = 0,
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

pub type Chunk = DynamicArray<u8>;
