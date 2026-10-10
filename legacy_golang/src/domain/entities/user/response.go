package entityuser

type UserLoginResponse struct {
	User  *UserEntity      `json:"user"`
	Token *UserTokenEntity `json:"token"`
}
