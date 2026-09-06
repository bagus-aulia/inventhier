package user

// User represents a user data from the user service
type User struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Position string `json:"position"`
}
