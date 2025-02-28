package types

type GCResponse[T any] struct {
	IsOK    bool   `json:"IsOK"`
	Message string `json:"Message"`
	Payload *T     `json:"Payload,omitempty"`
}

type Deck struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}
