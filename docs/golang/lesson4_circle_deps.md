# 

Updated my struct to include some extra information that is useful to passaround.  That is, I noticed often when I fetched goods there was a need to also know what system it was coming from:

```
type Good struct {
	Symbol               string `json:"symbol,omitempty"`
    . . .
	System               string
}
```

Made sure this had the write info:
```
# debugging 
fmt.Printf("good == %s; system == %s\n", good.Symbol, system)

# output
good == FUEL; system == OE-UC-OB

# debugging
fmt.Printf("systems %+v", systemGoods)

# output
systems [{PricePerUnit:80 PurchasePricePerUnit:91 QuantityAvailable:23886 SellPricePerUnit:69 Spread:11 Symbol:BIOMETRIC_FIREARMS VolumePerUnit:1 System:}, ...]
```

Something weird going on here - `good.system` returns empty string:

```golang
var goodsList []goods.Good
for system := range systems {
    systemGoods, _ := goods.List(good, system)
    for _, good := range systemGoods {
        good.System = system
    }

    goodsList = append(goodsList, systemGoods...)
}
```


Some searching for modifications while iterating turned out this [SO post](https://stackoverflow.com/a/15947177/5660197):

Leading me to:
```golang
for i := 0; i < len(systemGoods); i++ {
    good := &systemGoods[i]
    good.System = system
}
```