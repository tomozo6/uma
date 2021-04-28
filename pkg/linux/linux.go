package linux

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/tomozo6/uma/pkg/types"
)

// ----------------------------------------------------------------------------
// User 操作関連
// ----------------------------------------------------------------------------
func UserAdd(params *types.GroupsKeysForUser) error {
	userName := params.UserName

	// リストをカンマ区切りの文字列に変換
	// 例) [groupa, groupb] -> groupa,groupb
	groupNamesStr := strings.Join(params.GroupNames, ",")

	// ユーザー存在確認
	if err := exec.Command("id", userName).Run(); err != nil {
		// ユーザーが存在しなかったらユーザー作成
		if out, err := exec.Command("useradd", "-m", userName, "-c", "ManagedByUMA", "-G", groupNamesStr).CombinedOutput(); err != nil {
			fmt.Println("Command Exec Error.")
			fmt.Printf("useradd result: \n%s", string(out))
			return err
		}
	} else {
		// ユーザーが存在していたらユーザー情報変更
		if out, err := exec.Command("usermod", userName, "-c", "ManagedByUMA", "-G", groupNamesStr).CombinedOutput(); err != nil {
			fmt.Println("Command Exec Error.")
			fmt.Printf("useradd result: \n%s", string(out))
			return err
		}
	}

	// 作成したユーザー情報の取得
	u, err := user.Lookup(userName)
	if err != nil {
		fmt.Println("ユーザー情報の取得に失敗しました。")
		return err
	}

	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		fmt.Println("uidの string->int 型変換に失敗しました。")
		return err
	}

	gid, err := strconv.Atoi(u.Gid)
	if err != nil {
		fmt.Println("gidの string->int 型変換に失敗しました。")
		return err
	}

	// .sshディレクトリの作成
	sshDir := filepath.Join(u.HomeDir, ".ssh")
	if _, err := os.Stat(sshDir); os.IsNotExist(err) {
		if err := os.Mkdir(sshDir, 0700); err != nil {
			fmt.Println(".sshディレクトリの作成に失敗しました。")
			return err
		}
	}

	if err := os.Chown(sshDir, uid, gid); err != nil {
		fmt.Println(".sshディレクトリのchownに失敗しました。")
		return err
	}

	// authorized_keysの作成
	authkey := filepath.Join(u.HomeDir, ".ssh", "authorized_keys")

	f, err := os.Create(authkey)
	if err != nil {
		fmt.Println("authorized_keysの作成or編集に失敗しました。")
		return err
	}
	defer f.Close()

	if err := f.Chmod(0600); err != nil {
		fmt.Println("authorized_keysのchmodに失敗しました。")
		return err
	}

	if err := f.Chown(uid, gid); err != nil {
		fmt.Println("authorized_keysのchownに失敗しました。")
		return err
	}

	// 公開鍵書き込み
	for _, v := range params.SSHPublicKeyBodys {
		f.Write([]byte(v + "\n"))
	}

	return nil
}

func UserDel(userName string) error {
	// ユーザー存在確認
	if err := exec.Command("id", userName).Run(); err == nil {
		// ユーザーが存在していたらユーザー削除
		if out, err := exec.Command("userdel", "-r", userName).CombinedOutput(); err != nil {
			fmt.Printf("ユーザー削除に失敗しました: \n%s", string(out))
			return err
		}
	}
	return nil
}

func ListUMAUser() ([]string, error) {
	var s []string

	// ファイルオープン
	f, err := os.Open("/etc/passwd")
	if err != nil {
		fmt.Println("/etc/passwdのオープンに失敗しました")
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// 一行ずつ処理
	for scanner.Scan() {
		// コメント欄にManagedByUMAが入っているユーザー名をsliceに追加
		if strings.Contains(scanner.Text(), ":ManagedByUMA:") {
			s = append(s, strings.Split(scanner.Text(), ":")[0])
		}

	}
	if err = scanner.Err(); err != nil {
		fmt.Println("/etc/passwdの読み込みに失敗しました。")
		return nil, err
	}

	return s, nil
}

func Contains(s []string, e string) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
