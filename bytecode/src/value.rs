use crate::chunk::DynamicArray;

#[derive(Debug)]
#[allow(dead_code)]
pub struct Value(pub f64);

pub type ValueArray = DynamicArray<Value>;

pub fn print_value(value: &Value) {
    print!("{:?}", value.0);
}
