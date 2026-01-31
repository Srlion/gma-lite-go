# gma-lite-go

📦 Minimal Go library to read and write Garry's Mod Addon (.gma) archives.

## Installation

```sh
go get -u github.com/Srlion/gma-lite-go/gma
```

## Usage

### Write a GMA file

```go
package main

import (
    "os"
    "github.com/Srlion/gma-lite-go/gma"
)

func main() {
    builder := gma.NewBuilder("My Addon", 76561198261855442, 2)
    builder.SetAuthor("Srlion")
    builder.SetDescription("Example Garry's Mod addon")

    builder.FileFromBytes("lua/autorun/init.lua", []byte("print('Hello from gma-lite-go')"))
    builder.FileFromString("readme.txt", "This was packed using gma-lite-go!")

    f, err := os.Create("my_addon.gma")
    if err != nil {
        panic(err)
    }
    defer f.Close()

    if err := builder.WriteGMATo(f); err != nil {
        panic(err)
    }

    println("GMA written successfully")
}
```

### Read a GMA file

```go
package main

import (
    "bufio"
    "os"
    "github.com/Srlion/gma-lite-go/gma"
)

func main() {
    f, err := os.Open("my_addon.gma")
    if err != nil {
        panic(err)
    }
    defer f.Close()

    entries, err := gma.ReadGMA(bufio.NewReader(f))
    if err != nil {
        panic(err)
    }

    for _, e := range entries {
        println("File:", e.Name(), "Size:", e.Size())
        // To access file content: e.Content()
    }
}
```

## API

- **Entry:** struct representing an individual file in the archive.
    - `Name() string`
    - `Size() uint64`
    - `Content() []byte`
- **Builder:** used for creating GMA archives.
    - `NewBuilder(name string, steamid64 int64, size int) *Builder`
    - `SetAuthor(author string)`
    - `SetDescription(desc string)`
    - `FileFromBytes(name string, content []byte)`
    - `FileFromString(name string, content string)`
    - `WriteGMATo(w io.Writer) error`
- **Reading GMA:**
    - `ReadGMA(r *bufio.Reader) ([]Entry, error)`

## Format Constants

- `HEADER`: []byte("GMAD")
- `VERSION`: int8(3)

## License

MIT
