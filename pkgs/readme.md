Still trying to better understand when it is best to break out a collection of functions into their own package.  

The difficulty that comes with testing due to specific usernames makes me think a better project structure could be used here?

to be continuted...


## errors

```
$> go test
cli.go:8:2: found packages github (account.go) and account (account_test.go) in /Users/a114383/gotraders/pkgs/account
...
```

Solution: https://stackoverflow.com/questions/6997524/the-declared-package-does-not-match-the-expected-package/25220460

Turns out by default you need to CD into each directory to run a test else the compiled package script and package script_test will run into conflicts when it gets run.

```
$> go test ./...  # works
```
