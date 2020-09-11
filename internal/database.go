package internal

import (
	"cloud.google.com/go/datastore"
	"context"
	"google.golang.org/api/iterator"
	"net/http"
	"reflect"
	"strconv"
)

//クライアント作成
func NewClient(ctx context.Context) (*datastore.Client, bool) {
	client, err := datastore.NewClient(ctx, project_id)
	if err != nil {
		return nil, false
	}
	return client, true
}

//ログインチェック
func CheckUserLogin(user UserAccount, password string) (bool, UserAccount) {
	ctx := context.Background()
	client, flg := NewClient(ctx)
	if flg == false {
		return false, user
	}
	defer client.Close()

	query := datastore.NewQuery("user_account").Filter("Mail =", user.Mail)
	it := client.Run(ctx, query)
	var tmp UserAccount
	_, err := it.Next(&tmp)
	if err == nil {
		if CompareHashAndFiled(tmp.HashPassword, password) {
			return true, tmp
		}
		return false, user
	} else {
		return false, user
	}
}

//アカウント作成
func SaveUserAccount(user UserAccount) bool {
	ctx := context.Background()
	client, flg := NewClient(ctx)
	if flg == false {
		return false
	}
	defer client.Close()

	var keys []*datastore.Key
	key := datastore.IncompleteKey("user_account", nil)
	keys = append(keys, key)
	//有効なキーの取得
	keys, err := client.AllocateIDs(ctx, keys)
	if err != nil {
		//datastoreキー作成エラー
		return false
	} else {
		user.UserId = keys[0].ID
		key.ID = keys[0].ID
		_, err = client.Put(ctx, key, &user)
		if err != nil {
			//datastore格納エラー
			return false
		}
		return true
	}
}

//アカウント更新
func UpdateUserAccount(cookie *http.Cookie, updateAccount UserAccount) (bool, UserAccount) {
	ctx := context.Background()
	client, flg := NewClient(ctx)
	if flg == false {
		return false, updateAccount
	}
	defer client.Close()

	//テスト
	num, _ := strconv.ParseInt(cookie.Value, 10, 64)
	Key := datastore.IDKey("user_account", num, nil)
	tx, err := client.NewTransaction(ctx)
	if err != nil {
		return false, UserAccount{}
	}

	var tmp UserAccount
	if err := tx.Get(Key, &tmp); err != nil {
		return false,tmp
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

//アカウント削除
func DeleteUserAccount(cookie *http.Cookie) bool {
	ctx := context.Background()
	client, flg := NewClient(ctx)
	if flg == false {
		return false
	}
	defer client.Close()

	key := datastore.NameKey("user_account", cookie.Value, nil)
	err := client.Delete(ctx, key)
	if err != nil {
		return false
	}
	return true
}

//アカウント検索
func SelectAccount(user UserAccount) bool {
	ctx := context.Background()
	client, flg := NewClient(ctx)
	if flg == false {
		return false
	}
	defer client.Close()

	query := datastore.NewQuery("user_account").
		Filter("UserId =", user.UserId).
		Filter("Mail =", user.Mail)

	it := client.Run(ctx, query)
	var tmp UserAccount
	_, err := it.Next(&tmp)
	if err != nil {
		return false
	}
	return true
}

//問題集作成
func CreateWorkbook(book WorkbookContent) bool {
	//クライアント作成
	ctx := context.Background()
	client, flg := NewClient(ctx)
	if flg == false {
		return false
	}
	defer client.Close()

	//キー作成
	var keys []*datastore.Key
	parentKey := datastore.IDKey("user_account", book.UserId, nil)
	childKey := datastore.IncompleteKey("workbook", parentKey)
	keys = append(keys, childKey)

	//有効なキーの取得
	keys, err := client.AllocateIDs(ctx, keys)
	if err != nil {
		return false
	} else {
		book.BookId = keys[0].ID
		childKey.ID = keys[0].ID

		//book格納
		_, err := client.Put(ctx, childKey, &book)
		if err != nil {
			return false
		}
		return true
	}
}

//問題集検索
func SelectWorkbooks(id string) (bool, []WorkbookContent) {
	ctx := context.Background()
	client, flg := NewClient(ctx)
	if flg == false {
		return false, nil
	}
	defer client.Close()

	IntId, _ := strconv.Atoi(id)
	query := datastore.NewQuery("workbook").
		Filter("UserId =", IntId)
	var workbooks []WorkbookContent
	it := client.Run(ctx, query)
	for {
		var tmp WorkbookContent
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

func SelectWorkbook(id string) (bool, WorkbookContent) {
	ctx := context.Background()
	var workbook WorkbookContent
	client, flg := NewClient(ctx)
	if flg == false {
		return false, workbook
	}
	defer client.Close()

	bookId, _ := strconv.Atoi(id)
	query := datastore.NewQuery("workbook").
		Filter("BookId =", bookId)
	it := client.Run(ctx, query)
	_, err := it.Next(&workbook)
	if err != nil {
		return false, workbook
	}
	return true, workbook
}