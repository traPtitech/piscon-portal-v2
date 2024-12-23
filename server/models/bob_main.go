// Code generated by BobGen mysql v0.29.0. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"hash/maphash"
	"strings"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/clause"
	"github.com/stephenafamo/bob/dialect/mysql"
	"github.com/stephenafamo/bob/dialect/mysql/dialect"
)

var TableNames = struct {
	Teams string
	Users string
}{
	Teams: "teams",
	Users: "users",
}

var ColumnNames = struct {
	Teams teamColumnNames
	Users userColumnNames
}{
	Teams: teamColumnNames{
		ID:        "id",
		Name:      "name",
		CreatedAt: "created_at",
	},
	Users: userColumnNames{
		ID:      "id",
		Name:    "name",
		TeamID:  "team_id",
		IsAdmin: "is_admin",
	},
}

var (
	SelectWhere = Where[*dialect.SelectQuery]()
	UpdateWhere = Where[*dialect.UpdateQuery]()
	DeleteWhere = Where[*dialect.DeleteQuery]()
)

func Where[Q mysql.Filterable]() struct {
	Teams teamWhere[Q]
	Users userWhere[Q]
} {
	return struct {
		Teams teamWhere[Q]
		Users userWhere[Q]
	}{
		Teams: buildTeamWhere[Q](TeamColumns),
		Users: buildUserWhere[Q](UserColumns),
	}
}

var (
	SelectJoins = getJoins[*dialect.SelectQuery]()
	UpdateJoins = getJoins[*dialect.UpdateQuery]()
	DeleteJoins = getJoins[*dialect.DeleteQuery]()
)

type joinSet[Q interface{ aliasedAs(string) Q }] struct {
	InnerJoin Q
	LeftJoin  Q
	RightJoin Q
}

func (j joinSet[Q]) AliasedAs(alias string) joinSet[Q] {
	return joinSet[Q]{
		InnerJoin: j.InnerJoin.aliasedAs(alias),
		LeftJoin:  j.LeftJoin.aliasedAs(alias),
		RightJoin: j.RightJoin.aliasedAs(alias),
	}
}

type joins[Q dialect.Joinable] struct {
	Teams joinSet[teamJoins[Q]]
	Users joinSet[userJoins[Q]]
}

func buildJoinSet[Q interface{ aliasedAs(string) Q }, C any, F func(C, string) Q](c C, f F) joinSet[Q] {
	return joinSet[Q]{
		InnerJoin: f(c, clause.InnerJoin),
		LeftJoin:  f(c, clause.LeftJoin),
		RightJoin: f(c, clause.RightJoin),
	}
}

func getJoins[Q dialect.Joinable]() joins[Q] {
	return joins[Q]{
		Teams: buildJoinSet[teamJoins[Q]](TeamColumns, buildTeamJoins),
		Users: buildJoinSet[userJoins[Q]](UserColumns, buildUserJoins),
	}
}

type modAs[Q any, C interface{ AliasedAs(string) C }] struct {
	c C
	f func(C) bob.Mod[Q]
}

func (m modAs[Q, C]) Apply(q Q) {
	m.f(m.c).Apply(q)
}

func (m modAs[Q, C]) AliasedAs(alias string) bob.Mod[Q] {
	m.c = m.c.AliasedAs(alias)
	return m
}

func randInt() int64 {
	out := int64(new(maphash.Hash).Sum64())

	if out < 0 {
		return -out % 10000
	}

	return out % 10000
}

// ErrUniqueConstraint captures all unique constraint errors by explicitly leaving `s` empty.
var ErrUniqueConstraint = &errUniqueConstraint{s: ""}

type errUniqueConstraint struct {
	// s is a string uniquely identifying the constraint in the raw error message returned from the database.
	s string
}

func (e *errUniqueConstraint) Error() string {
	return e.s
}

func (e *errUniqueConstraint) Is(target error) bool {
	err, ok := target.(*mysqlDriver.MySQLError)
	if !ok {
		return false
	}
	return err.Number == 1062 && strings.Contains(err.Message, e.s)
}
