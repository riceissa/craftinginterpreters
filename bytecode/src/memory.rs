use std::{mem, ptr, process};
use std::alloc::{realloc, dealloc, Layout};

pub fn grow_capacity(capacity: isize) -> isize {
    if capacity < 8 { 8 } else { capacity * 2 }
}

pub fn grow_array<T>(pointer: *mut T, old_count: isize, new_count: isize) -> *mut T {
    let t_size: isize = mem::size_of::<T>() as isize;
    return reallocate::<T>(pointer, t_size * old_count, t_size * new_count);
}

pub fn free_array<T>(pointer: *mut T, old_count: isize) -> *mut T {
    let t_size: isize = mem::size_of::<T>() as isize;
    reallocate::<T>(pointer, t_size * old_count, 0)
}

fn reallocate<T>(pointer: *mut T, old_size: isize, new_size: isize) -> *mut T {
    if new_size == 0 {
        let layout = Layout::from_size_align(old_size as usize, mem::align_of::<T>()).unwrap();
        unsafe {
            dealloc(pointer as *mut u8, layout);
        }
        return ptr::null::<T>() as *mut T;
    }

    let layout = Layout::from_size_align(new_size as usize, mem::align_of::<T>()).unwrap();
    unsafe {
        let result: *mut u8 = realloc(pointer as *mut u8, layout, new_size as usize);
        if result.is_null() {
            process::exit(1);
        }
        return result as *mut T;
    }
}
