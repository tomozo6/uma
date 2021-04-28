package awsiam

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

// ----------------------------------------------------------------------------
// Groups
// ----------------------------------------------------------------------------
type Groups struct {
	GroupNames []string
	Client     *iam.Client
}

func NewGroups(groups []string) (*Groups, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := iam.NewFromConfig(cfg)

	g := &Groups{
		GroupNames: groups,
		Client:     client,
	}
	return g, nil
}

func (g *Groups) FilterIAMGroup() error {
	input := &iam.ListGroupsInput{}

	result, err := IAMListGroups(context.TODO(), g.Client, input)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var s []string
	for _, v := range result.Groups {
		for _, w := range g.GroupNames {
			if w == *v.GroupName {
				s = append(s, w)
			}
		}
	}
	g.GroupNames = s
	return nil
}

func (g *Groups) FilterOSGroup() error {
	var s []string

	for _, v := range g.GroupNames {
		if err := exec.Command("getent", "group", v).Run(); err == nil {
			s = append(s, v)
		}
	}

	g.GroupNames = s
	return nil
}

func (g *Groups) GetIAMUserNames() ([]string, error) {
	var s []string

	for _, v := range g.GroupNames {
		input := &iam.GetGroupInput{
			GroupName: &v,
		}

		result, err := IAMGetGroup(context.TODO(), g.Client, input)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		for _, w := range result.Users {
			s = append(s, *w.UserName)
		}
	}

	// Slice重複削除
	m := make(map[string]struct{})
	iamUserNames := make([]string, 0)

	for _, v := range s {
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			iamUserNames = append(iamUserNames, v)
		}
	}

	return iamUserNames, nil
}
