package dto

type AuthData struct {
	Email    string `json:"email" example:"shad@tinkoff.ru" doc:"user email"`
	Password string `json:"password" example:"v3RyH@RdPa$$w0rd" doc:"user password"`
}

type Token struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30"` // Token string itself
}
