# 极客时间 Go训练营作业

## 作业
1.总结几种 socket 粘包的解包方式：fix length/delimiter based/length field based frame decoder。尝试举例其应用。
## 答：
1. fix length based 每次发送固定长度的数据，并且不超过缓冲区，接收方每次按照固定长度去读取数据，会出现粘包、半包的问题
```
// 服务端
func server_handle(c net.Conn) {
    defer c.Close()
    for {
		var buf = make([]byte, 65536) 
		log.Println("start to read from conn")
		n, err := c.Read(buf)
		if err != nil {
			log.Printf("conn read %d bytes,  error: %s", n, err)
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				continue
			}
		}
		log.Printf("read %d bytes, content is %s\n", n, string(buf[:n]))
	}
}

// 客户端
func client() {
	conn, _ := net.Dial("tcp", ":8888")
	defer conn.Close()

	data := make([]byte, 65536)
	var total int
	for i := 0; i < 10; i++ {
		n, err := conn.Write(data)
		if err != nil {
			total += n
			log.Printf("write %d bytes, error:%s\n", n, err)
			break
		}
		total += n
		log.Printf("write %d bytes this time, %d bytes in total\n", n, total)
	}

	log.Printf("write %d bytes in total\n", total)
}
```

2. delimiter based 发送在数据包加入特殊分隔符，标识数据包开始或结束。如果数据量过大，查找定界符会消耗一些性能，且内容和定界符可能混乱。
```
// 服务端
func server_handle(c net.Conn) {
    ...
	n, err := c.ReadSlice('\n')
    ...
}

// 客户端
func client() {
    ...
    data + "..........\n"
	conn.Write(data)
    ...
}
```

3. length field based 在TCP协议头里面写入每次发送请求的长度。 客户端在协议头里面带入数据长度，服务器在接收到请求后，根据协议头里面的数据长度来决定接受多少数据

```
func (p *myFrameCodec) Encode(w io.Writer, framePayload FramePayload) error {
	var f = framePayload
	var totalLen int32 = int32(len(framePayload)) + 4

	err := binary.Write(w, binary.BigEndian, &totalLen)
	if err != nil {
		return err
	}

	n, err := w.Write([]byte(f)) // write the frame payload to outbound stream
	if err != nil {
		return err
	}

	if n != len(framePayload) {
		return ErrShortWrite
	}
	return nil
}

func (p *myFrameCodec) Decode(r io.Reader) (FramePayload, error) {
	var totalLen int32
	err := binary.Read(r, binary.BigEndian, &totalLen)
	if err != nil {
		return nil, err
	}
	if totalLen < 0 {
		return nil, errors.New("totalLen:" + strconv.Itoa(int(totalLen)))
	}
	
	buf := make([]byte, totalLen-4)
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}

	if n != int(totalLen-4) {
		return nil, ErrShortRead
	}
	return FramePayload(buf), nil
}

```

2.实现一个从 socket connection 中解码出 goim 协议的解码器。

## 答：
https://github.com/Terry-Mao/goim/blob/e742c99ad76e626d5f6df8b33bc47ca005501980/api/protocol/protocol.go

```
const (
	// size
	_packSize      = 4
	_headerSize    = 2
	_verSize       = 2
	_opSize        = 4
	_seqSize       = 4
	_heartSize     = 4
	_rawHeaderSize = _packSize + _headerSize + _verSize + _opSize + _seqSize
	_maxPackSize   = MaxBodySize + int32(_rawHeaderSize)
	// offset
	_packOffset   = 0
	_headerOffset = _packOffset + _packSize
	_verOffset    = _headerOffset + _headerSize
	_opOffset     = _verOffset + _verSize
	_seqOffset    = _opOffset + _opSize
	_heartOffset  = _seqOffset + _seqSize
)

func (p *Proto) ReadTCP(rr *bufio.Reader) (err error) {
	var (
		bodyLen   int
		headerLen int16
		packLen   int32
		buf       []byte
	)
	if buf, err = rr.Pop(_rawHeaderSize); err != nil {
		return
	}
	packLen = binary.BigEndian.Int32(buf[_packOffset:_headerOffset])
	headerLen = binary.BigEndian.Int16(buf[_headerOffset:_verOffset])
	p.Ver = int32(binary.BigEndian.Int16(buf[_verOffset:_opOffset]))
	p.Op = binary.BigEndian.Int32(buf[_opOffset:_seqOffset])
	p.Seq = binary.BigEndian.Int32(buf[_seqOffset:])
	if packLen > _maxPackSize {
		return ErrProtoPackLen
	}
	if headerLen != _rawHeaderSize {
		return ErrProtoHeaderLen
	}
	if bodyLen = int(packLen - int32(headerLen)); bodyLen > 0 {
		p.Body, err = rr.Pop(bodyLen)
	} else {
		p.Body = nil
	}
	return
}
```