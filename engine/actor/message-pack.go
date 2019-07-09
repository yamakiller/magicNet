package actor

type messageHeader map[string]string

func (m messageHeader) Get(key string) string {
	return m[key]
}

func (m messageHeader) Set(key string, value string) {
	m[key] = value
}

func (m messageHeader) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

func (m messageHeader) Length() int {
	return len(m)
}

func (m messageHeader) ToMap() map[string]string {
	mp := make(map[string]string)
	for k, v := range m {
		mp[k] = v
	}
	return mp
}

// ReadOnlyMessageHeader : 只读消息头对象
type ReadOnlyMessageHeader interface {
	Get(key string) string
	Keys() []string
	Length() int
	ToMap() map[string]string
}

// MessagePack : 消息报对象
type MessagePack struct {
	Header  messageHeader
	Message interface{}
	Sender  *PID
}

// GetHeader : 获取包的对象的包头信息 key -> string
func (mp *MessagePack) GetHeader(key string) string {
	if mp.Header == nil {
		return ""
	}

	return mp.Header.Get(key)
}

// SetHeader :  设置消息头信息 key -> value
func (mp *MessagePack) SetHeader(key string, value string) {
	if mp.Header == nil {
		mp.Header = make(map[string]string)
	}
	mp.Header.Set(key, value)
}

// DefaultMessageHeader : 默认消息头
var DefaultMessageHeader = make(messageHeader)

// WrapPack 消息打包
func WrapPack(message interface{}) *MessagePack {
	if e, ok := message.(*MessagePack); ok {
		return e
	}

	return &MessagePack{nil, message, nil}
}

// UnWrapPack : 消息包拆分返回 [消息头 | 消息 | 发送者]
func UnWrapPack(message interface{}) (ReadOnlyMessageHeader, interface{}, *PID) {
	if e, ok := message.(*MessagePack); ok {
		return e.Header, e.Message, e.Sender
	}
	return nil, message, nil
}

// UnWrapPackHeader : 消息包拆分返回 [消息头]
func UnWrapPackHeader(message interface{}) ReadOnlyMessageHeader {
	if e, ok := message.(*MessagePack); ok {
		return e.Header
	}

	return nil
}

// UnWrapPackMessage : 消息包拆分返回[消息]
func UnWrapPackMessage(message interface{}) interface{} {
	if e, ok := message.(*MessagePack); ok {
		return e.Message
	}

	return message
}

// UnWrapPackSender : 消息包拆分返回[发送者]
func UnWrapPackSender(message interface{}) *PID {
	if e, ok := message.(*MessagePack); ok {
		return e.Sender
	}

	return nil
}
