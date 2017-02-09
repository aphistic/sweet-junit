# sweet-junit #
A plugin for the Sweet testing framework to output JUnit files for test results

## Usage ##

Using this plugin with [sweet](https://www.github.com/aphistic/sweet) is pretty
straightforward:

```go
func Test(t *testing.T) {
    sweet.T(func(s *sweet.S) {
        s.RegisterPlugin(junit.NewPlugin())

        s.RunSuite(t, &mySuite{})
    })
}

```

Once the plugin is registered with sweet, you can specify the file to write the
output to by passing the `-junit.output` flag when running `go test`, such as:

```bash
$ go test -junit.output junit.xml
```
