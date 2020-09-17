package hw04_lru_cache //nolint:golint,stylecheck
import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool // Добавить значение в кэш по ключу
	Get(key Key) (interface{}, bool)     // Получить значение из кэша по ключу
	Clear()                              // Очистить кэш
}

type cacheItem struct {
	key   Key
	value interface{}
}

type lruCache struct {
	items    map[Key]*ListItem
	queue    List
	capacity int
	mutex    sync.Mutex
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	item, ok := l.items[key]
	if ok {
		cachedItem := getAsCachedItem(item)
		cachedItem.value = value
	} else {
		queueItem := l.queue.PushFront(&cacheItem{
			key:   key,
			value: value,
		})
		if l.queue.Len() > l.capacity {
			removedQueueItem := l.queue.Back()
			l.queue.Remove(removedQueueItem)

			delete(l.items, getAsCachedItem(removedQueueItem).key)
		}

		l.items[key] = queueItem
	}

	return ok
}

func getAsCachedItem(item *ListItem) *cacheItem {
	cachedItem, ok := item.Value.(*cacheItem)
	if !ok {
		panic("not expected item")
	}
	return cachedItem
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	item, ok := l.items[key]
	if ok {
		cachedItem := getAsCachedItem(item)
		l.queue.MoveToFront(item)
		return cachedItem.value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.items = make(map[Key]*ListItem)
	l.queue = NewList()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem),
	}
}
