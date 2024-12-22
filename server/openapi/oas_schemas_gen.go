// Code generated by ogen, DO NOT EDIT.

package openapi

import (
	"net/url"
	"time"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
)

type AdminAuth struct {
	APIKey string
}

// GetAPIKey returns the value of APIKey.
func (s *AdminAuth) GetAPIKey() string {
	return s.APIKey
}

// SetAPIKey sets the value of APIKey.
func (s *AdminAuth) SetAPIKey(val string) {
	s.APIKey = val
}

// DeleteTeamInstanceOK is response for DeleteTeamInstance operation.
type DeleteTeamInstanceOK struct{}

func (*DeleteTeamInstanceOK) deleteTeamInstanceRes() {}

// Ref: #/components/schemas/ErrorBadRequest
type ErrorBadRequest struct {
	Message OptString `json:"message"`
}

// GetMessage returns the value of Message.
func (s *ErrorBadRequest) GetMessage() OptString {
	return s.Message
}

// SetMessage sets the value of Message.
func (s *ErrorBadRequest) SetMessage(val OptString) {
	s.Message = val
}

func (*ErrorBadRequest) createTeamInstanceRes() {}
func (*ErrorBadRequest) patchTeamInstanceRes()  {}
func (*ErrorBadRequest) patchTeamRes()          {}
func (*ErrorBadRequest) postTeamRes()           {}

type Forbidden struct {
	Message OptString `json:"message"`
}

// GetMessage returns the value of Message.
func (s *Forbidden) GetMessage() OptString {
	return s.Message
}

// SetMessage sets the value of Message.
func (s *Forbidden) SetMessage(val OptString) {
	s.Message = val
}

func (*Forbidden) createTeamInstanceRes() {}
func (*Forbidden) deleteTeamInstanceRes() {}
func (*Forbidden) getInstancesRes()       {}
func (*Forbidden) getTeamInstancesRes()   {}
func (*Forbidden) patchTeamInstanceRes()  {}
func (*Forbidden) patchTeamRes()          {}

type GetInstancesOKApplicationJSON []Instance

func (*GetInstancesOKApplicationJSON) getInstancesRes() {}

type GetOauth2CallbackBadRequest struct {
	Message OptString `json:"message"`
}

// GetMessage returns the value of Message.
func (s *GetOauth2CallbackBadRequest) GetMessage() OptString {
	return s.Message
}

// SetMessage sets the value of Message.
func (s *GetOauth2CallbackBadRequest) SetMessage(val OptString) {
	s.Message = val
}

func (*GetOauth2CallbackBadRequest) getOauth2CallbackRes() {}

// GetOauth2CallbackOK is response for GetOauth2Callback operation.
type GetOauth2CallbackOK struct{}

func (*GetOauth2CallbackOK) getOauth2CallbackRes() {}

// GetOauth2CodeSeeOther is response for GetOauth2Code operation.
type GetOauth2CodeSeeOther struct {
	Location OptURI
}

// GetLocation returns the value of Location.
func (s *GetOauth2CodeSeeOther) GetLocation() OptURI {
	return s.Location
}

// SetLocation sets the value of Location.
func (s *GetOauth2CodeSeeOther) SetLocation(val OptURI) {
	s.Location = val
}

func (*GetOauth2CodeSeeOther) getOauth2CodeRes() {}

type GetTeamInstancesOKApplicationJSON []Instance

func (*GetTeamInstancesOKApplicationJSON) getTeamInstancesRes() {}

type GetTeamsOKApplicationJSON []Team

func (*GetTeamsOKApplicationJSON) getTeamsRes() {}

type GetUsersOKApplicationJSON []User

func (*GetUsersOKApplicationJSON) getUsersRes() {}

type IPAddress string

// サーバーのインスタンス.
// Ref: #/components/schemas/Instance
type Instance struct {
	ID     InstanceId `json:"id"`
	TeamId TeamId     `json:"teamId"`
	// サーバーのID。チームごとに1から始まる。表示用。
	// 新しくインスタンスを起動する場合は、1以上の現在使われていない整数が採用される。
	// インスタンスを削除したら、そのIDは使用されていないものとし、再利用できる。.
	ServerId         int            `json:"serverId"`
	PublicIPAddress  IPAddress      `json:"publicIPAddress"`
	PrivateIPAddress IPAddress      `json:"privateIPAddress"`
	Status           InstanceStatus `json:"status"`
	CreatedAt        time.Time      `json:"createdAt"`
}

// GetID returns the value of ID.
func (s *Instance) GetID() InstanceId {
	return s.ID
}

// GetTeamId returns the value of TeamId.
func (s *Instance) GetTeamId() TeamId {
	return s.TeamId
}

// GetServerId returns the value of ServerId.
func (s *Instance) GetServerId() int {
	return s.ServerId
}

// GetPublicIPAddress returns the value of PublicIPAddress.
func (s *Instance) GetPublicIPAddress() IPAddress {
	return s.PublicIPAddress
}

// GetPrivateIPAddress returns the value of PrivateIPAddress.
func (s *Instance) GetPrivateIPAddress() IPAddress {
	return s.PrivateIPAddress
}

// GetStatus returns the value of Status.
func (s *Instance) GetStatus() InstanceStatus {
	return s.Status
}

// GetCreatedAt returns the value of CreatedAt.
func (s *Instance) GetCreatedAt() time.Time {
	return s.CreatedAt
}

// SetID sets the value of ID.
func (s *Instance) SetID(val InstanceId) {
	s.ID = val
}

// SetTeamId sets the value of TeamId.
func (s *Instance) SetTeamId(val TeamId) {
	s.TeamId = val
}

// SetServerId sets the value of ServerId.
func (s *Instance) SetServerId(val int) {
	s.ServerId = val
}

// SetPublicIPAddress sets the value of PublicIPAddress.
func (s *Instance) SetPublicIPAddress(val IPAddress) {
	s.PublicIPAddress = val
}

// SetPrivateIPAddress sets the value of PrivateIPAddress.
func (s *Instance) SetPrivateIPAddress(val IPAddress) {
	s.PrivateIPAddress = val
}

// SetStatus sets the value of Status.
func (s *Instance) SetStatus(val InstanceStatus) {
	s.Status = val
}

// SetCreatedAt sets the value of CreatedAt.
func (s *Instance) SetCreatedAt(val time.Time) {
	s.CreatedAt = val
}

func (*Instance) createTeamInstanceRes() {}

type InstanceId uuid.UUID

// Ref: #/components/schemas/InstanceStatus
type InstanceStatus string

const (
	InstanceStatusBuilding InstanceStatus = "building"
	InstanceStatusRunning  InstanceStatus = "running"
	InstanceStatusStopped  InstanceStatus = "stopped"
)

// AllValues returns all InstanceStatus values.
func (InstanceStatus) AllValues() []InstanceStatus {
	return []InstanceStatus{
		InstanceStatusBuilding,
		InstanceStatusRunning,
		InstanceStatusStopped,
	}
}

// MarshalText implements encoding.TextMarshaler.
func (s InstanceStatus) MarshalText() ([]byte, error) {
	switch s {
	case InstanceStatusBuilding:
		return []byte(s), nil
	case InstanceStatusRunning:
		return []byte(s), nil
	case InstanceStatusStopped:
		return []byte(s), nil
	default:
		return nil, errors.Errorf("invalid value: %q", s)
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *InstanceStatus) UnmarshalText(data []byte) error {
	switch InstanceStatus(data) {
	case InstanceStatusBuilding:
		*s = InstanceStatusBuilding
		return nil
	case InstanceStatusRunning:
		*s = InstanceStatusRunning
		return nil
	case InstanceStatusStopped:
		*s = InstanceStatusStopped
		return nil
	default:
		return errors.Errorf("invalid value: %q", data)
	}
}

type InternalServerError struct {
	Message OptString `json:"message"`
}

// GetMessage returns the value of Message.
func (s *InternalServerError) GetMessage() OptString {
	return s.Message
}

// SetMessage sets the value of Message.
func (s *InternalServerError) SetMessage(val OptString) {
	s.Message = val
}

func (*InternalServerError) createTeamInstanceRes() {}
func (*InternalServerError) deleteTeamInstanceRes() {}
func (*InternalServerError) getInstancesRes()       {}
func (*InternalServerError) getMeRes()              {}
func (*InternalServerError) getOauth2CallbackRes()  {}
func (*InternalServerError) getOauth2CodeRes()      {}
func (*InternalServerError) getTeamInstancesRes()   {}
func (*InternalServerError) getTeamRes()            {}
func (*InternalServerError) getTeamsRes()           {}
func (*InternalServerError) getUsersRes()           {}
func (*InternalServerError) patchTeamInstanceRes()  {}
func (*InternalServerError) patchTeamRes()          {}
func (*InternalServerError) postOauth2LogoutRes()   {}
func (*InternalServerError) postTeamRes()           {}

type NotFound struct {
	Message OptString `json:"message"`
}

// GetMessage returns the value of Message.
func (s *NotFound) GetMessage() OptString {
	return s.Message
}

// SetMessage sets the value of Message.
func (s *NotFound) SetMessage(val OptString) {
	s.Message = val
}

func (*NotFound) createTeamInstanceRes() {}
func (*NotFound) deleteTeamInstanceRes() {}
func (*NotFound) getTeamInstancesRes()   {}
func (*NotFound) getTeamRes()            {}
func (*NotFound) patchTeamInstanceRes()  {}
func (*NotFound) patchTeamRes()          {}

// NewOptInstanceStatus returns new OptInstanceStatus with value set to v.
func NewOptInstanceStatus(v InstanceStatus) OptInstanceStatus {
	return OptInstanceStatus{
		Value: v,
		Set:   true,
	}
}

// OptInstanceStatus is optional InstanceStatus.
type OptInstanceStatus struct {
	Value InstanceStatus
	Set   bool
}

// IsSet returns true if OptInstanceStatus was set.
func (o OptInstanceStatus) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptInstanceStatus) Reset() {
	var v InstanceStatus
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptInstanceStatus) SetTo(v InstanceStatus) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptInstanceStatus) Get() (v InstanceStatus, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptInstanceStatus) Or(d InstanceStatus) InstanceStatus {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptPatchTeamReq returns new OptPatchTeamReq with value set to v.
func NewOptPatchTeamReq(v PatchTeamReq) OptPatchTeamReq {
	return OptPatchTeamReq{
		Value: v,
		Set:   true,
	}
}

// OptPatchTeamReq is optional PatchTeamReq.
type OptPatchTeamReq struct {
	Value PatchTeamReq
	Set   bool
}

// IsSet returns true if OptPatchTeamReq was set.
func (o OptPatchTeamReq) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptPatchTeamReq) Reset() {
	var v PatchTeamReq
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptPatchTeamReq) SetTo(v PatchTeamReq) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptPatchTeamReq) Get() (v PatchTeamReq, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptPatchTeamReq) Or(d PatchTeamReq) PatchTeamReq {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptPostTeamReq returns new OptPostTeamReq with value set to v.
func NewOptPostTeamReq(v PostTeamReq) OptPostTeamReq {
	return OptPostTeamReq{
		Value: v,
		Set:   true,
	}
}

// OptPostTeamReq is optional PostTeamReq.
type OptPostTeamReq struct {
	Value PostTeamReq
	Set   bool
}

// IsSet returns true if OptPostTeamReq was set.
func (o OptPostTeamReq) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptPostTeamReq) Reset() {
	var v PostTeamReq
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptPostTeamReq) SetTo(v PostTeamReq) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptPostTeamReq) Get() (v PostTeamReq, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptPostTeamReq) Or(d PostTeamReq) PostTeamReq {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptString returns new OptString with value set to v.
func NewOptString(v string) OptString {
	return OptString{
		Value: v,
		Set:   true,
	}
}

// OptString is optional string.
type OptString struct {
	Value string
	Set   bool
}

// IsSet returns true if OptString was set.
func (o OptString) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptString) Reset() {
	var v string
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptString) SetTo(v string) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptString) Get() (v string, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptString) Or(d string) string {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptTeamId returns new OptTeamId with value set to v.
func NewOptTeamId(v TeamId) OptTeamId {
	return OptTeamId{
		Value: v,
		Set:   true,
	}
}

// OptTeamId is optional TeamId.
type OptTeamId struct {
	Value TeamId
	Set   bool
}

// IsSet returns true if OptTeamId was set.
func (o OptTeamId) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptTeamId) Reset() {
	var v TeamId
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptTeamId) SetTo(v TeamId) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptTeamId) Get() (v TeamId, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptTeamId) Or(d TeamId) TeamId {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptTeamName returns new OptTeamName with value set to v.
func NewOptTeamName(v TeamName) OptTeamName {
	return OptTeamName{
		Value: v,
		Set:   true,
	}
}

// OptTeamName is optional TeamName.
type OptTeamName struct {
	Value TeamName
	Set   bool
}

// IsSet returns true if OptTeamName was set.
func (o OptTeamName) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptTeamName) Reset() {
	var v TeamName
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptTeamName) SetTo(v TeamName) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptTeamName) Get() (v TeamName, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptTeamName) Or(d TeamName) TeamName {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptURI returns new OptURI with value set to v.
func NewOptURI(v url.URL) OptURI {
	return OptURI{
		Value: v,
		Set:   true,
	}
}

// OptURI is optional url.URL.
type OptURI struct {
	Value url.URL
	Set   bool
}

// IsSet returns true if OptURI was set.
func (o OptURI) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptURI) Reset() {
	var v url.URL
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptURI) SetTo(v url.URL) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptURI) Get() (v url.URL, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptURI) Or(d url.URL) url.URL {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// PatchTeamInstanceOK is response for PatchTeamInstance operation.
type PatchTeamInstanceOK struct{}

func (*PatchTeamInstanceOK) patchTeamInstanceRes() {}

type PatchTeamInstanceReq struct {
	Status OptInstanceStatus `json:"status"`
}

// GetStatus returns the value of Status.
func (s *PatchTeamInstanceReq) GetStatus() OptInstanceStatus {
	return s.Status
}

// SetStatus sets the value of Status.
func (s *PatchTeamInstanceReq) SetStatus(val OptInstanceStatus) {
	s.Status = val
}

type PatchTeamReq struct {
	Name OptTeamName `json:"name"`
	// チームに所属させる部員のID.
	Members []UserId `json:"members"`
}

// GetName returns the value of Name.
func (s *PatchTeamReq) GetName() OptTeamName {
	return s.Name
}

// GetMembers returns the value of Members.
func (s *PatchTeamReq) GetMembers() []UserId {
	return s.Members
}

// SetName sets the value of Name.
func (s *PatchTeamReq) SetName(val OptTeamName) {
	s.Name = val
}

// SetMembers sets the value of Members.
func (s *PatchTeamReq) SetMembers(val []UserId) {
	s.Members = val
}

type PostOauth2LogoutBadRequest struct {
	Message OptString `json:"message"`
}

// GetMessage returns the value of Message.
func (s *PostOauth2LogoutBadRequest) GetMessage() OptString {
	return s.Message
}

// SetMessage sets the value of Message.
func (s *PostOauth2LogoutBadRequest) SetMessage(val OptString) {
	s.Message = val
}

func (*PostOauth2LogoutBadRequest) postOauth2LogoutRes() {}

// PostOauth2LogoutOK is response for PostOauth2Logout operation.
type PostOauth2LogoutOK struct{}

func (*PostOauth2LogoutOK) postOauth2LogoutRes() {}

type PostTeamReq struct {
	Name TeamName `json:"name"`
	// チームに所属させる部員のID.
	Members []UserId `json:"members"`
}

// GetName returns the value of Name.
func (s *PostTeamReq) GetName() TeamName {
	return s.Name
}

// GetMembers returns the value of Members.
func (s *PostTeamReq) GetMembers() []UserId {
	return s.Members
}

// SetName sets the value of Name.
func (s *PostTeamReq) SetName(val TeamName) {
	s.Name = val
}

// SetMembers sets the value of Members.
func (s *PostTeamReq) SetMembers(val []UserId) {
	s.Members = val
}

// チーム.
// Ref: #/components/schemas/Team
type Team struct {
	ID   TeamId   `json:"id"`
	Name TeamName `json:"name"`
	// チームに所属している部員のID.
	Members   []UserId  `json:"members"`
	CreatedAt time.Time `json:"createdAt"`
}

// GetID returns the value of ID.
func (s *Team) GetID() TeamId {
	return s.ID
}

// GetName returns the value of Name.
func (s *Team) GetName() TeamName {
	return s.Name
}

// GetMembers returns the value of Members.
func (s *Team) GetMembers() []UserId {
	return s.Members
}

// GetCreatedAt returns the value of CreatedAt.
func (s *Team) GetCreatedAt() time.Time {
	return s.CreatedAt
}

// SetID sets the value of ID.
func (s *Team) SetID(val TeamId) {
	s.ID = val
}

// SetName sets the value of Name.
func (s *Team) SetName(val TeamName) {
	s.Name = val
}

// SetMembers sets the value of Members.
func (s *Team) SetMembers(val []UserId) {
	s.Members = val
}

// SetCreatedAt sets the value of CreatedAt.
func (s *Team) SetCreatedAt(val time.Time) {
	s.CreatedAt = val
}

func (*Team) getTeamRes()   {}
func (*Team) patchTeamRes() {}
func (*Team) postTeamRes()  {}

type TeamAuth struct {
	APIKey string
}

// GetAPIKey returns the value of APIKey.
func (s *TeamAuth) GetAPIKey() string {
	return s.APIKey
}

// SetAPIKey sets the value of APIKey.
func (s *TeamAuth) SetAPIKey(val string) {
	s.APIKey = val
}

type TeamId uuid.UUID

type TeamName string

type Unauthorized struct {
	Message OptString `json:"message"`
}

// GetMessage returns the value of Message.
func (s *Unauthorized) GetMessage() OptString {
	return s.Message
}

// SetMessage sets the value of Message.
func (s *Unauthorized) SetMessage(val OptString) {
	s.Message = val
}

func (*Unauthorized) createTeamInstanceRes() {}
func (*Unauthorized) deleteTeamInstanceRes() {}
func (*Unauthorized) getInstancesRes()       {}
func (*Unauthorized) getMeRes()              {}
func (*Unauthorized) getTeamInstancesRes()   {}
func (*Unauthorized) getTeamRes()            {}
func (*Unauthorized) getTeamsRes()           {}
func (*Unauthorized) getUsersRes()           {}
func (*Unauthorized) patchTeamInstanceRes()  {}
func (*Unauthorized) patchTeamRes()          {}
func (*Unauthorized) postOauth2LogoutRes()   {}
func (*Unauthorized) postTeamRes()           {}

// 部員.
// Ref: #/components/schemas/User
type User struct {
	ID     UserId    `json:"id"`
	Name   UserName  `json:"name"`
	TeamId OptTeamId `json:"teamId"`
	// 管理者権限を持っているかどうか.
	IsAdmin bool `json:"isAdmin"`
}

// GetID returns the value of ID.
func (s *User) GetID() UserId {
	return s.ID
}

// GetName returns the value of Name.
func (s *User) GetName() UserName {
	return s.Name
}

// GetTeamId returns the value of TeamId.
func (s *User) GetTeamId() OptTeamId {
	return s.TeamId
}

// GetIsAdmin returns the value of IsAdmin.
func (s *User) GetIsAdmin() bool {
	return s.IsAdmin
}

// SetID sets the value of ID.
func (s *User) SetID(val UserId) {
	s.ID = val
}

// SetName sets the value of Name.
func (s *User) SetName(val UserName) {
	s.Name = val
}

// SetTeamId sets the value of TeamId.
func (s *User) SetTeamId(val OptTeamId) {
	s.TeamId = val
}

// SetIsAdmin sets the value of IsAdmin.
func (s *User) SetIsAdmin(val bool) {
	s.IsAdmin = val
}

func (*User) getMeRes() {}

type UserAuth struct {
	APIKey string
}

// GetAPIKey returns the value of APIKey.
func (s *UserAuth) GetAPIKey() string {
	return s.APIKey
}

// SetAPIKey sets the value of APIKey.
func (s *UserAuth) SetAPIKey(val string) {
	s.APIKey = val
}

type UserId uuid.UUID

type UserName string