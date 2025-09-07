#[derive(Debug, Copy, Clone)]
#[allow(dead_code)]
pub struct Value(pub f64);

pub type ValueArray = Vec<Value>;

pub fn print_value(value: &Value) {
    print!("{:?}", value.0);
}
