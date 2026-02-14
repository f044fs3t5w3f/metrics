package pool

type resetable interface {
	Reset()
}

type Pool[T resetable] struct {
	objects []T
	newFn   func() T
}

func New[T resetable](newFn func() T) *Pool[T] {
	return &Pool[T]{
		objects: make([]T, 0),
		newFn:   newFn,
	}
}

func (p *Pool[T]) Get() T {
	len := len(p.objects)
	if len == 0 {
		return p.newFn()
	}

	obj := p.objects[len-1]
	p.objects = p.objects[:len-1]
	return obj
}

func (p *Pool[T]) Put(obj T) {
	obj.Reset()
	p.objects = append(p.objects, obj)
}
