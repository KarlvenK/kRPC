package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

/*
客户端与服务端的通信需要协商一些内容，例如 HTTP 报文，分为 header 和 body 2 部分，
body 的格式和长度通过 header 中的 Content-Type 和 Content-Length 指定，服务端
通过解析 header 就能够知道如何从 body 中读取需要的信息。
对于 RPC 协议来说，这部分协商是需要自主设计的。为了提升性能，一般在报文的最开始会规划
固定的字节，来协商相关的信息。比如第1个字节用来表示序列化方式，第2个字节表示压缩方式，
第3-6字节表示 header 的长度，7-10 字节表示 body 的长度。
*/

type GobCodec struct {
	//conn 是由构建函数传入，通常是通过 TCP 或者 Unix 建立 socket 时得到的链接实例
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *gob.Decoder
	enc  *gob.Encoder
}

var _ Codec = (*GobCodec)(nil)

func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(buf),
	}
}

func (c *GobCodec) ReadHeader(head *Header) error {
	return c.dec.Decode(head)
}

func (c *GobCodec) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

func (c *GobCodec) Write(head *Header, body interface{}) (err error) {
	defer func() {
		_ = c.buf.Flush()
		if err != nil {
			_ = c.Close()
		}
	}()
	if err = c.enc.Encode(head); err != nil {
		log.Println("rpc codec: gob error encoding header")
		return
	}
	if err = c.enc.Encode(body); err != nil {
		log.Println("rpc codec: gob error encoding body")
		return
	}
	return
}

func (c *GobCodec) Close() error {
	return c.conn.Close()
}
