package painter

import (
	"image"
	"sync"

	"golang.org/x/exp/shiny/screen"
)

// Receiver отримує текстуру, яка була підготовлена в результаті виконання команд у циклі подій.
type Receiver interface {
	Update(t screen.Texture)
}

// Loop реалізує цикл подій для формування текстури отриманої через виконання операцій отриманих з внутрішньої черги.
type Loop struct {
	Receiver Receiver

	next screen.Texture // текстура, яка зараз формується
	prev screen.Texture // текстура, яка була відправлення останнього разу у Receiver

	mq messageQueue

	stop    chan struct{}
	stopReq bool
}

var size = image.Pt(400, 400)

// Start запускає цикл подій. Цей метод потрібно запустити до того, як викликати на ньому будь-які інші методи.
func (l *Loop) Start(s screen.Screen) {
	l.next, _ = s.NewTexture(size)
	l.prev, _ = s.NewTexture(size)

	l.stop = make(chan struct{})

	go func() {
		for !(l.stopReq && l.mq.empty()) {
			op := l.mq.pull()
			update := op.Do(l.next)
			if update {
				l.Receiver.Update(l.next)
				l.next, l.prev = l.prev, l.next
			}
		}
		close(l.stop)
	}()
}

// Post додає нову операцію у внутрішню чергу.
func (l *Loop) Post(op Operation) {
	if update := op.Do(l.next); update {
		l.Receiver.Update(l.next)
		l.next, l.prev = l.prev, l.next
	}
}

// StopAndWait сигналізує про необхідність завершити цикл та блокується до моменту його повної зупинки.
func (l *Loop) StopAndWait() {
	l.mq.push(OperationFunc(func(t screen.Texture) {
		l.stopReq = true
	}))
}

type messageQueue struct {
	ops []Operation
	mu  sync.Mutex

	blocked chan struct{}
}

func (mq *messageQueue) push(op Operation) {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	mq.ops = append(mq.ops, op)

	if mq.blocked != nil {
		close(mq.blocked)
		mq.blocked = nil
	}
}

func (mq *messageQueue) pull() Operation {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	for len(mq.ops) == 0 {
		if mq.blocked == nil {
			mq.blocked = make(chan struct{})
		}
		mq.mu.Unlock()
		<-mq.blocked
		mq.mu.Lock()
	}

	op := mq.ops[0]
	mq.ops[0] = nil
	mq.ops = mq.ops[1:]
	return op
}

func (mq *messageQueue) empty() bool {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	return len(mq.ops) == 0
}
