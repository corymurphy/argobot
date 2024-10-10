package events

// type Event interface {
// 	Action() string
// }

type Event struct {
	Action *string `json:"action,omitempty"`
}
