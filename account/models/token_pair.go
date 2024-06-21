package models

type TokenPair struct {
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
}
