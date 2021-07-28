# gotraders

Playing with go and spacetraders api


## software thinnking

### lesson 0: variable declaration

From initially running the CLI you'll see that the printed information isn't that helpful.

```
git checkout 799d44815d9369f879158e46c535b31baca48c0e
a114383@C02XV0P9JHD4 gotraders % go run . status
Account info == {Credits:0 JoinedAt:2021-07-24T17:25:24.754Z ShipCount:0 StructureCount:0 Username:bru}% 
```

A command should only do one thing, and do it well.

If you're hoping to get account information then this should return a useful nicely formatted status box.  That display formatting should be seperate from the fucntionality to construct the object so that other functions can access it's members rather than re-format a string back into struct.

First worked on creating a function to return an Account struct, I ran into an issue that wasn't clear.  VSCode tells me `creds declared but unused`:

```golang
func GetAccount(token string, username string) (Credentials, error) {
	if token != "" || username != "" {
		// try to load credentials based on ~/.spacetravlers/credentials file
		creds, _ := GetUsernameAndToken()
	} else {
		// should we validate these are correct?
		creds := Credentials{
			Username: username,
			Token:    token,
		}
	}
    
    . . .
    return creds
}
```

Until I read [this note on variable scope](https://stackoverflow.com/a/21481424/5660197), the important point being the difference between `:=` and `=`.  After switching to declaring the variable ahead of time and using the `=` operator everything cleared up:
```golang
func GetAccount(token string, username string) (Credentials, error) {
    var creds Credentials

	if token != "" || username != "" {
		// try to load credentials based on ~/.spacetravlers/credentials file
		creds, _ = GetUsernameAndToken()
	} else {
		// should we validate these are correct?
		creds = Credentials{
			Username: username,
			Token:    token,
		}
	}
    
    . . .
    return creds
}
```


### lesson 1: trailling "%"


Trailing "%" found when using the CLI, looks like `fmt.Printf()` [does not pad a trailing new line](https://stackoverflow.com/a/59094048/5660197).

```
a114383@C02XV0P9JHD4 gotraders % go run . init -u brew    
Username & Token written to ~/.spacetravels/credentials% 
```

Fixed with 1162a134284fd963970f6659ce9e146b76de970a


### lesson 2: ovrloading struct string method

Allows a cleaner interface to the original problem, now our funnction `GetAccount()` returns an `Account` type which the caller can either user or pretty print by default using the default `fmt.Printf("%+v", GetAccount())`.


### Learning goals for week of 07/25/2021

1. the amount if if err != nil in the code base is too damn high
When to use `panic`: "the main thing is that it should be utilized in cases where you don't feel like your program can recover from the state you're in" -sajan

2. There has got to be a better way to keep track of all the json payloads

3. Anotheer todo, find out why POST request requires a bytes object for the third parameter

 I might have complicated things a bit, [this post](https://blog.logrocket.com/making-http-requests-in-go/) seems to suggest using ioutils.ReadAll() over bytes.  Also when sending query params its easiest to 

4. betterr understand reflect package
i'd like to do something to the equiv of python's `getaddr(obj, key) == value`

