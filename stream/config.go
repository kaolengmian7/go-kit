package stream

// 配置相关
type option struct {
	workerNum int // 协程数
}

type OptionFunc func(op *option)

func WithWorkerNum(num int) OptionFunc {
	return func(op *option) {
		op.workerNum = num
	}
}

func buildOption(opts ...OptionFunc) (opt *option) {
	opt = buildDefaultOption()
	for _, fn := range opts {
		fn(opt)
	}
	return
}

// 默认配置
func buildDefaultOption() (opt *option) {
	opt = &option{
		workerNum: 1, // 1 个协程
	}

	return
}
