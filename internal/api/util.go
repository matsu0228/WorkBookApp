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
	"io"
	"log"
	"mime/multipart"
	"os"
	"time"
)

//クライアント
type Client struct {
	CloudStorage *storage.Client
}

//Oauth認証(Google)
func GoogleGetConnect() *oauth2.Config {
	config := &oauth2.Config{
		ClientID:     "907386733410-poe2om4820040g0vbuug7iiin6jajfjt.apps.googleusercontent.com",
		ClientSecret: "My6OndPCv4v6lx2igjxpVrZV",
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:8080/account/create/google",
	}
	return config
}

//Oauth認証(FaceBook)
func FacebookGetConnect() *oauth2.Config {
	config := &oauth2.Config{
		ClientID:     "641206823488638",
		ClientSecret: "f3053b7fe1d41fe7acbb682268d2ed02",
		Scopes:       []string{"email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.facebook.com/dialog/oauth",
			TokenURL: "https://graph.facebook.com/oauth/access_token",
		},
		RedirectURL: "http://localhost:8080/account/create/facebook",
	}
	return config
}

//Oauth認証(github)
func GithubGetConnect() *oauth2.Config {
	config := &oauth2.Config{
		ClientID:     "Iv1.503404aa1ff2e2e6",
		ClientSecret: "112aefef6df4e5bd8d48c675390f6d2caad3b124",
		Scopes:       []string{"email"},
		Endpoint:     github.Endpoint,
		RedirectURL:  "http://localhost:8080/account/create/github",
	}
	return config
}

/*CreationUrl()*/
func CreationUrl() map[string]string {
	c := GoogleGetConnect()
	GoogleUrl := c.AuthCodeURL(oauthStateString, oauth2.AccessTypeOffline)
	c = FacebookGetConnect()
	FacebookUrl := c.AuthCodeURL(oauthStateString, oauth2.AccessTypeOffline)
	c = GithubGetConnect()
	GithubUrl := c.AuthCodeURL(oauthStateString, oauth2.AccessTypeOffline)
	return map[string]string{"google": GoogleUrl, "facebook": FacebookUrl, "github": GithubUrl}
}

//cloud storage クライアント作成
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
	defer func() {
		if err := writer.Close(); err != nil {
			log.Println(err)
		}
	}()

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
ハッシュ値とパスワード（文字列）を比較して、真偽値とエラーを返す関数*/
func CompareHashAndFiled(hash []byte, filed string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, []byte(filed))
	if err != nil {
		return false, err
	}
	return true, nil
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
