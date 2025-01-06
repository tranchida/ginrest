package message

type Message struct {
    Id      string            `json:"id"`
    Content string            `json:"body"`
    Headers map[string]string `json:"headers"`
}

// AddHeader adds or updates a header
func (m *Message) AddHeader(key, value string) {
    if m.Headers == nil {
        m.Headers = make(map[string]string)
    }
    m.Headers[key] = value
}

// GetHeader returns header value for given key
func (m *Message) GetHeader(key string) (string, bool) {
    if m.Headers == nil {
        return "", false
    }
    val, exists := m.Headers[key]
    return val, exists
}

// RemoveHeader deletes a header
func (m *Message) RemoveHeader(key string) {
    if m.Headers != nil {
        delete(m.Headers, key)
    }
}
