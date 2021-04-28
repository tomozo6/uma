package awsiam

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/tomozo6/uma/pkg/types"
)

// ----------------------------------------------------------------------------
// Users
// ----------------------------------------------------------------------------
type Users struct {
	UserNames []string
	Client    *iam.Client
}

type GetGroupsKeysForUsersOutput struct {
	GroupsKeysForUsers []*types.GroupsKeysForUser
}

// type GroupsKeysForUser struct {
// 	UserName          string
// 	GroupNames        []string
// 	SSHPublicKeyBodys []string
// }

func NewUsers(userNames []string) (*Users, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	client := iam.NewFromConfig(cfg)

	u := &Users{
		UserNames: userNames,
		Client:    client,
	}
	return u, nil
}

func (u *Users) GetGroupsKeysForUsers(inputGroupList []string) (*GetGroupsKeysForUsersOutput, error) {

	// return用の構造体
	r := new(GetGroupsKeysForUsersOutput)

	for _, userName := range u.UserNames {
		// --------------------------------------------------------------------
		// 鍵情報を取得
		// --------------------------------------------------------------------
		input_k := &iam.ListSSHPublicKeysInput{
			UserName: &userName,
		}

		resp_k, err := IAMListSSHPublicKeys(context.TODO(), u.Client, input_k)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		// 鍵が無いユーザーは除外
		if len(resp_k.SSHPublicKeys) == 0 {
			continue
		}

		// キーID毎に鍵の実体を取得
		var SSHPublicKeyBodys []string
		for _, key := range resp_k.SSHPublicKeys {

			input := &iam.GetSSHPublicKeyInput{
				Encoding:       "SSH",
				SSHPublicKeyId: key.SSHPublicKeyId,
				UserName:       &userName,
			}

			resp, err := IAMGetSSHPublicKey(context.TODO(), u.Client, input)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			SSHPublicKeyBodys = append(SSHPublicKeyBodys, *resp.SSHPublicKey.SSHPublicKeyBody)
		}

		// --------------------------------------------------------------------
		// グループ情報取得
		// --------------------------------------------------------------------
		input_g := &iam.ListGroupsForUserInput{
			UserName: &userName,
		}

		resp_g, err := IAMListGroupsForUser(context.TODO(), u.Client, input_g)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		// 指定されたグループに属していないユーザーは除外
		if len(resp_g.Groups) == 0 {
			continue
		}

		var GroupNames []string
		for _, group := range resp_g.Groups {

			for _, inputGroup := range inputGroupList {
				// 入力されたグループ名と一致した場合Sliceに追加
				if *group.GroupName == inputGroup {
					GroupNames = append(GroupNames, *group.GroupName)
				}
			}
		}

		// --------------------------------------------------------------------
		// リターン用構造体に情報を格納
		// --------------------------------------------------------------------
		// IAMユーザー名から@以降の文字列を削除
		uName := strings.Split(userName, "@")[0]

		l := types.GroupsKeysForUser{
			UserName:          uName,
			GroupNames:        GroupNames,
			SSHPublicKeyBodys: SSHPublicKeyBodys,
		}

		r.GroupsKeysForUsers = append(r.GroupsKeysForUsers, &l)
	}
	return r, nil
}
