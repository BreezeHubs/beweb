package beweb

import "net/http"

func (c *Context) Render(templateName string, data any) error {
	bytes, err := c.templateEngine.Render(c.Req.Context(), templateName, data)
	if err != nil {
		c.ResponseStatus = http.StatusInternalServerError
		return err
	}

	c.Response(http.StatusOK, bytes)
	return nil
}
