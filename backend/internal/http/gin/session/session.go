package session

type Session struct {
	UserID      int    `json:"user_id"`
	RedirectURL string `json:"redirect_url"`
	State       string `json:"state"`
}

func NewSession() Session {
	return Session{}
}
