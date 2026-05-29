package activities

import "context"

type ActivitiesService struct {
	repo ActivityRepository
}

func NewActivitiesService(repo ActivityRepository) *ActivitiesService {
	return &ActivitiesService{repo: repo}
}

func (s *ActivitiesService) CreateActivity(ctx context.Context, activity *Activity) error {
	if err := activity.Validate(); err != nil {
		return err
	}
	return s.repo.Create(ctx, activity)
}

func (s *ActivitiesService) FindActivityByID(ctx context.Context, id string) (*Activity, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *ActivitiesService) FindActivities(ctx context.Context) ([]*Activity, error) {
	return s.repo.FindAll(ctx)
}

func (s *ActivitiesService) FindActivitiesByHandle(ctx context.Context, handle string) ([]*Activity, error) {
	return s.repo.FindByHandle(ctx, handle)
}

func (s *ActivitiesService) UpdateActivity(ctx context.Context, activity *Activity) error {
	if err := activity.Validate(); err != nil {
		return err
	}
	return s.repo.Update(ctx, activity)
}

func (s *ActivitiesService) DeleteActivity(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
