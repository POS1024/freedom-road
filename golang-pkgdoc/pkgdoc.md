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


# 4. bufio — 缓存IO
bufio 包实现了缓存IO。它包装了 io.Reader 和 io.Writer 对象，创建了另外的Reader和Writer对象，它们也实现了 io.Reader 和 io.Writer 接口，不过它们是有缓存的。

## 4.1 bufio.Reader 类型
``` go
    type Reader struct {
		buf          []byte		// 缓存
		rd           io.Reader	// 底层的io.Reader
		// r:从buf中读走的字节（偏移）；w:buf中填充内容的偏移；
		// w - r 是buf中可被读的长度（缓存数据的大小），也是Buffered()方法的返回值
		r, w         int
		err          error		// 读过程中遇到的错误
		lastByte     int		// 最后一次读到的字节（ReadByte/UnreadByte)
		lastRuneSize int		// 最后一次读到的Rune的大小 (ReadRune/UnreadRune)
	}
	
	func NewReader(rd io.Reader) *Reader {
		// 默认缓存大小：defaultBufSize=4096
		return NewReaderSize(rd, defaultBufSize)
	}
	
	func NewReaderSize(rd io.Reader, size int) *Reader {
		// 已经是bufio.Reader类型，且缓存大小不小于 size，则直接返回
		b, ok := rd.(*Reader)
		if ok && len(b.buf) >= size {
			return b
		}
		// 缓存大小不会小于 minReadBufferSize （16字节）
		if size < minReadBufferSize {
			size = minReadBufferSize
		}
		// 构造一个bufio.Reader实例
		return &Reader{
			buf:          make([]byte, size),
			rd:           rd,
			lastByte:     -1,
			lastRuneSize: -1,
		}
	}
```
* 自定义的Reader类型，提供缓冲。

## 4.2 ReadSlice、ReadBytes、ReadString 和 ReadLine 方法
``` go
    func (b *Reader) ReadSlice(delim byte) (line []byte, err error)
```
* `ReadSlice` 从输入中读取，直到遇到第一个 `界定符（delim）` 为止，返回一个指向缓存中字节的 `slice` ，在下次调用 `读操作（read）` 时，这些字节会无效。
* 如果 `ReadSlice` 在找到界定符之前遇到了 `error` ，它就会返回缓存中所有的数据和错误本身（经常是 `io.EOF` ）。
* 如果在 `找到界定符之前` 缓存已经满了，`ReadSlice` 会返回 `bufio.ErrBufferFull` 错误。
* 当且仅当返回的结果（line）`没有以界定符结束` 的时候，`ReadSlice` 返回 `err != nil` ，也就是说，如果 `ReadSlice` 返回的结果 `line` 不是以界定符 `delim` 结尾，那么返回的 `err` 也一定不等于 `nil`（可能是 `bufio.ErrBufferFull` 或 `io.EOF` ）。

``` go
    func (b *Reader) ReadBytes(delim byte) (line []byte, err error)
```
* `ReadBytes` 从输入中读取直到遇到 `界定符（delim）` 为止，返回的 `slice` 包含了从当前到界定符的内容 （`包括界定符`）。
* 如果 `ReadBytes` 在遇到界定符之前就捕获到一个错误，它会返回遇到错误之前已经读取的数据，和这个捕获到的错误（经常是 `io.EOF` ）。
* 如果 `ReadBytes` 返回的结果 `line` 不是以界定符 `delim` 结尾，那么返回的 `err` 也一定不等于 nil（可能是 `bufio.ErrBufferFull` 或 `io.EOF` ）。

``` go
    func (b *Reader) ReadString(delim byte) (line string, err error) {
		bytes, err := b.ReadBytes(delim)
		return string(bytes), err
	}
```
* 调用了 ReadBytes 方法，并将结果的 []byte 转为 string 类型。

``` go
    func (b *Reader) ReadLine() (line []byte, isPrefix bool, err error)
```
* ReadLine 尝试返回单独的行，不包括行尾的换行符。
* 如果一行大于缓存，isPrefix 会被设置为 true，同时返回该行的开始部分（等于缓存大小的部分）。
* 该行剩余的部分就会在下次调用的时候返回。当下次调用返回该行剩余部分时，isPrefix 将会是 false 。
* 跟 ReadSlice 一样，返回的 line 只是 buffer 的引用，在下次执行IO操作时，line 会无效。
* 返回值中，要么 line 不是 nil，要么 err 非 nil，两者不会同时非 nil。

## 4.3 Scanner 类型和方法
``` go
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
	    fmt.Println(scanner.Text()) // Println will add back the final '\n'
	}
	if err := scanner.Err(); err != nil {
	    fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
```
* 更容易的处理如按行读取输入序列或空格分隔单词等。
* `Split` 分词
* `bufio.ScanWords` 返回通过 `空格` 分词的单词。这里的 `空格` 是 `unicode.IsSpace()`，即包括：`'\t', '\n', '\v', '\f', '\r', ' ', U+0085 (NEL), U+00A0 (NBSP)`。
* `bufio.ScanBytes` 返回单个字节作为一个 `token`。
* `bufio.ScanRunes` 返回单个 `UTF-8` 编码的 `rune` 作为一个 `token` 。对于 `无效的 UTF-8 编码` 会解释为 `U+FFFD = "\xef\xbf\xbd"` 。
* `bufio.ScanLines` 返回一行文本，`不包括` 行尾的换行符。这里的换行包括了Windows下的 `"\r\n"` 和Unix下的 `"\n"` 。

## 4.4 Writer 类型和方法
``` go
    type Writer struct {
		err error		// 写过程中遇到的错误
		buf []byte		// 缓存
		n   int			// 当前缓存中的字节数
		wr  io.Writer	// 底层的 io.Writer 对象
	}
	func NewWriter(wr io.Writer) *Writer {
		// 默认缓存大小：defaultBufSize=4096
		return NewWriterSize(wr, defaultBufSize)
	}
	func NewWriterSize(wr io.Writer, size int) *Writer {
		// 已经是 bufio.Writer 类型，且缓存大小不小于 size，则直接返回
		b, ok := wr.(*Writer)
		if ok && len(b.buf) >= size {
			return b
		}
		if size <= 0 {
			size = defaultBufSize
		}
		return &Writer{
			buf: make([]byte, size),
			wr:  w,
		}
	}
```
* `bufio.Writer` 结构包装了一个 `io.Writer` 对象，提供缓存功能，同时实现了 `io.Writer` 接口。
* `Available` 方法获取缓存中还未使用的字节数（缓存大小 - 字段 n 的值）。
* `Buffered` 方法获取写入当前缓存中的字节数（字段 n 的值）。
* `Flush` 方法将缓存中的所有数据写入底层的 io.Writer 对象中。

``` go
    // 实现了 io.ReaderFrom 接口
	func (b *Writer) ReadFrom(r io.Reader) (n int64, err error)
	
	// 实现了 io.Writer 接口
	func (b *Writer) Write(p []byte) (nn int, err error)
	
	// 实现了 io.ByteWriter 接口
	func (b *Writer) WriteByte(c byte) error
	
	// io 中没有该方法的接口，它用于写入单个 Unicode 码点，返回写入的字节数（码点占用的字节），内部实现会根据当前 rune 的范围调用 WriteByte 或 WriteString
	func (b *Writer) WriteRune(r rune) (size int, err error)
	
	// 写入字符串，如果返回写入的字节数比 len(s) 小，返回的error会解释原因
	func (b *Writer) WriteString(s string) (int, error)
```
* 这些写方法在缓存满了时会调用 Flush 方法。
* 只要写的过程中遇到了错误，再次调用写操作会直接返回该错误。

## 4.5 ReadWriter 类型和实例化
``` go
	type ReadWriter struct {
		*Reader
		*Writer
	}
	func NewReadWriter(r *Reader, w *Writer) *ReadWriter {
		return &ReadWriter{r, w}
	}
```
* 可以使用 `bufio.Reader` 和 `bufio.Writer` 所有的函数。

# 5. 待定