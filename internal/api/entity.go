package api

import (
	"cloud.google.com/go/storage"
	"mime/multipart"
	"net/http"
)

//DBとの接続情報（リポジトリ）
type Repository interface {
	/**/
	CheckUserLogin(user UserAccount, password string) (error, UserAccount)

	//SaveUserAccount はユーザIDを主キーにして、データを登録する関数
	InsertUserAccount(user UserAccount) error

	/**/
	InsertOauthAccount(user UserAccount, kind string) error

	/**/
	SelectOauthAccount(user UserAccount, kind string) error

	/*SelectAccountMailは*/
	SelectUserAccountMail(searchTarget string) (error, int64)

	/*UpdateAccountPasswordは*/
	UpdateUserAccountPassword(mail string, password string, userId string) error

	//UpdateUserAccount　は登録されているユーザ情報を更新する関数
	UpdateUserAccount(cookie *http.Cookie, updateAccount UserAccount) (error, UserAccount)

	//DeleteUserAccount　は登録しているユーザ情報を削除する関数
	DeleteUserAccount(cookie *http.Cookie) error

	//CreateWorkbook　は４択問題集をbookIDを主キーにデータを登録する関数
	InsertWorkbook(book WorkbookContent) error

	//SelectWorkbooks は問題集のタイトル,IDを取得して,boolean,構造体の配列を返す関数
	SelectWorkbooks(id string) (error, []WorkbookContent)

	//SelectWorkbook は
	SelectWorkbook(id string) (error, WorkbookContent)

	//InsertShareWorkbook は
	InsertShareWorkbook(bookId string) error

	//SelectShareWorkbook　は
	SelectShareWorkbooks()
}

//Cloud Storageクライアント
type Client struct {
	CloudStorage *storage.Client
}

//ストレージとの接続情報（リポジトリ）
type Storage interface {
	//
	UploadImg(file multipart.File, fileHeader *multipart.FileHeader) error
}

//アプリケーション情報
type App struct {
	//データベース情報
	DB Repository
	//ストレージ情報
	ST Storage
}

//ユーザアカウント
type UserAccount struct {
	Id           int64
	Name         string
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
