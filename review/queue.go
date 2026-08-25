package review

import (
	"sort"

	"subsidy11/domain"
	"subsidy11/storage"
)

type Queue struct{ store *storage.Store }

func NewQueue(store *storage.Store) *Queue { return &Queue{store: store} }

func (q *Queue) Pending() ([]domain.Record, error) {
	items, err := q.store.ListRecords()
	if err != nil {
		return nil, err
	}
	result := make([]domain.Record, 0, len(items))
	for _, item := range items {
		if item.Status == domain.StatusPendingReview || item.Status == domain.StatusValidated {
			result = append(result, item)
		}
	}
	sort.SliceStable(result, func(i, j int) bool { return result[i].CreatedAt.Before(result[j].CreatedAt) })
	return result, nil
}

func (q *Queue) FindForOfficer(officer string) ([]domain.Record, error) {
	items, err := q.Pending()
	if err != nil {
		return nil, err
	}
	result := make([]domain.Record, 0, len(items))
	for _, item := range items {
		if !item.HasOfficer(officer) {
			result = append(result, item)
		}
	}
	return result, nil
}
