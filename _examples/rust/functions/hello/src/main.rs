extern crate rust_apex;
extern crate serde_json;

use std::error::Error;
use std::collections::BTreeMap;
use std::fmt::{Display, Formatter, Error as FmtError};

use serde_json::{Value, to_value};

#[derive(Debug)]
struct DummyError;

impl Display for DummyError {
    fn fmt(&self, f: &mut Formatter) -> Result<(), FmtError> {
        write!(f, "{:?}", self)
    }
}

impl Error for DummyError {
    fn description(&self) -> &str {
        "dummy"
    }
}

fn main() {
    rust_apex::run::<_, _, DummyError, _>(|input: Value, c: rust_apex::Context| {
        let mut bt = BTreeMap::new();
        bt.insert("c", to_value(&c).unwrap());
        bt.insert("i", input);
        Ok(bt)
    });
}
