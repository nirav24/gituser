package user

type User struct {
	Username   string
	Email      string
	SigningKey string
	ShouldSign bool
}

func (u User) FilterValue() string {
	return u.Username
}
