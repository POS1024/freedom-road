# 1. io — 基本的 IO 接口
在 io 包中最重要的是两个接口：Reader 和 Writer 接口。只要满足这两个接口，它就可以使用 IO 包的功能。

## 1.1 Reader 接口
``` go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```
* `Read` 将 `len(p)` 个字节读取到 `p` 中。它返回读取的字节数以及任何遇到的错误。
* 即使 `Read` 返回的 `n < len(p)`，它也会在调用过程中占用 `len(p)` 个字节作为暂存空间。
* 若可读取的数据不到 `len(p)` 个字节，`Read` 会返回可用数据，而不是等待更多数据。
* 返回的错误注意 `io.EOF` 类型与其他类型。

## 1.2 Writer 接口
``` go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```
* `Write` 将 `len(p)` 个字节从 `p` 中写入到基本数据流中。
* 它返回从 `p` 中被写入的字节数以及任何遇到的引起写入提前停止的错误。
* 若 `Write` 返回的 `n < len(p)` ，它就必须返回一个 `非nil` 的错误。

## 1.3 Closer 接口
``` go
type Closer interface {
    Close() error
}
```
* 用于关闭数据流
* 一些需要手动关闭的资源最好实现 `Closer` 接口


# 2. ioutil — 方便的IO操作函数集
提供了一些常用、方便的IO操作函数。

## 2.1 ioutil.ReadAll 函数
``` go
func ReadAll(r io.Reader) ([]byte, error)
```
* 用来从 `io.Reader` 中一次读取所有数据。
* 该函数成功调用后会返回 `err == nil` 而不是 `err == EOF` 。

## 2.2 ioutil.ReadDir 函数
``` go
fileInfos, err := ioutil.ReadDir("")
if err == nil {
    for _,fileInfo := range fileInfos {
        # fileInfo fs.FileInfo
        if fileInfo.IsDir() {
            # DIR
        }else{
            # FILE
            fileName := fileInfo.Name()
            fmt.Println(fileName)
        }
    }
}
```
* 输出目录下的 `文件`（包含 `文件目录` ）。
* 遍历为 `fs.FileInfo` 类型，`IsDir` 判断是否是文件目录，`Name` 得到文件名。

## 2.3 ioutil.ReadFile 函数
``` go
func ReadFile(filename string) ([]byte, error)
```
* `ReadFile` 从 `filename` 指定的文件中读取数据并返回文件的内容。成功调用返回的 `err` 为 `nil` 而非 `EOF`。
* `ReadFile` 会先判断文件的大小，给 `bytes.Buffer` 一个预定义容量，避免额外分配内存。

## 2.4 ioutil.WriteFile 函数
``` go
func WriteFile(filename string, data []byte, perm os.FileMode) error
```
* `WriteFile` 将 `data` 写入 `filename` 文件中，当文件不存在时会根据 `perm` 指定的权限进行创建一个,文件存在时会先清空文件内容。


# 3. fmt - 格式化IO
占位符格式化转换、格式化输出、普通输出等。

## 3.1 Stringer 接口
``` go
	type Stringer interface {
	    String() string
	}
	type OwnStringer type {
	    Name string
	    Age int
	}
	func (o *OwnStringer) String() string {
		buffer := bytes.NewBufferString("name : ")
		buffer.WriteString(this.Name + ", ")
		buffer.WriteString("age : ")
		buffer.WriteString(strconv.Itoa(this.Age))
		buffer.WriteString(" years old.")
		return buffer.String()
	}
	o := &OwnStringer{"own", 18}
	fmt.Println(o)
```
* 某结构体如果实现了 `Stringer` 接口，`fmt` 包中的打印函数打印这个结构体实例时，将调用该结构体的 `String` 方法。


# 4. 待定