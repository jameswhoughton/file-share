package file_share

type User struct {
	Id       int
	Password string
	Email    string
	ApiKey   string
}

type UserService interface {
	Get(id int) (User, error)
	GetFromSessionId(sessionId string) (User, error)
	GetWithCredentials(email, password string) (User, error)
	Add(user User) (User, error)
}
