package timeslice

import (
	"errors"
	"fmt"
	"time"
)

type (
	Item[T any] struct {
		value   *T
		startAt time.Time
		endAt   time.Time
	}
)

var (
	ErrItemStartAtEmpty      = errors.New("Item: startAt cannot be empty")
	ErrItemEndAtEmpty        = errors.New("Item: endAt cannot be empty")
	ErrItemInvalid           = errors.New("Item: invalid")
	ErrCollectionUnexpection = errors.New("Collection: unexpection")
)

func NewItem[T any](value *T, startAt, endAt time.Time) (*Item[T], error) {
	if startAt.IsZero() {
		return nil, ErrItemStartAtEmpty
	}
	if endAt.IsZero() {
		return nil, ErrItemEndAtEmpty
	}
	if startAt.Compare(endAt) > 0 {
		return nil, errors.Join(ErrItemInvalid, fmt.Errorf(" want startAt: %s <= endAt: %s", startAt.Format(time.RFC3339), endAt.Format(time.RFC3339)))
	}
	return &Item[T]{
		value:   value,
		startAt: startAt,
		endAt:   endAt,
	}, nil
}

func (i Item[T]) Value() T {
	return *i.value
}

func (i Item[T]) StartAt() time.Time {
	return i.startAt
}

func (i Item[T]) EndAt() time.Time {
	return i.endAt
}

func (i Item[T]) Contains(t time.Time) bool {
	return i.startAt.Compare(t) <= 0 && i.endAt.Compare(t) >= 0
}

// 既存期間と追加期間が重複している場合は、追加期間を優先して調整した既存期間を返す
func (i Item[T]) Adjust(adding *Item[T]) ([]*Item[T], error) {
	// 追加期間に対し、既存期間は重複しない前方に位置するため、そのまま返す
	if adding.startAt.Compare(i.endAt) > 0 {
		return []*Item[T]{&i}, nil
	}

	// 追加期間に対し、既存期間は重複しない後方に位置するため、そのまま返す
	if adding.endAt.Compare(i.startAt) < 0 {
		return []*Item[T]{&i}, nil
	}

	// 追加期間が既存期間を包含するため、nilで返す
	if adding.startAt.Compare(i.startAt) <= 0 && adding.endAt.Compare(i.endAt) >= 0 {
		return []*Item[T]{}, nil
	}

	// 追加期間が既存期間に包含されるため、追加期間を優先して前方と後方に分割して返す
	if adding.startAt.Compare(i.startAt) > 0 && adding.endAt.Compare(i.endAt) < 0 {

		// 分割した既存期間の前方は、開始日時を調整して返す
		foward, err := NewItem(i.value, i.startAt, adding.startAt.Add(-1*time.Second))
		if err != nil {
			return nil, errors.Join(ErrCollectionUnexpection, err)
		}

		// 分割した既存期間の後方は、終了日時を調整して返す
		backward, err := NewItem(i.value, adding.endAt.Add(1*time.Second), i.endAt)
		if err != nil {
			return nil, errors.Join(ErrCollectionUnexpection, err)
		}

		return []*Item[T]{foward, backward}, nil
	}

	// 追加期間に対し、既存期間の終了日時が重複するため、既存期間の終了日時を調整して返す
	if adding.startAt.Compare(i.endAt) < 0 && adding.endAt.Compare(i.endAt) > 0 {
		foward, err := NewItem[T](i.value, i.startAt, adding.startAt.Add(-1*time.Second))
		if err != nil {
			return nil, errors.Join(ErrCollectionUnexpection, err)
		}
		return []*Item[T]{foward}, nil
	}

	// 追加期間に対し、既存期間の開始日時が重複するため、既存期間の開始日時を調整して返す
	if adding.startAt.Compare(i.startAt) < 0 && adding.endAt.Compare(i.endAt) < 0 {
		backward, err := NewItem[T](i.value, adding.endAt.Add(1*time.Second), i.endAt)
		if err != nil {
			return nil, errors.Join(ErrCollectionUnexpection, err)
		}
		return []*Item[T]{backward}, nil
	}

	return nil, ErrCollectionUnexpection
}
