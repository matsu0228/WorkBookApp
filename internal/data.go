package internal

import (
	"github.com/astaxie/beego"
	"html/template"
	"net/http"
)

//アカウントテンプレート
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

//オプション
type Option struct {
	NumberOfQuestions string
	Shuffle           bool
}

//設問
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

// CallbackController コールバックコントローラ
type CallbackController struct {
	beego.Controller
}

// CallbackRequest コールバックリクエスト
type CallbackRequest struct {
	Code  string `form:"code"`
	State string `form:"state"`
}

type Cookies []*http.Cookie

const (

	//

	//HTML各パーツ
	head                      = "web/template/includeParts/head.html"
	script                    = "web/template/includeParts/script.html"
	home_content              = "web/template/includeParts/home_content.html"
	account_edit_content      = "web/template/includeParts/account_edit_content.html"
	sidebar                   = "web/template/includeParts/sidebar.html"
	nav                       = "web/template/includeParts/nav.html"
	footer                    = "web/template/includeParts/footer.html"
	login                     = "web/template/login.html"
	account_create            = "web/template/account_create.html"
	home                      = "web/template/home.html"
	workbook_create_content   = "web/template/includeParts/workbook_create_content.html"
	workbook_folder           = "web/template/includeParts/workbook_folder.html"
	workbook_learning_content = "web/template/includeParts/workbook_learning_content.html"
	workbook_share            = "web/template/includeParts/workbook_share.html"
	workbook_question         = "web/template/includeParts/workbook_question.html"

	//showdata用(HTML名前)
	show_login          = "login"
	show_account_create = "account_create"
	show_home           = "home"

	//バケット名
	BucketName = "gompei/"

	//表示メッセージ
	succes_account_create_message       = "アカウント作成が完了しました。ログインして下さい。"
	succes_account_delete_message       = "アカウントの削除が完了しました"
	succes_logout_message               = "ログアウトが完了しました"
	error_database_message              = "システムエラーが起きました。管理者にご連絡下さい(エラー内容:データベース登録)"
	error_htmlfile_message              = "システムエラーが起きました。管理者にご連絡下さい(エラー内容:htmlファイル読み込み)"
	error_cookie_cannot_confirm_message = "セッション情報が確認出来ませんでした。ログインをもう一度お願いします。"
	error_login_failed_message          = "メールアドレス又はパスワードが違います"

	//switch文用文字列（スペルミスを減らすため）
	s_title             = "title"
	s_numberOfQuestions = "numberOfQuestions"
	s_shuffle           = "shuffle"

	//from取得用文字列
	f_user_id   = "userId"
	f_user_name = "userName"
	f_email     = "email"
	f_password  = "password"
	f_image     = "image"
	f_book_id   = "bookId"

	//datastore関係
	project_id  = "quiz-283808"
	project_key = "./quiz-283808-424072cf23ec.json"
)

var (
	//ページ本体
	pageLogin            = []string{head, script, login}
	pageAccountCreate    = []string{head, script, account_create}
	pageHome             = []string{head, script, home, home_content, sidebar, nav, footer}
	pageAccountEdit      = []string{head, script, home, account_edit_content, sidebar, nav, footer}
	pageWorkbookCreate   = []string{head, script, home, workbook_create_content, sidebar, nav, footer}
	pageWorkbookFolder   = []string{head, script, home, workbook_folder, sidebar, nav, footer}
	pageWorkbookLearning = []string{head, script, home, workbook_learning_content, sidebar, nav, footer}
	pageWorkbookShare    = []string{head, script, home, workbook_share, sidebar, nav, footer}
	pageWorkbookQuestion = []string{head, script, home, workbook_question, sidebar, nav, footer}

	//独自関数（フロント側）
	FuncMap = template.FuncMap{
		"Increment": func(i int) int {
			return i + 1
		},
	}
)
