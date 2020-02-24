package utility

type UserObject struct {
	ObjectId string `json:"objectid"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TodoObject struct {
	ObjectId string `json:"objectid"`
	// CreatedAt time.Time `json:"created_at"`
	Todo      string `json:"todo"`
	CreatedBy string `json:"created_by"`
}

// https://yourbasic.org/golang/iota/
// enum for querying item attributes
type QueryCondition int

const (
	CREATED_AT = iota
	CREATED_BY
	EMAIL
)

func (q QueryCondition) String() string {
	return [...]string{"CREATED_BY", "CREATED_AT"}[q]
}

type Todo struct {
	UserID  string `json:"userId"`
	TodoID  string `json:"noteId"`
	Content string `json:"content"`
}
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
