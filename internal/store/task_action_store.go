package store

type TypeAction string

const (
	Type1 TypeAction = "type_1"
	Type2 TypeAction = "type_2"
	Type3 TypeAction = "type_3"
)

type ActionTask struct {
	ID          int        `json:"id"`
	Type        TypeAction `json:"type"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
}
