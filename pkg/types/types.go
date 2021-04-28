package types

type GroupsKeysForUser struct {
	UserName          string
	GroupNames        []string
	SSHPublicKeyBodys []string
}
