package user

import (
	"fmt"
	"strings"
)

// User represents Git configuration
type User struct {
	Username    string `json:"username" cmd:"user.name"`
	Email       string `json:"email" cmd:"user.email"`
	SigningKey  string `json:"signingKey" cmd:"user.signingkey"`
	SignCommits *bool  `json:"gpgSign" cmd:"commit.gpgsign"`
	GpgFormat   string `json:"gpgFormat" cmd:"gpg.format"`
}

// FilterValue satisfies list.Item interface
func (u User) FilterValue() string {
	return u.Username
}

// Title ...
func (u User) Title() string {
	return fmt.Sprintf("%s | %s", u.Username, u.Email)
}

// Description ...
func (u User) Description() string {
	var b strings.Builder
	if u.SigningKey != "" {
		b.WriteString(fmt.Sprintf("SigningKey %s", u.SigningKey))
	}
	if b.Len() != 0 {
		b.WriteString(" | ")
	}

	if u.SignCommits != nil && *u.SignCommits {
		b.WriteString("Auto sign commits")
	} else {
		b.WriteString("Do not sign commits")
	}

	if u.GpgFormat != "" {
		if b.Len() != 0 {
			b.WriteString(" | ")
		}
		b.WriteString(fmt.Sprintf("GPG Format %s", u.GpgFormat))
	}

	return b.String()
}
