package utility

type UserObject struct {
	ObjectId string `json:"objectid"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type TodoObject struct {
	ObjectId  string `json:"objectid"`
	CreatedAt string `json:"created_at"`
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
