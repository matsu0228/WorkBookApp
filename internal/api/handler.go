package api

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	oauthapi "google.golang.org/api/oauth2/v2"
	"html/template"
	"net/http"
	"strconv"
)

//NewAppは
func NewApp(d Repository, s Storage) *App {
	return &App{
		DB: d,
		ST: s,
	}
}

//IndexShowPage　は
func (a *App) IndexShowPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	ReadTemplate(w, PageIndex, "index", nil)
}

//ShowLoginPage は
func (a *App) ShowLoginPage(w http.ResponseWriter, r *http.Request) {
	ReadTemplate(w, PageLogin, Show_login, nil)
}

//ShowAccountCreatePage　は
func (a *App) ShowAccountCreatePage(w http.ResponseWriter, r *http.Request) {
	conf := GoogleGetConnect()
	GoogleUrl := conf.AuthCodeURL("yourStateUUID", oauth2.AccessTypeOffline)
	ReadTemplate(w, PageAccountCreate, Show_account_create, GoogleUrl)
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
			ErrorHandling("")
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
			ErrorHandling("")
		}
	}
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
		ErrorHandling("login")
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

}

//Oauth認証（ログイン）(Facebook)
func (a *App) LoginFacebookAccount(w http.ResponseWriter, r *http.Request) {

}

//Oauth認証（作成）(Twitter)
func (a *App) ExternalAuthenticationTwitter(w http.ResponseWriter, r *http.Request) {

}

//Oauth認証(ログイン)(Twitter)
func (a *App) LoginTwitterAccount(w http.ResponseWriter, r *http.Request) {

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
			ErrorHandling(nil)
		}
	}
}

//画像変更
func (a *App) ImageUpload(w http.ResponseWriter, r *http.Request) {
	if ok, cookies := ConfirmationCookie(w, r); ok {
		//画像ファイル取得
		file, fileHeader, err := r.FormFile(F_image)
		if err != nil {
			ErrorHandling(err)
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
			ErrorHandling(nil)
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
			ErrorHandling(nil)
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
			ErrorHandling(err)
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
			ErrorHandling(nil)
		}
	}
}

//ReadTemplate　は
func ReadTemplate(w http.ResponseWriter, files []string, showPage string, date interface{}) {
	//
	t, err := template.ParseFiles(files...)
	if err != nil {
		ErrorHandling(err)
		return
	}
	t.ExecuteTemplate(w, showPage, date)
}

// ReadTemplateToIncludeFunction は
func ReadTemplateToIncludeFunction(w http.ResponseWriter, files []string, showPage string, date interface{}, funcMap template.FuncMap) {
	//
	t, err := template.New(showPage).Funcs(funcMap).ParseFiles(files...)
	if err != nil {
		return
	}
	t.ExecuteTemplate(w, showPage, date)
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
		ErrorHandling(err)
		ReadTemplate(w, PageLogin, "login", Error_cookie_cannot_confirm_message)
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
