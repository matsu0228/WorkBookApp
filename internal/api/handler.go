package api

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	fb "github.com/huandu/facebook"
	oauthapi "google.golang.org/api/oauth2/v2"
	"html/template"
	"net/http"
	"strconv"
)

//NewApp は
func NewApp(d Repository, s Storage) *App {
	return &App{
		DB: d,
		ST: s,
	}
}

//IndexShowPage　は
func (a *App) IndexShowPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		a.Show404Page(w, r, http.StatusNotFound)
		return
	}
	ReadTemplate(w, PageIndex, "index", nil)
}

//ShowLoginPage は
func (a *App) ShowLoginPage(w http.ResponseWriter, r *http.Request) {
	ReadTemplate(w, PageLogin, Show_login, nil)
}

/*ShowForgotPasswordPage*/
func (a *App) ShowForgotPasswordPage(w http.ResponseWriter, r *http.Request) {
	ReadTemplate(w, PageForgotPassword, Show_forgot_password, nil)
}

/*ShowRecoverPasswordPageは*/
func (a *App) ShowRecoverPasswordPage(w http.ResponseWriter, r *http.Request) {
	//
	if _, err := r.Cookie("email"); err != nil {
		ErrorHandling(err, nil)
		return
	}

	ReadTemplate(w, PageRecoverPassword, Show_recover_password, nil)
}

//ShowAccountCreatePageは
func (a *App) ShowAccountCreatePage(w http.ResponseWriter, r *http.Request) {
	urls := CreationUrl()
	ReadTemplate(w, PageAccountCreate, Show_account_create, urls)
}

//HOMEページ
func (a *App) ShowHomePage(w http.ResponseWriter, r *http.Request) {
	if ok, _ := ConfirmationCookie(w, r); ok {
		ReadTemplate(w, PageHome, Show_home, nil)
	}
}

//ShowEditPage　は
func (a *App) ShowEditPage(w http.ResponseWriter, r *http.Request) {
	if ok, _ := ConfirmationCookie(w, r); ok {
		ReadTemplate(w, PageAccountEdit, Show_home, nil)
	}
}

//ShowWorkbookPage　は
func (a *App) ShowWorkbookPage(w http.ResponseWriter, r *http.Request) {
	if ok, _ := ConfirmationCookie(w, r); ok {
		ReadTemplate(w, PageWorkbookCreate, Show_home, nil)
	}
}

//学習フォルダページ
func (a *App) ShowWorkbookFolderPage(w http.ResponseWriter, r *http.Request) {
	if ok, cookies := ConfirmationCookie(w, r); ok {
		if ok, workbooks := a.DB.SelectWorkbooks(cookies.UserID.Value); ok {
			ReadTemplateToIncludeFunction(w, PageWorkbookFolder, Show_home, workbooks, FuncMap)
		} else {
			ErrorHandling(nil, nil)
		}
	}
}

//問題集共有ページ
func (a *App) ShowWorkbookSharePage(w http.ResponseWriter, r *http.Request) {
	if ok, _ := ConfirmationCookie(w, r); ok {
		ReadTemplate(w, PageWorkbookShare, Show_home, nil)
	}
}

//問題質問ページ
func (a *App) ShowWorkbookQuestion(w http.ResponseWriter, r *http.Request) {
	if ok, _ := ConfirmationCookie(w, r); ok {
		ReadTemplate(w, PageWorkbookQuestion, Show_home, nil)
	}
}

//各学習ページ
func (a *App) LearningWorkbook(w http.ResponseWriter, r *http.Request) {
	if ok, _ := ConfirmationCookie(w, r); ok {
		bookId := r.FormValue(F_book_id)
		if ok, workbook := a.DB.SelectWorkbook(bookId); ok {
			ReadTemplateToIncludeFunction(w, PageWorkbookLearning, Show_home, workbook, FuncMap)
		} else {
			ErrorHandling(nil, nil)
		}
	}
}

//Show404Page
func (a *App) Show404Page(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		//ReadTemplate()
	} else {

	}
}

//Show500Page は
func (a *App) Show500Page() {

}

//ログアウト
func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	if DiscardCookie(w, r) {
		ReadTemplate(w, PageLogin, Show_login, Succes_logout_message)
	}
}

/*アカウント関係ハンドラ*/
//ログイン
func (a *App) ValidateLoginData(w http.ResponseWriter, r *http.Request) {
	//from情報取得
	user := UserAccount{
		Mail: r.FormValue(F_email),
	}
	password := r.FormValue(F_password)

	//データベースに問い合わせ
	if ok, user := a.DB.CheckUserLogin(user, password); ok {
		if ok, c := CreateCookie(w, r, user); ok {
			ReadTemplate(w, PageHome, Show_home, c)
		}
	} else {
		ReadTemplate(w, PageLogin, Show_login, nil)
		ErrorHandling("login", nil)
	}
}

//アカウント作成
func (a *App) CreateAccount(w http.ResponseWriter, r *http.Request) {
	//fromからデータ取得
	user := UserAccount{
		UserName:     r.FormValue(F_user_name),
		Mail:         r.FormValue(F_email),
		HashPassword: HashFiled(r.FormValue(F_password)),
		ProfileImg:   BucketName + "no_image_square.jpg",
	}

	//アカウント作成
	if a.DB.SaveUserAccount(user) {
		ReadTemplate(w, PageLogin, Show_login, Succes_account_create_message)

	} else {
		ReadTemplate(w, PageAccountCreate, Show_account_create, Error_database_message)
	}
}

//Oauth認証（登録）(Google)
func (a *App) ExternalAuthenticationGoogle(w http.ResponseWriter, r *http.Request) {
	conf := GoogleGetConnect()
	code := r.URL.Query()["code"]
	if code == nil || len(code) == 0 {
		fmt.Fprint(w, "Invalid Parameter")
	}
	ctx := context.Background()
	tok, err := conf.Exchange(ctx, code[0])
	if err != nil {
		fmt.Fprintf(w, "OAuth Error:%v", err)
	}
	client := conf.Client(ctx, tok)

	svr, err := oauthapi.New(client)
	ui, err := svr.Userinfo.Get().Do()
	if err != nil {
		fmt.Fprintf(w, "OAuth Error:%v", err)
	} else {
		fmt.Fprintf(w, "Your are logined as : %s", ui.Email)
		fmt.Fprintf(w, "Your are logined as : %s", ui.Name)
		fmt.Fprintf(w, "Your are logined as : %s", ui.Id)
	}
}

//Oauth認証（ログイン）(Google)
func (a *App) LoginGoogleAccount(w http.ResponseWriter, r *http.Request) {

}

//Oauth認証(登録)(FaceBook)
func (a *App) ExternalAuthenticationFaceBook(w http.ResponseWriter, r *http.Request) {
	conf := FacebookGetConnect()
	code := r.URL.Query()["code"]
	if code == nil || len(code) == 0 {
		fmt.Fprint(w, "Invalid Parameter")
	}
	ctx := context.Background()
	tok, err := conf.Exchange(ctx, code[0])
	if err != nil {
		fmt.Fprintf(w, "OAuth Error:%v", err)
	}
	client := conf.Client(ctx, tok)

	session := fb.Session{
		Version:    "v2.8",
		HttpClient: client,
	}

	res, err := session.Get("/me?fields=id,name,email", nil)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "Your are logined as : %s", res["id"])
	fmt.Fprintf(w, "Your are logined as : %s", res["name"])
	fmt.Fprintf(w, "Your are logined as : %s", res["email"])

}

//Oauth認証（ログイン）(Facebook)
func (a *App) LoginFacebookAccount(w http.ResponseWriter, r *http.Request) {

}

//Oauth認証（作成）(Github)
func (a *App) ExternalAuthenticationGithub(w http.ResponseWriter, r *http.Request) {
	conf := GithubGetConnect()
	code := r.URL.Query()["code"]
	if code == nil || len(code) == 0 {
		fmt.Fprint(w, "Invalid Parameter")
	}
	ctx := context.Background()
	tok, err := conf.Exchange(ctx, code[0])
	if err != nil {
		fmt.Fprintf(w, "OAuth Error:%v", err)
	}
	client := conf.Client(ctx, tok)
	c := github.NewClient(client)

	user, _, err := c.Users.Get(ctx, "")
	if err != nil {
		fmt.Println("aaa")
	}
	fmt.Fprintf(w, "Your are logined as : %s", *user.Name)
	fmt.Fprintf(w, "Your are logined as : %s", *user.ID)
	fmt.Fprintf(w, "Your are logined as : %s", *user.Email)
}

//Oauth認証(ログイン)(Github)
func (a *App) LoginGithubAccount(w http.ResponseWriter, r *http.Request) {

}

//アカウント情報変更
func (a *App) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	if ok, cookies := ConfirmationCookie(w, r); ok {
		user := UserAccount{
			UserName: r.FormValue(F_user_name),
			Mail:     r.FormValue(F_email),
		}
		if r.FormValue(F_password) != "" {
			user.HashPassword = HashFiled(r.FormValue(F_password))
		}

		if ok, tmp := a.DB.UpdateUserAccount(cookies.Image, user); ok {
			//クッキー作成
			if ok, _ := CreateCookie(w, r, tmp); ok {
				ReadTemplate(w, PageLogin, Show_home, nil)
			}

		} else {
			ErrorHandling(nil, nil)
		}
	}
}

//画像変更
func (a *App) ImageUpload(w http.ResponseWriter, r *http.Request) {
	if ok, cookies := ConfirmationCookie(w, r); ok {
		//画像ファイル取得
		file, fileHeader, err := r.FormFile(F_image)
		if err != nil {
			ErrorHandling(err, nil)
		}
		user := UserAccount{
			ProfileImg: BucketName + fileHeader.Filename,
		}

		//データベース更新
		if ok, user = a.DB.UpdateUserAccount(cookies.Image, user); ok {

			//画像アップロード
			err := a.ST.UploadImg(file, fileHeader)
			if err != nil {
				//TODO:エラー処理は今後修正
			}

			//クッキー作成
			if ok, _ := CreateCookie(w, r, user); ok {
				ReadTemplate(w, PageHome, Show_home, nil)
			}

		} else {
			ErrorHandling(nil, nil)
		}
	}
}

//アカウント削除
func (a *App) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	if ok, cookies := ConfirmationCookie(w, r); ok {
		//アカウント、クッキー削除
		if a.DB.DeleteUserAccount(cookies.Image) && DiscardCookie(w, r) {
			ReadTemplate(w, PageLogin, Show_login, Succes_account_delete_message)

		} else {
			ErrorHandling(nil, nil)
		}
	}
}

/*問題集関係ハンドラ*/
//問題集作成
func (a *App) CreateWorkBook(w http.ResponseWriter, r *http.Request) {
	//クッキー確認
	if ok, cookies := ConfirmationCookie(w, r); ok {
		//from全情報取得
		if err := r.ParseForm(); err != nil {
			ErrorHandling(err, nil)
		}

		//
		workbook := WorkbookContent{}
		workbook.Contents = make([]Content, 0)
		option := Option{}
		//なぜかエラーハンドリングができない
		workbook.UserId, _ = strconv.ParseInt(cookies.UserID.Value, 10, 64)

		//
		for k, v := range r.Form {
			switch k {
			case S_title:
				workbook.Title = v[0]
			case S_numberOfQuestions:
				option.NumberOfQuestions = v[0]
			case S_shuffle:
				option.Shuffle, _ = strconv.ParseBool(v[0])
			default:
				content := Content{
					ProblemNumber:    v[0],
					ProblemStatement: v[1],
					Choice1:          v[2],
					Choice2:          v[3],
					Choice3:          v[4],
					Choice4:          v[5],
					Answer:           v[6],
					Explanation:      v[7],
				}
				workbook.Contents = append(workbook.Contents, content)
			}
		}

		//問題集保存
		if a.DB.CreateWorkbook(workbook) {
			ReadTemplate(w, PageHome, Show_home, nil)
		} else {
			ErrorHandling(nil, nil)
		}
	}
}

//WorkbookUpload は問題集を共有（カインド）に登録する関数
func (a *App) WorkbookUpload(w http.ResponseWriter, r *http.Request) {
	if ok, _ := ConfirmationCookie(w, r); ok {
		bookId := r.FormValue(F_book_id)
		if a.DB.InsertWorkbookShare(bookId) {

		}
	}
}

//ReadTemplate　は
func ReadTemplate(w http.ResponseWriter, files []string, showPage string, date interface{}) {
	//
	t, err := template.ParseFiles(files...)
	if err != nil {
		ErrorLogOutput(err)
		return
	}
	if err := t.ExecuteTemplate(w, showPage, date); err != nil {
		ErrorLogOutput(err)
	}
}

// ReadTemplateToIncludeFunction は
func ReadTemplateToIncludeFunction(w http.ResponseWriter, files []string, showPage string, date interface{}, funcMap template.FuncMap) {
	//
	t, err := template.New(showPage).Funcs(funcMap).ParseFiles(files...)
	if err != nil {
		ErrorLogOutput(err)
		return
	}
	if err := t.ExecuteTemplate(w, showPage, date); err != nil {
		ErrorLogOutput(err)
	}
}

//ConfirmationCookie は
func ConfirmationCookie(w http.ResponseWriter, r *http.Request) (bool, Cookies) {
	//
	cookie1, err := r.Cookie(F_user_name)
	if err != nil {
		return false, Cookies{}
	}
	cookie2, _ := r.Cookie(F_image)
	cookie3, _ := r.Cookie(F_user_id)

	//
	c := Cookies{
		UserName: cookie1,
		Image:    cookie2,
		UserID:   cookie3,
	}

	return true, c
}

//CreateCookie は
func CreateCookie(w http.ResponseWriter, r *http.Request, user UserAccount) (bool, Cookies) {
	//
	cookie1 := http.Cookie{
		Name:  F_user_name,
		Value: user.UserName,
	}
	cookie2 := http.Cookie{
		Name:  F_image,
		Value: user.ProfileImg,
	}
	cookie3 := http.Cookie{
		Name:  F_user_id,
		Value: strconv.Itoa(int(user.UserId)),
	}

	//
	http.SetCookie(w, &cookie1)
	http.SetCookie(w, &cookie2)
	http.SetCookie(w, &cookie3)
	c := Cookies{
		UserName: &cookie1,
		Image:    &cookie2,
		UserID:   &cookie3,
	}

	return true, c
}

//DiscardCookie は
func DiscardCookie(w http.ResponseWriter, r *http.Request) bool {
	//
	cookie1, err := r.Cookie(F_user_name)
	if err != nil {
		ErrorHandling(err, nil)
		return false
	}
	cookie2, _ := r.Cookie(F_image)
	cookie3, _ := r.Cookie(F_user_id)

	//
	cookie1.MaxAge = -1
	cookie2.MaxAge = -1
	cookie3.MaxAge = -1

	//
	http.SetCookie(w, cookie1)
	http.SetCookie(w, cookie2)
	http.SetCookie(w, cookie3)

	return true
}

//SendMail は
func (a *App) SendReissueEmail(w http.ResponseWriter, r *http.Request) {
	//
	searchMail := r.FormValue(F_email)
	if ok, userId := a.DB.SelectAccountMail(searchMail); !ok {
		ErrorHandling(nil, w)
		return

	} else {
		if err := gmailSend(searchMail); err != nil {
			ErrorHandling(err, nil)
		} else {
			cookie1 := http.Cookie{
				Name:  "email",
				Value: searchMail,
			}
			cookie2 := http.Cookie{
				Name:  F_user_id,
				Value: strconv.Itoa(int(userId)),
			}
			http.SetCookie(w, &cookie1)
			http.SetCookie(w, &cookie2)
			ReadTemplate(w, PageLogin, Show_login, "メールを送信しました！")
		}
	}
}

/*RecoverPasswordは*/
func (a *App) RecoverPassword(w http.ResponseWriter, r *http.Request) {
	//
	mail, err := r.Cookie("mail")
	if err != nil {
		ErrorHandling(err, w)
		return
	}

	//
	password := r.FormValue(F_password)
	userId, err := r.Cookie(F_user_id)
	if err != nil {
		ErrorHandling(err, w)
		return
	}

	//
	if ok := a.DB.UpdateAccountPassword(mail.Value, password, userId.Value); !ok {
		ErrorHandling(nil, w)
		return
	}
	ReadTemplate(w, PageLogin, Show_login, "パスワードの再発行が完了しました")
}

/*ErrorHandlingは*/
func ErrorHandling(err interface{}, w http.ResponseWriter) {
	ErrorLogOutput(err)
	ReadTemplate(w, PageLogin, Show_home, "エラー")
}
