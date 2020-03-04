package util

import (
	"fmt"
	"reflect"
	"sync"
)

type Watcher interface {
	Has(topic string) bool

	Pub(topic string, args ...interface{})

	PubAs(topicValue ...interface{})

	PubTo(topicType reflect.Type, args ...interface{})

	Sub(topic string, function interface{}) Handler

	SubAs(function interface{}) Handler

	SubTo(topicType reflect.Type, function interface{}) Handler
}

type Handler interface {
	UnSub()
}

func NewWatcher() Watcher {
	return &watcher{
		locker: sync.Mutex{},
		topics: make(map[string][]*handler),
	}
}

type watcher struct {
	locker sync.Mutex
	topics map[string][]*handler
}

func (w *watcher) Has(topic string) bool {
	handlers, contains := w.topics[topic]

	return contains && len(handlers) > 0
}

func (w *watcher) Pub(topic string, args ...interface{}) {
	if handlers, contains := w.topics[topic]; contains && len(handlers) > 0 {

		callArgs := make([]reflect.Value, 0)
		for _, arg := range args {
			callArgs = append(callArgs, reflect.ValueOf(arg))
		}

		for _, handler := range handlers {
			if handler.function.Type().NumIn() != len(callArgs) {
				continue
			}

			handler.function.Call(callArgs)
		}
	}
}

func (w *watcher) PubAs(topicValue ...interface{}) {
	w.PubTo(reflect.TypeOf(topicValue[0]), topicValue...)
}

func (w *watcher) PubTo(topicType reflect.Type, args ...interface{}) {
	w.Pub(topicType.String(), args...)
}

func (w *watcher) Sub(topic string, function interface{}) Handler {
	w.locker.Lock()
	defer w.locker.Unlock()

	// check if function is actual a function
	if reflect.TypeOf(function).Kind() != reflect.Func {
		panic(fmt.Errorf("cannot sub with %v, must be reflect.Func", reflect.TypeOf(function).Kind()))
	}

	handler := &handler{
		topic:    topic,
		watch:    w,
		function: reflect.ValueOf(function),
	}

	// append watch handler to topic list
	w.topics[topic] = append(w.topics[topic], handler)

	return handler
}

func (w *watcher) SubAs(function interface{}) Handler {
	funcType := reflect.TypeOf(function)

	// check if function is actual a function
	if funcType.Kind() != reflect.Func {
		panic(fmt.Errorf("cannot sub with %v, must be reflect.Func", funcType.Kind()))
	}

	if funcType.NumIn() == 0 {
		panic(fmt.Errorf("cannot sub with %v, must have at least 1 input parameter", funcType))
	}

	return w.SubTo(funcType.In(0), function)
}

func (w *watcher) SubTo(topicType reflect.Type, function interface{}) Handler {
	return w.Sub(topicType.String(), function)
}

type handler struct {
	topic string
	watch *watcher

	function reflect.Value
}

func (h *handler) UnSub() {
	handlers := h.watch.topics[h.topic]
	if handlers == nil {
		return
	}

	for i, elem := range handlers {
		if elem == h {
			h.watch.topics[h.topic] = append(h.watch.topics[h.topic][:i], h.watch.topics[h.topic][i+1:]...)
		}
	}
}
