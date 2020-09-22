package api

import "testing"

func TestAuthGoogle(t *testing.T) {
	invalidCode := "hoge"
	if _, err := AuthGoogle(invalidCode); err == nil {
		t.Errorf("不適切なコードからは認証が通らないこと %v err:%v", invalidCode, err)
	}

	/* 認証が通る場合のテスト

	   dummyCode := "hoge"       //TODO: 認証がとおるテスト用のコードがあればいれる
	   wantUser := UserAccount{} //TODO: 上記ユーザーから取得できるはずのデータを用意

	   gotUser, err := AuthGoogle(dummyCode)
	   if err != nil {
	           t.Errorf("cant auth from %v err:%v", dummyCode, err)
	   }
	   if diff := pretty.Compare(gotUser, wantUser); diff != "" {
	           t.Errorf("invalid userInfo: %s", pretty.Compare(gotUser, wantUser))
	   }
	*/
}
