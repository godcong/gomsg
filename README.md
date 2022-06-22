# GOMSG #

## About ##
gomsg is an asynchronous notification processing implemented by gorutine

## How to use ##

- Create the monitor
```go
    m := New[datatype](size)
```

- Add waiting
```go
    m.WaitFor(key,handlefunctions)
```

- Send data
```go
    m.Send(key,value)
```

## Example ##
```go
type data struct {
    s string
    i int
}
m := New[*data](1024)

m.WaitFor("test", func(key string, val *data) {
    fmt.Println("called", key, val.s, val.i)
})
data := &data{
    s: "str",
    i: 1,
}
m.Send("test",data)
```