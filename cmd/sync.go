package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tomozo6/uma/pkg/awsiam"
	"github.com/tomozo6/uma/pkg/linux"
)

var (
	targetIAMGroups []string
	skipDelete      bool
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronize AWS IAM users with Linux users.",
	Long: `UMA (User Management Automation)

Synchronize AWS IAM users with Linux users.`,
	Run: func(cmd *cobra.Command, args []string) {
		// ------------------------------------------------------------------------
		// AWS情報取得
		// ------------------------------------------------------------------------
		// Groups構造体を生成
		g, _ := awsiam.NewGroups(targetIAMGroups)

		// 引数で与えられたグループ名のうち、AWS上に存在しないグループをスライスから削除
		if err := g.FilterIAMGroup(); err != nil {
			fmt.Println(err)
		}

		// 引数で与えられたグループ名のうち、サーバー上に存在しないグループをサーバー上に作成
		if err := g.AddOSGroup(); err != nil {
			fmt.Println(err)
		}

		// 各IAMグループに属するIAMユーザー名をマージ&ユニークしたスライスを作成
		iamUserNames, err := g.GetIAMUserNames()
		if err != nil {
			fmt.Println(err)
		}

		// Users構造体を生成
		u, _ := awsiam.NewUsers(iamUserNames)
		// Usersリストを取得
		usersInfo, err := u.GetGroupsKeysForUsers(targetIAMGroups)

		// ------------------------------------------------------------------------
		// ユーザー作成
		// ------------------------------------------------------------------------
		for _, v := range usersInfo.GroupsKeysForUsers {

			fmt.Println("add or mod user:", v.UserName)
			err := linux.UserAdd(v)
			if err != nil {
				fmt.Println("Error: Failed to add or mod user.")
				fmt.Println(err)
			}
		}

		// ------------------------------------------------------------------------
		// ユーザー削除
		// ------------------------------------------------------------------------
		if skipDelete == false {
			// 対象IAMユーザー名のSliceを作成
			// またIAMユーザー名から@以降の文字列を削除している
			var targetUserNames []string
			for _, v := range usersInfo.GroupsKeysForUsers {
				targetUserNames = append(targetUserNames, strings.Split(v.UserName, "@")[0])
			}

			// Linux上でUMAに管理されているユーザー名のSliceを作成
			linuxUserNames, err := linux.ListUMAUser()
			if err != nil {
				fmt.Println("Error: Failed to get the LinuxUMAUsersList.")
				fmt.Println(err)
			}

			// 対象ユーザーとなっていないUMA管理Linuxユーザーを削除
			for _, v := range linuxUserNames {
				if linux.Contains(targetUserNames, v) == false {
					// Linuxユーザー削除
					fmt.Println("delete user:", v)
					err := linux.UserDel(v)
					if err != nil {
						fmt.Printf("Error: Failed to delete user.")
						fmt.Println(err)
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	syncCmd.Flags().StringSliceVarP(&targetIAMGroups, "groups", "g", nil, "target IAM groups(ex: sre,director)")
	syncCmd.MarkFlagRequired("groups")

	syncCmd.Flags().Bool("skip-delete", false, "skip delete users")

}
