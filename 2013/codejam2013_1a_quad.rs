// using Mozilla Rust
// rust run prog.rs < prog.in > prog.out

extern mod std;

use std::bigint;
use core::io::ReaderUtil;
use core::from_str::FromStr::from_str;

// Newton's method square root.
fn bigroot(N: &bigint::BigInt) -> ~bigint::BigInt {
  // Any bigger than this is overflow.
  let mut x = ~bigint::BigInt::from_uint(10000000000000000000);
  let two = ~bigint::BigInt::from_uint(2);
  loop {
    let y = ~N.div(x).add(x).div(two);
    if y >= x {
      return x;
    }
    x = y;
  }
}

fn main() {
  // Lambda: convert string to integer/float; unsafe.
  let atoi = |num_str: &str| from_str::<int>(num_str).get();

  let reader = io::stdin();
  // Get number of cases.
  let T = atoi(reader.read_line());

  // Loop through cases.
  for int::range(0, T) |c| {
    // Read entire line into int array.
    let line1 = reader.read_line();
    let mut rt = ~[];
    for str::each_word(line1) |word| { rt.push(atoi(word)); }

    let r = rt[0];
    let t = rt[1];

    // Using the quadratic equation. a=2, b=(2r-1), c=-t

    // Coefficient `b`.
    let b: int = 2*r-1;
    // BigInt `b`.
    let bigb = ~bigint::BigInt::from_uint(b.to_uint());
    // b^2 - 4ac
    let bsquared_minus4ac = ~bigint::BigInt::from_uint(t.to_uint())
        .mul(~bigint::BigInt::from_uint(8)).add(~bigb.mul(bigb));
    // -b + sqrt(b^2 - 4ac) / 2 * a
    let x = (bigroot(bsquared_minus4ac).to_uint().to_int() - b) / 4;

    io::println(fmt!("Case #%i: %i", c+1, x));
  }
}

