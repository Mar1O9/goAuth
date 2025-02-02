package inputs

type InputUser struct {
	Username        string `json:"username"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}
type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
