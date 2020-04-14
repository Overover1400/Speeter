package structs


type SendMessage struct {
	Msg string
	State bool
	FolderName string
	UserInfo EmailPass

}

type EmailPass struct {
	UserId string
	Email    string
	UserName string
	Password string
	FolderName string
	Folders []string
	Coin   int
	Status bool
}

type SendValue struct {
	MusicName string
}