package message

type MemoryMessageStore struct {
	messages map[string]Message
}

func NewMemoryMessageStore() (*MemoryMessageStore, error) {
	return &MemoryMessageStore{
		messages: make(map[string]Message),
	}, nil
}

func (s *MemoryMessageStore) Add(id string, message Message) error {
	s.messages[id] = message
	return nil
}

func (s *MemoryMessageStore) Get(id string) (Message, error) {
	msg, exists := s.messages[id]
	if !exists {
		return Message{}, ErrMessageNotFound
	}
	return msg, nil
}	

func (s *MemoryMessageStore) Update(id string, message Message) error {
	if _, exists := s.messages[id]; !exists {
		return ErrMessageNotFound
	}
	s.messages[id] = message
	return nil
}

func (s *MemoryMessageStore) Remove(id string) error {
	if _, exists := s.messages[id]; !exists {
		return ErrMessageNotFound
	}
	delete(s.messages, id)
	return nil
}

func (s *MemoryMessageStore) List() (map[string]Message, error) {
	return s.messages, nil
}

