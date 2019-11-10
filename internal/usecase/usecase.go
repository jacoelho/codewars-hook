package usecase

import (
	"context"
	"fmt"

	"github.com/jacoelho/codewars/internal/notifier"
	"github.com/jacoelho/codewars/internal/user"
)

type UserHonorUpdated func(context.Context, user.Event) error

func UserHonorUpdatedCase(notify notifier.Notifier, repo user.Repository) UserHonorUpdated {
	return func(ctx context.Context, event user.Event) error {
		if event.HonorUpgraded == nil {
			return fmt.Errorf("received empty event")
		}

		u, err := repo.GetUserByID(ctx, event.HonorUpgraded.ID)
		if err != nil {
			return err
		}

		return notify.Notify(
			ctx,
			u.Username,
			fmt.Sprintf(":arrow_up:  Honor went up to %d", event.HonorUpgraded.Honor),
		)
	}
}
