package notifier

import "context"

// Notifier sends a message
type Notifier interface {
	Notify(ctx context.Context, user, message string) error
}
