package github

type AccessToken struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}
