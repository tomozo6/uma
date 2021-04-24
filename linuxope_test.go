package uma

import (
	"testing"
)

func TestDelLinuxUser(t *testing.T) {
	err := DelLinuxUser("aiueo")

	if err != nil {
		t.Fatal("なにかがおこった！")
	}

}

func TestListLinuxUMAUser(t *testing.T) {
	_, err := ListLinuxUMAUser()
	if err != nil {
		t.Fatal("ユーザーリストの作成に失敗しました")
	}
}

func TestContains(t *testing.T) {
	fruits := []string{"apple", "orange", "greap"}

	if Contains(fruits, "apple") == false {
		t.Fatal("おかしい")
	}

	if Contains(fruits, "dog") {
		t.Fatal("おかしい")
	}
}
