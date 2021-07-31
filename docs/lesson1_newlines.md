### lesson 1: trailling "%"


Trailing "%" found when using the CLI, looks like `fmt.Printf()` [does not pad a trailing new line](https://stackoverflow.com/a/59094048/5660197).

```
a114383@C02XV0P9JHD4 gotraders % go run . init -u brew    
Username & Token written to ~/.spacetravels/credentials% 
```

To fix this you'll want to append the "\n" to any `fmt.Printf` and if you do not need formatted strings you can use the `fmt.Println` to do this automatically for  you.


Fixed with 1162a134284fd963970f6659ce9e146b76de970a
