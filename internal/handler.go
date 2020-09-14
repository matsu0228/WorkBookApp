package internal

import (
	"Workbook/internal/pkg"
	"context"
	"fmt"
	"golang.org/x/oauth2"
	oauthapi "google.golang.org/api/oauth2/v2"
	"html/template"
	"net/http"
	"strconv"
)

/*ページ表示用ハンドラ*/
//アプリ紹介ページ
func IndexShowPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	ReadTemplate(w, pkg.PageIndex, "index", nil)
}

//ログインページ
func ShowLoginPage(w http.ResponseWriter, r *http.Request) {
	ReadTemplate(w, pkg.PageLogin, pkg.Show_login, nil)
}

//アカウント作成ページ
func ShowAccountCreatePage(w http.ResponseWriter, r *http.Request) {
	conf := pkg.GoogleGetConnect()
	GoogleUrl := conf.AuthCodeURL("yourStateUUID", oauth2.AccessTypeOffline)
	ReadTemplate(w, pkg.PageAccountCreate, pkg.Show_account_create, GoogleUrl)
}

//HOMEページ
func ShowHomePage(w http.ResponseWriter, r *http.Request) {
	if flg, _ := ConfirmationCookie(w, r); flg {
		ReadTemplate(w, pkg.PageHome, pkg.Show_home, nil)
	}
}

//設定ページ
func ShowEditPage(w http.ResponseWriter, r *http.Request) {
	if flg, _ := ConfirmationCookie(w, r); flg {
		ReadTemplate(w, pkg.PageAccountEdit, pkg.Show_home, nil)
	}
}

//問題作成ページ
func ShowWorkbookPage(w http.ResponseWriter, r *http.Request) {
	if flg, _ := ConfirmationCookie(w, r); flg {
		ReadTemplate(w, pkg.PageWorkbookCreate, pkg.Show_home, nil)
	}
}

//学習フォルダページ
func ShowWorkbookFolderPage(w http.ResponseWriter, r *http.Request) {
	if flg, cookies := ConfirmationCookie(w, r); flg {
		if flg, workbooks := pkg.SelectWorkbooks(cookies[2].Value); flg {
			ReadTemplateToIncludeFunction(w, pkg.PageWorkbookFolder, pkg.Show_home, workbooks, pkg.FuncMap)
		} else {
			pkg.ErrorHandling("")
		}
	}
}

//問題集共有ページ
func ShowWorkbookSharePage(w http.ResponseWriter, r *http.Request) {
	if flg, _ := ConfirmationCookie(w, r); flg {
		ReadTemplate(w, pkg.PageWorkbookShare, pkg.Show_home, nil)
	}
}

//問題質問ページ
func ShowWorkbookQuestion(w http.ResponseWriter, r *http.Request) {
	if flg, _ := ConfirmationCookie(w, r); flg {
		ReadTemplate(w, pkg.PageWorkbookQuestion, pkg.Show_home, nil)
	}
}

//各学習ページ
func LearningWorkbook(w http.ResponseWriter, r *http.Request) {
	if flg, _ := ConfirmationCookie(w, r); flg {
		bookId := r.FormValue(pkg.F_book_id)
		if flg, workbook := pkg.SelectWorkbook(bookId); flg {
			ReadTemplateToIncludeFunction(w, pkg.PageWorkbookLearning, pkg.Show_home, workbook, pkg.FuncMap)
		} else {
			pkg.ErrorHandling("")
		}
	}
}

//ログアウト
func Logout(w http.ResponseWriter, r *http.Request) {
	if DiscardCookie(w, r) {
		ReadTemplate(w, pkg.PageLogin, pkg.Show_login, pkg.Succes_logout_message)
	}
}

/*アカウント関係ハンドラ*/
//ログイン
func ValidateLoginData(w http.ResponseWriter, r *http.Request) {
	//from情報取得
	user := pkg.UserAccount{
		Mail: r.FormValue(pkg.F_email),
	}
	password := r.FormValue(pkg.F_password)

	//データベースに問い合わせ
	if flg, user := pkg.CheckUserLogin(user, password); flg {
		if flg, c := CreateCookie(w, r, user); flg {
			ReadTemplate(w, pkg.PageHome, pkg.Show_home, c)
		}
	} else {
		ReadTemplate(w, pkg.PageLogin, pkg.Show_login, nil)
		pkg.ErrorHandling("login")
	}
}

//アカウント作成
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	//fromからデータ取得
	user := pkg.UserAccount{
		UserName:     r.FormValue(pkg.F_user_name),
		Mail:         r.FormValue(pkg.F_email),
		HashPassword: pkg.HashFiled(r.FormValue(pkg.F_password)),
		ProfileImg:   pkg.BucketName + "no_image_square.jpg",
	}

	//アカウント作成
	if pkg.SaveUserAccount(user) {
		ReadTemplate(w, pkg.PageLogin, pkg.Show_login, pkg.Succes_account_create_message)

	} else {
		ReadTemplate(w, pkg.PageAccountCreate, pkg.Show_account_create, pkg.Error_database_message)
	}
}

//Oauth認証（登録）(Google)
func ExternalAuthenticationGoogle(w http.ResponseWriter, r *http.Request) {
	conf := pkg.GoogleGetConnect()
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
func LoginGoogleAccount(w http.ResponseWriter, r *http.Request) {

}

//Oauth認証(登録)(FaceBook)
func ExternalAuthenticationFaceBook(w http.ResponseWriter, r *http.Request) {

}

//Oauth認証（ログイン）(Facebook)
func LoginFacebookAccount(w http.ResponseWriter, r *http.Request) {

}

//Oauth認証（作成）(Twitter)
func ExternalAuthenticationTwitter(w http.ResponseWriter, r *http.Request) {

}

//Oauth認証(ログイン)(Twitter)
func LoginTwitterAccount(w http.ResponseWriter, r *http.Request) {

}

//アカウント情報変更
func UpdateAccount(w http.ResponseWriter, r *http.Request) {
	if flg, cookies := ConfirmationCookie(w, r); flg {
		user := pkg.UserAccount{
			UserName: r.FormValue(pkg.F_user_name),
			Mail:     r.FormValue(pkg.F_email),
		}
		if r.FormValue(pkg.F_password) != "" {
			user.HashPassword = pkg.HashFiled(r.FormValue(pkg.F_password))
		}

		if flg, tmp := pkg.UpdateUserAccount(cookies[2], user); flg {
			//クッキー作成
			if flg, _ := CreateCookie(w, r, tmp); flg {
				ReadTemplate(w, pkg.PageLogin, pkg.Show_home, nil)
			}

		} else {
			pkg.ErrorHandling(nil)
		}
	}
}

//画像変更
func ImageUpload(w http.ResponseWriter, r *http.Request) {
	if flg, cookies := ConfirmationCookie(w, r); flg {
		//画像ファイル取得
		file, fileHeader, err := r.FormFile(pkg.F_image)
		if err != nil {
			pkg.ErrorHandling(err)
		}
		user := pkg.UserAccount{
			ProfileImg: pkg.BucketName + fileHeader.Filename,
		}

		//データベース更新
		if flg, user = pkg.UpdateUserAccount(cookies[2], user); flg {

			//画像アップロード
			pkg.UploadImg(file, fileHeader)

			//クッキー作成
			if flg, _ := CreateCookie(w, r, user); flg {
				ReadTemplate(w, pkg.PageHome, pkg.Show_home, nil)
			}

		} else {
			pkg.ErrorHandling(nil)
		}
	}
}

//アカウント削除
func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	if flg, cookies := ConfirmationCookie(w, r); flg {
		//アカウント、クッキー削除
		if pkg.DeleteUserAccount(cookies[2]) && DiscardCookie(w, r) {
			ReadTemplate(w, pkg.PageLogin, pkg.Show_login, pkg.Succes_account_delete_message)

		} else {
			pkg.ErrorHandling(nil)
		}
	}
}

/*問題集関係ハンドラ*/
//問題集作成
func CreateWorkBook(w http.ResponseWriter, r *http.Request) {
	//クッキー確認
	if flg, cookies := ConfirmationCookie(w, r); flg {
		//from全情報取得
		if err := r.ParseForm(); err != nil {
			pkg.ErrorHandling(err)
		}

		//
		workbook := pkg.WorkbookContent{}
		workbook.Contents = make([]pkg.Content, 0)
		option := pkg.Option{}
		//なぜかエラーハンドリングができない
		workbook.UserId, _ = strconv.ParseInt(cookies[2].Value, 10, 64)

		//
		for k, v := range r.Form {
			switch k {
			case pkg.S_title:
				workbook.Title = v[0]
			case pkg.S_numberOfQuestions:
				option.NumberOfQuestions = v[0]
			case pkg.S_shuffle:
				option.Shuffle, _ = strconv.ParseBool(v[0])
			default:
				content := pkg.Content{
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
		if pkg.CreateWorkbook(workbook) {
			ReadTemplate(w, pkg.PageHome, pkg.Show_home, nil)
		} else {
			pkg.ErrorHandling(nil)
		}
	}
}

/*utilハンドラ*/
//HTMLファイル読み込み(独自Func無し)
func ReadTemplate(w http.ResponseWriter, files []string, showPage string, date interface{}) {
	t, err := template.ParseFiles(files...)
	if err != nil {
		pkg.ErrorHandling(err)
		return
	}
	t.ExecuteTemplate(w, showPage, date)
}

//HTMLファイル読み込み
func ReadTemplateToIncludeFunction(w http.ResponseWriter, files []string, showPage string, date interface{}, funcMap template.FuncMap) {
	t, err := template.New(showPage).Funcs(funcMap).ParseFiles(files...)
	if err != nil {
		pkg.ErrorHandling(err)
		return
	}
	t.ExecuteTemplate(w, showPage, date)
}

//クッキー確認util
func ConfirmationCookie(w http.ResponseWriter, r *http.Request) (bool, []*http.Cookie) {
	cookie1, err := r.Cookie(pkg.F_user_name)
	if err != nil {
		pkg.ErrorHandling(err)
		ReadTemplate(w, pkg.PageLogin, "login", pkg.Error_cookie_cannot_confirm_message)
		return false, nil
	}
	cookie2, _ := r.Cookie(pkg.F_image)
	cookie3, _ := r.Cookie(pkg.F_user_id)
	var c pkg.Cookies
	c = append(c, cookie1, cookie2, cookie3)
	return true, c
}

//クッキー作成
func CreateCookie(w http.ResponseWriter, r *http.Request, user pkg.UserAccount) (bool, pkg.Cookies) {
	cookie1 := http.Cookie{
		Name:  pkg.F_user_name,
		Value: user.UserName,
	}
	cookie2 := http.Cookie{
		Name:  pkg.F_image,
		Value: user.ProfileImg,
	}
	cookie3 := http.Cookie{
		Name:  pkg.F_user_id,
		Value: strconv.Itoa(int(user.UserId)),
	}
	http.SetCookie(w, &cookie1)
	http.SetCookie(w, &cookie2)
	http.SetCookie(w, &cookie3)
	c := pkg.Cookies{}
	c = append(c, &cookie1, &cookie2, &cookie3)
	return true, c
}

//クッキー破棄
func DiscardCookie(w http.ResponseWriter, r *http.Request) bool {
	cookie1, err := r.Cookie(pkg.F_user_name)
	if err != nil {
		pkg.ErrorHandling(err)
		ReadTemplate(w, pkg.PageLogin, "login", pkg.Error_cookie_cannot_confirm_message)
		return false
	}
	cookie2, _ := r.Cookie(pkg.F_image)
	cookie3, _ := r.Cookie(pkg.F_user_id)
	cookie1.MaxAge = -1
	cookie2.MaxAge = -1
	cookie3.MaxAge = -1
	http.SetCookie(w, cookie1)
	http.SetCookie(w, cookie2)
	http.SetCookie(w, cookie3)
	return true
}
