package closer

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

type MockFunc func(ctx context.Context) error

func (m MockFunc) Execute(ctx context.Context) error {
	return m(ctx)
}

func TestCloser_Add(t *testing.T) {
	c := Closer{}
	c.Add((func(ctx context.Context) error { return nil }))
	if len(c.funcs) != 1 {
		t.Error("Expected 1 function added, got", len(c.funcs))
	}
}

func TestCloser_Close(t *testing.T) {
	t.Run("Successful shutdown", func(t *testing.T) {
		c := Closer{}
		var executed bool
		c.Add(func(ctx context.Context) error {
			executed = true
			return nil
		})

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if err := c.Close(ctx); err != nil {
			t.Error("Expected nil error, got", err)
		}
		if !executed {
			t.Error("Expected the function to be executed")
		}
	})

	t.Run("Error during shutdown", func(t *testing.T) {
		c := Closer{}
		c.Add(func(ctx context.Context) error {
			return errors.New("error during shutdown")
		})

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		expectedErr := "shutdown finished with error(s): \n[!] error during shutdown"
		if err := c.Close(ctx); err == nil || err.Error() != expectedErr {
			t.Errorf("Expected error: '%s', got: '%v'", expectedErr, err)
		}
	})

	t.Run("Cancellation", func(t *testing.T) {
		c := Closer{}
		c.Add(func(ctx context.Context) error {
			<-ctx.Done()
			return nil
		})

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		expectedErr := "shutdown cancelled: context canceled"
		if err := c.Close(ctx); err == nil || err.Error() != expectedErr {
			t.Errorf("Expected error: '%s', got: '%v'", expectedErr, err)
		}
	})
}

func TestCloser_Close_Concurrent(t *testing.T) {
	t.Run("Concurrent shutdown", func(t *testing.T) {
		c := Closer{}
		var executed int
		numFuncs := 10

		var wg sync.WaitGroup
		for i := 0; i < numFuncs; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				c.Add(func(ctx context.Context) error {
					executed++
					return nil
				})
			}()
		}
		wg.Wait()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if err := c.Close(ctx); err != nil {
			t.Error("Expected nil error, got", err)
		}
		if executed != numFuncs {
			t.Errorf("Expected %d functions to be executed, got %d", numFuncs, executed)
		}
	})
}
