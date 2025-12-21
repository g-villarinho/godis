package ports

import "context"

type Persistence interface {
	Append(ctx context.Context, command string, args []string) error
	Replay(ctx context.Context, store KeyValueRepository) error
	Close() error
}
