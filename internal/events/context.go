package events

type Handler func(*ConsumerCtx) error
type ConsumerCtx struct {
	handlers []Handler
	message  []byte
	values   map[string]any
	pivot    int
}

func (cc *ConsumerCtx) SetValue(key string, value any) {
	cc.values[key] = value
}
func (cc *ConsumerCtx) GetValue(key string) any {
	value, ok := cc.values[key]
	_ = ok
	return value
}
func (cc *ConsumerCtx) GetMessage() []byte {
	return cc.message
}
func (cc *ConsumerCtx) Next() error {
	if len(cc.handlers) > cc.pivot {
		cc.pivot++
		return cc.handlers[cc.pivot-1](cc)
	}
	return nil
}
