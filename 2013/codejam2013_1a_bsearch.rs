// using Mozilla Rust
// rust run prog.rs < prog.in > prog.out

use core::io::ReaderUtil;
use core::from_str::FromStr::from_str;

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

    let mut high = 1;
    let mut low = 0;

    // Double the ring size till we surpass total paint.
    while t > computeTotalPaint(r, high) { 
      low = high;
      high *= 2; 
    }

    let mut rings = 0;
    // Answer now bounded between low/high. Binary search.
    while low <= high {
      let mut mid = (high + low)/2;
      match (t, computeTotalPaint(r, mid)) {
        (x, y) if x > y => { low = mid+1; }
        (x, y) if x < y => { high = mid-1; }
        (_, _) => { rings = mid; break; }
      }
      rings = high;
    }

    io::println(fmt!("Case #%i: %i", c+1, rings));
  }
}

fn computeTotalPaint(radius: int, rings: int) -> int {
  rings * (2*radius - 1 + 2*rings)
}

