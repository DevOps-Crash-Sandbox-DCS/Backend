package scenarios

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]ScenarioResponse, error) {
	items, err := s.repo.GetAllActive(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]ScenarioResponse, 0, len(items))

	for _, item := range items {
		response = append(response, ToScenarioResponse(item))
	}

	return response, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*ScenarioResponse, error) {
	scenario, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := ToScenarioResponse(*scenario)

	return &response, nil
}
