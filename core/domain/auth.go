package domain

const (
	Issuer      = "golang-api-hexagonal"
	SuperRole   = "admin"
	ClaimsKey   = "claims"
	JwtTokenKey = "jwt"
)

type Auth struct {
	Username string `json:"username" validate:"required,not_blank,min=2,max=256"`
	Password string `json:"password" validate:"required,not_blank,min=8,max=256"`
}

type AuthClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}
