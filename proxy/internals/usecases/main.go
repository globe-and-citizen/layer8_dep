package usecases

import "globe-and-citizen/layer8/proxy/internals/repository"

type UseCase struct {
	Repo repository.Repository
}

func NewUseCase(repo repository.Repository) *UseCase {
	return &UseCase{
		Repo: repo,
	}
}
