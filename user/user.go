package user

type User struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	SigningKey string `json:"signingKey"`
}

// FilterValue satisfies list.Item interface
func (u User) FilterValue() string {
	return u.Username
}
