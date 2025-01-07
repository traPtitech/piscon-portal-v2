// Code generated by BobGen mysql v0.29.0. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package factory

import (
	"context"
	"testing"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/jaswdr/faker/v2"
	"github.com/stephenafamo/bob"
	models "github.com/traPtitech/piscon-portal-v2/server/repository/db/models"
)

type SessionMod interface {
	Apply(*SessionTemplate)
}

type SessionModFunc func(*SessionTemplate)

func (f SessionModFunc) Apply(n *SessionTemplate) {
	f(n)
}

type SessionModSlice []SessionMod

func (mods SessionModSlice) Apply(n *SessionTemplate) {
	for _, f := range mods {
		f.Apply(n)
	}
}

// SessionTemplate is an object representing the database table.
// all columns are optional and should be set by mods
type SessionTemplate struct {
	ID        func() string
	UserID    func() string
	CreatedAt func() time.Time
	ExpiredAt func() time.Time

	r sessionR
	f *Factory
}

type sessionR struct {
	User *sessionRUserR
}

type sessionRUserR struct {
	o *UserTemplate
}

// Apply mods to the SessionTemplate
func (o *SessionTemplate) Apply(mods ...SessionMod) {
	for _, mod := range mods {
		mod.Apply(o)
	}
}

// toModel returns an *models.Session
// this does nothing with the relationship templates
func (o SessionTemplate) toModel() *models.Session {
	m := &models.Session{}

	if o.ID != nil {
		m.ID = o.ID()
	}
	if o.UserID != nil {
		m.UserID = o.UserID()
	}
	if o.CreatedAt != nil {
		m.CreatedAt = o.CreatedAt()
	}
	if o.ExpiredAt != nil {
		m.ExpiredAt = o.ExpiredAt()
	}

	return m
}

// toModels returns an models.SessionSlice
// this does nothing with the relationship templates
func (o SessionTemplate) toModels(number int) models.SessionSlice {
	m := make(models.SessionSlice, number)

	for i := range m {
		m[i] = o.toModel()
	}

	return m
}

// setModelRels creates and sets the relationships on *models.Session
// according to the relationships in the template. Nothing is inserted into the db
func (t SessionTemplate) setModelRels(o *models.Session) {
	if t.r.User != nil {
		rel := t.r.User.o.toModel()
		rel.R.Sessions = append(rel.R.Sessions, o)
		o.UserID = rel.ID
		o.R.User = rel
	}
}

// BuildSetter returns an *models.SessionSetter
// this does nothing with the relationship templates
func (o SessionTemplate) BuildSetter() *models.SessionSetter {
	m := &models.SessionSetter{}

	if o.ID != nil {
		m.ID = omit.From(o.ID())
	}
	if o.UserID != nil {
		m.UserID = omit.From(o.UserID())
	}
	if o.CreatedAt != nil {
		m.CreatedAt = omit.From(o.CreatedAt())
	}
	if o.ExpiredAt != nil {
		m.ExpiredAt = omit.From(o.ExpiredAt())
	}

	return m
}

// BuildManySetter returns an []*models.SessionSetter
// this does nothing with the relationship templates
func (o SessionTemplate) BuildManySetter(number int) []*models.SessionSetter {
	m := make([]*models.SessionSetter, number)

	for i := range m {
		m[i] = o.BuildSetter()
	}

	return m
}

// Build returns an *models.Session
// Related objects are also created and placed in the .R field
// NOTE: Objects are not inserted into the database. Use SessionTemplate.Create
func (o SessionTemplate) Build() *models.Session {
	m := o.toModel()
	o.setModelRels(m)

	return m
}

// BuildMany returns an models.SessionSlice
// Related objects are also created and placed in the .R field
// NOTE: Objects are not inserted into the database. Use SessionTemplate.CreateMany
func (o SessionTemplate) BuildMany(number int) models.SessionSlice {
	m := make(models.SessionSlice, number)

	for i := range m {
		m[i] = o.Build()
	}

	return m
}

func ensureCreatableSession(m *models.SessionSetter) {
	if m.ID.IsUnset() {
		m.ID = omit.From(random_string(nil))
	}
	if m.UserID.IsUnset() {
		m.UserID = omit.From(random_string(nil))
	}
	if m.ExpiredAt.IsUnset() {
		m.ExpiredAt = omit.From(random_time_Time(nil))
	}
}

// insertOptRels creates and inserts any optional the relationships on *models.Session
// according to the relationships in the template.
// any required relationship should have already exist on the model
func (o *SessionTemplate) insertOptRels(ctx context.Context, exec bob.Executor, m *models.Session) (context.Context, error) {
	var err error

	return ctx, err
}

// Create builds a session and inserts it into the database
// Relations objects are also inserted and placed in the .R field
func (o *SessionTemplate) Create(ctx context.Context, exec bob.Executor) (*models.Session, error) {
	_, m, err := o.create(ctx, exec)
	return m, err
}

// MustCreate builds a session and inserts it into the database
// Relations objects are also inserted and placed in the .R field
// panics if an error occurs
func (o *SessionTemplate) MustCreate(ctx context.Context, exec bob.Executor) *models.Session {
	_, m, err := o.create(ctx, exec)
	if err != nil {
		panic(err)
	}
	return m
}

// CreateOrFail builds a session and inserts it into the database
// Relations objects are also inserted and placed in the .R field
// It calls `tb.Fatal(err)` on the test/benchmark if an error occurs
func (o *SessionTemplate) CreateOrFail(ctx context.Context, tb testing.TB, exec bob.Executor) *models.Session {
	tb.Helper()
	_, m, err := o.create(ctx, exec)
	if err != nil {
		tb.Fatal(err)
		return nil
	}
	return m
}

// create builds a session and inserts it into the database
// Relations objects are also inserted and placed in the .R field
// this returns a context that includes the newly inserted model
func (o *SessionTemplate) create(ctx context.Context, exec bob.Executor) (context.Context, *models.Session, error) {
	var err error
	opt := o.BuildSetter()
	ensureCreatableSession(opt)

	var rel0 *models.User
	if o.r.User == nil {
		var ok bool
		rel0, ok = userCtx.Value(ctx)
		if !ok {
			SessionMods.WithNewUser().Apply(o)
		}
	}
	if o.r.User != nil {
		ctx, rel0, err = o.r.User.o.create(ctx, exec)
		if err != nil {
			return ctx, nil, err
		}
	}
	opt.UserID = omit.From(rel0.ID)

	m, err := models.Sessions.Insert(opt).One(ctx, exec)
	if err != nil {
		return ctx, nil, err
	}
	ctx = sessionCtx.WithValue(ctx, m)

	m.R.User = rel0

	ctx, err = o.insertOptRels(ctx, exec, m)
	return ctx, m, err
}

// CreateMany builds multiple sessions and inserts them into the database
// Relations objects are also inserted and placed in the .R field
func (o SessionTemplate) CreateMany(ctx context.Context, exec bob.Executor, number int) (models.SessionSlice, error) {
	_, m, err := o.createMany(ctx, exec, number)
	return m, err
}

// MustCreateMany builds multiple sessions and inserts them into the database
// Relations objects are also inserted and placed in the .R field
// panics if an error occurs
func (o SessionTemplate) MustCreateMany(ctx context.Context, exec bob.Executor, number int) models.SessionSlice {
	_, m, err := o.createMany(ctx, exec, number)
	if err != nil {
		panic(err)
	}
	return m
}

// CreateManyOrFail builds multiple sessions and inserts them into the database
// Relations objects are also inserted and placed in the .R field
// It calls `tb.Fatal(err)` on the test/benchmark if an error occurs
func (o SessionTemplate) CreateManyOrFail(ctx context.Context, tb testing.TB, exec bob.Executor, number int) models.SessionSlice {
	tb.Helper()
	_, m, err := o.createMany(ctx, exec, number)
	if err != nil {
		tb.Fatal(err)
		return nil
	}
	return m
}

// createMany builds multiple sessions and inserts them into the database
// Relations objects are also inserted and placed in the .R field
// this returns a context that includes the newly inserted models
func (o SessionTemplate) createMany(ctx context.Context, exec bob.Executor, number int) (context.Context, models.SessionSlice, error) {
	var err error
	m := make(models.SessionSlice, number)

	for i := range m {
		ctx, m[i], err = o.create(ctx, exec)
		if err != nil {
			return ctx, nil, err
		}
	}

	return ctx, m, nil
}

// Session has methods that act as mods for the SessionTemplate
var SessionMods sessionMods

type sessionMods struct{}

func (m sessionMods) RandomizeAllColumns(f *faker.Faker) SessionMod {
	return SessionModSlice{
		SessionMods.RandomID(f),
		SessionMods.RandomUserID(f),
		SessionMods.RandomCreatedAt(f),
		SessionMods.RandomExpiredAt(f),
	}
}

// Set the model columns to this value
func (m sessionMods) ID(val string) SessionMod {
	return SessionModFunc(func(o *SessionTemplate) {
		o.ID = func() string { return val }
	})
}

// Set the Column from the function
func (m sessionMods) IDFunc(f func() string) SessionMod {
	return SessionModFunc(func(o *SessionTemplate) {
		o.ID = f
	})
}

// Clear any values for the column
func (m sessionMods) UnsetID() SessionMod {
	return SessionModFunc(func(o *SessionTemplate) {
		o.ID = nil
	})
}

// Generates a random value for the column using the given faker
// if faker is nil, a default faker is used
func (m sessionMods) RandomID(f *faker.Faker) SessionMod {
	return SessionModFunc(func(o *SessionTemplate) {
		o.ID = func() string {
			return random_string(f)
		}
	})
}

// Set the model columns to this value
func (m sessionMods) UserID(val string) SessionMod {
	return SessionModFunc(func(o *SessionTemplate) {
		o.UserID = func() string { return val }
	})
}

// Set the Column from the function
func (m sessionMods) UserIDFunc(f func() string) SessionMod {
	return SessionModFunc(func(o *SessionTemplate) {
		o.UserID = f
	})
}

// Clear any values for the column
func (m sessionMods) UnsetUserID() SessionMod {
	return SessionModFunc(func(o *SessionTemplate) {
		o.UserID = nil
	})
}

// Generates a random value for the column using the given faker
// if faker is nil, a default faker is used
func (m sessionMods) RandomUserID(f *faker.Faker) SessionMod {
	return SessionModFunc(func(o *SessionTemplate) {
		o.UserID = func() string {
			return random_string(f)
		}
	})
}

// Set the model columns to this value
func (m sessionMods) CreatedAt(val time.Time) SessionMod {
	return SessionModFunc(func(o *SessionTemplate) {
		o.CreatedAt = func() time.Time { return val }
	})
}

// Set the Column from the function
func (m sessionMods) CreatedAtFunc(f func() time.Time) SessionMod {
	return SessionModFunc(func(o *SessionTemplate) {
		o.CreatedAt = f
	})
}

// Clear any values for the column
func (m sessionMods) UnsetCreatedAt() SessionMod {
	return SessionModFunc(func(o *SessionTemplate) {
		o.CreatedAt = nil
	})
}

// Generates a random value for the column using the given faker
// if faker is nil, a default faker is used
func (m sessionMods) RandomCreatedAt(f *faker.Faker) SessionMod {
	return SessionModFunc(func(o *SessionTemplate) {
		o.CreatedAt = func() time.Time {
			return random_time_Time(f)
		}
	})
}

// Set the model columns to this value
func (m sessionMods) ExpiredAt(val time.Time) SessionMod {
	return SessionModFunc(func(o *SessionTemplate) {
		o.ExpiredAt = func() time.Time { return val }
	})
}

// Set the Column from the function
func (m sessionMods) ExpiredAtFunc(f func() time.Time) SessionMod {
	return SessionModFunc(func(o *SessionTemplate) {
		o.ExpiredAt = f
	})
}

// Clear any values for the column
func (m sessionMods) UnsetExpiredAt() SessionMod {
	return SessionModFunc(func(o *SessionTemplate) {
		o.ExpiredAt = nil
	})
}

// Generates a random value for the column using the given faker
// if faker is nil, a default faker is used
func (m sessionMods) RandomExpiredAt(f *faker.Faker) SessionMod {
	return SessionModFunc(func(o *SessionTemplate) {
		o.ExpiredAt = func() time.Time {
			return random_time_Time(f)
		}
	})
}

func (m sessionMods) WithUser(rel *UserTemplate) SessionMod {
	return SessionModFunc(func(o *SessionTemplate) {
		o.r.User = &sessionRUserR{
			o: rel,
		}
	})
}

func (m sessionMods) WithNewUser(mods ...UserMod) SessionMod {
	return SessionModFunc(func(o *SessionTemplate) {
		related := o.f.NewUser(mods...)

		m.WithUser(related).Apply(o)
	})
}

func (m sessionMods) WithoutUser() SessionMod {
	return SessionModFunc(func(o *SessionTemplate) {
		o.r.User = nil
	})
}
