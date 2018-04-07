package App

type User struct {
	Uid      int
	Username string
	Password string
}

type Tweet struct {
	Tid        int
	Created_by User
	Content    string
	Time       string
}
