package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/jacoelho/codewars/internal/user"
	"github.com/jacoelho/codewars/internal/user/api"
)

func TestFetchFromAPI(t *testing.T) {
	const payload = `{
		"username":"someUsername",
		"name":"someName",
		"honor":324,
		"clan":"",
		"leaderboardPosition":59031,
		"skills":[
	 
		],
		"ranks":{
		   "overall":{
			  "rank":-5,
			  "name":"5 kyu",
			  "color":"yellow",
			  "score":307
		   },
		   "languages":{
			  "python":{
				 "rank":-6,
				 "name":"6 kyu",
				 "color":"yellow",
				 "score":209
			  },
			  "javascript":{
				 "rank":-6,
				 "name":"6 kyu",
				 "color":"yellow",
				 "score":100
			  }
		   }
		},
		"codeChallenges":{
		   "totalAuthored":0,
		   "totalCompleted":46
		}
	 }`

	var calls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++

		w.WriteHeader(http.StatusOK)
		//nolint:errcheck
		w.Write([]byte(payload))
	}))
	defer server.Close()

	r := api.New(http.DefaultClient)
	r.Endpoint = server.URL + "/"

	_, err := r.GetUserByID(context.Background(), "123")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	// get again to check calls
	u, err := r.GetUserByID(context.Background(), "123")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	expected := user.User{
		Name:     "someName",
		Username: "someUsername",
		Honor:    324,
	}

	if !reflect.DeepEqual(expected, u) {
		t.Fatalf("expected %v, got %v", expected, u)
	}

	if calls != 1 {
		t.Fatalf("unexpected number of calls %d", calls)
	}
}
