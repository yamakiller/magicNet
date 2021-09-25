package netbytes

import "io"

//WriterBytes 写入所有数据
func WriterBytes(w io.Writer, b []byte) (err error) {

	n := 0
	length := len(b)
	offset := 0

	for {
		if n, err = w.Write(b[offset:]); err != nil {
			return
		}

		offset += n
		if offset >= length {
			break
		}
	}
	return
}

//ReadBytes 读取所有数据到缓冲区
func ReadBytes(r io.Reader, out []byte) (err error) {
	n := 0
	length := len(out)
	offset := 0

	for {
		if n, err = r.Read(out[offset:]); err != nil {
			return
		}

		offset += n
		if offset >= length {
			break
		}
	}
	return
}
