package users

type GetMeResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	OauthType string `json:"oauth_type"`
	Verified  bool   `json:"verified"`
}
