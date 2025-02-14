// Code generated by BobGen mysql v0.30.0. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"database/sql"
	"database/sql/driver"

	"github.com/stephenafamo/bob"
)

// Make sure the type BenchmarkLog runs hooks after queries
var _ bob.HookableType = &BenchmarkLog{}

// Make sure the type Benchmark runs hooks after queries
var _ bob.HookableType = &Benchmark{}

// Make sure the type Instance runs hooks after queries
var _ bob.HookableType = &Instance{}

// Make sure the type Session runs hooks after queries
var _ bob.HookableType = &Session{}

// Make sure the type Team runs hooks after queries
var _ bob.HookableType = &Team{}

// Make sure the type User runs hooks after queries
var _ bob.HookableType = &User{}

// Make sure the type BenchmarksStatus satisfies database/sql.Scanner
var _ sql.Scanner = (*BenchmarksStatus)(nil)

// Make sure the type BenchmarksStatus satisfies database/sql/driver.Valuer
var _ driver.Valuer = *new(BenchmarksStatus)

// Make sure the type BenchmarksResult satisfies database/sql.Scanner
var _ sql.Scanner = (*BenchmarksResult)(nil)

// Make sure the type BenchmarksResult satisfies database/sql/driver.Valuer
var _ driver.Valuer = *new(BenchmarksResult)

// Make sure the type InstancesStatus satisfies database/sql.Scanner
var _ sql.Scanner = (*InstancesStatus)(nil)

// Make sure the type InstancesStatus satisfies database/sql/driver.Valuer
var _ driver.Valuer = *new(InstancesStatus)
