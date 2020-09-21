package api

import (
	"context"
	"errors"
	"github.com/google/go-github/github"
	"github.com/huandu/facebook"
	"google.golang.org/api/oauth2/v1"
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
		Show404Page(w, r, http.StatusNotFound)
		return
	}
	ReadTemplate(w, PageIndex, Show_index, nil)
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
		ErrorHandling(err, w, r)
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
	err, _ := ConfirmationCookie(w, r)
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}
	ReadTemplate(w, PageHome, Show_home, nil)
}

//ShowEditPage　は
func (a *App) ShowEditPage(w http.ResponseWriter, r *http.Request) {
	err, _ := ConfirmationCookie(w, r)
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}
	ReadTemplate(w, PageAccountEdit, Show_home, nil)
}

//ShowWorkbookPage　は
func (a *App) ShowWorkbookPage(w http.ResponseWriter, r *http.Request) {
	err, _ := ConfirmationCookie(w, r)
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}
	ReadTemplate(w, PageWorkbookCreate, Show_home, nil)
}

//学習フォルダページ
func (a *App) ShowWorkbookFolderPage(w http.ResponseWriter, r *http.Request) {
	err, cookies := ConfirmationCookie(w, r)
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}
	//
	err, workbooks := a.DB.SelectWorkbooks(cookies.UserID.Value)
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}
	ReadTemplateToIncludeFunction(w, PageWorkbookFolder, Show_home, workbooks, FuncMap)
}

//問題集共有ページ
func (a *App) ShowWorkbookSharePage(w http.ResponseWriter, r *http.Request) {
	err, _ := ConfirmationCookie(w, r)
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}
	ReadTemplate(w, PageWorkbookShare, Show_home, nil)
}

//問題質問ページ
func (a *App) ShowWorkbookQuestion(w http.ResponseWriter, r *http.Request) {
	err, _ := ConfirmationCookie(w, r)
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}
	ReadTemplate(w, PageWorkbookQuestion, Show_home, nil)
}

//各学習ページ
func (a *App) ShowLearningWorkbook(w http.ResponseWriter, r *http.Request) {
	err, _ := ConfirmationCookie(w, r)
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}

	bookId := r.FormValue(F_book_id)
	err, workbook := a.DB.SelectWorkbook(bookId)
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}
	ReadTemplateToIncludeFunction(w, PageWorkbookLearning, Show_home, workbook, FuncMap)
}

//Show404Page
func Show404Page(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		ReadTemplate(w, PageIndex, Show_404, nil)
	} else {
		ReadTemplate(w, PageIndex, Show_index, nil)
	}
}

//Show500Page は
func Show500Page(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusInternalServerError {
		ReadTemplate(w, PageIndex, Show_500, nil)
	} else {
		ReadTemplate(w, PageIndex, Show_index, nil)
	}
}

//ログアウト
func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	err := DiscardCookie(w, r)
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}
	ReadTemplate(w, PageLogin, Show_login, Succes_logout_message)
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
	err, user := a.DB.CheckUserLogin(user, password)
	if err != nil {
		ReadTemplate(w, PageLogin, Show_login, nil)
		return
	}
	c := CreateCookies(w, r, user)
	ReadTemplate(w, PageHome, Show_home, c)

}

//アカウント作成
func (a *App) CreateAccount(w http.ResponseWriter, r *http.Request) {
	//fromからデータ取得
	hash, err := HashFiled(r.FormValue(F_password))
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}

	user := UserAccount{
		Name:         r.FormValue(F_user_name),
		Mail:         r.FormValue(F_email),
		HashPassword: hash,
		ProfileImg:   BucketName + "no_image_square.jpg",
	}

	//アカウント作成
	err = a.DB.InsertUserAccount(user)
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}

	ReadTemplate(w, PageLogin, Show_login, Succes_account_create_message)
}

//SendMail は
func (a *App) SendReissueEmail(w http.ResponseWriter, r *http.Request) {
	//
	searchMail := r.FormValue(F_email)
	err, id := a.DB.SelectUserAccountMail(searchMail)
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}
	err = gmailSend(searchMail)
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}

	cookie1 := http.Cookie{
		Name:  "email",
		Value: searchMail,
	}
	cookie2 := http.Cookie{
		Name:  F_user_id,
		Value: strconv.Itoa(int(id)),
	}
	http.SetCookie(w, &cookie1)
	http.SetCookie(w, &cookie2)
	ReadTemplate(w, PageLogin, Show_login, "メールを送信しました！")
}

/*RecoverPasswordは*/
func (a *App) RecoverPassword(w http.ResponseWriter, r *http.Request) {
	//
	_, err := r.Cookie("mail")
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}

	//
	password := r.FormValue(F_password)
	userId, err := r.Cookie(F_user_id)
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}

	//
	err = a.DB.UpdateUserAccountPassword(password, userId.Value)
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}
	ReadTemplate(w, PageLogin, Show_login, "パスワードの再発行が完了しました")
}

//ExternalAuthenticationGoogleは
func (a *App) ExternalAuthenticationGoogle(w http.ResponseWriter, r *http.Request) {
	conf := GoogleGetConnect()
	code := r.URL.Query()["code"]
	if code == nil || len(code) == 0 {
		err := errors.New("外部認証エラー")
		ErrorHandling(err, w, r)
	}
	ctx := context.Background()
	tok, err := conf.Exchange(ctx, code[0])
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}
	client := conf.Client(ctx, tok)

	svr, err := oauth2.New(client)
	ui, err := svr.Userinfo.Get().Do()
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}

	id, err := strconv.ParseInt(ui.Id, 10, 64)
	if err != nil {
		ErrorHandling(err, w, r)
	}
	user := UserAccount{
		Id:         id,
		Name:       ui.Name,
		Mail:       ui.Email,
		ProfileImg: ui.Picture,
	}

	err = a.DB.SelectOauthAccount(user, "google_account")
	if err == nil {
		c := CreateCookies(w, r, user)
		ReadTemplate(w, PageHome, Show_home, c)
		return
	}

	err = a.DB.InsertOauthAccount(user, "google_account")
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}
	c := CreateCookies(w, r, user)
	ReadTemplate(w, PageHome, Show_home, c)
}

//Oauth認証(登録)(FaceBook)
func (a *App) ExternalAuthenticationFaceBook(w http.ResponseWriter, r *http.Request) {
	conf := FacebookGetConnect()
	code := r.URL.Query()["code"]
	if code == nil || len(code) == 0 {
		err := errors.New("外部認証エラー")
		ErrorHandling(err, w, r)
		return
	}
	ctx := context.Background()
	tok, err := conf.Exchange(ctx, code[0])
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}
	client := conf.Client(ctx, tok)

	session := facebook.Session{
		Version:    "v2.8",
		HttpClient: client,
	}

	res, err := session.Get("/me?fields=id,name,email,picture", nil)
	if err != nil {
		ErrorHandling(err, w, r)
	}
	//主キーはキャストしておく
	id, ok := res["id"].(int64)
	if !ok {
		err = errors.New("インタフェースキャスト失敗")
		ErrorHandling(err, w, r)
		return
	}

	picture := res["picture"].([]interface{})[0].(map[string]interface{})["url"].(string)
	if picture == "" {
		err = errors.New("facebookユーザー画像url取得失敗")
		ErrorHandling(err, w, r)
		return
	}

	user := UserAccount{
		Id:         id,
		Name:       res["name"].(string),
		Mail:       res["email"].(string),
		ProfileImg: picture,
	}

	err = a.DB.SelectOauthAccount(user, "facebook_account")
	if err == nil {
		c := CreateCookies(w, r, user)
		ReadTemplate(w, PageHome, Show_home, c)
		return
	}

	err = a.DB.InsertOauthAccount(user, "facebook_account")
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}
	c := CreateCookies(w, r, user)
	ReadTemplate(w, PageHome, Show_home, c)
}

//Oauth認証（作成）(Github)
func (a *App) ExternalAuthenticationGithub(w http.ResponseWriter, r *http.Request) {
	conf := GithubGetConnect()
	code := r.URL.Query()["code"]
	if code == nil || len(code) == 0 {
		err := errors.New("外部認証エラー")
		ErrorHandling(err, w, r)
		return
	}
	ctx := context.Background()
	tok, err := conf.Exchange(ctx, code[0])
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}
	client := conf.Client(ctx, tok)
	c := github.NewClient(client)

	tmp, _, err := c.Users.Get(ctx, "")
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}

	user := UserAccount{
		Id:         *tmp.ID,
		Name:       *tmp.Name,
		Mail:       *tmp.Email,
		ProfileImg: *tmp.AvatarURL,
	}

	err = a.DB.SelectOauthAccount(user, "github_account")
	if err == nil {
		c := CreateCookies(w, r, user)
		ReadTemplate(w, PageHome, Show_home, c)
		return
	}

	err = a.DB.InsertOauthAccount(user, "github_account")
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}

	cookies := CreateCookies(w, r, user)
	ReadTemplate(w, PageHome, Show_home, cookies)
}

//アカウント情報変更
func (a *App) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	err, cookies := ConfirmationCookie(w, r)
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}

	user := UserAccount{
		Name: r.FormValue(F_user_name),
		Mail: r.FormValue(F_email),
	}
	if r.FormValue(F_password) != "" {
		hash, err := HashFiled(r.FormValue(F_password))
		user.HashPassword = hash
		if err != nil {
			ErrorHandling(err, w, r)
			return
		}
	}

	err, tmp := a.DB.UpdateUserAccount(cookies.Image, user)
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}

	//クッキー作成
	CreateCookies(w, r, tmp)
	ReadTemplate(w, PageLogin, Show_home, nil)
}

//画像変更
func (a *App) ImageUpload(w http.ResponseWriter, r *http.Request) {
	err, cookies := ConfirmationCookie(w, r)
	if err != nil {
		ErrorHandling(err, w, r)
	}

	//画像ファイル取得
	file, fileHeader, err := r.FormFile(F_image)
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}
	user := UserAccount{
		ProfileImg: BucketName + fileHeader.Filename,
	}

	err, user = a.DB.UpdateUserAccount(cookies.Image, user)
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}

	//画像アップロード
	err = a.ST.UploadImg(file, fileHeader)
	if err != nil {
		ErrorHandling(err, w, r)
	}

	//クッキー作成
	CreateCookies(w, r, user)
	ReadTemplate(w, PageHome, Show_home, nil)

}

//アカウント削除
func (a *App) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	err, cookies := ConfirmationCookie(w, r)
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}

	err = a.DB.DeleteUserAccount(cookies.Image)
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}
	err = DiscardCookie(w, r)
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}

	ReadTemplate(w, PageLogin, Show_login, Succes_account_delete_message)
}

/*問題集関係ハンドラ*/
//問題集作成
func (a *App) CreateWorkBook(w http.ResponseWriter, r *http.Request) {
	err, cookies := ConfirmationCookie(w, r)
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}

	//from全情報取得
	err = r.ParseForm()
	if err != nil {
		ErrorHandling(err, w, r)
	}

	//
	workbook := WorkbookContent{}
	workbook.Contents = make([]Content, 0)
	option := Option{}
	workbook.UserId, err = strconv.ParseInt(cookies.UserID.Value, 10, 64)
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}

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

	err = a.DB.InsertWorkbook(workbook)
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}
	ReadTemplate(w, PageHome, Show_home, nil)
}

//WorkbookUpload は問題集を共有（カインド）に登録する関数
func (a *App) WorkbookUpload(w http.ResponseWriter, r *http.Request) {
	err, _ := ConfirmationCookie(w, r)
	if err != nil {
		ErrorHandling(err, w, r)
	}

	bookId := r.FormValue(F_book_id)
	err = a.DB.InsertShareWorkbook(bookId)
	if err != nil {
		ErrorHandling(err, w, r)
		return
	}
	ReadTemplate(w, PageHome, Show_home, nil)
}
