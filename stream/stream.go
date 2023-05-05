package stream

import (
	"sync"
)

// Stream 包实现了流式数据处理，使用 Stream 的好处是 worker 可以专注于写业务逻辑。
// 定时任务一般的流程都是：获取数据 -> 过滤数据 -> 处理数据 -> 输出结果
// 这个流程大约可以简化成： FROM -> FILTER -> MAP -> DONE
// 此外，FILTER、MAP 流程支持并发。
type Stream struct {
	src <-chan interface{}
}

func New(src <-chan interface{}) Stream {
	return Stream{src: src}
}

// 数据源迭代器，负责获取数据、数据分组
type srcIteratorFunc func(ch chan<- interface{})

func From(fn srcIteratorFunc) Stream {
	src := make(chan interface{})

	go func(src chan<- interface{}) {
		fn(src)
		close(src)
	}(src)

	return New(src)
}

type walkFunc func(val interface{}, ch chan<- interface{})

// Walk 是 stream 包的底层函数，提供基础的能力：
// 从 stream.src 读取数据 val，调用 walkFunc 处理数据后输出到 chan，支持并发。
func (s Stream) Walk(fn walkFunc, opts ...OptionFunc) Stream {
	ch := make(chan interface{})
	opt := buildOption(opts...)

	go func() {
		var wg sync.WaitGroup
		for i := 0; i < opt.workerNum; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for val := range s.src {
					fn(val, ch)
				}
			}()
		}

		wg.Wait()
		close(ch)
	}()

	return New(ch)
}

type MapFunc func(val interface{}) interface{}

// Map 是对 Walk 的基础封装，执行 MapFunc 函数对`数据val`做处理。
func (s Stream) Map(fn MapFunc, opts ...OptionFunc) Stream {
	return s.Walk(func(val interface{}, ch chan<- interface{}) {
		ch <- fn(val)
	}, opts...)
}

type FilterFunc func(val interface{}) bool

// Filter 会根据 FilterFunc 的返回值做过滤。
func (s Stream) Filter(fn FilterFunc, opts ...OptionFunc) Stream {
	return s.Walk(func(val interface{}, ch chan<- interface{}) {
		if fn(val) {
			ch <- val
		}
	}, opts...)
}

type done func(src <-chan interface{}) (result interface{}, err error)

// Done 直接读取 channel 的所有数据，
func (s Stream) Done(fn done) (result interface{}, err error) {
	return fn(s.src)
}
