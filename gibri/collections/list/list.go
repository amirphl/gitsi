package list

// TODO concurrency
// TODO similar to copyonwritearraylist

type Filter func(item interface{}) bool

type List interface {
	Add(interface{})
	Remove(idx int)
	RetainAll(Filter)
}

type list struct {
	items []interface{}
}

func (c *list) Add(item interface{}) {
	c.items = append(c.items, item)
}

func (c *list) Remove(idx int) {
	lastIdx := len(c.items) - 1
	c.items[idx] = c.items[lastIdx]
	c.items = c.items[:lastIdx]
}

func (c *list) RetainAll(filter Filter) {
	for idx, item := range c.items {
		if filter(item) {
			c.Remove(idx)
		}
	}
}

func New() List {
	return &list{}
}
