package beweb

func (c *Context) Response(code int, value []byte) {
	c.ResponseStatus = code
	c.ResponseContent = value
}
