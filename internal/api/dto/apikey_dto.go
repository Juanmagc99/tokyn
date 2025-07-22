package dto

type APIKeyCreateResponse struct {
	ID    string `json:"id"`
	Token string `json:"token"`
	Name  string `json:"name"`
}
