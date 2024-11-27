# go-pointer

Originally forked from https://github.com/mattn/go-pointer and adapted to fit the needs of go-gst

## Usage

https://github.com/golang/proposal/blob/master/design/12416-cgo-pointers.md

In go 1.6, cgo argument can't be passed Go pointer.

```
var s string
C.pass_pointer(pointer.Save(&s))
v := *(pointer.Restore(C.get_from_pointer()).(*string))
```

## Installation

```
go get github.com/go-gst/go-pointer
```

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a mattn)
