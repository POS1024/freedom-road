# io — 基本的 IO 接口
在 io 包中最重要的是两个接口：Reader 和 Writer 接口。只要满足这两个接口，它就可以使用 IO 包的功能。

## Reader 接口
```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```
* Read 将 len(p) 个字节读取到 p 中。它返回读取的字节数以及任何遇到的错误。
* 即使 Read 返回的 n < len(p)，它也会在调用过程中占用 len(p) 个字节作为暂存空间。
* 若可读取的数据不到 len(p) 个字节，Read 会返回可用数据，而不是等待更多数据。
* 返回的错误注意io.EOF类型与其他类型。

## Writer 接口
``` go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```
* Write 将 len(p) 个字节从 p 中写入到基本数据流中。
* 它返回从 p 中被写入的字节数以及任何遇到的引起写入提前停止的错误。
* 若 Write 返回的 n < len(p)，它就必须返回一个 非nil 的错误。

## Closer 接口
``` go
type Closer interface {
    Close() error
}
```
* 用于关闭数据流
* 一些需要手动关闭的资源最好实现Closer接口