package main

import (
	"Workbook/internal"
	"net/http"
)

func main() {

	//静的ファイル読み込み
	http.Handle("/web/js/", http.StripPrefix("/web/js/", http.FileServer(http.Dir("web/js/"))))
	http.Handle("/web/css/", http.StripPrefix("/web/css/", http.FileServer(http.Dir("web/css/"))))
	http.Handle("/web/img/", http.StripPrefix("/web/img/", http.FileServer(http.Dir("web/img/"))))

	/*ハンドラ登録*/
	//ログイン関係
	http.HandleFunc("/", internal.ShowLoginPage)
	http.HandleFunc("/login", internal.ValidateLoginData)
	http.HandleFunc("/logout",internal.Logout)

	//アカウント関係
	http.HandleFunc("/account_create_page",internal.ShowAccountCreatePage)
	http.HandleFunc("/account_home_page",internal.ShowHomePage)
	http.HandleFunc("/account_edit_page",internal.ShowEditPage)
	http.HandleFunc("/account_create",internal.CreateAccount)
	http.HandleFunc("/account_update",internal.UpdateAccount)
	http.HandleFunc("/account_delete",internal.DeleteAccount)
	http.HandleFunc("/image_upload",internal.ImageUpload)

	//問題集関係
	http.HandleFunc("/workbook_create_page",internal.ShowWorkbookPage)
	http.HandleFunc("/workbook_folder_page",internal.ShowWorkbookFolderPage)
	http.HandleFunc("/workbook_learning_page",internal.LearningWorkbook)
	http.HandleFunc("/workbook_create",internal.CreateWorkBook)

	//サーバー起動
	server:= http.Server{
		Addr: "127.0.0.1:8080",
	}
	server.ListenAndServe()
}
