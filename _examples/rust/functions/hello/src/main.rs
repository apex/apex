extern crate rust_apex;
#[macro_use]
extern crate serde_json;
extern crate failure;

use failure::{Compat, Error};
use serde_json::{to_value, Value};
use rust_apex::Context;

fn main() {
    rust_apex::run::<_, _, Compat<Error>, _>(|input: Value, c: Context| {
        Ok(json!({
            "name": to_value(&c).unwrap(),
            "age": input,
            "phones": [
                "+44 1234567",
                "+44 2345678"
            ]
        }))
    });
}
