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

	//サーバー起動
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	server.ListenAndServe()
}

//func main() {
//	mux := http.NewServeMux()
//	mux.HandleFunc("/login", LoginHandler)
//	mux.HandleFunc("/loginr", LoginRHandler)
//	log.Println("Server has started")
//	http.ListenAndServe(":5001", mux)
//}
//
//func LoginHandler (w http.ResponseWriter, r *http.Request) {
//	var url = conf.AuthCodeURL("yourStateUUID", oauth2.AccessTypeOffline)
//	fmt.Fprintf(w, "Visit here : %s", url)
//}
//
//func LoginRHandler(w http.ResponseWriter, r *http.Request) {
//	code := r.URL.Query()["code"]
//	if code == nil ||  len(code) == 0 {
//		fmt.Fprint(w,"Invalid Parameter")
//	}
//	ctx := context.Background()
//	tok, err := conf.Exchange(ctx, code[0])
//	if err != nil {
//		fmt.Fprintf(w,"OAuth Error:%v", err)
//	}
//	client := conf.Client(ctx, tok)
//	svr, err := oauthapi.New(client)
//	ui, err := svr.Userinfo.Get().Do()
//	if err != nil {
//		fmt.Fprintf(w,"OAuth Error:%v", err)
//	} else {
//		fmt.Fprintf(w, "Your are logined as : %s",  ui.Email)
//		fmt.Fprintf(w, "Your are logined as : %s",  ui.Name)
//		fmt.Fprintf(w, "Your are logined as : %s",  ui.Id)
//	}
//}
//var conf =internal.GoogleGetConnect()