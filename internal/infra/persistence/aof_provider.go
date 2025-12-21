package persistence

import "github.com/g-villarinho/godis/internal/core/ports/repositories"

type AOFProviderOptions struct {
	EnableAOF bool
	Filepath  string
}

func NewAOFProvider(opt AOFProviderOptions) (repositories.PersistenceRepository, error) {
	if !opt.EnableAOF {
		return nil, nil
	}

	return nil, nil
}
