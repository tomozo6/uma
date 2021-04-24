package main

import (
	"fmt"
	"strings"
	"uma"

	flag "github.com/spf13/pflag"
)

var version = "0.0.1"

// ----------------------------------------------------------------------------
// main
// ----------------------------------------------------------------------------
// groupadd IMJ_Infra -g 2000 ; groupadd IMJ_Server -g 2001 ; groupadd IMJ_PMO -g 2002 ; groupadd IMJ_App -g 2003

func main() {
	// ----------------------------------------------------------------------------
	// フラグハンドリング
	// ----------------------------------------------------------------------------
	var showHelp bool
	flag.BoolVarP(&showHelp, "help", "h", false, "show help")

	var showVersion bool
	flag.BoolVarP(&showVersion, "version", "v", false, "show version")

	var targetIAMGroups []string
	flag.StringSliceVarP(&targetIAMGroups, "target", "t", nil, "target iam groups")

	flag.Parse()

	if showHelp {
		flag.PrintDefaults()
		return
	}

	if showVersion {
		fmt.Println("version:", version)
		return
	}

	if len(targetIAMGroups) == 0 {
		label := "The number of targetIAMGroups specified is 0. Therefore, all OS users managed by UMA will be deleted. Are you sure you want to do it?"
		if uma.YesNo(label) == false {
			return
		}
	}

	// ------------------------------------------------------------------------
	// AWS情報取得
	// ------------------------------------------------------------------------
	// Groups構造体を生成
	fmt.Println("targetIAMGroups: ", targetIAMGroups)
	g, _ := uma.NewGroups(targetIAMGroups)

	// 引数で与えられたグループ名のうち、サーバー上に存在しないグループをスライスから削除
	if err := g.FilterOSGroup(); err != nil {
		fmt.Println(err)
	}

	// 引数で与えられたグループ名のうち、AWS上に存在しないグループをスライスから削除
	if err := g.FilterIAMGroup(); err != nil {
		fmt.Println(err)
	}

	// 各IAMグループに属するIAMユーザー名をマージ&ユニークしたスライスを作成
	iamUserNames, err := g.GetIAMUserNames()
	if err != nil {
		fmt.Println(err)
	}

	// Users構造体を生成
	u, _ := uma.NewUsers(iamUserNames)
	// Usersリストを取得
	usersInfo, err := u.GetGroupsKeysForUsers(targetIAMGroups)
	fmt.Println("usersInfo: ", usersInfo)

	// ------------------------------------------------------------------------
	// ユーザー作成
	// ------------------------------------------------------------------------
	for _, v := range usersInfo.GroupsKeysForUsers {
		err := uma.AddLinuxUser(v)
		if err != nil {
			fmt.Println("ユーザー作成に失敗しました。")
			fmt.Println(err)
		}
	}

	// ------------------------------------------------------------------------
	// ユーザー削除
	// ------------------------------------------------------------------------
	// 対象IAMユーザー名のSliceを作成
	// またIAMユーザー名から@以降の文字列を削除している
	var targetUserNames []string
	for _, v := range usersInfo.GroupsKeysForUsers {
		targetUserNames = append(targetUserNames, strings.Split(v.UserName, "@")[0])
	}

	// Linux上でUMAに管理されているユーザー名Sliceを作成
	linuxUserNames, err := uma.ListLinuxUMAUser()
	if err != nil {
		fmt.Println("linuxUserListの取得に失敗しました。")
		fmt.Println(err)
	}

	// 対象ユーザーとなっていないUMA管理Linuxユーザーを削除
	for _, v := range linuxUserNames {
		if uma.Contains(targetUserNames, v) == false {
			// Linuxユーザー削除
			err := uma.DelLinuxUser(v)
			if err != nil {
				fmt.Printf("ユーザー(%s)の削除に失敗しました。", v)
				fmt.Println(err)
			}
		}
	}
}
