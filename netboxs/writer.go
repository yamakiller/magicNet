package netboxs

import "bufio"

//WriteToBuffer Write data to writer/buffer TODO: remove
func WriteToBuffer(writer *bufio.Writer, buffer []byte) error {
	length := len(buffer)
	seek := 0
	for {
		if writer.Available() == 0 {
			if err := writer.Flush(); err != nil {
				return err
			}
		}

		n, err := writer.Write(buffer[seek:])
		if err != nil {
			return err
		}

		seek += n
		if seek >= length {
			break
		}
	}

	if writer.Buffered() > 0 {
		if err := writer.Flush(); err != nil {
			return err
		}
	}
	return nil
}
