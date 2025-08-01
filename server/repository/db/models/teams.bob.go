// Code generated by BobGen mysql v0.38.0. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"time"

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

// Team is an object representing the database table.
type Team struct {
	ID        string    `db:"id,pk" `
	Name      string    `db:"name" `
	CreatedAt time.Time `db:"created_at" `

	R teamR `db:"-" `
}

// TeamSlice is an alias for a slice of pointers to Team.
// This should almost always be used instead of []*Team.
type TeamSlice []*Team

// Teams contains methods to work with the teams table
var Teams = mysql.NewTablex[*Team, TeamSlice, *TeamSetter]("teams", []string{"id"})

// TeamsQuery is a query on the teams table
type TeamsQuery = *mysql.ViewQuery[*Team, TeamSlice]

// teamR is where relationships are stored.
type teamR struct {
	Users UserSlice // users_ibfk_1
}

type teamColumnNames struct {
	ID        string
	Name      string
	CreatedAt string
}

var TeamColumns = buildTeamColumns("teams")

type teamColumns struct {
	tableAlias string
	ID         mysql.Expression
	Name       mysql.Expression
	CreatedAt  mysql.Expression
}

func (c teamColumns) Alias() string {
	return c.tableAlias
}

func (teamColumns) AliasedAs(alias string) teamColumns {
	return buildTeamColumns(alias)
}

func buildTeamColumns(alias string) teamColumns {
	return teamColumns{
		tableAlias: alias,
		ID:         mysql.Quote(alias, "id"),
		Name:       mysql.Quote(alias, "name"),
		CreatedAt:  mysql.Quote(alias, "created_at"),
	}
}

type teamWhere[Q mysql.Filterable] struct {
	ID        mysql.WhereMod[Q, string]
	Name      mysql.WhereMod[Q, string]
	CreatedAt mysql.WhereMod[Q, time.Time]
}

func (teamWhere[Q]) AliasedAs(alias string) teamWhere[Q] {
	return buildTeamWhere[Q](buildTeamColumns(alias))
}

func buildTeamWhere[Q mysql.Filterable](cols teamColumns) teamWhere[Q] {
	return teamWhere[Q]{
		ID:        mysql.Where[Q, string](cols.ID),
		Name:      mysql.Where[Q, string](cols.Name),
		CreatedAt: mysql.Where[Q, time.Time](cols.CreatedAt),
	}
}

var TeamErrors = &teamErrors{
	ErrUniquePrimary: &UniqueConstraintError{
		schema:  "",
		table:   "teams",
		columns: []string{"id"},
		s:       "PRIMARY",
	},
}

type teamErrors struct {
	ErrUniquePrimary *UniqueConstraintError
}

// TeamSetter is used for insert/upsert/update operations
// All values are optional, and do not have to be set
// Generated columns are not included
type TeamSetter struct {
	ID        *string    `db:"id,pk" `
	Name      *string    `db:"name" `
	CreatedAt *time.Time `db:"created_at" `
}

func (s TeamSetter) SetColumns() []string {
	vals := make([]string, 0, 3)
	if s.ID != nil {
		vals = append(vals, "id")
	}

	if s.Name != nil {
		vals = append(vals, "name")
	}

	if s.CreatedAt != nil {
		vals = append(vals, "created_at")
	}

	return vals
}

func (s TeamSetter) Overwrite(t *Team) {
	if s.ID != nil {
		t.ID = *s.ID
	}
	if s.Name != nil {
		t.Name = *s.Name
	}
	if s.CreatedAt != nil {
		t.CreatedAt = *s.CreatedAt
	}
}

func (s *TeamSetter) Apply(q *dialect.InsertQuery) {
	q.AppendHooks(func(ctx context.Context, exec bob.Executor) (context.Context, error) {
		return Teams.BeforeInsertHooks.RunHooks(ctx, exec, s)
	})

	q.AppendValues(
		bob.ExpressionFunc(func(ctx context.Context, w io.Writer, d bob.Dialect, start int) ([]any, error) {
			if s.ID == nil {
				return mysql.Raw("DEFAULT").WriteSQL(ctx, w, d, start)
			}
			return mysql.Arg(s.ID).WriteSQL(ctx, w, d, start)
		}), bob.ExpressionFunc(func(ctx context.Context, w io.Writer, d bob.Dialect, start int) ([]any, error) {
			if s.Name == nil {
				return mysql.Raw("DEFAULT").WriteSQL(ctx, w, d, start)
			}
			return mysql.Arg(s.Name).WriteSQL(ctx, w, d, start)
		}), bob.ExpressionFunc(func(ctx context.Context, w io.Writer, d bob.Dialect, start int) ([]any, error) {
			if s.CreatedAt == nil {
				return mysql.Raw("DEFAULT").WriteSQL(ctx, w, d, start)
			}
			return mysql.Arg(s.CreatedAt).WriteSQL(ctx, w, d, start)
		}))
}

func (s TeamSetter) UpdateMod() bob.Mod[*dialect.UpdateQuery] {
	return um.Set(s.Expressions("teams")...)
}

func (s TeamSetter) Expressions(prefix ...string) []bob.Expression {
	exprs := make([]bob.Expression, 0, 3)

	if s.ID != nil {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			mysql.Quote(append(prefix, "id")...),
			mysql.Arg(s.ID),
		}})
	}

	if s.Name != nil {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			mysql.Quote(append(prefix, "name")...),
			mysql.Arg(s.Name),
		}})
	}

	if s.CreatedAt != nil {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			mysql.Quote(append(prefix, "created_at")...),
			mysql.Arg(s.CreatedAt),
		}})
	}

	return exprs
}

// FindTeam retrieves a single record by primary key
// If cols is empty Find will return all columns.
func FindTeam(ctx context.Context, exec bob.Executor, IDPK string, cols ...string) (*Team, error) {
	if len(cols) == 0 {
		return Teams.Query(
			SelectWhere.Teams.ID.EQ(IDPK),
		).One(ctx, exec)
	}

	return Teams.Query(
		SelectWhere.Teams.ID.EQ(IDPK),
		sm.Columns(Teams.Columns().Only(cols...)),
	).One(ctx, exec)
}

// TeamExists checks the presence of a single record by primary key
func TeamExists(ctx context.Context, exec bob.Executor, IDPK string) (bool, error) {
	return Teams.Query(
		SelectWhere.Teams.ID.EQ(IDPK),
	).Exists(ctx, exec)
}

// AfterQueryHook is called after Team is retrieved from the database
func (o *Team) AfterQueryHook(ctx context.Context, exec bob.Executor, queryType bob.QueryType) error {
	var err error

	switch queryType {
	case bob.QueryTypeSelect:
		ctx, err = Teams.AfterSelectHooks.RunHooks(ctx, exec, TeamSlice{o})
	case bob.QueryTypeInsert:
		ctx, err = Teams.AfterInsertHooks.RunHooks(ctx, exec, TeamSlice{o})
	case bob.QueryTypeUpdate:
		ctx, err = Teams.AfterUpdateHooks.RunHooks(ctx, exec, TeamSlice{o})
	case bob.QueryTypeDelete:
		ctx, err = Teams.AfterDeleteHooks.RunHooks(ctx, exec, TeamSlice{o})
	}

	return err
}

// primaryKeyVals returns the primary key values of the Team
func (o *Team) primaryKeyVals() bob.Expression {
	return mysql.Arg(o.ID)
}

func (o *Team) pkEQ() dialect.Expression {
	return mysql.Quote("teams", "id").EQ(bob.ExpressionFunc(func(ctx context.Context, w io.Writer, d bob.Dialect, start int) ([]any, error) {
		return o.primaryKeyVals().WriteSQL(ctx, w, d, start)
	}))
}

// Update uses an executor to update the Team
func (o *Team) Update(ctx context.Context, exec bob.Executor, s *TeamSetter) error {
	_, err := Teams.Update(s.UpdateMod(), um.Where(o.pkEQ())).Exec(ctx, exec)
	if err != nil {
		return err
	}

	s.Overwrite(o)

	return nil
}

// Delete deletes a single Team record with an executor
func (o *Team) Delete(ctx context.Context, exec bob.Executor) error {
	_, err := Teams.Delete(dm.Where(o.pkEQ())).Exec(ctx, exec)
	return err
}

// Reload refreshes the Team using the executor
func (o *Team) Reload(ctx context.Context, exec bob.Executor) error {
	o2, err := Teams.Query(
		SelectWhere.Teams.ID.EQ(o.ID),
	).One(ctx, exec)
	if err != nil {
		return err
	}
	o2.R = o.R
	*o = *o2

	return nil
}

// AfterQueryHook is called after TeamSlice is retrieved from the database
func (o TeamSlice) AfterQueryHook(ctx context.Context, exec bob.Executor, queryType bob.QueryType) error {
	var err error

	switch queryType {
	case bob.QueryTypeSelect:
		ctx, err = Teams.AfterSelectHooks.RunHooks(ctx, exec, o)
	case bob.QueryTypeInsert:
		ctx, err = Teams.AfterInsertHooks.RunHooks(ctx, exec, o)
	case bob.QueryTypeUpdate:
		ctx, err = Teams.AfterUpdateHooks.RunHooks(ctx, exec, o)
	case bob.QueryTypeDelete:
		ctx, err = Teams.AfterDeleteHooks.RunHooks(ctx, exec, o)
	}

	return err
}

func (o TeamSlice) pkIN() dialect.Expression {
	if len(o) == 0 {
		return mysql.Raw("NULL")
	}

	return mysql.Quote("teams", "id").In(bob.ExpressionFunc(func(ctx context.Context, w io.Writer, d bob.Dialect, start int) ([]any, error) {
		pkPairs := make([]bob.Expression, len(o))
		for i, row := range o {
			pkPairs[i] = row.primaryKeyVals()
		}
		return bob.ExpressSlice(ctx, w, d, start, pkPairs, "", ", ", "")
	}))
}

// copyMatchingRows finds models in the given slice that have the same primary key
// then it first copies the existing relationships from the old model to the new model
// and then replaces the old model in the slice with the new model
func (o TeamSlice) copyMatchingRows(from ...*Team) {
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
func (o TeamSlice) UpdateMod() bob.Mod[*dialect.UpdateQuery] {
	return bob.ModFunc[*dialect.UpdateQuery](func(q *dialect.UpdateQuery) {
		q.AppendHooks(func(ctx context.Context, exec bob.Executor) (context.Context, error) {
			return Teams.BeforeUpdateHooks.RunHooks(ctx, exec, o)
		})

		q.AppendLoader(bob.LoaderFunc(func(ctx context.Context, exec bob.Executor, retrieved any) error {
			var err error
			switch retrieved := retrieved.(type) {
			case *Team:
				o.copyMatchingRows(retrieved)
			case []*Team:
				o.copyMatchingRows(retrieved...)
			case TeamSlice:
				o.copyMatchingRows(retrieved...)
			default:
				// If the retrieved value is not a Team or a slice of Team
				// then run the AfterUpdateHooks on the slice
				_, err = Teams.AfterUpdateHooks.RunHooks(ctx, exec, o)
			}

			return err
		}))

		q.AppendWhere(o.pkIN())
	})
}

// DeleteMod modifies an delete query with "WHERE primary_key IN (o...)"
func (o TeamSlice) DeleteMod() bob.Mod[*dialect.DeleteQuery] {
	return bob.ModFunc[*dialect.DeleteQuery](func(q *dialect.DeleteQuery) {
		q.AppendHooks(func(ctx context.Context, exec bob.Executor) (context.Context, error) {
			return Teams.BeforeDeleteHooks.RunHooks(ctx, exec, o)
		})

		q.AppendLoader(bob.LoaderFunc(func(ctx context.Context, exec bob.Executor, retrieved any) error {
			var err error
			switch retrieved := retrieved.(type) {
			case *Team:
				o.copyMatchingRows(retrieved)
			case []*Team:
				o.copyMatchingRows(retrieved...)
			case TeamSlice:
				o.copyMatchingRows(retrieved...)
			default:
				// If the retrieved value is not a Team or a slice of Team
				// then run the AfterDeleteHooks on the slice
				_, err = Teams.AfterDeleteHooks.RunHooks(ctx, exec, o)
			}

			return err
		}))

		q.AppendWhere(o.pkIN())
	})
}

func (o TeamSlice) UpdateAll(ctx context.Context, exec bob.Executor, vals TeamSetter) error {
	_, err := Teams.Update(vals.UpdateMod(), o.UpdateMod()).Exec(ctx, exec)

	for i := range o {
		vals.Overwrite(o[i])
	}

	return err
}

func (o TeamSlice) DeleteAll(ctx context.Context, exec bob.Executor) error {
	if len(o) == 0 {
		return nil
	}

	_, err := Teams.Delete(o.DeleteMod()).Exec(ctx, exec)
	return err
}

func (o TeamSlice) ReloadAll(ctx context.Context, exec bob.Executor) error {
	if len(o) == 0 {
		return nil
	}

	o2, err := Teams.Query(sm.Where(o.pkIN())).All(ctx, exec)
	if err != nil {
		return err
	}

	o.copyMatchingRows(o2...)

	return nil
}

type teamJoins[Q dialect.Joinable] struct {
	typ   string
	Users modAs[Q, userColumns]
}

func (j teamJoins[Q]) aliasedAs(alias string) teamJoins[Q] {
	return buildTeamJoins[Q](buildTeamColumns(alias), j.typ)
}

func buildTeamJoins[Q dialect.Joinable](cols teamColumns, typ string) teamJoins[Q] {
	return teamJoins[Q]{
		typ: typ,
		Users: modAs[Q, userColumns]{
			c: UserColumns,
			f: func(to userColumns) bob.Mod[Q] {
				mods := make(mods.QueryMods[Q], 0, 1)

				{
					mods = append(mods, dialect.Join[Q](typ, Users.Name().As(to.Alias())).On(
						to.TeamID.EQ(cols.ID),
					))
				}

				return mods
			},
		},
	}
}

// Users starts a query for related objects on users
func (o *Team) Users(mods ...bob.Mod[*dialect.SelectQuery]) UsersQuery {
	return Users.Query(append(mods,
		sm.Where(UserColumns.TeamID.EQ(mysql.Arg(o.ID))),
	)...)
}

func (os TeamSlice) Users(mods ...bob.Mod[*dialect.SelectQuery]) UsersQuery {
	PKArgSlice := make([]bob.Expression, len(os))
	for i, o := range os {
		PKArgSlice[i] = mysql.ArgGroup(o.ID)
	}
	PKArgExpr := mysql.Group(PKArgSlice...)

	return Users.Query(append(mods,
		sm.Where(mysql.Group(UserColumns.TeamID).OP("IN", PKArgExpr)),
	)...)
}

func (o *Team) Preload(name string, retrieved any) error {
	if o == nil {
		return nil
	}

	switch name {
	case "Users":
		rels, ok := retrieved.(UserSlice)
		if !ok {
			return fmt.Errorf("team cannot load %T as %q", retrieved, name)
		}

		o.R.Users = rels

		for _, rel := range rels {
			if rel != nil {
				rel.R.Team = o
			}
		}
		return nil
	default:
		return fmt.Errorf("team has no relationship %q", name)
	}
}

type teamPreloader struct{}

func buildTeamPreloader() teamPreloader {
	return teamPreloader{}
}

type teamThenLoader[Q orm.Loadable] struct {
	Users func(...bob.Mod[*dialect.SelectQuery]) orm.Loader[Q]
}

func buildTeamThenLoader[Q orm.Loadable]() teamThenLoader[Q] {
	type UsersLoadInterface interface {
		LoadUsers(context.Context, bob.Executor, ...bob.Mod[*dialect.SelectQuery]) error
	}

	return teamThenLoader[Q]{
		Users: thenLoadBuilder[Q](
			"Users",
			func(ctx context.Context, exec bob.Executor, retrieved UsersLoadInterface, mods ...bob.Mod[*dialect.SelectQuery]) error {
				return retrieved.LoadUsers(ctx, exec, mods...)
			},
		),
	}
}

// LoadUsers loads the team's Users into the .R struct
func (o *Team) LoadUsers(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) error {
	if o == nil {
		return nil
	}

	// Reset the relationship
	o.R.Users = nil

	related, err := o.Users(mods...).All(ctx, exec)
	if err != nil {
		return err
	}

	for _, rel := range related {
		rel.R.Team = o
	}

	o.R.Users = related
	return nil
}

// LoadUsers loads the team's Users into the .R struct
func (os TeamSlice) LoadUsers(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) error {
	if len(os) == 0 {
		return nil
	}

	users, err := os.Users(mods...).All(ctx, exec)
	if err != nil {
		return err
	}

	for _, o := range os {
		o.R.Users = nil
	}

	for _, o := range os {
		for _, rel := range users {
			if o.ID != rel.TeamID.V {
				continue
			}

			rel.R.Team = o

			o.R.Users = append(o.R.Users, rel)
		}
	}

	return nil
}

func insertTeamUsers0(ctx context.Context, exec bob.Executor, users1 []*UserSetter, team0 *Team) (UserSlice, error) {
	for i := range users1 {
		users1[i].TeamID = func() *sql.Null[string] {
			v := sql.Null[string]{V: team0.ID, Valid: true}
			return &v
		}()
	}

	ret, err := Users.Insert(bob.ToMods(users1...)).All(ctx, exec)
	if err != nil {
		return ret, fmt.Errorf("insertTeamUsers0: %w", err)
	}

	return ret, nil
}

func attachTeamUsers0(ctx context.Context, exec bob.Executor, count int, users1 UserSlice, team0 *Team) (UserSlice, error) {
	setter := &UserSetter{
		TeamID: func() *sql.Null[string] {
			v := sql.Null[string]{V: team0.ID, Valid: true}
			return &v
		}(),
	}

	err := users1.UpdateAll(ctx, exec, *setter)
	if err != nil {
		return nil, fmt.Errorf("attachTeamUsers0: %w", err)
	}

	return users1, nil
}

func (team0 *Team) InsertUsers(ctx context.Context, exec bob.Executor, related ...*UserSetter) error {
	if len(related) == 0 {
		return nil
	}

	var err error

	users1, err := insertTeamUsers0(ctx, exec, related, team0)
	if err != nil {
		return err
	}

	team0.R.Users = append(team0.R.Users, users1...)

	for _, rel := range users1 {
		rel.R.Team = team0
	}
	return nil
}

func (team0 *Team) AttachUsers(ctx context.Context, exec bob.Executor, related ...*User) error {
	if len(related) == 0 {
		return nil
	}

	var err error
	users1 := UserSlice(related)

	_, err = attachTeamUsers0(ctx, exec, len(related), users1, team0)
	if err != nil {
		return err
	}

	team0.R.Users = append(team0.R.Users, users1...)

	for _, rel := range related {
		rel.R.Team = team0
	}

	return nil
}
