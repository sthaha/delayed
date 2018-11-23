# Delayed Methods


## usage


```go

fn := delayed.Call(500 * time.Millisecond, func() {
  fmt.Prinln("hello, world.")
})


fn.Reset(100 * time.Millisecond, func() {
  fmt.Prinln("hello, world. This should execute")
})

```

### Cancel a delayed function

`fn.Cancel()` cancels the delayed execution and return `true` or `false` to
indicate if cancel was required/happened.

### Reset delay and func

`fn.Reset(newDelay, func(){...})` cancels the delayed execution and starts a
new one with the new delay and the new `func`

### Reset delay

`fn.ResetDelay(newDelay)` cancels the delayed execution and starts a new one
with the new delay.


### Reset Func

`fn.ResetFunc(func(){ ...  })` cancels the delayed execution and starts a new
one to execute the new `func`

## Hacks
###  Enabled debug logs

use `-tags debug` to enable debug logs
e.g. `go test -race -tags debug -run Example`
