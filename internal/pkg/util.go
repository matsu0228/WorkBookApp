package pkg

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"io"
	"log"
	"mime/multipart"
	"os"
	"unsafe"
)

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
func CloudStoreNewClient(ctx context.Context) *storage.Client {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("./quiz-283808-424072cf23ec.json"))
	if err != nil {
		log.Fatal(err)
	}
	return client
}

//画像アップロード（CloudStore）
func UploadImg(file multipart.File, fileHeader *multipart.FileHeader) {
	ctx := context.Background()
	client := CloudStoreNewClient(ctx)
	defer client.Close()

	// オブジェクトのReaderを作成
	bucketName := "gompei"
	objectName := fileHeader.Filename

	//
	writer := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	writer.ContentType = fileHeader.Header.Get("Content-Type")
	_, _ = io.Copy(writer, file)
	_ = writer.Close()
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

//バイト配列　→　文字列変換(多用禁止)
func Benchmark_Unsafe(data []byte) string {
	hashMail := *(*string)(unsafe.Pointer(&data))
	return hashMail
}

//エラーハンドリング
func ErrorHandling(err interface{}) {
	file, _ := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer file.Close()

	log.Println(err)
	fmt.Fprintln(file, err)
}
