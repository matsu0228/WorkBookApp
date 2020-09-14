package database

import (
	"WorkBookApp/internal/api"
	"cloud.google.com/go/datastore"
	"context"
	"google.golang.org/api/iterator"
	"net/http"
	"reflect"
	"strconv"
)

//クライアント
var Client *datastore.Client

/*DataStore関連関数*/
//NewClient　はDataStoreのクライアントを生成する関数
func DataStoreNewClient(ctx context.Context) (*datastore.Client, error) {
	var err error
	Client, err = datastore.NewClient(ctx, api.Project_id)
	if err != nil {
		return nil, err
	}
	return Client, nil
}

//CheckUserLogin はメールアドレスとパスワードを比較して、booleanとユーザアカウント情報を返す関数
func CheckUserLogin(user api.UserAccount, password string) (bool, api.UserAccount) {
	ctx := context.Background()
	query := datastore.NewQuery("user_account").Filter("Mail =", user.Mail)
	it := Client.Run(ctx, query)
	var tmp api.UserAccount
	_, err := it.Next(&tmp)
	if err == nil {
		if api.CompareHashAndFiled(tmp.HashPassword, password) {
			return true, tmp
		}
		return false, user
	} else {
		return false, user
	}
}

//SaveUserAccount はユーザIDを主キーにして、データを登録する関数
func SaveUserAccount(user api.UserAccount) bool {
	ctx := context.Background()

	var keys []*datastore.Key
	key := datastore.IncompleteKey("user_account", nil)
	keys = append(keys, key)
	//有効なキーの取得
	keys, err := Client.AllocateIDs(ctx, keys)
	if err != nil {
		//datastoreキー作成エラー
		return false
	} else {
		user.UserId = keys[0].ID
		key.ID = keys[0].ID
		_, err = Client.Put(ctx, key, &user)
		if err != nil {
			//datastore格納エラー
			return false
		}
		return true
	}
}

//UpdateUserAccount　は登録されているユーザ情報を更新する関数
func UpdateUserAccount(cookie *http.Cookie, updateAccount api.UserAccount) (bool, api.UserAccount) {
	ctx := context.Background()

	//テスト
	num, _ := strconv.ParseInt(cookie.Value, 10, 64)
	Key := datastore.IDKey("user_account", num, nil)
	tx, err := Client.NewTransaction(ctx)
	if err != nil {
		return false, api.UserAccount{}
	}

	var tmp api.UserAccount
	if err := tx.Get(Key, &tmp); err != nil {
		return false, tmp
	}

	rv := reflect.ValueOf(updateAccount)
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		filed := rt.Field(i)
		if rv.FieldByName(filed.Name).Interface() != "" {
			switch filed.Name {
			case "UserName":
				tmp.UserName = rv.FieldByName(filed.Name).Interface().(string)
			case "Mail":
				tmp.Mail = rv.FieldByName(filed.Name).Interface().(string)
			case "HashPassword":
				if len(rv.FieldByName(filed.Name).Interface().([]byte)) > 0 {
					tmp.HashPassword = rv.FieldByName(filed.Name).Interface().([]byte)
				}
			case "ProfileImg":
				tmp.ProfileImg = rv.FieldByName(filed.Name).Interface().(string)
			}
		}
	}
	if _, err := tx.Put(Key, &tmp); err != nil {
		return false, updateAccount
	}
	if _, err := tx.Commit(); err != nil {
		return false, updateAccount
	}

	return true, tmp
}

//DeleteUserAccount　は登録しているユーザ情報を削除する関数
func DeleteUserAccount(cookie *http.Cookie) bool {
	ctx := context.Background()

	key := datastore.NameKey("user_account", cookie.Value, nil)
	err := Client.Delete(ctx, key)
	if err != nil {
		return false
	}
	return true
}

//CreateWorkbook　は４択問題集をbookIDを主キーにデータを登録する関数
func CreateWorkbook(book api.WorkbookContent) bool {
	//クライアント作成
	ctx := context.Background()

	//キー作成
	var keys []*datastore.Key
	parentKey := datastore.IDKey("user_account", book.UserId, nil)
	childKey := datastore.IncompleteKey("workbook", parentKey)
	keys = append(keys, childKey)

	//有効なキーの取得
	keys, err := Client.AllocateIDs(ctx, keys)
	if err != nil {
		return false
	} else {
		book.BookId = keys[0].ID
		childKey.ID = keys[0].ID

		//book格納
		_, err := Client.Put(ctx, childKey, &book)
		if err != nil {
			return false
		}
		return true
	}
}

//SelectWorkbooks は問題集のタイトル,IDを取得して,boolean,構造体の配列を返す関数
func SelectWorkbooks(id string) (bool, []api.WorkbookContent) {
	Context := context.Background()

	IntId, _ := strconv.Atoi(id)
	query := datastore.NewQuery("workbook").
		Filter("UserId =", IntId)
	var workbooks []api.WorkbookContent
	it := Client.Run(Context, query)
	for {
		var tmp api.WorkbookContent
		_, err := it.Next(&tmp)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return false, nil
		}
		workbooks = append(workbooks, tmp)
	}
	return true, workbooks
}

//TODO:記載予定
func SelectWorkbook(id string) (bool, api.WorkbookContent) {
	Context := context.Background()
	var workbook api.WorkbookContent

	bookId, _ := strconv.Atoi(id)
	query := datastore.NewQuery("workbook").
		Filter("BookId =", bookId)
	it := Client.Run(Context, query)
	_, err := it.Next(&workbook)
	if err != nil {
		return false, workbook
	}
	return true, workbook
}
