package activities

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"
)

type ActivityRepository interface {
	Create(ctx context.Context, activity *Activity) error
	FindByID(ctx context.Context, id string) (*Activity, error)
	FindAll(ctx context.Context) ([]*Activity, error)
	FindByHandle(ctx context.Context, handle string) ([]*Activity, error)
	Update(ctx context.Context, activity *Activity) error
	Delete(ctx context.Context, id string) error
}

// bd mocado para testes, usando um array em memória.
type MockActivitiesRepo struct {
	mu    sync.RWMutex
	data  map[string]*Activity
	order []string
	seq   int64
}

func NewMockActivitiesRepo(initialActivities ...*Activity) *MockActivitiesRepo {
	repo := &MockActivitiesRepo{
		data:  make(map[string]*Activity, len(initialActivities)),
		order: make([]string, 0, len(initialActivities)),
	}

	for _, activity := range initialActivities {
		if activity == nil {
			continue
		}
		cloned := cloneActivity(activity)
		if cloned.UUID == "" {
			repo.seq++
			cloned.UUID = fmt.Sprintf("mock-activity-%d", repo.seq)
		}
		repo.data[cloned.UUID] = cloned
		repo.order = append(repo.order, cloned.UUID)
	}

	return repo
}

func cloneActivity(activity *Activity) *Activity {
	if activity == nil {
		return nil
	}

	cloned := *activity
	if activity.ReplyToActivityUUID != nil {
		replyTo := *activity.ReplyToActivityUUID
		cloned.ReplyToActivityUUID = &replyTo
	}
	if len(activity.Replies) > 0 {
		cloned.Replies = make([]Activity, len(activity.Replies))
		copy(cloned.Replies, activity.Replies)
	}
	if activity.ExpiresAt != nil {
		expiresAt := *activity.ExpiresAt
		cloned.ExpiresAt = &expiresAt
	}

	return &cloned
}

func (r *MockActivitiesRepo) nextID() string {
	r.seq++
	return fmt.Sprintf("mock-activity-%d", r.seq)
}

func (r *MockActivitiesRepo) Create(ctx context.Context, activity *Activity) error {
	if activity == nil {
		return errors.New("activity required")
	}
	if err := activity.Validate(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	stored := cloneActivity(activity)
	if stored.UUID == "" {
		stored.UUID = r.nextID()
	}
	now := time.Now().UTC()
	if stored.CreatedAt.IsZero() {
		stored.CreatedAt = now
	}
	stored.UpdatedAt = now
	stored.DeletedAt = gorm.DeletedAt{}

	if r.data == nil {
		r.data = make(map[string]*Activity)
	}
	r.data[stored.UUID] = stored
	if stored.UUID != "" {
		r.order = append(r.order, stored.UUID)
	}
	activity.UUID = stored.UUID
	activity.CreatedAt = stored.CreatedAt
	activity.UpdatedAt = stored.UpdatedAt
	activity.DeletedAt = stored.DeletedAt
	return nil
}

func (r *MockActivitiesRepo) FindByID(ctx context.Context, id string) (*Activity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if activity, ok := r.data[id]; ok && activity != nil {
		return cloneActivity(activity), nil
	}

	return nil, errors.New("activity not found")
}

func (r *MockActivitiesRepo) FindAll(ctx context.Context) ([]*Activity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	activities := make([]*Activity, 0, len(r.data))
	for _, id := range r.order {
		if activity, ok := r.data[id]; ok && activity != nil {
			activities = append(activities, cloneActivity(activity))
		}
	}

	return activities, nil
}

func (r *MockActivitiesRepo) FindByHandle(ctx context.Context, handle string) ([]*Activity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	activities := make([]*Activity, 0)
	for _, id := range r.order {
		activity := r.data[id]
		if activity != nil && activity.Handle == handle {
			activities = append(activities, cloneActivity(activity))
		}
	}

	return activities, nil
}

func (r *MockActivitiesRepo) Update(ctx context.Context, activity *Activity) error {
	if activity == nil {
		return errors.New("activity required")
	}
	if activity.UUID == "" {
		return errors.New("activity id required")
	}
	if err := activity.Validate(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	current, ok := r.data[activity.UUID]
	if !ok || current == nil {
		return errors.New("activity not found")
	}

	updated := cloneActivity(activity)
	updated.CreatedAt = current.CreatedAt
	updated.UpdatedAt = time.Now().UTC()
	r.data[activity.UUID] = updated

	activity.CreatedAt = updated.CreatedAt
	activity.UpdatedAt = updated.UpdatedAt
	activity.DeletedAt = updated.DeletedAt
	return nil
}

func (r *MockActivitiesRepo) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[id]; ok {
		delete(r.data, id)
		for index, storedID := range r.order {
			if storedID == id {
				r.order = append(r.order[:index], r.order[index+1:]...)
				break
			}
		}
		return nil
	}

	return errors.New("activity not found")
}
