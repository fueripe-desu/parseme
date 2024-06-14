package parseme

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestObserver struct {
	test_notified bool
}

func (o *TestObserver) OnUpdate(info ErrorInfo) {
	o.test_notified = true
}

type TestObserver2 struct{}

func (o *TestObserver2) OnUpdate(info ErrorInfo) {
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
		assert := assert.New(t)
		pool := &ErrorPool{}
		observer := &TestObserver{}
		data := errorData{name: "random-error", level: Warning, message: "Some message"}

		pool.Subscribe(observer)
		pool.Error(data, nil)

		assert.Equal(pool.errorStack.Peek(), data)
		assert.Equal(observer.test_notified, true)
	})
}

func Test_Notify(t *testing.T) {
	t.Run("notify last error", func(t *testing.T) {
		assert := assert.New(t)
		pool := &ErrorPool{}
		observer := &TestObserver{}
		data1 := errorData{name: "random-error", level: Warning, message: "Some message"}
		data2 := errorData{name: "random-error2", level: Warning, message: "Some message2"}

		pool.Subscribe(observer)
		pool.AddError(data1, nil)
		pool.AddError(data2, nil)
		pool.Notify()

		assert.Equal(observer.test_notified, true)
	})

	t.Run("notify empty pool", func(t *testing.T) {
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
	testcases := []struct {
		name         string
		inputData    errorData
		args         []string
		expectedData errorData
	}{
		{
			"add error",
			errorData{name: "random-error", level: Warning, message: "Some message"},
			nil,
			errorData{name: "random-error", level: Warning, message: "Some message"},
		},
		{
			"add masked error",
			errorData{name: "random-error", level: Warning, message: "Name: %{0}, Age: %{1}"},
			[]string{"John", "40"},
			errorData{name: "random-error", level: Warning, message: "Name: John, Age: 40"},
		},

		{
			"masked but nil args",
			errorData{name: "masked error", level: Warning, message: "Name: %{0}"},
			nil,
			errorData{name: "masked error", level: Warning, message: "Name: %{0}"},
		},
		{
			"masked but empty args",
			errorData{name: "masked error", level: Warning, message: "Name: %{0}"},
			[]string{},
			errorData{name: "masked error", level: Warning, message: "Name: %{0}"},
		},
		{
			"args but no mask",
			errorData{name: "masked error", level: Warning, message: "Name: John"},
			[]string{"Carl", "Paul"},
			errorData{name: "masked error", level: Warning, message: "Name: John"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			pool := &ErrorPool{}
			pool.AddError(tc.inputData, tc.args)

			assert.Equal(pool.errorStack.Peek(), tc.expectedData)
		})

	}
}

func Test_CleanErrors(t *testing.T) {
	t.Run("clean errors", func(t *testing.T) {
		assert := assert.New(t)
		pool := &ErrorPool{}
		data1 := errorData{name: "random-error", level: Warning, message: "Some message"}
		data2 := errorData{name: "random-error2", level: Warning, message: "Some message2"}

		pool.AddError(data1, nil)
		pool.AddError(data2, nil)
		pool.ClearErrors()

		assert.Equal(pool.errorStack.IsEmpty(), true)
	})
}

func Test_HasErrors(t *testing.T) {
	t.Run("not empty pool", func(t *testing.T) {
		assert := assert.New(t)
		pool := &ErrorPool{}
		data := errorData{name: "random-error", level: Warning, message: "Some message"}

		pool.AddError(data, nil)

		assert.Equal(pool.HasErrors(), true)
	})

	t.Run("empty bool", func(t *testing.T) {
		assert := assert.New(t)
		pool := &ErrorPool{}

		assert.Equal(pool.HasErrors(), false)
	})
}

func Test_consumeMask(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected int
		endIndex int
	}{
		{
			"normal mask",
			"22}",
			22,
			2,
		},
		{
			"empty mask",
			"}",
			-1,
			0,
		},
		{
			"alphabetic mask",
			"something}",
			-1,
			9,
		},
		{
			"alphanumeric mask",
			"abc1def2}",
			12,
			8,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			pool := &ErrorPool{}
			bytes := []byte(tc.input)
			index, argIndex := pool.consumeMask(0, &bytes)
			assert.Equal(argIndex, tc.expected)
			assert.Equal(index, tc.endIndex)
		})
	}
}

func Test_replaceMasks(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		args     []string
		expected string
	}{
		{
			"normal mask",
			"Hello! My name is %{0}.",
			[]string{"Paul"},
			"Hello! My name is Paul.",
		},
		{
			"empty mask",
			"This %{}should not be replaced.",
			[]string{"something"},
			"This should not be replaced.",
		},
		{
			"alphabetic mask",
			"This %{something}should not be replaced.",
			[]string{"Carl"},
			"This should not be replaced.",
		},
		{
			"alphanumeric mask",
			"I'm %{abc2} years old!",
			[]string{"18", "21", "34"},
			"I'm 34 years old!",
		},
		{
			"nil args",
			"I'm %{0} years old!",
			nil,
			"I'm  years old!",
		},
		{
			"empty args",
			"I'm %{0} years old!",
			[]string{},
			"I'm  years old!",
		},
		{
			"index out of range",
			"I'm %{33} years old!",
			[]string{"22"},
			"I'm  years old!",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			pool := &ErrorPool{}
			newString := pool.replaceMasks(tc.input, tc.args)
			assert.Equal(newString, tc.expected)
		})
	}
}

func Test_precompileError(t *testing.T) {
	t.Run("precompile error", func(t *testing.T) {
		assert := assert.New(t)
		pool := &ErrorPool{}
		data := errorData{name: "some error", level: Warning, message: "Name: %{0}, Age: %{1}"}
		newData := pool.precompileError(data, []string{"Carl", "40"})
		assert.Equal(newData.message, "Name: Carl, Age: 40")
	})
}
