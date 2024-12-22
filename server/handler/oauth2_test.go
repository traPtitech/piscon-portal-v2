package handler_test

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/mysql"
	"github.com/stephenafamo/bob/dialect/mysql/sm"
	"github.com/traPtitech/piscon-portal-v2/server/models"
)

func TestAuthorize(t *testing.T) {
	client := NewClient()
	userName := t.Name()

	// login and create user
	if err := Login(t, client, userName); err != nil {
		t.FailNow()
	}

	// check if the user is created
	exists, err := models.Users.Query(models.SelectWhere.Users.Name.EQ(userName)).Exists(context.TODO(), bob.NewDB(db))
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("user not created")
	}

	// logout
	res, err := client.Post(ServerURL+"/api/oauth2/logout", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(res.Body)
		t.Fatalf("status code is %d: %s", res.StatusCode, msg)
	}

	// logout after logout should be unauthorized
	res, err = client.Post(server.URL+"/api/oauth2/logout", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		msg, _ := io.ReadAll(res.Body)
		t.Fatalf("status code is %d: %s", res.StatusCode, msg)
	}

	// login again
	if err := Login(t, client, userName); err != nil {
		t.FailNow()
	}
}

func TestExpiredSession(t *testing.T) {
	client := NewClient()
	userName := t.Name()

	if err := Login(t, client, userName); err != nil {
		t.FailNow()
	}
	// get session id
	session, err := models.Sessions.Query(
		models.SelectJoins.Sessions.InnerJoin.User(context.Background()),
		sm.Where(mysql.Quote("users", "name").EQ(mysql.S(userName))),
	).One(context.TODO(), bob.NewDB(db))
	if err != nil {
		t.Fatal(err)
	}

	// update expiredAt to past
	_, err = models.Sessions.Update(
		models.UpdateWhere.Sessions.ID.EQ(session.ID),
		models.SessionSetter{ExpiredAt: omit.From(time.Now().Add(-time.Hour))}.UpdateMod(),
	).Exec(context.Background(), bob.NewDB(db))
	if err != nil {
		t.Fatal(err)
	}

	// now, logout should be unauthorized
	res, err := client.Post(server.URL+"/api/oauth2/logout", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		msg, _ := io.ReadAll(res.Body)
		t.Fatalf("status code is %d: %s", res.StatusCode, msg)
	}
}
