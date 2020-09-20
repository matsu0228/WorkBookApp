package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

func Route(app *App) *mux.Router {
	//
	r := mux.NewRouter()

	//静的ファイルの読み込み
	r.PathPrefix("/web/js/").Handler(http.StripPrefix("/web/js/", http.FileServer(http.Dir("web/js/"))))
	r.PathPrefix("/web/css/").Handler(http.StripPrefix("/web/css/", http.FileServer(http.Dir("web/css/"))))
	r.PathPrefix("/web/img/").Handler(http.StripPrefix("/web/img/", http.FileServer(http.Dir("web/img/"))))

	//ハンドラ登録
	//アプリ紹介ページ
	r.HandleFunc("/", app.IndexShowPage).Methods("GET")
	//ユーザ認証関
	r.HandleFunc("/login/page", app.ShowLoginPage).Methods("GET")
	r.HandleFunc("/login/forgot-password", app.ShowForgotPasswordPage).Methods("GET")
	r.HandleFunc("/login/PasswordReissue", app.SendReissueEmail).Methods("POST")
	r.HandleFunc("/login/recover-password/page", app.ShowRecoverPasswordPage).Methods("GET")
	r.HandleFunc("/login/recover-password", app.RecoverPassword).Methods("POST")

	r.HandleFunc("/login", app.ValidateLoginData).Methods("POST")
	r.HandleFunc("/login/google", app.ShowLoginPage).Methods("GET")
	r.HandleFunc("/login/facebook", app.ValidateLoginData).Methods("GET")
	r.HandleFunc("/login/github", app.Logout).Methods("GET")
	r.HandleFunc("/logout", app.Logout).Methods("GET")

	//アカウント関係
	r.HandleFunc("/account/create/page", app.ShowAccountCreatePage).Methods("GET")
	r.HandleFunc("/account/home/page", app.ShowHomePage).Methods("GET")
	r.HandleFunc("/account/edit/page", app.ShowEditPage).Methods("GET")
	r.HandleFunc("/account/create", app.CreateAccount).Methods("POST")
	r.HandleFunc("/account/create/google", app.ExternalAuthenticationGoogle).Methods("GET")
	r.HandleFunc("/account/create/facebook", app.ExternalAuthenticationFaceBook).Methods("GET")
	r.HandleFunc("/account/create/github", app.ExternalAuthenticationGithub).Methods("GET")
	r.HandleFunc("/account/update", app.UpdateAccount).Methods("POST")
	r.HandleFunc("/account/delete", app.DeleteAccount).Methods("POST")
	r.HandleFunc("/account/imageUpload", app.ImageUpload).Methods("POST")

	//問題集関係
	r.HandleFunc("/workbook/create/page", app.ShowWorkbookPage).Methods("GET")
	r.HandleFunc("/workbook/folder/page", app.ShowWorkbookFolderPage).Methods("GET")
	r.HandleFunc("/workbook/share/page", app.ShowWorkbookSharePage).Methods("GET")
	r.HandleFunc("/workbook/learning/page", app.LearningWorkbook).Methods("POST")
	r.HandleFunc("/workbook/share", app.WorkbookUpload).Methods("POST")
	r.HandleFunc("/workbook/create", app.CreateWorkBook).Methods("POST")

	// CORS対応（省略）
	return r
}
