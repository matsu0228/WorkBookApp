package pkg

import (
	"cloud.google.com/go/storage"
	"context"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/option"
	"io"
	"log"
	"mime/multipart"
	"unsafe"
)

//cloud storage クライアント作成
func NewClient(ctx context.Context) *storage.Client {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("./quiz-283808-424072cf23ec.json"))
	if err != nil {
		log.Fatal(err)
	}
	return client
}

//画像アップロード
func UploadImg(file multipart.File, fileHeader *multipart.FileHeader){
	ctx := context.Background()
	client := NewClient(ctx)
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

//画像ダウンロード(廃棄予定)
//func DownloadImg() {
//	ctx := context.Background()
//	client := NewClient(ctx)
//	defer client.Close()
//
//	// GCSオブジェクトを書き込むファイルの作成
//	f, err := os.Create("sample.png")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// オブジェクトのReaderを作成
//	bucketName := "gompei"
//	objectPath := "test"
//	obj := client.Bucket(bucketName).Object(objectPath)
//	reader, err := obj.NewReader(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer reader.Close()
//
//	// 書き込み
//	tee := io.TeeReader(reader, f)
//	s := bufio.NewScanner(tee)
//	for s.Scan() {
//	}
//	if err := s.Err(); err != nil {
//		log.Fatal(err)
//	}
//}

//ハッシュ化
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
