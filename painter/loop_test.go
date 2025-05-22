package painter

import (
	"image"
	"image/color"
	"image/draw"
	"reflect"
	"sync"
	"testing"
	"time"

	"golang.org/x/exp/shiny/screen"
)

func TestLoop_Post(t *testing.T) {
	var (
		l  Loop
		tr testReceiver
	)
	l.Receiver = &tr

	var testOps []string

	l.Start(mockScreen{})
	l.Post(logOp(t, "do white fill", WhiteFill))
	l.Post(logOp(t, "do green fill", GreenFill))
	l.Post(UpdateOp)

	for i := 0; i < 3; i++ {
		go l.Post(logOp(t, "do green fill", GreenFill))
	}

	l.Post(OperationFunc(func(screen.Texture) {
		testOps = append(testOps, "op 1")
		l.Post(OperationFunc(func(screen.Texture) {
			testOps = append(testOps, "op 2")
		}))
	}))
	l.Post(OperationFunc(func(screen.Texture) {
		testOps = append(testOps, "op 3")
	}))

	l.StopAndWait()

	if tr.lastTexture == nil {
		t.Fatal("Texture was not updated")
	}
	mt, ok := tr.lastTexture.(*mockTexture)
	if !ok {
		t.Fatal("Unexpected texture", tr.lastTexture)
	}
	if mt.Colors[0] != color.White {
		t.Error("First color is not white:", mt.Colors)
	}
	if len(mt.Colors) != 2 {
		t.Error("Unexpected size of colors:", mt.Colors)
	}

	if !reflect.DeepEqual(testOps, []string{"op 1", "op 2", "op 3"}) {
		t.Error("Bad order:", testOps)
	}
}

func logOp(t *testing.T, msg string, op OperationFunc) OperationFunc {
	return func(tx screen.Texture) {
		t.Log(msg)
		op(tx)
	}
}

type testReceiver struct {
	lastTexture screen.Texture
}

func (tr *testReceiver) Update(t screen.Texture) {
	tr.lastTexture = t
}

type mockScreen struct{}

func (m mockScreen) NewBuffer(size image.Point) (screen.Buffer, error) {
	panic("implement me")
}

func (m mockScreen) NewTexture(size image.Point) (screen.Texture, error) {
	return new(mockTexture), nil
}

func (m mockScreen) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	panic("implement me")
}

type mockTexture struct {
	Colors []color.Color
}

func (m *mockTexture) Release() {}

func (m *mockTexture) Size() image.Point { return size }

func (m *mockTexture) Bounds() image.Rectangle {
	return image.Rectangle{Max: m.Size()}
}

func (m *mockTexture) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {}
func (m *mockTexture) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	m.Colors = append(m.Colors, src)
}

type mockOp struct {
	id int
}

func (m mockOp) Do(t screen.Texture) bool {
	return true
}

func TestMessageQueuePush(t *testing.T) {
	var mq messageQueue

	op := &mockOp{}
	mq.push(op)

	if len(mq.ops) != 1 {
		t.Errorf("expected 1 operation, got %d", len(mq.ops))
	}
}

func TestMessageQueuePull(t *testing.T) {
	var mq messageQueue

	op1 := &mockOp{}
	op2 := &mockOp{}
	mq.push(op1)
	mq.push(op2)

	op := mq.pull()
	if op != op1 {
		t.Errorf("expected op1, got %+v", op)
	}
	if len(mq.ops) != 1 {
		t.Errorf("expected 1 operation left, got %d", len(mq.ops))
	}
}

func TestMessageQueuePullBlocksUntilPush(t *testing.T) {
	var mq messageQueue
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		op := mq.pull()
		if mop, ok := op.(*mockOp); !ok || mop.id != 99 {
			t.Errorf("expected operation id=99, got %+v", op)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	mq.push(&mockOp{id: 99})
	wg.Wait()
}

func TestMessageQueueEmpty(t *testing.T) {
	var mq messageQueue

	if !mq.empty() {
		t.Errorf("expected queue to be empty")
	}

	mq.push(&mockOp{})

	if mq.empty() {
		t.Errorf("expected queue to not be empty")
	}
}
