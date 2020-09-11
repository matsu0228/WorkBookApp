package main

import (
	"Workbook/internal"
	"net/http"
	"google.golang.org/appengine"
)

func main() {
	//静的ファイル読み込み
	http.Handle("/web/js/", http.StripPrefix("/web/js/", http.FileServer(http.Dir("web/js/"))))
	http.Handle("/web/css/", http.StripPrefix("/web/css/", http.FileServer(http.Dir("web/css/"))))
	http.Handle("/web/img/", http.StripPrefix("/web/img/", http.FileServer(http.Dir("web/img/"))))

	//ハンドラ登録
	//ユーザ認証関係
	http.HandleFunc("/", internal.ShowLoginPage)
	http.HandleFunc("/login", internal.ValidateLoginData)
	http.HandleFunc("/login/google", internal.ShowLoginPage)
	http.HandleFunc("/login/facebook", internal.ValidateLoginData)
	http.HandleFunc("/login/twitter", internal.Logout)
	http.HandleFunc("/logout", internal.Logout)
	//アカウント関係
	http.HandleFunc("/account/create/page", internal.ShowAccountCreatePage)
	http.HandleFunc("/account/home/page", internal.ShowHomePage)
	http.HandleFunc("/account/edit/page", internal.ShowEditPage)
	http.HandleFunc("/account/create", internal.CreateAccount)
	http.HandleFunc("/account/create/google", internal.ExternalAuthenticationGoogle)
	http.HandleFunc("/account/create/facebook", internal.ValidateLoginData)
	http.HandleFunc("/account/create/twitter", internal.Logout)
	http.HandleFunc("/account/update", internal.UpdateAccount)
	http.HandleFunc("/account/delete", internal.DeleteAccount)
	http.HandleFunc("/account/imageUpload", internal.ImageUpload)
	//問題集関係
	http.HandleFunc("/workbook/create/page", internal.ShowWorkbookPage)
	http.HandleFunc("/workbook/folder/page", internal.ShowWorkbookFolderPage)
	http.HandleFunc("/workbook/share", internal.ShowWorkbookSharePage)
	http.HandleFunc("/workbook/question", internal.ShowWorkbookQuestion)
	http.HandleFunc("/workbook/learning/page", internal.LearningWorkbook)
	http.HandleFunc("/workbook/create", internal.CreateWorkBook)

	//サーバー起動(ローカルテスト時)
	//server := http.Server{
	//	Addr: "127.0.0.1:8080",
	//}
	//server.ListenAndServe()

	//GAEでのプログラム実行
	appengine.Main()
}