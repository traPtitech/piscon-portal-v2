// Code generated by ogen, DO NOT EDIT.

package openapi

// CreateTeamInstanceParams is parameters of createTeamInstance operation.
type CreateTeamInstanceParams struct {
	// チームID.
	TeamId TeamId
}

// DeleteTeamInstanceParams is parameters of deleteTeamInstance operation.
type DeleteTeamInstanceParams struct {
	// チームID.
	TeamId TeamId
	// インスタンスID.
	InstanceId InstanceId
}

// GetOauth2CallbackParams is parameters of getOauth2Callback operation.
type GetOauth2CallbackParams struct {
	// TraQからのAuthorization Code.
	Code string
}

// GetTeamParams is parameters of getTeam operation.
type GetTeamParams struct {
	// チームID.
	TeamId TeamId
}

// GetTeamInstancesParams is parameters of getTeamInstances operation.
type GetTeamInstancesParams struct {
	// チームID.
	TeamId TeamId
}

// PatchTeamParams is parameters of patchTeam operation.
type PatchTeamParams struct {
	// チームID.
	TeamId TeamId
}

// PatchTeamInstanceParams is parameters of patchTeamInstance operation.
type PatchTeamInstanceParams struct {
	// チームID.
	TeamId TeamId
	// インスタンスID.
	InstanceId InstanceId
}
