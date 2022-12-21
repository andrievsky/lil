package main

type Content interface {
	Path() Path
	Data() string
}

type ContentImpl struct {
	path Path
	data string
}

func NewContent(path Path, data string) Content {
	return &ContentImpl{
		path,
		data,
	}
}

func (c *ContentImpl) Path() Path {
	return c.path
}

func (c *ContentImpl) Data() string {
	return c.data
}
