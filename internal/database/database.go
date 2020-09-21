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
type Client struct {
	DataStore *datastore.Client
}

//NewClient　はDataStoreのクライアントを生成する関数
func NewClient(ctx context.Context) (*Client, error) {
	client, err := datastore.NewClient(ctx, api.Project_id)
	//client, err := datastore.NewClient(ctx, api.Project_id, option.WithCredentialsFile("./apptestgo0000-bef404e886bb.json"))
	if err != nil {
		return nil, err
	}
	return &Client{
		DataStore: client,
	}, nil
}

//CheckUserLogin はメールアドレスとパスワードを比較して、booleanとユーザアカウント情報を返す関数
func (c *Client) CheckUserLogin(user api.UserAccount, password string) (error, api.UserAccount) {
	//
	ctx := context.Background()

	//
	query := datastore.NewQuery("user_account").Filter("Mail =", user.Mail)

	//
	it := c.DataStore.Run(ctx, query)
	var tmp api.UserAccount

	//
	_, err := it.Next(&tmp)
	if err != nil {
		return err, tmp
	}

	//
	err = api.CompareHashAndFiled(tmp.HashPassword, password)
	if err != nil {
		return err, tmp
	}

	return nil, tmp
}

//InsertUserAccount はユーザIDを主キーにして、データを登録する関数
func (c *Client) InsertUserAccount(user api.UserAccount) error {
	//
	ctx := context.Background()

	//
	var keys []*datastore.Key
	key := datastore.IncompleteKey("user_account", nil)
	keys = append(keys, key)

	//有効なキーの取得
	keys, err := c.DataStore.AllocateIDs(ctx, keys)
	if err != nil {
		return err
	}

	//
	user.Id = keys[0].ID
	key.ID = keys[0].ID
	_, err = c.DataStore.Put(ctx, key, &user)
	if err != nil {
		return err
	}

	return nil
}

/**/
func (c *Client) InsertOauthAccount(user api.UserAccount, kind string) error {
	//
	ctx := context.Background()

	//
	key := datastore.IDKey(kind, user.Id, nil)

	//
	_, err := c.DataStore.Put(ctx, key, &user)
	if err != nil {
		return err
	}

	return nil
}

/**/
func (c *Client) SelectOauthAccount(user api.UserAccount, kind string) error {
	ctx := context.Background()

	//
	query := datastore.NewQuery(kind).Filter("Id =", user.Id)

	//
	it := c.DataStore.Run(ctx, query)
	var tmp api.UserAccount

	//
	_, err := it.Next(&tmp)
	if err != nil {
		return err
	}

	return nil
}

/*SelectAccountMailは*/
func (c *Client) SelectUserAccountMail(searchTarget string) (error, int64) {
	ctx := context.Background()

	//
	query := datastore.NewQuery("user_account").Filter("Mail =", searchTarget)

	//
	it := c.DataStore.Run(ctx, query)
	var tmp api.UserAccount

	//
	_, err := it.Next(&tmp)
	if err != nil {
		return err, 0
	}

	return nil, tmp.Id
}

/*UpdateAccountPasswordは*/
func (c *Client) UpdateUserAccountPassword(password string, userId string) error {
	ctx := context.Background()

	//クッキーにあるユーザIDを元に更新
	id, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		return err
	}

	//
	Key := datastore.IDKey("user_account", id, nil)
	tx, err := c.DataStore.NewTransaction(ctx)
	if err != nil {
		return err
	}

	//
	var tmp api.UserAccount
	if err := tx.Get(Key, &tmp); err != nil {
		return err
	}

	//パスワードの詰め替え
	hash, err := api.HashFiled(password)
	if err != nil {
		return err
	}
	tmp.HashPassword = hash
	if _, err := tx.Put(Key, &tmp); err != nil {
		return err
	}
	if _, err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

//UpdateUserAccount　は登録されているユーザ情報を更新する関数
func (c *Client) UpdateUserAccount(cookie *http.Cookie, updateAccount api.UserAccount) (error, api.UserAccount) {
	ctx := context.Background()

	//クッキーにあるユーザIDを元に更新
	userId, err := strconv.ParseInt(cookie.Value, 10, 64)
	if err != nil {
		return err, api.UserAccount{}
	}

	Key := datastore.IDKey("user_account", userId, nil)
	tx, err := c.DataStore.NewTransaction(ctx)
	if err != nil {
		return err, api.UserAccount{}
	}

	var tmp api.UserAccount
	if err := tx.Get(Key, &tmp); err != nil {
		return err, tmp
	}

	rv := reflect.ValueOf(updateAccount)
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		filed := rt.Field(i)
		if rv.FieldByName(filed.Name).Interface() != "" {
			switch filed.Name {
			case "Name":
				tmp.Name = rv.FieldByName(filed.Name).Interface().(string)
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
		return err, updateAccount
	}
	if _, err := tx.Commit(); err != nil {
		return err, updateAccount
	}

	return nil, tmp
}

//DeleteUserAccount　は登録しているユーザ情報を削除する関数
func (c *Client) DeleteUserAccount(cookie *http.Cookie) error {
	//
	ctx := context.Background()

	//
	key := datastore.NameKey("user_account", cookie.Value, nil)
	err := c.DataStore.Delete(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

//CreateWorkbook　は４択問題集をbookIDを主キーにデータを登録する関数
func (c *Client) InsertWorkbook(book api.WorkbookContent) error {
	//クライアント作成
	ctx := context.Background()

	//キー作成
	var keys []*datastore.Key
	parentKey := datastore.IDKey("user_account", book.UserId, nil)
	childKey := datastore.IncompleteKey("workbook", parentKey)
	keys = append(keys, childKey)

	//有効なキーの取得
	keys, err := c.DataStore.AllocateIDs(ctx, keys)
	if err != nil {
		return err
	}

	book.BookId = keys[0].ID
	childKey.ID = keys[0].ID

	//book格納
	_, err = c.DataStore.Put(ctx, childKey, &book)
	if err != nil {
		return err
	}
	return nil

}

//SelectWorkbooks は問題集のタイトル,IDを取得して,boolean,構造体の配列を返す関数
func (c *Client) SelectWorkbooks(id string) (error, []api.WorkbookContent) {
	Context := context.Background()

	IntId, _ := strconv.Atoi(id)
	query := datastore.NewQuery("workbook").Filter("UserId =", IntId)
	var WorkbookContents []api.WorkbookContent
	it := c.DataStore.Run(Context, query)
	for {
		var tmp api.WorkbookContent
		_, err := it.Next(&tmp)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err, nil
		}
		WorkbookContents = append(WorkbookContents, tmp)
	}
	return nil, WorkbookContents
}

//SelectWorkbook は
func (c *Client) SelectWorkbook(id string) (error, api.WorkbookContent) {
	Context := context.Background()
	var workbook api.WorkbookContent

	bookId, _ := strconv.Atoi(id)
	query := datastore.NewQuery("workbook").Filter("BookId =", bookId)
	it := c.DataStore.Run(Context, query)
	_, err := it.Next(&workbook)
	if err != nil {
		return err, workbook
	}
	return nil, workbook
}

//InsertShareWorkbook は
func (c *Client) InsertShareWorkbook(bookId string) error {
	ctx := context.Background()
	err, book := c.SelectWorkbook(bookId)
	if err != nil {
		return err
	}

	//キー作成
	parentKey := datastore.IDKey("user_account", book.UserId, nil)
	childKey := datastore.IncompleteKey("share_workbook", parentKey)

	//
	_, err = c.DataStore.Put(ctx, childKey, &book)
	if err != nil {
		return err
	}
	return nil
}

//SelectShareWorkbook　は
func (c *Client) SelectShareWorkbooks() {

}
