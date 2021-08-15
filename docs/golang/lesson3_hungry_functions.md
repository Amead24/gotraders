## Lesson 3: Watchout for Hungry Functions

I spent hours trying to debug why my http response wasn't properly getting marshalled into a struct, instead everything was returning as default values when curling he endpoint was returning real values:

```zsh
gotraders % go run . goods buy -s ckrmxi90r925411ds6hoveg0oz -g FUEL -q 1
Printing goodsBuy struct:
{Credits:0 Order:{Good: PricePerUnit:0 Quantity:0 Total:0} Ship:{Id: Location: X:0 Y:0 Cargo:[] SpaceAvailable:0 Type: Class: MaxCargo:0 LoadingSpeed:0 Speed:0 Manufacturer: Plating:0 Weapons:0}}
```

When you run out of the obvious solutions you start grasping for wild ones, and once you're tired of bashing your head against the keyboard you start looking for ways to rule out pieces of code to know what to focus your attention on.  You know, what I should have been doing from the start.

First thing I wanted to do was stop switching between curl and the cli tool.  Searching around the interwebs turned up this little nugget on converting the http response body into a string:

```golang
bodyBytes, err := ioutil.ReadAll(resp.Body)
if err != nil { return err }
fmt.Println(string(bodyBytes))
```

This way I could see what the payload was while iterating which helped me narrow down the fact that one of my struct fields was named wrong.  Double, triple, quadruple checked... This still wasn't fixing the issue:

```zsh
gotraders % go run . goods buy -s ckrmxi90r925411ds6hoveg0oz -g FUEL -q 1
Response body contents:
{"credits":136460,"order":{"good":"FUEL","quantity":1,"pricePerUnit":3,"total":3},"ship":{"id":"ckrmxi90r925411ds6hoveg0oz","location":"OE-PM-TR","x":14,"y":18,"cargo":[{"good":"FUEL","quantity":5,"totalVolume":5}],"spaceAvailable":45,"type":"JW-MK-I","class":"MK-I","maxCargo":50,"loadingSpeed":25,"speed":1,"manufacturer":"Jackshaw","plating":5,"weapons":5}}

Printing goodsBuy struct:
{Credits:0 Order:{Good: PricePerUnit:0 Quantity:0 Total:0} Ship:{Id: Location: X:0 Y:0 Cargo:[] SpaceAvailable:0 Type: Class: MaxCargo:0 LoadingSpeed:0 Speed:0 Manufacturer: Plating:0 Weapons:0}}
```

Eventually I landed on some helpful code for debugging json marshalling.  Tying in nicely with my goals for learning better error handling let's break this down a bit:

```golang
//https://stackoverflow.com/a/42624472/5660197
package main

import (
    "bytes"
    "fmt"
    "encoding/json"
)

func main() {
    buff := bytes.NewBufferString("{\"bar\": -123}")
    decoder := json.NewDecoder(buff)

    var foo struct{
        Bar uint32
    }
    if err := decoder.Decode(&foo); err != nil {
        if terr, ok := err.(*json.UnmarshalTypeError); ok {
                fmt.Printf("Failed to unmarshal field %s \n", terr.Field)
        } else {
            fmt.Printf("Not an unmarhsalling problem: %s\n", err)
        }
    } else {
        fmt.Println(foo.Bar)
    }
}
```

Still no luck though, when I slapped this into my code the terminal spit back `Not an unmarhsalling problem: EOF`.  Putting a small subset into [the go playground](https://play.golang.org/p/io7XGUMJM7p) returned the healthy, beautiful, and happy struct I always wanted.  I threw my hands up and went to bed.  The next morning comparing the playground versus my code the different pieces were:
1. The making of an http request
2. `defer resp.Body.Close()`
3. Printing of the string before the marshalling

Ruling these out, `curl` had me pretty confident it wasn't (1). I went and re-read the docs for `defer` and `resp.Body.Close()` but nothing seemed to jump out at me as the culprit.  The last seemed pretty inplasable, how could printing a value break something?

If you've ever played with Rust the concept of [ownership](https://doc.rust-lang.org/book/ch04-01-what-is-ownership.html) can be a hard adjustment at first.  The dea being that once you assign a variable a value, the original variable no longer holds a referance to that value.

Searching for `ioutils modifies request body` lead me right to [the answer](https://stackoverflow.com/questions/43021058/golang-read-request-body) I was looking for.  Turns out that there was really a bug with the the struct naming but by adding those logging statements I re-introduced a bug and conflated the two issues.  Does Go have ownership then?

No... But what I didn't realize (and as great as the Go integration to VS code is unfortunatetly there's no docs on this popping up while writing the code) that [bytes.Buffer consumes the bytes that are read](https://golang.org/src/bytes/buffer.go#L292).  Removing the first read from my code returned the result I was looking for:

```zsh
gotraders % go run . goods buy -s ckrmxi90r925411ds6hoveg0oz -g FUEL -q 1
Printing goodsBuy struct:
{Credits:136457 Order:{Good:FUEL PricePerUnit:3 Quantity:1 Total:3} Ship:{Id:ckrmxi90r925411ds6hoveg0oz Location:OE-PM-TR X:14 Y:18 Cargo:[{Good:FUEL Quantity:6 TotalVolume:6}] SpaceAvailable:44 Type:JW-MK-I Class:MK-I MaxCargo:50 LoadingSpeed:25 Speed:1 Manufacturer:Jackshaw Plating:5 Weapons:5}}
```

The moral of the story?  It feels wrong to say approach all functions with caution!  Instead coming from Python it's important to remind myself that referances and pointers are alot more prevelant here.  Slowing down and paying more attention to the types and memory will help reduce these bugs in the first place.

....And that I suck at writing conclusions.
