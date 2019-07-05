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

func (m messageHeader) Length() int  {
  return len(m)
}

func (m messageHeader) ToMap() map[string]string {
  mp := make(map[string]string)
  for k, v := range m {
    mp[k] = v
  }
  return mp
}

type ReadOnlyMessageHeader interface {
  Get(key string) string
  Keys() []string
  Length() int
  ToMap() map[string]string
}

type MessagePack struct {
  Header  messageHeader
  Message interface{}
  Sender  *PID
}

func (mp *MessagePack) GetHeader(key string) string {
  if mp.Header == nil {
    return ""
  }

  return mp.Header.Get(key)
}

func (mp *MessagePack) SetHeader(key string, value string) {
  if mp.Header == nil {
    mp.Header = make(map[string]string)
  }
  mp.Header.Set(key, value)
}

var EmptyMessageHeader = make(messageHeader)

func WrapPack(message interface{}) *MessagePack {
  if e, ok := message.(*MessagePack); ok {
    return e
  }
  return &MessagePack{nil, message, nil}
}

func UnWrapPack(message interface{}) (ReadOnlyMessageHeader, interface{}, *PID) {
  if e, ok := message.(*MessagePack);ok {
    return e.Header, e.Message, e.Sender
  }
  return nil, message, nil
}

func UnWrapPackHeader(message interface{}) ReadOnlyMessageHeader {
   if e, ok := message.(*MessagePack); ok {
     return e.Header
   }

   return nil
}

func UnWrapPackMessage(message interface{}) interface{} {
  if e, ok := message.(*MessagePack); ok {
    return e.Message
  }

  return nil
}

func UnWrapPackSender(message interface{}) *PID {
   if e, ok := message.(*MessagePack); ok {
     return e.Sender
   }

   return nil
}
