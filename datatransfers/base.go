package datatransfers

type FindAllParams struct {
	Page    int
	Limit   int
	UserID  int
	Offset  int
	Status  string
	Email   string
	Name    string
	Search  string
	Token   string
	UserIDs []int
}
