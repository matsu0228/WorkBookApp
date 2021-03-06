package api

import (
	"html/template"
)

const (
	//HTML各パーツ
	index                     = "web/template/index.html"
	head                      = "web/template/includeParts/head.html"
	script                    = "web/template/includeParts/script.html"
	home_content              = "web/template/includeParts/home_content.html"
	account_edit_content      = "web/template/includeParts/account_edit_content.html"
	sidebar                   = "web/template/includeParts/sidebar.html"
	nav                       = "web/template/includeParts/nav.html"
	footer                    = "web/template/includeParts/footer.html"
	login                     = "web/template/login.html"
	forgot_password           = "web/template/forgot_password.html"
	recover_password          = "web/template/recover_password.html"
	account_create            = "web/template/account_create.html"
	home                      = "web/template/home.html"
	workbook_create_content   = "web/template/includeParts/workbook_create_content.html"
	workbook_folder           = "web/template/includeParts/workbook_folder.html"
	workbook_learning_content = "web/template/includeParts/workbook_learning_content.html"
	workbook_share            = "web/template/includeParts/workbook_share.html"
	workbook_question         = "web/template/includeParts/workbook_question.html"
	not_found                 = "web/template/404.html"
	error500                  = "web/template/500.html"

	//ShowData用{{define "ここの名前"}}
	Show_index            = "index"
	Show_login            = "login"
	Show_forgot_password  = "forgot_password"
	Show_recover_password = "recover_password"
	Show_account_create   = "account_create"
	Show_home             = "home"
	Show_404              = "404"
	Show_500              = "500"
	//表示メッセージ
	Succes_account_create_message = "アカウント作成が完了しました。ログインして下さい。"
	Succes_account_delete_message = "アカウントの削除が完了しました"
	Succes_logout_message         = "ログアウトが完了しました"
	Error_database_message        = "システムエラーが起きました。管理者にご連絡下さい(エラー内容:データベース登録)"
	//Error_htmlfile_message              = "システムエラーが起きました。管理者にご連絡下さい(エラー内容:htmlファイル読み込み)"
	Error_cookie_cannot_confirm_message = "セッション情報が確認出来ませんでした。ログインをもう一度お願いします。"
	//Error_login_failed_message          = "メールアドレス又はパスワードが違います"

	//switch文用文字列（スペルミスを減らすため）
	S_title             = "title"
	S_numberOfQuestions = "numberOfQuestions"
	S_shuffle           = "shuffle"

	//from取得用文字列（スペルミスを減らすため）
	F_user_id   = "userId"
	F_user_name = "userName"
	F_email     = "email"
	F_password  = "password"
	F_image     = "image"
	F_book_id   = "bookId"

	//GCPプロジェクトID
	Project_id = "apptestgo0000"
	//CloudStoreバケット名
	BucketName = "gompei/"
)

var (
	//
	oauthStateString = "thisshouldberandom"

	//ページ本体
	PageIndex            = []string{index}
	Page404              = []string{head, script, not_found}
	Page500              = []string{head, script, error500}
	PageLogin            = []string{head, script, login}
	PageForgotPassword   = []string{head, script, forgot_password}
	PageRecoverPassword  = []string{head, script, recover_password}
	PageAccountCreate    = []string{head, script, account_create}
	PageHome             = []string{head, script, home, home_content, sidebar, nav, footer}
	PageAccountEdit      = []string{head, script, home, account_edit_content, sidebar, nav, footer}
	PageWorkbookCreate   = []string{head, script, home, workbook_create_content, sidebar, nav, footer}
	PageWorkbookFolder   = []string{head, script, home, workbook_folder, sidebar, nav, footer}
	PageWorkbookLearning = []string{head, script, home, workbook_learning_content, sidebar, nav, footer}
	PageWorkbookShare    = []string{head, script, home, workbook_share, sidebar, nav, footer}
	PageWorkbookQuestion = []string{head, script, home, workbook_question, sidebar, nav, footer}

	//フロント側に組み込む関数
	FuncMap = template.FuncMap{
		"Increment": func(i int) int {
			return i + 1
		},
	}
)
