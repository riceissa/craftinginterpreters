#[derive(Debug)]
#[allow(dead_code)]
pub struct Value(pub f64);

// impl Clone for Value {
//     fn clone(&self) -> Self {
//         Value(self.0)
//     }
// }

pub type ValueArray = Vec<Value>;

pub fn print_value(value: &Value) {
    print!("{:?}", value.0);
}
