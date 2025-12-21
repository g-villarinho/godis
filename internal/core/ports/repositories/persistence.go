package repositories

import "context"

type PersistenceRepository interface {
	Append(ctx context.Context, command string, args []string) error
	Replay(ctx context.Context, store KeyValueRepository) error
	Close() error
}
