package users

type GetMeResponse struct {
	ID        uint     `json:"id"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles"`
	OauthType string   `json:"oauth_type"`
	Verified  bool     `json:"verified"`
}
