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
	r.HandleFunc("/", app.IndexShowPage)
	//ユーザ認証関
	r.HandleFunc("/login/page", app.ShowLoginPage)
	r.HandleFunc("/login", app.ValidateLoginData)
	r.HandleFunc("/login/google", app.ShowLoginPage)
	r.HandleFunc("/login/facebook", app.ValidateLoginData)
	r.HandleFunc("/login/twitter", app.Logout)
	r.HandleFunc("/logout", app.Logout)
	//アカウント関係
	r.HandleFunc("/account/create/page", app.ShowAccountCreatePage)
	r.HandleFunc("/account/home/page", app.ShowHomePage)
	r.HandleFunc("/account/edit/page", app.ShowEditPage)
	r.HandleFunc("/account/create", app.CreateAccount)
	r.HandleFunc("/account/create/google", app.ExternalAuthenticationGoogle)
	r.HandleFunc("/account/create/facebook", app.ValidateLoginData)
	r.HandleFunc("/account/create/twitter", app.Logout)
	r.HandleFunc("/account/update", app.UpdateAccount)
	r.HandleFunc("/account/delete", app.DeleteAccount)
	r.HandleFunc("/account/imageUpload", app.ImageUpload)
	//問題集関係
	r.HandleFunc("/workbook/create/page", app.ShowWorkbookPage)
	r.HandleFunc("/workbook/folder/page", app.ShowWorkbookFolderPage)
	r.HandleFunc("/workbook/share", app.ShowWorkbookSharePage)
	r.HandleFunc("/workbook/question", app.ShowWorkbookQuestion)
	r.HandleFunc("/workbook/learning/page", app.LearningWorkbook)
	r.HandleFunc("/workbook/create", app.CreateWorkBook)

	// CORS対応（省略）
	return r

	//r := mux.NewRouter()
	//
	//// 単純なハンドラ
	//r.HandleFunc("/", YourHandler)
	//
	//// パスに変数を埋め込み
	//r.HandleFunc("/hello/{name}", VarsHandler)
	//
	//// パス変数で正規表現を使用
	//r.HandleFunc("/hello/{name}/{age:[0-9]+}", RegexHandler)
	//
	//// クエリ文字列の取得
	//r.HandleFunc("/hi/", QueryStringHandler)
	//
	//// 静的ファイルの提供
	//// $PROROOT/assets/about.html が http://localhost:8080/assets/about.html でアクセスできる
	//r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	//
	//// リダイレクト
	//r.HandleFunc("/moved", RedirectHandler)
	//
	//// マッチするパスがない場合のハンドラ
	//r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	//
	//// 複数のハンドラで共通の処理を実行する
	//// 今回はcontextのセットとゲットを試しているが、同じパターンでDBの初期化や認証処理やログ書き出しなどにも応用できる
	//// ハンドラを引き渡すには http.Handler 型を使い func(http.ResponseWriter, *http.Request) から http.Handler への変換には http.HandlerFunc を利用する
	//// http.Handler をハンドラとして登録する場合は Router.Handle を利用する
	//r.Handle("/some1", UseContext(http.HandlerFunc(SomeHandler1)))
	//r.Handle("/some2", UseContext(http.HandlerFunc(SomeHandler2)))
	//
	//// http://localhost:8080 でサービスを行う
	//http.ListenAndServe(":8080", r)

	//router.HandleFunc("/api/todos", app.getTodos).Methods("GET")
	//router.HandleFunc("api/todos", addTodos).Method("POST")
	//router.HandleFunc("api/todos/{id}", updateTodos).Method("POST")
	//router.HandleFunc("api/todos/{id}", deleteTodos).Method("DELETE")

	//静的ファイル読み込み
	//http.Handle("/web/js/", http.StripPrefix("/web/js/", http.FileServer(http.Dir("web/js/"))))
	//http.Handle("/web/css/", http.StripPrefix("/web/css/", http.FileServer(http.Dir("web/css/"))))
	//http.Handle("/web/img/", http.StripPrefix("/web/img/", http.FileServer(http.Dir("web/img/"))))
}
