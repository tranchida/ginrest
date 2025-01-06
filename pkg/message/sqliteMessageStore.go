package message

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dataSourceName string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Create messages table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (
        id TEXT PRIMARY KEY,
        content TEXT
    )`)
	if err != nil {
		return nil, err
	}

	// Create headers table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS message_headers (
        message_id TEXT,
        key TEXT,
        value TEXT,
        PRIMARY KEY (message_id, key),
        FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE
    )`)
	if err != nil {
		return nil, err
	}

	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Add(id string, message Message) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert message
	_, err = tx.Exec("INSERT INTO messages (id, content) VALUES (?, ?)",
		id, message.Content)
	if err != nil {
		return err
	}

	// Insert headers
	for key, value := range message.Headers {
		_, err = tx.Exec("INSERT INTO message_headers (message_id, key, value) VALUES (?, ?, ?)",
			id, key, value)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *SQLiteStore) Get(id string) (Message, error) {
	// Get message content
	row := s.db.QueryRow("SELECT content FROM messages WHERE id = ?", id)
	var content string
	err := row.Scan(&content)
	if err == sql.ErrNoRows {
		return Message{}, ErrMessageNotFound
	} else if err != nil {
		return Message{}, err
	}

	// Get headers
	headers := make(map[string]string)
	rows, err := s.db.Query("SELECT key, value FROM message_headers WHERE message_id = ?", id)
	if err != nil {
		return Message{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return Message{}, err
		}
		headers[key] = value
	}

	return Message{Id: id, Content: content, Headers: headers}, nil
}

func (s *SQLiteStore) List() (map[string]Message, error) {
	messages := make(map[string]Message)

	// Get all messages
	rows, err := s.db.Query("SELECT id, content FROM messages")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, content string
		if err := rows.Scan(&id, &content); err != nil {
			return nil, err
		}
		messages[id] = Message{Id: id, Content: content, Headers: make(map[string]string)}
	}

	// Get all headers
	rows, err = s.db.Query("SELECT message_id, key, value FROM message_headers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var messageId, key, value string
		if err := rows.Scan(&messageId, &key, &value); err != nil {
			return nil, err
		}
		if msg, exists := messages[messageId]; exists {
			msg.Headers[key] = value
			messages[messageId] = msg
		}
	}

	return messages, nil
}

func (s *SQLiteStore) Update(id string, message Message) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update message content
	_, err = tx.Exec("UPDATE messages SET content = ? WHERE id = ?", message.Content, id)
	if err != nil {
		return err
	}

	// Delete existing headers
	_, err = tx.Exec("DELETE FROM message_headers WHERE message_id = ?", id)
	if err != nil {
		return err
	}

	// Insert new headers
	for key, value := range message.Headers {
		_, err = tx.Exec("INSERT INTO message_headers (message_id, key, value) VALUES (?, ?, ?)",
			id, key, value)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *SQLiteStore) Remove(id string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete headers
	_, err = tx.Exec("DELETE FROM message_headers WHERE message_id = ?", id)
	if err != nil {
		return err
	}

	// Delete message
	_, err = tx.Exec("DELETE FROM messages WHERE id = ?", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}
