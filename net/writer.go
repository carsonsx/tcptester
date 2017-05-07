package net


type Writer interface {
	Write(*TCPClient, interface{}) (data []byte, err error)
}

type StringWriter struct {

}

func (w *StringWriter) Write(c *TCPClient, v interface{}) (data []byte, err error) {
	str := v.(string)
	return c.Write([]byte(str))
}

type ProtobufWriter struct {

}

func (w *ProtobufWriter) Write(c *TCPClient, v interface{}) (data []byte, err error) {
	str := v.(string)
	return c.Write([]byte(str))
}