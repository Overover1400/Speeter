package structs


type SendMessage struct {
	Msg string
	State bool
	FolderName string
	UserInfo EmailPass

}

type EmailPass struct {
	Email    string
	UserName string
	Password string
}