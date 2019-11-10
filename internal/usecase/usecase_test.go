package usecase_test

import (
	"context"
	"testing"

	"github.com/jacoelho/codewars/internal/usecase"
	"github.com/jacoelho/codewars/internal/user"
)

type fakeNotifier func(user, message string) error

func (f fakeNotifier) Notify(ctx context.Context, user, message string) error {
	return f(user, message)
}

type fakeRepo func(id string) (user.User, error)

func (f fakeRepo) GetUserByID(ctx context.Context, id string) (user.User, error) {
	return f(id)
}

func TestUserHonorUpdatedEmptyEvent(t *testing.T) {
	f := fakeNotifier(func(u, m string) error {
		return nil
	})

	r := fakeRepo(func(id string) (user.User, error) {
		return user.User{
			Username: "test",
		}, nil
	})

	err := usecase.UserHonorUpdatedCase(f, r)(context.Background(), user.Event{})

	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestUserHonorUpdated(t *testing.T) {
	var (
		recipient string
		message   string
	)

	f := fakeNotifier(func(u, m string) error {
		recipient = u
		message = m
		return nil
	})

	r := fakeRepo(func(id string) (user.User, error) {
		return user.User{
			Username: "userFromRepo",
		}, nil
	})

	err := usecase.UserHonorUpdatedCase(f, r)(
		context.Background(),
		user.Event{
			HonorUpgraded: &user.HonorUpgraded{
				ID:         "123",
				Honor:      100,
				HonorDelta: 2,
			},
		})

	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if want := "userFromRepo"; want != recipient {
		t.Fatalf("unexpected user, want %s, got %s", want, recipient)
	}

	if want := ":arrow_up:  Honor went up to 100"; want != message {
		t.Fatalf("unexpected message, want %s, got %s", want, message)
	}
}
