package parseme

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var test_notified bool

type TestObserver struct{}

func (o *TestObserver) OnUpdate(name string, level ErrorLevel, message string) {
	test_notified = true
}

type TestObserver2 struct{}

func (o *TestObserver2) OnUpdate(name string, level ErrorLevel, message string) {
	return
}

func Test_Subscribe(t *testing.T) {
	t.Run("add observer", func(t *testing.T) {
		assert := assert.New(t)
		pool := &ErrorPool{}
		observer := &TestObserver{}

		err := pool.Subscribe(observer)

		if err != nil {
			t.Log("Subscribe() returned error: " + err.Error())
			t.FailNow()
		}
		assert.Equal(pool.observers, []ErrorObserver{observer})
	})

	t.Run("add more than one observer", func(t *testing.T) {
		assert := assert.New(t)
		pool := &ErrorPool{}
		observer1 := &TestObserver{}
		observer2 := &TestObserver2{}

		err1 := pool.Subscribe(observer1)

		if err1 != nil {
			t.Log("Subscribe() returned error: " + err1.Error())
			t.FailNow()
		}

		err2 := pool.Subscribe(observer2)

		if err2 != nil {
			t.Log("Subscribe() returned error: " + err2.Error())
			t.FailNow()
		}

		assert.Equal(pool.observers, []ErrorObserver{observer1, observer2})
	})

	t.Run("add duplicate observer", func(t *testing.T) {
		assert := assert.New(t)
		pool := &ErrorPool{}
		observer := &TestObserver{}

		firstErr := pool.Subscribe(observer)

		if firstErr != nil {
			t.Log("Subscribe() returned unexepected error: " + firstErr.Error())
			t.FailNow()
		}

		secondErr := pool.Subscribe(observer)
		assert.Equal(secondErr, &ObserverDuplicateError{})
		assert.EqualError(secondErr, (&ObserverDuplicateError{}).Error())
	})
}

func Test_Unsubscribe(t *testing.T) {
	t.Run("unsubscribe existent observer", func(t *testing.T) {
		assert := assert.New(t)
		pool := &ErrorPool{}
		observer := &TestObserver{}

		pool.Subscribe(observer)
		err := pool.Unsubscribe(observer)

		if err != nil {
			t.Log("Unsubscribe() returned an unexpected error: " + err.Error())
			t.FailNow()
		}

		assert.Equal(pool.observers, []ErrorObserver{})
	})

	t.Run("unsubscribe non existent observer", func(t *testing.T) {
		assert := assert.New(t)
		pool := &ErrorPool{}
		observer := &TestObserver{}

		err := pool.Unsubscribe(observer)

		assert.Equal(err, &ObserverNotFoundError{})
		assert.EqualError(err, (&ObserverNotFoundError{}).Error())
	})

}

func Test_UnsubscribeAll(t *testing.T) {
	t.Run("unsubscribe all", func(t *testing.T) {
		assert := assert.New(t)
		pool := &ErrorPool{}
		observer1 := &TestObserver{}
		observer2 := &TestObserver2{}

		pool.Subscribe(observer1)
		pool.Subscribe(observer2)

		pool.UnsubscribeAll()

		assert.Equal(pool.observers, []ErrorObserver{})
	})
}

func Test_Error(t *testing.T) {
	t.Run("add error", func(t *testing.T) {
		t.Cleanup(func() {
			test_notified = false
		})
		assert := assert.New(t)
		pool := &ErrorPool{}
		observer := &TestObserver{}
		data := ErrorData{name: "random-error", level: Warning, message: "Some message"}

		pool.Subscribe(observer)
		pool.Error(data)

		assert.Equal(pool.errorStack.Peek(), data)
		assert.Equal(test_notified, true)
	})
}

func Test_Notify(t *testing.T) {
	t.Run("notify last error", func(t *testing.T) {
		t.Cleanup(func() {
			test_notified = false
		})
		assert := assert.New(t)
		pool := &ErrorPool{}
		observer := &TestObserver{}
		data1 := ErrorData{name: "random-error", level: Warning, message: "Some message"}
		data2 := ErrorData{name: "random-error2", level: Warning, message: "Some message2"}

		pool.Subscribe(observer)
		pool.AddError(data1)
		pool.AddError(data2)
		pool.Notify()

		assert.Equal(test_notified, true)
	})

	t.Run("notify empty pool", func(t *testing.T) {
		t.Cleanup(func() {
			test_notified = false
		})
		assert := assert.New(t)
		pool := &ErrorPool{}
		observer := &TestObserver{}

		pool.Subscribe(observer)
		err := pool.Notify()

		assert.Equal(err, &NotifyObserverError{})
		assert.EqualError(err, (&NotifyObserverError{}).Error())
	})
}

func Test_AddError(t *testing.T) {
	t.Run("add error", func(t *testing.T) {
		assert := assert.New(t)
		pool := &ErrorPool{}
		data := ErrorData{name: "random-error", level: Warning, message: "Some message"}

		pool.AddError(data)

		assert.Equal(pool.errorStack.Peek(), data)
		assert.Equal(test_notified, false)
	})
}

func Test_CleanErrors(t *testing.T) {
	t.Run("clean errors", func(t *testing.T) {
		assert := assert.New(t)
		pool := &ErrorPool{}
		data1 := ErrorData{name: "random-error", level: Warning, message: "Some message"}
		data2 := ErrorData{name: "random-error2", level: Warning, message: "Some message2"}

		pool.AddError(data1)
		pool.AddError(data2)
		pool.ClearErrors()

		assert.Equal(pool.errorStack.IsEmpty(), true)
	})
}

func Test_HasErrors(t *testing.T) {
	t.Run("not empty pool", func(t *testing.T) {
		assert := assert.New(t)
		pool := &ErrorPool{}
		data := ErrorData{name: "random-error", level: Warning, message: "Some message"}

		pool.AddError(data)

		assert.Equal(pool.HasErrors(), true)
	})

	t.Run("empty bool", func(t *testing.T) {
		assert := assert.New(t)
		pool := &ErrorPool{}

		assert.Equal(pool.HasErrors(), false)
	})
}
