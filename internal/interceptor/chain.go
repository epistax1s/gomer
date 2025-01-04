package interceptor

type ChainBuilder struct {
	interceptors [] Interceptor
}

func NewChainBuilder() *ChainBuilder {
	return &ChainBuilder{}
}

func (c *ChainBuilder) Add(i Interceptor) *ChainBuilder {
	c.interceptors = append(c.interceptors, i)
	return c
}

func (c *ChainBuilder) Build() Interceptor {
	for i := 0; i < len(c.interceptors) - 1; i++ {
		c.interceptors[i].SetNext(c.interceptors[i+1])
	}
	return c.interceptors[0]
}