// Code generated by BobGen mysql v0.29.0. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package factory

type Factory struct {
	baseTeamMods TeamModSlice
	baseUserMods UserModSlice
}

func New() *Factory {
	return &Factory{}
}

func (f *Factory) NewTeam(mods ...TeamMod) *TeamTemplate {
	o := &TeamTemplate{f: f}

	if f != nil {
		f.baseTeamMods.Apply(o)
	}

	TeamModSlice(mods).Apply(o)

	return o
}

func (f *Factory) NewUser(mods ...UserMod) *UserTemplate {
	o := &UserTemplate{f: f}

	if f != nil {
		f.baseUserMods.Apply(o)
	}

	UserModSlice(mods).Apply(o)

	return o
}

func (f *Factory) ClearBaseTeamMods() {
	f.baseTeamMods = nil
}

func (f *Factory) AddBaseTeamMod(mods ...TeamMod) {
	f.baseTeamMods = append(f.baseTeamMods, mods...)
}

func (f *Factory) ClearBaseUserMods() {
	f.baseUserMods = nil
}

func (f *Factory) AddBaseUserMod(mods ...UserMod) {
	f.baseUserMods = append(f.baseUserMods, mods...)
}
