package auth

// User represents the User object.
type User struct {
	UID         string `json:"uid,omitempty"`
	DisplayName string `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	Email       string `json:"email,omitempty"`
	Token       Token  `json:"stsTokenManager,omitempty"`
}

// Token stores the Oauth2 Token.
type Token struct {
	APIKey         string `json:"apiKey,omitempty" yaml:"apiKey,omitempty"`
	RefreshToken   string `json:"refreshToken,omitempty" yaml:"refreshToken,omitempty"`
	AccessToken    string `json:"accessToken,omitempty" yaml:"accessToken,omitempty"`
	ExpirationTime int    `json:"expirationTime,omitempty" yaml:"expirationTime,omitempty"`
}
