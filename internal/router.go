package internal

import (
	"github.com/gorilla/mux"
	"net/http"
)

func Route() *mux.Router {
	router := mux.NewRouter()

	//router.HandleFunc("/api/todos", app.getTodos).Methods("GET")
	//router.HandleFunc("api/todos", addTodos).Method("POST")
	//router.HandleFunc("api/todos/{id}", updateTodos).Method("POST")
	//router.HandleFunc("api/todos/{id}", deleteTodos).Method("DELETE")

	//静的ファイル読み込み
	router.Handle("/web/js/", http.StripPrefix("/web/js/", http.FileServer(http.Dir("web/js/"))))
	router.Handle("/web/css/", http.StripPrefix("/web/css/", http.FileServer(http.Dir("web/css/"))))
	router.Handle("/web/img/", http.StripPrefix("/web/img/", http.FileServer(http.Dir("web/img/"))))

	//ハンドラ登録
	//アプリ紹介ページ
	router.HandleFunc("/", IndexShowPage)
	//ユーザ認証関
	router.HandleFunc("/login/page", ShowLoginPage)
	router.HandleFunc("/login", ValidateLoginData)
	router.HandleFunc("/login/google", ShowLoginPage)
	router.HandleFunc("/login/facebook", ValidateLoginData)
	router.HandleFunc("/login/twitter", Logout)
	router.HandleFunc("/logout", Logout)
	//アカウント関係
	router.HandleFunc("/account/create/page", ShowAccountCreatePage)
	router.HandleFunc("/account/home/page", ShowHomePage)
	router.HandleFunc("/account/edit/page", ShowEditPage)
	router.HandleFunc("/account/create", CreateAccount)
	router.HandleFunc("/account/create/google", ExternalAuthenticationGoogle)
	router.HandleFunc("/account/create/facebook", ValidateLoginData)
	router.HandleFunc("/account/create/twitter", Logout)
	router.HandleFunc("/account/update", UpdateAccount)
	router.HandleFunc("/account/delete", DeleteAccount)
	router.HandleFunc("/account/imageUpload", ImageUpload)
	//問題集関係
	router.HandleFunc("/workbook/create/page", ShowWorkbookPage)
	router.HandleFunc("/workbook/folder/page", ShowWorkbookFolderPage)
	router.HandleFunc("/workbook/share", ShowWorkbookSharePage)
	router.HandleFunc("/workbook/question", ShowWorkbookQuestion)
	router.HandleFunc("/workbook/learning/page", LearningWorkbook)
	router.HandleFunc("/workbook/create", CreateWorkBook)

	// CORS対応（省略）
	return router
}
