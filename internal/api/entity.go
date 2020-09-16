package api

import (
	"mime/multipart"
	"net/http"
)

//DBとStorageの置き場所
type Repository interface {
	CheckUserLogin(user UserAccount, password string) (bool, UserAccount)
	SaveUserAccount(user UserAccount) bool
	UpdateUserAccount(cookie *http.Cookie, updateAccount UserAccount) (bool, UserAccount)
	DeleteUserAccount(cookie *http.Cookie) bool
	CreateWorkbook(book WorkbookContent) bool
	SelectWorkbooks(id string) (bool, []WorkbookContent)
	SelectWorkbook(id string) (bool, WorkbookContent)
}

type Storage interface {
	UploadImg(file multipart.File, fileHeader *multipart.FileHeader) error
}

//どうしても実装できなった（interfaceなしだと）
type App struct {
	//DB *datastore.Client
	//ST *storage.Client
	DB Repository
	ST Storage
}

//ユーザアカウント
type UserAccount struct {
	UserId       int64
	UserName     string
	Mail         string
	HashPassword []byte
	ProfileImg   string
}

//問題集（ユーザーID+タイトル+オプション+設問の配列）
type WorkbookContent struct {
	UserId   int64
	BookId   int64
	Title    string
	Options  Option
	Contents []Content
}

//オプション（問題集）
type Option struct {
	NumberOfQuestions string
	Shuffle           bool
}

//設問（問題集）
type Content struct {
	ProblemNumber    string
	ProblemStatement string
	Choice1          string
	Choice2          string
	Choice3          string
	Choice4          string
	Answer           string
	Explanation      string
}

//クッキーで名前、ユーザー画像、IDを保持
type Cookies struct {
	UserName *http.Cookie
	Image    *http.Cookie
	UserID   *http.Cookie
}
