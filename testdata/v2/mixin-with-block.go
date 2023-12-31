// Code generated by "jade.go"; DO NOT EDIT.

package jade

import (
	"bytes"

	pool "github.com/valyala/bytebufferpool"
)

func Jade_mixinwithblock(buffer *pool.ByteBuffer) {

	{
		var block []byte
		buffer.WriteString(`<div></div>`)

		if len(block) > 0 {
			buffer.Write(block)
		}
	}

	{
		var block []byte
		{
			buffer := new(bytes.Buffer)
			buffer.WriteString(`<p>test</p>`)

			block = buffer.Bytes()
		}

		buffer.WriteString(`<div></div>`)

		if len(block) > 0 {
			buffer.Write(block)
		}
	}

}
