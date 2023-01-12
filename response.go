package beweb

func (c *Context) Response(code int, value []byte) error {
	c.Resp.WriteHeader(code)
	_, err := c.Resp.Write(value)
	return err
}
