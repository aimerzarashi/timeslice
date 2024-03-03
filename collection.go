package timeslice

import (
	"errors"
	"sort"
	"time"
)

type (
	Collection[T any] struct {
		items []*Item[T]
	}
)

var (
	ErrCollectionInvalid  = errors.New("Collection: invalid")
	ErrCollectionNotFound = errors.New("Collection: not found")
)

func NewCollection[T any](initItems ...*Item[T]) (*Collection[T], error) {
	// startAtで昇順ソート
	sort.Slice(initItems, func(i, j int) bool {
		return initItems[i].StartAt().Before(initItems[j].StartAt())
	})

	// 期間が重複していないか確認
	for i := 0; i < len(initItems)-1; i++ {
		if initItems[i].EndAt().Compare(initItems[i+1].StartAt()) >= 0 {
			return nil, ErrCollectionInvalid
		}
	}

	var items []*Item[T]
	if len(initItems) == 0 {
		items = make([]*Item[T], len(initItems))
	} else {
		items = initItems
	}

	return &Collection[T]{
		items: items,
	}, nil
}

func (c Collection[T]) Items() []*Item[T] {
	return c.items
}

func (d Collection[T]) Find(criteria time.Time) (Item[T], error) {
	for _, v := range d.items {
		if v.Contains(criteria) {
			return *v, nil
		}
	}

	return Item[T]{}, ErrCollectionNotFound
}

func (d *Collection[T]) Add(adding *Item[T]) (*Collection[T], error) {
	buffer := make([]*Item[T], 0)
	buffer = append(buffer, adding)

	// 追加する期間が重複している場合は、追加する期間を優先して既存の期間を調整する
	for _, v := range d.items {
		adjusted, err := v.Adjust(adding)
		if err != nil {
			return nil, err
		}

		buffer = append(buffer, adjusted...)
	}

	// startAtで昇順ソートする
	sort.Slice(buffer, func(i, j int) bool {
		return buffer[i].StartAt().Before(buffer[j].StartAt())
	})

	// 調整済みの期間に置き換える
	items, err := NewCollection(buffer...)
	if err != nil {
		return nil, err
	}
	return items, nil
}
