package storage

type Collection []interface{}

func NewCollection() *Collection {
	return &Collection{}
}

func (c *Collection) Save() {
}

func (c *Collection) Find() {
}

func (c *Collection) Clear() {
}
