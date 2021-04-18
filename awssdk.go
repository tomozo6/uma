package uma

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iam"
)

// ----------------------------------------------------------------------------
// IAM ListGroups
// ----------------------------------------------------------------------------
type IAMListGroupsAPI interface {
	ListGroups(ctx context.Context, params *iam.ListGroupsInput, optFns ...func(*iam.Options)) (*iam.ListGroupsOutput, error)
}

func IAMListGroups(c context.Context, api IAMListGroupsAPI, input *iam.ListGroupsInput) (*iam.ListGroupsOutput, error) {
	return api.ListGroups(c, input)
}

// ----------------------------------------------------------------------------
// IAM GetGroup
// ----------------------------------------------------------------------------
type IAMGetGroupAPI interface {
	GetGroup(ctx context.Context, params *iam.GetGroupInput, optFns ...func(*iam.Options)) (*iam.GetGroupOutput, error)
}

func IAMGetGroup(c context.Context, api IAMGetGroupAPI, input *iam.GetGroupInput) (*iam.GetGroupOutput, error) {
	return api.GetGroup(c, input)
}

// ----------------------------------------------------------------------------
// IAM ListSSHPublicKeys
// ----------------------------------------------------------------------------
type IAMListSSHPublicKeysAPI interface {
	ListSSHPublicKeys(ctx context.Context, params *iam.ListSSHPublicKeysInput, optFns ...func(*iam.Options)) (*iam.ListSSHPublicKeysOutput, error)
}

func IAMListSSHPublicKeys(c context.Context, api IAMListSSHPublicKeysAPI, input *iam.ListSSHPublicKeysInput) (*iam.ListSSHPublicKeysOutput, error) {
	return api.ListSSHPublicKeys(c, input)
}

// ----------------------------------------------------------------------------
// IAM GetSSHPublicKey
// ----------------------------------------------------------------------------
type IAMGetSSHPublicKeyAPI interface {
	GetSSHPublicKey(ctx context.Context, params *iam.GetSSHPublicKeyInput, optFns ...func(*iam.Options)) (*iam.GetSSHPublicKeyOutput, error)
}

func IAMGetSSHPublicKey(c context.Context, api IAMGetSSHPublicKeyAPI, input *iam.GetSSHPublicKeyInput) (*iam.GetSSHPublicKeyOutput, error) {
	return api.GetSSHPublicKey(c, input)
}

// ----------------------------------------------------------------------------
// IAM ListGroupsForUser
// ----------------------------------------------------------------------------
type IAMListGroupsForUserAPI interface {
	ListGroupsForUser(ctx context.Context, params *iam.ListGroupsForUserInput, optFns ...func(*iam.Options)) (*iam.ListGroupsForUserOutput, error)
}

func IAMListGroupsForUser(c context.Context, api IAMListGroupsForUserAPI, input *iam.ListGroupsForUserInput) (*iam.ListGroupsForUserOutput, error) {
	return api.ListGroupsForUser(c, input)
}
