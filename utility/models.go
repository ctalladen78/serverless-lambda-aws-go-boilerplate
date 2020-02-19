package utility

type UserObject struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TodoObject struct {
	Id        string `json:"id"`
	CreatedAt string `json:"created_at"`
	Todo      string `json:"todo"`
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
