use std::ops::Neg;

#[derive(Debug, Copy, Clone)]
#[allow(dead_code)]
pub struct Value(pub f64);

impl Neg for Value {
    type Output = Self;

    fn neg(mut self) -> Self::Output {
        self.0 = -self.0;
        self
    }
}

pub type ValueArray = Vec<Value>;

pub fn print_value(value: &Value) {
    print!("{:?}", value.0);
}
