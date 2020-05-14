package pool

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"testing"
	"time"
)

type DemoCloser struct {
	Name     string
	activeAt time.Time
}

func (p *DemoCloser) Close() error {
	fmt.Println(p.Name, "closed")
	return nil
}

func (p *DemoCloser) GetActiveTime() time.Time {
	return p.activeAt
}

func TestNewGenericPool(t *testing.T) {
	_, err := NewGenericPool(0, 10, time.Minute*10, func() (Poolable, error) {
		time.Sleep(time.Second)
		return &DemoCloser{Name: "test"}, nil
	})
	if err != nil {
		t.Error(err)
	}
}

func TestGenericPool_Get(t *testing.T) {
	pool, err := NewGenericPool(0, 5, time.Minute*10, func() (Poolable, error) {
		time.Sleep(time.Second)
		name := strconv.FormatInt(time.Now().Unix(), 10)
		log.Printf("%s created", name)
		// TODO: FIXME &DemoCloser{Name: name}后，pool.Get陷入死循环
		return &DemoCloser{Name: name, activeAt: time.Now()}, nil
	})
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i < 10; i++ {
		s, err := pool.Get()
		if err != nil {
			t.Error(err)
			return
		}
		_ = pool.Put(s)
	}
}

func TestGenericPool_Len(t *testing.T) {
	pool, err := NewGenericPool(0, 5, time.Minute*10, func() (Poolable, error) {
		time.Sleep(time.Second)
		name := strconv.FormatInt(time.Now().Unix(), 10)
		log.Printf("%s created", name)
		// TODO: FIXME &DemoCloser{Name: name}后，pool.Get陷入死循环
		return &DemoCloser{Name: name, activeAt: time.Now()}, nil
	})
	if err != nil {
		t.Error(err)
		return
	}
	if i := pool.Len(); reflect.TypeOf(i).String() != "int" {
		t.Error(err)
	}
}

func TestGenericPool_Shutdown(t *testing.T) {
	pool, err := NewGenericPool(0, 10, time.Minute*10, func() (Poolable, error) {
		time.Sleep(time.Second)
		return &DemoCloser{Name: "test"}, nil
	})
	if err != nil {
		t.Error(err)
		return
	}
	if err := pool.Shutdown(); err != nil {
		t.Error(err)
		return
	}
	if _, err := pool.Get(); err != ErrPoolClosed {
		t.Error(err)
	}
}

func TestGenericPool_Close(t *testing.T) {
	pool, err := NewGenericPool(1, 10, time.Minute*10, func() (Poolable, error) {
		time.Sleep(time.Second)
		return &DemoCloser{Name: "test close", activeAt: time.Now()}, nil
	})
	if err != nil {
		t.Error(err)
		return
	}
	p, err := pool.Get()
	if err != nil {
		t.Error(err)
	}
	if err := pool.Close(p); err != nil {
		t.Error(err)
	}
}

