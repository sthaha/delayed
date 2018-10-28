# Delayed Methods


## usage


```go

fn := delayed.Call(500 * time.Millisecond, func() {
  fmt.Prinln("hello, world.")
})


fn.Call(100 * time.Millisecond, func(){
  fmt.Prinln("hello, world. This should execute")
})

```

### Cancel a delayed function

`fn.Cancel()` cancels the delayed execution and return `true` or `false` to
indicate if cancel was required/happened.
