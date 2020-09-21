package api

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

/*GoogleGetConnectはOauth認証に必要な設定情報を生成して、返す関数*/
func GoogleGetConnect() *oauth2.Config {
	config := &oauth2.Config{
		ClientID:     "xxxx",
		ClientSecret: "xxxx",
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:8080/login/google",
	}
	return config
}

//Oauth認証(FaceBook)
func FacebookGetConnect() *oauth2.Config {
	config := &oauth2.Config{
		ClientID:     "xxxx",
		ClientSecret: "xxxx",
		Scopes:       []string{"email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.facebook.com/dialog/oauth",
			TokenURL: "https://graph.facebook.com/oauth/access_token",
		},
		RedirectURL: "http://localhost:8080/login/facebook",
	}
	return config
}

//Oauth認証(github)
func GithubGetConnect() *oauth2.Config {
	config := &oauth2.Config{
		ClientID:     "xxxx",
		ClientSecret: "xxxx",
		Scopes:       []string{"email"},
		Endpoint:     github.Endpoint,
		RedirectURL:  "http://localhost:8080/login/github",
	}
	return config
}

/*CreationUrlは*/
func CreationUrl() map[string]string {
	c := GoogleGetConnect()
	googleUrl := c.AuthCodeURL(oauthStateString, oauth2.AccessTypeOffline)
	c = FacebookGetConnect()
	facebookUrl := c.AuthCodeURL(oauthStateString, oauth2.AccessTypeOffline)
	c = GithubGetConnect()
	githubUrl := c.AuthCodeURL(oauthStateString, oauth2.AccessTypeOffline)
	return map[string]string{"google": googleUrl, "facebook": facebookUrl, "github": githubUrl}
}

//NewClient はCloudStorageのクライアントを生成して、返す関数
func NewClient(ctx context.Context) (*Client, error) {
	//client, err := storage.NewClient(ctx)
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("./apptestgo0000-bef404e886bb.json"))
	if err != nil {
		return nil, err
	}
	return &Client{
		CloudStorage: client,
	}, nil
}

//画像アップロード（CloudStore）
func (c *Client) UploadImg(file multipart.File, fileHeader *multipart.FileHeader) error {
	ctx := context.Background()

	// オブジェクトのReaderを作成
	bucketName := "gompei"
	objectName := fileHeader.Filename

	writer := c.CloudStorage.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	defer writer.Close()

	writer.ContentType = fileHeader.Header.Get("Content-Type")
	_, err := io.Copy(writer, file)

	return err
}

/*HashFiledは
パスワード（文字列）をハッシュ化して、バイト配列とエラーを返す関数*/
func HashFiled(filed string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(filed), 12)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

/*CompareHashAndFiledは
ハッシュ値とパスワード（文字列）を比較して、エラーを返す関数*/
func CompareHashAndFiled(hash []byte, filed string) error {
	err := bcrypt.CompareHashAndPassword(hash, []byte(filed))
	if err != nil {
		return err
	}
	return nil
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
		return
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
		return
	}
}

//ConfirmationCookie は
func ConfirmationCookie(w http.ResponseWriter, r *http.Request) (error, Cookies) {
	//
	cookie1, err := r.Cookie(F_user_name)
	if err != nil {
		return err, Cookies{}
	}
	cookie2, err := r.Cookie(F_image)
	if err != nil {
		return err, Cookies{}
	}
	cookie3, err := r.Cookie(F_user_id)
	if err != nil {
		return err, Cookies{}
	}

	//
	c := Cookies{
		UserName: cookie1,
		Image:    cookie2,
		UserID:   cookie3,
	}

	return nil, c
}

//CreateCookie は
func CreateCookies(w http.ResponseWriter, r *http.Request, user UserAccount) Cookies {
	//
	cookie1 := http.Cookie{
		Name:  F_user_name,
		Value: user.Name,
	}
	cookie2 := http.Cookie{
		Name:  F_image,
		Value: user.ProfileImg,
	}
	cookie3 := http.Cookie{
		Name:  F_user_id,
		Value: strconv.Itoa(int(user.Id)),
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

	return c
}

//DiscardCookie は
func DiscardCookie(w http.ResponseWriter, r *http.Request) error {
	//
	cookie1, err := r.Cookie(F_user_name)
	if err != nil {
		return err
	}
	cookie2, err := r.Cookie(F_image)
	if err != nil {
		return err
	}
	cookie3, err := r.Cookie(F_user_id)
	if err != nil {
		return err
	}

	//
	cookie1.MaxAge = -1
	cookie2.MaxAge = -1
	cookie3.MaxAge = -1

	//
	http.SetCookie(w, cookie1)
	http.SetCookie(w, cookie2)
	http.SetCookie(w, cookie3)

	return nil
}

/*ErrorHandlingは
log.txtに現在日付とエラー内容を書き出す関数*/
func ErrorLogOutput(error interface{}) {
	//エラーログファイルオープン
	file, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}

	//エラー内容書き込み
	_, err = fmt.Fprint(file, time.Now())
	if err != nil {
		log.Println(err)
	}
	_, err = fmt.Fprintln(file, error)
	if err != nil {
		log.Println(err)
	}

	//ファイルを閉じる
	err = file.Close()
	if err != nil {
		log.Println(err)
	}
}

/*ErrorHandlingは*/
func ErrorHandling(err interface{}, w http.ResponseWriter, r *http.Request) {
	ErrorLogOutput(err)
	Show500Page(w, r, http.StatusInternalServerError)
}
