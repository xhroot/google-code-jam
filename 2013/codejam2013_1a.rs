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

    let mut paint_used = 0;
    let mut rings_drawn = 0;
    
    loop {
      let current_r = 2 * rings_drawn + r;
      let paint_needed = 2 * current_r + 1;
      if (paint_needed + paint_used > t) {
        break;
      }
      paint_used += paint_needed;
      rings_drawn += 1;
    }

    io::println(fmt!("Case #%i: %i", c+1, rings_drawn));
  }
}

