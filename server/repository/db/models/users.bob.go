// Code generated by BobGen mysql v0.34.2. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"

	"github.com/aarondl/opt/null"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/mysql"
	"github.com/stephenafamo/bob/dialect/mysql/dialect"
	"github.com/stephenafamo/bob/dialect/mysql/dm"
	"github.com/stephenafamo/bob/dialect/mysql/sm"
	"github.com/stephenafamo/bob/dialect/mysql/um"
	"github.com/stephenafamo/bob/expr"
	"github.com/stephenafamo/bob/mods"
	"github.com/stephenafamo/bob/orm"
)

// User is an object representing the database table.
type User struct {
	ID      string           `db:"id,pk" `
	Name    string           `db:"name" `
	TeamID  null.Val[string] `db:"team_id" `
	IsAdmin bool             `db:"is_admin" `

	R userR `db:"-" `
}

// UserSlice is an alias for a slice of pointers to User.
// This should almost always be used instead of []*User.
type UserSlice []*User

// Users contains methods to work with the users table
var Users = mysql.NewTablex[*User, UserSlice, *UserSetter]("users", []string{"name"}, []string{"id"})

// UsersQuery is a query on the users table
type UsersQuery = *mysql.ViewQuery[*User, UserSlice]

// userR is where relationships are stored.
type userR struct {
	Sessions SessionSlice // sessions_ibfk_1
	Team     *Team        // users_ibfk_1
}

type userColumnNames struct {
	ID      string
	Name    string
	TeamID  string
	IsAdmin string
}

var UserColumns = buildUserColumns("users")

type userColumns struct {
	tableAlias string
	ID         mysql.Expression
	Name       mysql.Expression
	TeamID     mysql.Expression
	IsAdmin    mysql.Expression
}

func (c userColumns) Alias() string {
	return c.tableAlias
}

func (userColumns) AliasedAs(alias string) userColumns {
	return buildUserColumns(alias)
}

func buildUserColumns(alias string) userColumns {
	return userColumns{
		tableAlias: alias,
		ID:         mysql.Quote(alias, "id"),
		Name:       mysql.Quote(alias, "name"),
		TeamID:     mysql.Quote(alias, "team_id"),
		IsAdmin:    mysql.Quote(alias, "is_admin"),
	}
}

type userWhere[Q mysql.Filterable] struct {
	ID      mysql.WhereMod[Q, string]
	Name    mysql.WhereMod[Q, string]
	TeamID  mysql.WhereNullMod[Q, string]
	IsAdmin mysql.WhereMod[Q, bool]
}

func (userWhere[Q]) AliasedAs(alias string) userWhere[Q] {
	return buildUserWhere[Q](buildUserColumns(alias))
}

func buildUserWhere[Q mysql.Filterable](cols userColumns) userWhere[Q] {
	return userWhere[Q]{
		ID:      mysql.Where[Q, string](cols.ID),
		Name:    mysql.Where[Q, string](cols.Name),
		TeamID:  mysql.WhereNull[Q, string](cols.TeamID),
		IsAdmin: mysql.Where[Q, bool](cols.IsAdmin),
	}
}

var UserErrors = &userErrors{
	ErrUniqueName: &UniqueConstraintError{s: "name"},

	ErrUniquePrimary: &UniqueConstraintError{s: "PRIMARY"},
}

type userErrors struct {
	ErrUniqueName *UniqueConstraintError

	ErrUniquePrimary *UniqueConstraintError
}

// UserSetter is used for insert/upsert/update operations
// All values are optional, and do not have to be set
// Generated columns are not included
type UserSetter struct {
	ID      omit.Val[string]     `db:"id,pk" `
	Name    omit.Val[string]     `db:"name" `
	TeamID  omitnull.Val[string] `db:"team_id" `
	IsAdmin omit.Val[bool]       `db:"is_admin" `
}

func (s UserSetter) SetColumns() []string {
	vals := make([]string, 0, 4)
	if !s.ID.IsUnset() {
		vals = append(vals, "id")
	}

	if !s.Name.IsUnset() {
		vals = append(vals, "name")
	}

	if !s.TeamID.IsUnset() {
		vals = append(vals, "team_id")
	}

	if !s.IsAdmin.IsUnset() {
		vals = append(vals, "is_admin")
	}

	return vals
}

func (s UserSetter) Overwrite(t *User) {
	if !s.ID.IsUnset() {
		t.ID, _ = s.ID.Get()
	}
	if !s.Name.IsUnset() {
		t.Name, _ = s.Name.Get()
	}
	if !s.TeamID.IsUnset() {
		t.TeamID, _ = s.TeamID.GetNull()
	}
	if !s.IsAdmin.IsUnset() {
		t.IsAdmin, _ = s.IsAdmin.Get()
	}
}

func (s *UserSetter) Apply(q *dialect.InsertQuery) {
	q.AppendHooks(func(ctx context.Context, exec bob.Executor) (context.Context, error) {
		return Users.BeforeInsertHooks.RunHooks(ctx, exec, s)
	})

	q.AppendValues(
		bob.ExpressionFunc(func(ctx context.Context, w io.Writer, d bob.Dialect, start int) ([]any, error) {
			if s.ID.IsUnset() {
				return mysql.Raw("DEFAULT").WriteSQL(ctx, w, d, start)
			}
			return mysql.Arg(s.ID).WriteSQL(ctx, w, d, start)
		}), bob.ExpressionFunc(func(ctx context.Context, w io.Writer, d bob.Dialect, start int) ([]any, error) {
			if s.Name.IsUnset() {
				return mysql.Raw("DEFAULT").WriteSQL(ctx, w, d, start)
			}
			return mysql.Arg(s.Name).WriteSQL(ctx, w, d, start)
		}), bob.ExpressionFunc(func(ctx context.Context, w io.Writer, d bob.Dialect, start int) ([]any, error) {
			if s.TeamID.IsUnset() {
				return mysql.Raw("DEFAULT").WriteSQL(ctx, w, d, start)
			}
			return mysql.Arg(s.TeamID).WriteSQL(ctx, w, d, start)
		}), bob.ExpressionFunc(func(ctx context.Context, w io.Writer, d bob.Dialect, start int) ([]any, error) {
			if s.IsAdmin.IsUnset() {
				return mysql.Raw("DEFAULT").WriteSQL(ctx, w, d, start)
			}
			return mysql.Arg(s.IsAdmin).WriteSQL(ctx, w, d, start)
		}))
}

func (s UserSetter) UpdateMod() bob.Mod[*dialect.UpdateQuery] {
	return um.Set(s.Expressions("users")...)
}

func (s UserSetter) Expressions(prefix ...string) []bob.Expression {
	exprs := make([]bob.Expression, 0, 4)

	if !s.ID.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			mysql.Quote(append(prefix, "id")...),
			mysql.Arg(s.ID),
		}})
	}

	if !s.Name.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			mysql.Quote(append(prefix, "name")...),
			mysql.Arg(s.Name),
		}})
	}

	if !s.TeamID.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			mysql.Quote(append(prefix, "team_id")...),
			mysql.Arg(s.TeamID),
		}})
	}

	if !s.IsAdmin.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			mysql.Quote(append(prefix, "is_admin")...),
			mysql.Arg(s.IsAdmin),
		}})
	}

	return exprs
}

// FindUser retrieves a single record by primary key
// If cols is empty Find will return all columns.
func FindUser(ctx context.Context, exec bob.Executor, IDPK string, cols ...string) (*User, error) {
	if len(cols) == 0 {
		return Users.Query(
			SelectWhere.Users.ID.EQ(IDPK),
		).One(ctx, exec)
	}

	return Users.Query(
		SelectWhere.Users.ID.EQ(IDPK),
		sm.Columns(Users.Columns().Only(cols...)),
	).One(ctx, exec)
}

// UserExists checks the presence of a single record by primary key
func UserExists(ctx context.Context, exec bob.Executor, IDPK string) (bool, error) {
	return Users.Query(
		SelectWhere.Users.ID.EQ(IDPK),
	).Exists(ctx, exec)
}

// AfterQueryHook is called after User is retrieved from the database
func (o *User) AfterQueryHook(ctx context.Context, exec bob.Executor, queryType bob.QueryType) error {
	var err error

	switch queryType {
	case bob.QueryTypeSelect:
		ctx, err = Users.AfterSelectHooks.RunHooks(ctx, exec, UserSlice{o})
	case bob.QueryTypeInsert:
		ctx, err = Users.AfterInsertHooks.RunHooks(ctx, exec, UserSlice{o})
	case bob.QueryTypeUpdate:
		ctx, err = Users.AfterUpdateHooks.RunHooks(ctx, exec, UserSlice{o})
	case bob.QueryTypeDelete:
		ctx, err = Users.AfterDeleteHooks.RunHooks(ctx, exec, UserSlice{o})
	}

	return err
}

// PrimaryKeyVals returns the primary key values of the User
func (o *User) PrimaryKeyVals() bob.Expression {
	return mysql.Arg(o.ID)
}

func (o *User) pkEQ() dialect.Expression {
	return mysql.Quote("users", "id").EQ(bob.ExpressionFunc(func(ctx context.Context, w io.Writer, d bob.Dialect, start int) ([]any, error) {
		return o.PrimaryKeyVals().WriteSQL(ctx, w, d, start)
	}))
}

// Update uses an executor to update the User
func (o *User) Update(ctx context.Context, exec bob.Executor, s *UserSetter) error {
	_, err := Users.Update(s.UpdateMod(), um.Where(o.pkEQ())).Exec(ctx, exec)
	if err != nil {
		return err
	}

	s.Overwrite(o)

	return nil
}

// Delete deletes a single User record with an executor
func (o *User) Delete(ctx context.Context, exec bob.Executor) error {
	_, err := Users.Delete(dm.Where(o.pkEQ())).Exec(ctx, exec)
	return err
}

// Reload refreshes the User using the executor
func (o *User) Reload(ctx context.Context, exec bob.Executor) error {
	o2, err := Users.Query(
		SelectWhere.Users.ID.EQ(o.ID),
	).One(ctx, exec)
	if err != nil {
		return err
	}
	o2.R = o.R
	*o = *o2

	return nil
}

// AfterQueryHook is called after UserSlice is retrieved from the database
func (o UserSlice) AfterQueryHook(ctx context.Context, exec bob.Executor, queryType bob.QueryType) error {
	var err error

	switch queryType {
	case bob.QueryTypeSelect:
		ctx, err = Users.AfterSelectHooks.RunHooks(ctx, exec, o)
	case bob.QueryTypeInsert:
		ctx, err = Users.AfterInsertHooks.RunHooks(ctx, exec, o)
	case bob.QueryTypeUpdate:
		ctx, err = Users.AfterUpdateHooks.RunHooks(ctx, exec, o)
	case bob.QueryTypeDelete:
		ctx, err = Users.AfterDeleteHooks.RunHooks(ctx, exec, o)
	}

	return err
}

func (o UserSlice) pkIN() dialect.Expression {
	if len(o) == 0 {
		return mysql.Raw("NULL")
	}

	return mysql.Quote("users", "id").In(bob.ExpressionFunc(func(ctx context.Context, w io.Writer, d bob.Dialect, start int) ([]any, error) {
		pkPairs := make([]bob.Expression, len(o))
		for i, row := range o {
			pkPairs[i] = row.PrimaryKeyVals()
		}
		return bob.ExpressSlice(ctx, w, d, start, pkPairs, "", ", ", "")
	}))
}

// copyMatchingRows finds models in the given slice that have the same primary key
// then it first copies the existing relationships from the old model to the new model
// and then replaces the old model in the slice with the new model
func (o UserSlice) copyMatchingRows(from ...*User) {
	for i, old := range o {
		for _, new := range from {
			if new.ID != old.ID {
				continue
			}
			new.R = old.R
			o[i] = new
			break
		}
	}
}

// UpdateMod modifies an update query with "WHERE primary_key IN (o...)"
func (o UserSlice) UpdateMod() bob.Mod[*dialect.UpdateQuery] {
	return bob.ModFunc[*dialect.UpdateQuery](func(q *dialect.UpdateQuery) {
		q.AppendHooks(func(ctx context.Context, exec bob.Executor) (context.Context, error) {
			return Users.BeforeUpdateHooks.RunHooks(ctx, exec, o)
		})

		q.AppendLoader(bob.LoaderFunc(func(ctx context.Context, exec bob.Executor, retrieved any) error {
			var err error
			switch retrieved := retrieved.(type) {
			case *User:
				o.copyMatchingRows(retrieved)
			case []*User:
				o.copyMatchingRows(retrieved...)
			case UserSlice:
				o.copyMatchingRows(retrieved...)
			default:
				// If the retrieved value is not a User or a slice of User
				// then run the AfterUpdateHooks on the slice
				_, err = Users.AfterUpdateHooks.RunHooks(ctx, exec, o)
			}

			return err
		}))

		q.AppendWhere(o.pkIN())
	})
}

// DeleteMod modifies an delete query with "WHERE primary_key IN (o...)"
func (o UserSlice) DeleteMod() bob.Mod[*dialect.DeleteQuery] {
	return bob.ModFunc[*dialect.DeleteQuery](func(q *dialect.DeleteQuery) {
		q.AppendHooks(func(ctx context.Context, exec bob.Executor) (context.Context, error) {
			return Users.BeforeDeleteHooks.RunHooks(ctx, exec, o)
		})

		q.AppendLoader(bob.LoaderFunc(func(ctx context.Context, exec bob.Executor, retrieved any) error {
			var err error
			switch retrieved := retrieved.(type) {
			case *User:
				o.copyMatchingRows(retrieved)
			case []*User:
				o.copyMatchingRows(retrieved...)
			case UserSlice:
				o.copyMatchingRows(retrieved...)
			default:
				// If the retrieved value is not a User or a slice of User
				// then run the AfterDeleteHooks on the slice
				_, err = Users.AfterDeleteHooks.RunHooks(ctx, exec, o)
			}

			return err
		}))

		q.AppendWhere(o.pkIN())
	})
}

func (o UserSlice) UpdateAll(ctx context.Context, exec bob.Executor, vals UserSetter) error {
	_, err := Users.Update(vals.UpdateMod(), o.UpdateMod()).Exec(ctx, exec)

	for i := range o {
		vals.Overwrite(o[i])
	}

	return err
}

func (o UserSlice) DeleteAll(ctx context.Context, exec bob.Executor) error {
	if len(o) == 0 {
		return nil
	}

	_, err := Users.Delete(o.DeleteMod()).Exec(ctx, exec)
	return err
}

func (o UserSlice) ReloadAll(ctx context.Context, exec bob.Executor) error {
	if len(o) == 0 {
		return nil
	}

	o2, err := Users.Query(sm.Where(o.pkIN())).All(ctx, exec)
	if err != nil {
		return err
	}

	o.copyMatchingRows(o2...)

	return nil
}

type userJoins[Q dialect.Joinable] struct {
	typ      string
	Sessions func(context.Context) modAs[Q, sessionColumns]
	Team     func(context.Context) modAs[Q, teamColumns]
}

func (j userJoins[Q]) aliasedAs(alias string) userJoins[Q] {
	return buildUserJoins[Q](buildUserColumns(alias), j.typ)
}

func buildUserJoins[Q dialect.Joinable](cols userColumns, typ string) userJoins[Q] {
	return userJoins[Q]{
		typ:      typ,
		Sessions: usersJoinSessions[Q](cols, typ),
		Team:     usersJoinTeam[Q](cols, typ),
	}
}

func usersJoinSessions[Q dialect.Joinable](from userColumns, typ string) func(context.Context) modAs[Q, sessionColumns] {
	return func(ctx context.Context) modAs[Q, sessionColumns] {
		return modAs[Q, sessionColumns]{
			c: SessionColumns,
			f: func(to sessionColumns) bob.Mod[Q] {
				mods := make(mods.QueryMods[Q], 0, 1)

				{
					mods = append(mods, dialect.Join[Q](typ, Sessions.Name().As(to.Alias())).On(
						to.UserID.EQ(from.ID),
					))
				}

				return mods
			},
		}
	}
}

func usersJoinTeam[Q dialect.Joinable](from userColumns, typ string) func(context.Context) modAs[Q, teamColumns] {
	return func(ctx context.Context) modAs[Q, teamColumns] {
		return modAs[Q, teamColumns]{
			c: TeamColumns,
			f: func(to teamColumns) bob.Mod[Q] {
				mods := make(mods.QueryMods[Q], 0, 1)

				{
					mods = append(mods, dialect.Join[Q](typ, Teams.Name().As(to.Alias())).On(
						to.ID.EQ(from.TeamID),
					))
				}

				return mods
			},
		}
	}
}

// Sessions starts a query for related objects on sessions
func (o *User) Sessions(mods ...bob.Mod[*dialect.SelectQuery]) SessionsQuery {
	return Sessions.Query(append(mods,
		sm.Where(SessionColumns.UserID.EQ(mysql.Arg(o.ID))),
	)...)
}

func (os UserSlice) Sessions(mods ...bob.Mod[*dialect.SelectQuery]) SessionsQuery {
	PKArgs := make([]bob.Expression, len(os))
	for i, o := range os {
		PKArgs[i] = mysql.ArgGroup(o.ID)
	}

	return Sessions.Query(append(mods,
		sm.Where(mysql.Group(SessionColumns.UserID).In(PKArgs...)),
	)...)
}

// Team starts a query for related objects on teams
func (o *User) Team(mods ...bob.Mod[*dialect.SelectQuery]) TeamsQuery {
	return Teams.Query(append(mods,
		sm.Where(TeamColumns.ID.EQ(mysql.Arg(o.TeamID))),
	)...)
}

func (os UserSlice) Team(mods ...bob.Mod[*dialect.SelectQuery]) TeamsQuery {
	PKArgs := make([]bob.Expression, len(os))
	for i, o := range os {
		PKArgs[i] = mysql.ArgGroup(o.TeamID)
	}

	return Teams.Query(append(mods,
		sm.Where(mysql.Group(TeamColumns.ID).In(PKArgs...)),
	)...)
}

func (o *User) Preload(name string, retrieved any) error {
	if o == nil {
		return nil
	}

	switch name {
	case "Sessions":
		rels, ok := retrieved.(SessionSlice)
		if !ok {
			return fmt.Errorf("user cannot load %T as %q", retrieved, name)
		}

		o.R.Sessions = rels

		for _, rel := range rels {
			if rel != nil {
				rel.R.User = o
			}
		}
		return nil
	case "Team":
		rel, ok := retrieved.(*Team)
		if !ok {
			return fmt.Errorf("user cannot load %T as %q", retrieved, name)
		}

		o.R.Team = rel

		if rel != nil {
			rel.R.Users = UserSlice{o}
		}
		return nil
	default:
		return fmt.Errorf("user has no relationship %q", name)
	}
}

func ThenLoadUserSessions(queryMods ...bob.Mod[*dialect.SelectQuery]) mysql.Loader {
	return mysql.Loader(func(ctx context.Context, exec bob.Executor, retrieved any) error {
		loader, isLoader := retrieved.(interface {
			LoadUserSessions(context.Context, bob.Executor, ...bob.Mod[*dialect.SelectQuery]) error
		})
		if !isLoader {
			return fmt.Errorf("object %T cannot load UserSessions", retrieved)
		}

		err := loader.LoadUserSessions(ctx, exec, queryMods...)

		// Don't cause an issue due to missing relationships
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return err
	})
}

// LoadUserSessions loads the user's Sessions into the .R struct
func (o *User) LoadUserSessions(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) error {
	if o == nil {
		return nil
	}

	// Reset the relationship
	o.R.Sessions = nil

	related, err := o.Sessions(mods...).All(ctx, exec)
	if err != nil {
		return err
	}

	for _, rel := range related {
		rel.R.User = o
	}

	o.R.Sessions = related
	return nil
}

// LoadUserSessions loads the user's Sessions into the .R struct
func (os UserSlice) LoadUserSessions(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) error {
	if len(os) == 0 {
		return nil
	}

	sessions, err := os.Sessions(mods...).All(ctx, exec)
	if err != nil {
		return err
	}

	for _, o := range os {
		o.R.Sessions = nil
	}

	for _, o := range os {
		for _, rel := range sessions {
			if o.ID != rel.UserID {
				continue
			}

			rel.R.User = o

			o.R.Sessions = append(o.R.Sessions, rel)
		}
	}

	return nil
}

func PreloadUserTeam(opts ...mysql.PreloadOption) mysql.Preloader {
	return mysql.Preload[*Team, TeamSlice](orm.Relationship{
		Name: "Team",
		Sides: []orm.RelSide{
			{
				From: TableNames.Users,
				To:   TableNames.Teams,
				FromColumns: []string{
					ColumnNames.Users.TeamID,
				},
				ToColumns: []string{
					ColumnNames.Teams.ID,
				},
			},
		},
	}, Teams.Columns().Names(), opts...)
}

func ThenLoadUserTeam(queryMods ...bob.Mod[*dialect.SelectQuery]) mysql.Loader {
	return mysql.Loader(func(ctx context.Context, exec bob.Executor, retrieved any) error {
		loader, isLoader := retrieved.(interface {
			LoadUserTeam(context.Context, bob.Executor, ...bob.Mod[*dialect.SelectQuery]) error
		})
		if !isLoader {
			return fmt.Errorf("object %T cannot load UserTeam", retrieved)
		}

		err := loader.LoadUserTeam(ctx, exec, queryMods...)

		// Don't cause an issue due to missing relationships
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return err
	})
}

// LoadUserTeam loads the user's Team into the .R struct
func (o *User) LoadUserTeam(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) error {
	if o == nil {
		return nil
	}

	// Reset the relationship
	o.R.Team = nil

	related, err := o.Team(mods...).One(ctx, exec)
	if err != nil {
		return err
	}

	related.R.Users = UserSlice{o}

	o.R.Team = related
	return nil
}

// LoadUserTeam loads the user's Team into the .R struct
func (os UserSlice) LoadUserTeam(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) error {
	if len(os) == 0 {
		return nil
	}

	teams, err := os.Team(mods...).All(ctx, exec)
	if err != nil {
		return err
	}

	for _, o := range os {
		for _, rel := range teams {
			if o.TeamID.GetOrZero() != rel.ID {
				continue
			}

			rel.R.Users = append(rel.R.Users, o)

			o.R.Team = rel
			break
		}
	}

	return nil
}

func insertUserSessions0(ctx context.Context, exec bob.Executor, sessions1 []*SessionSetter, user0 *User) (SessionSlice, error) {
	for i := range sessions1 {
		sessions1[i].UserID = omit.From(user0.ID)
	}

	ret, err := Sessions.Insert(bob.ToMods(sessions1...)).All(ctx, exec)
	if err != nil {
		return ret, fmt.Errorf("insertUserSessions0: %w", err)
	}

	return ret, nil
}

func attachUserSessions0(ctx context.Context, exec bob.Executor, count int, sessions1 SessionSlice, user0 *User) (SessionSlice, error) {
	setter := &SessionSetter{
		UserID: omit.From(user0.ID),
	}

	err := sessions1.UpdateAll(ctx, exec, *setter)
	if err != nil {
		return nil, fmt.Errorf("attachUserSessions0: %w", err)
	}

	return sessions1, nil
}

func (user0 *User) InsertSessions(ctx context.Context, exec bob.Executor, related ...*SessionSetter) error {
	if len(related) == 0 {
		return nil
	}

	var err error

	sessions1, err := insertUserSessions0(ctx, exec, related, user0)
	if err != nil {
		return err
	}

	user0.R.Sessions = append(user0.R.Sessions, sessions1...)

	for _, rel := range sessions1 {
		rel.R.User = user0
	}
	return nil
}

func (user0 *User) AttachSessions(ctx context.Context, exec bob.Executor, related ...*Session) error {
	if len(related) == 0 {
		return nil
	}

	var err error
	sessions1 := SessionSlice(related)

	_, err = attachUserSessions0(ctx, exec, len(related), sessions1, user0)
	if err != nil {
		return err
	}

	user0.R.Sessions = append(user0.R.Sessions, sessions1...)

	for _, rel := range related {
		rel.R.User = user0
	}

	return nil
}

func attachUserTeam0(ctx context.Context, exec bob.Executor, count int, user0 *User, team1 *Team) (*User, error) {
	setter := &UserSetter{
		TeamID: omitnull.From(team1.ID),
	}

	err := user0.Update(ctx, exec, setter)
	if err != nil {
		return nil, fmt.Errorf("attachUserTeam0: %w", err)
	}

	return user0, nil
}

func (user0 *User) InsertTeam(ctx context.Context, exec bob.Executor, related *TeamSetter) error {
	team1, err := Teams.Insert(related).One(ctx, exec)
	if err != nil {
		return fmt.Errorf("inserting related objects: %w", err)
	}

	_, err = attachUserTeam0(ctx, exec, 1, user0, team1)
	if err != nil {
		return err
	}

	user0.R.Team = team1

	team1.R.Users = append(team1.R.Users, user0)

	return nil
}

func (user0 *User) AttachTeam(ctx context.Context, exec bob.Executor, team1 *Team) error {
	var err error

	_, err = attachUserTeam0(ctx, exec, 1, user0, team1)
	if err != nil {
		return err
	}

	user0.R.Team = team1

	team1.R.Users = append(team1.R.Users, user0)

	return nil
}
