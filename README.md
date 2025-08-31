## GMA (only writer)

### Usage

```go
builder := gma.NewBuilder("test", 76561198261855442, 30)
builder.SetAuthor("Srlion")
builder.SetDescription("fuck")

builder.FileFromBytes("yes.lua", []byte("xd"))
builder.FileFromString("yes2.lua", "xd")

f, err := os.Create("test.gma")
if err != nil {
    panic(err)
}
defer f.Close()

err = builder.WriteGMATo(f)
if err != nil {
    panic(err)
}
```
