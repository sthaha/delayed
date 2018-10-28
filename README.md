# Delayed Methods

```go

fn := &delayed.Fn{}

fn.Call(500 * time.Millisecond, func() {
  fmt.Prinln("hello, world.")
})


```

### Cancel a delayed function

```go
fn.Cancel() // return true | false to indicate if cancel was required

```
