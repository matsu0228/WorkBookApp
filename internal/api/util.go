package api

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"log"
	"mime/multipart"
	"os"
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

//Oauth認証(Twitter)

//cloud storage クライアント作成
func NewClient(ctx context.Context) (*Client, error) {
	client, err := storage.NewClient(ctx)
	//client, err := storage.NewClient(ctx, option.WithCredentialsFile("./apptestgo0000-bef404e886bb.json"))
	if err != nil {
		return nil, err
	}
	return &Client{
		CloudStorage: client,
	}, nil
	//return nil, nil
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
			log.Println("can't close", err)
		}
	}()

	writer.ContentType = fileHeader.Header.Get("Content-Type")
	_, err := io.Copy(writer, file)

	return err
}

//文字列ハッシュ化
func HashFiled(filed string) []byte {
	hash, err := bcrypt.GenerateFromPassword([]byte(filed), 12)
	if err != nil {
		log.Fatal("")
	}
	return hash
}

//ハッシュ値と文字列比較
func CompareHashAndFiled(hash []byte, filed string) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(filed))
	if err != nil {
		return false
	}
	return true
}

//エラーハンドリング
func ErrorHandling(err interface{}) {
	file, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}

	log.Println(err)
	_, err = fmt.Fprintln(file, err)
	if err != nil {
		log.Println(err)
	}

	err = file.Close()
	if err != nil {
		log.Println(err)
	}
}
