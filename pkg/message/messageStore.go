package message

import "errors"

var (
	ErrMessageNotFound = errors.New("not found")
)

type MessageStore interface {
	Add(id string, message Message) error
	Get(id string) (Message, error)
	Update(id string, message Message) error
	Remove(id string) error
	List() (map[string]Message, error)
}
