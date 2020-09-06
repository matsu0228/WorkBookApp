package internal

import (
	"Workbook/internal/pkg"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

/*ページ表示用ハンドラ*/
//ログインページ
func ShowLoginPage(w http.ResponseWriter, r *http.Request) {
	ReadTemplate(w, pageLogin, show_login, nil)
}

//アカウント作成ページ
func ShowAccountCreatePage(w http.ResponseWriter, r *http.Request) {
	ReadTemplate(w, pageAccountCreate, show_account_create, nil)
}

//HOMEページ
func ShowHomePage(w http.ResponseWriter, r *http.Request) {
	if flg, _ := ConfirmationCookie(w, r); flg {
		ReadTemplate(w, pageHome, show_home, nil)
	}
}

//設定ページ
func ShowEditPage(w http.ResponseWriter, r *http.Request) {
	if flg, _ := ConfirmationCookie(w, r); flg {
		ReadTemplate(w, pageAccountEdit, show_home, nil)
	}
}

//問題作成ページ
func ShowWorkbookPage(w http.ResponseWriter, r *http.Request) {
	if flg, _ := ConfirmationCookie(w, r); flg {
		ReadTemplate(w, pageWorkbookCreate, show_home, nil)
	}
}

//学習フォルダページ
func ShowWorkbookFolderPage(w http.ResponseWriter, r *http.Request) {
	if flg, cookies := ConfirmationCookie(w, r); flg {
		if flg, workbooks := SelectWorkbooks(cookies[2].Value); flg {
			ReadTemplate(w, pageWorkbookFolder, show_home, workbooks)
		} else {
			ErrorHandling(nil)
		}
	}
}

//各学習ページ
func LearningWorkbook(w http.ResponseWriter, r *http.Request) {
	if flg, _ := ConfirmationCookie(w, r); flg {
		bookId := r.FormValue(f_book_id)
		if flg, workbook := SelectWorkbook(bookId); flg {
			ReadTemplateToIncludeFunction(w, pageWorkbookLearning, show_home, workbook, FuncMap)
		} else {
			ErrorHandling(nil)
		}
	}
}

//ログアウト
func Logout(w http.ResponseWriter, r *http.Request) {
	if DiscardCookie(w, r) {
		ReadTemplate(w, pageLogin, show_login, succes_logout_message)
	}
}

/*アカウント関係ハンドラ*/
//ログイン
func ValidateLoginData(w http.ResponseWriter, r *http.Request) {

	//from情報取得
	user := UserAccount{
		Mail: r.FormValue(f_email),
	}
	password := r.FormValue(f_password)

	//データベースに問い合わせ
	if flg, user := CheckUserLogin(user, password); flg {
		if flg, c := CreateCookie(w, r, user); flg {
			ReadTemplate(w, pageHome, show_home, c)
		}
	} else {
		ErrorHandling(nil)
	}
}

//作成
func CreateAccount(w http.ResponseWriter, r *http.Request) {

	//fromからデータ取得
	//最低限の入力チェックはフロントでやる
	//今後入力正規表現でチェックする
	rand.Seed(time.Now().UnixNano())
	user := UserAccount{
		UserName:     r.FormValue(f_user_name),
		Mail:         r.FormValue(f_email),
		HashPassword: pkg.HashFiled(r.FormValue(f_password)),
		ProfileImg:   BucketName + "no_image_square.jpg",
	}

	//アカウント作成
	if SaveUserAccount(user) {
		ReadTemplate(w, pageLogin, show_login, succes_account_create_message)

	} else {
		ReadTemplate(w, pageAccountCreate, show_account_create, error_database_message)
	}
}

//アカウント情報変更
func UpdateAccount(w http.ResponseWriter, r *http.Request) {
	if flg, cookies := ConfirmationCookie(w, r); flg {
		user := UserAccount{
			UserName: r.FormValue(f_user_name),
			Mail:     r.FormValue(f_email),
		}
		if r.FormValue(f_password) != "" {
			user.HashPassword = pkg.HashFiled(r.FormValue(f_password))
		}

		if flg, tmp := UpdateUserAccount(cookies[2], user); flg {
			//クッキー作成
			if flg, _ := CreateCookie(w, r, tmp); flg {
				ReadTemplate(w, pageLogin, show_home, nil)
			}

		} else {
			ErrorHandling(nil)
		}
	}
}

//画像変更
func ImageUpload(w http.ResponseWriter, r *http.Request) {
	if flg, cookies := ConfirmationCookie(w, r); flg {

		//画像ファイル取得
		file, fileHeader, err := r.FormFile(f_image)
		if err != nil {
			ErrorHandling(err)
		}
		user := UserAccount{
			ProfileImg: BucketName + fileHeader.Filename,
		}

		//データベース更新
		if flg, user = UpdateUserAccount(cookies[2], user); flg {

			//画像アップロード
			pkg.UploadImg(file, fileHeader)

			//クッキー作成
			if flg, _ := CreateCookie(w, r, user); flg {
				ReadTemplate(w, pageHome, show_home, nil)
			}

		} else {
			ErrorHandling(nil)
		}
	}
}

//アカウント削除
func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	if flg, cookies := ConfirmationCookie(w, r); flg {
		//アカウント、クッキー削除
		if DeleteUserAccount(cookies[2]) && DiscardCookie(w, r) {
			ReadTemplate(w, pageLogin, show_login, succes_account_delete_message)

		} else {
			ErrorHandling(nil)
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
			ErrorHandling(err)
		}

		//
		workbook := WorkbookContent{}
		workbook.Contents = make([]Content, 0)
		option := Option{}
		//なぜかエラーハンドリングができない
		workbook.UserId, _ = strconv.ParseInt(cookies[2].Value, 10, 64)

		//
		for k, v := range r.Form {
			switch k {
			case s_title:
				workbook.Title = v[0]
			case s_numberOfQuestions:
				option.NumberOfQuestions = v[0]
			case s_shuffle:
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
		if CreateWorkbook(workbook) {
			ReadTemplate(w, pageHome, show_home, nil)
		} else {
			ErrorHandling(nil)
		}
	}
}

/*utilハンドラ*/
//HTMLファイル読み込み(独自Func無し)
func ReadTemplate(w http.ResponseWriter, files []string, showPage string, date interface{}) {
	t, err := template.ParseFiles(files...)
	if err != nil {
		ErrorHandling(err)
	}
	t.ExecuteTemplate(w, showPage, date)
}

//HTMLファイル読み込み
func ReadTemplateToIncludeFunction(w http.ResponseWriter, files []string, showPage string, date interface{}, funcMap template.FuncMap) {
	t, err := template.New(showPage).Funcs(funcMap).ParseFiles(files...)
	if err != nil {
		ErrorHandling(err)
	}
	t.ExecuteTemplate(w, showPage, date)
}

//クッキー確認util
func ConfirmationCookie(w http.ResponseWriter, r *http.Request) (bool, []*http.Cookie) {
	cookie1, err := r.Cookie(f_user_name)
	if err != nil {
		//ReadTemplate(w, pageLogin, "login", error_cookie_cannot_confirm_message)
		ErrorHandling(err)
		return false, nil
	}
	cookie2, _ := r.Cookie(f_image)
	cookie3, _ := r.Cookie(f_book_id)
	var c Cookies
	c = append(c, cookie1, cookie2, cookie3)
	return true, c
}

//クッキー作成
func CreateCookie(w http.ResponseWriter, r *http.Request, user UserAccount) (bool, Cookies) {
	cookie1 := http.Cookie{
		Name:  f_user_name,
		Value: user.UserName,
	}
	cookie2 := http.Cookie{
		Name:  f_image,
		Value: user.ProfileImg,
	}
	cookie3 := http.Cookie{
		Name:  f_user_id,
		Value: strconv.Itoa(int(user.UserId)),
	}
	http.SetCookie(w, &cookie1)
	http.SetCookie(w, &cookie2)
	http.SetCookie(w, &cookie3)
	c := Cookies{}
	c = append(c, &cookie1, &cookie2, &cookie3)
	return true, c
}

//クッキー破棄
func DiscardCookie(w http.ResponseWriter, r *http.Request) bool {
	cookie1, err := r.Cookie(f_user_name)
	if err != nil {
		//ReadTemplate(w, pageLogin, "login", error_cookie_cannot_confirm_message)
		ErrorHandling(err)
		return false
	}
	cookie2, _ := r.Cookie(f_image)
	cookie3, _ := r.Cookie(f_user_id)
	cookie1.MaxAge = -1
	cookie2.MaxAge = -1
	cookie3.MaxAge = -1
	http.SetCookie(w, cookie1)
	http.SetCookie(w, cookie2)
	http.SetCookie(w, cookie3)
	return true
}

//エラーハンドリング
func ErrorHandling(err interface{}) {

}
