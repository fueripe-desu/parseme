package errors

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

type TestFuncObserver struct {
	observerFunc func(info ErrorInfo)
}

func (o *TestFuncObserver) OnUpdate(info ErrorInfo) {
	o.observerFunc(info)
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

func Test_error(t *testing.T) {
	t.Run("add error", func(t *testing.T) {
		assert := assert.New(t)
		pool := &ErrorPool{}
		observer := &TestObserver{}
		data := &errorData{name: "Fandom error", message: "Some message", code: "ABC", fix: "Fix some bugs"}

		pool.Subscribe(observer)
		pool.error(Warning, "Testing", data, nil)

		assert.Equal(pool.errorStack.peek(), *data)
		assert.Equal(observer.test_notified, true)
	})

	t.Run("nil error data", func(t *testing.T) {
		assert := assert.New(t)
		pool := &ErrorPool{}
		var data *errorData = nil

		assert.PanicsWithError("Error data must not be nil.", func() {
			pool.error(Error, "Testing", data, nil)
		})
	})

	t.Run("add non-recursive error", func(t *testing.T) {
		assert := assert.New(t)
		pool := &ErrorPool{}
		var resultInfo *ErrorInfo
		observer := &TestFuncObserver{
			observerFunc: func(info ErrorInfo) {
				resultInfo = &info
			},
		}
		var data *errorData = NewErrorData(
			"Test error",
			"This is a test message.",
			"EC1",
			"Try fixing those stupid bugs.",
		)

		pool.Subscribe(observer)
		pool.error(Error, "Test", data, nil)

		expectedInfo := ErrorInfo{
			Name:           "Test error",
			Message:        "This is a test message.",
			Code:           "EC1",
			Module:         "Test",
			Fix:            "Try fixing those stupid bugs.",
			CallerFuncName: "(*ErrorPool).error",
			CallerFilename: "error_pool.go",
			CallerLine:     103,
			ErrorFuncName:  "Test_error.func3",
			ErrorFilename:  "error_pool_test.go",
			ErrorLine:      174,
			ErrorSite:      "pool.error(Error, \"Test\", data, nil)",
			Level:          Error,
			Timestamp:      (*resultInfo).Timestamp,
		}

		assert.Equal(*resultInfo, expectedInfo)
	})
}

func Test_notify(t *testing.T) {
	t.Run("notify last error", func(t *testing.T) {
		assert := assert.New(t)
		pool := &ErrorPool{}
		observer := &TestObserver{}
		data1 := &errorData{name: "Random error", message: "Some message", code: "ABC", fix: "Some fix"}
		data2 := &errorData{name: "Random error two", message: "Some message2", code: "ABC", fix: "Some fix"}

		pool.Subscribe(observer)
		pool.addError(data1, nil)
		pool.addError(data2, nil)
		pool.notify(Warning, "Testing")

		assert.Equal(observer.test_notified, true)
	})

	t.Run("notify empty pool", func(t *testing.T) {
		assert := assert.New(t)
		pool := &ErrorPool{}
		observer := &TestObserver{}

		pool.Subscribe(observer)

		assert.PanicsWithError("Cannot notify when there are no errors.", func() {
			pool.notify(Error, "Testing")
		})
	})
}

func Test_addError(t *testing.T) {
	testcases := []struct {
		name         string
		inputData    *errorData
		args         []string
		expectedData errorData
	}{
		{
			"add error",
			&errorData{name: "Random error", message: "Some message"},
			nil,
			errorData{name: "Random error", message: "Some message"},
		},
		{
			"add masked error",
			&errorData{name: "Random error", message: "Name: %{0}, Age: %{1}"},
			[]string{"John", "40"},
			errorData{name: "Random error", message: "Name: John, Age: 40"},
		},

		{
			"masked but nil args",
			&errorData{name: "Masked error", message: "Name: %{0}"},
			nil,
			errorData{name: "Masked error", message: "Name: %{0}"},
		},
		{
			"masked but empty args",
			&errorData{name: "Masked error", message: "Name: %{0}"},
			[]string{},
			errorData{name: "Masked error", message: "Name: %{0}"},
		},
		{
			"args but no mask",
			&errorData{name: "Masked error", message: "Name: John"},
			[]string{"Carl", "Paul"},
			errorData{name: "Masked error", message: "Name: John"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			pool := &ErrorPool{}
			pool.addError(tc.inputData, tc.args)

			assert.Equal(pool.errorStack.peek(), tc.expectedData)
		})

	}
}

func Test_CleanErrors(t *testing.T) {
	t.Run("clean errors", func(t *testing.T) {
		assert := assert.New(t)
		pool := &ErrorPool{}
		data1 := &errorData{name: "random-error", message: "Some message"}
		data2 := &errorData{name: "random-error2", message: "Some message2"}

		pool.addError(data1, nil)
		pool.addError(data2, nil)
		pool.ClearErrors()

		assert.Equal(pool.errorStack.isEmpty(), true)
	})
}

func Test_HasErrors(t *testing.T) {
	t.Run("not empty pool", func(t *testing.T) {
		assert := assert.New(t)
		pool := &ErrorPool{}
		data := &errorData{name: "random-error", message: "Some message"}

		pool.addError(data, nil)

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
		data := errorData{name: "some error", message: "Name: %{0}, Age: %{1}"}
		newData := pool.precompileError(data, []string{"Carl", "40"})
		assert.Equal(newData.message, "Name: Carl, Age: 40")
	})
}

func Test_getStackIndex(t *testing.T) {
	testcases := []struct {
		name     string
		caller   bool
		expected int
	}{
		{
			"with caller",
			true,
			2,
		},
		{
			"without caller",
			false,
			3,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			pool := &ErrorPool{}
			result := pool.getStackIndex(tc.caller)
			assert.Equal(result, tc.expected)
		})
	}
}

func Test_extractFuncName(t *testing.T) {
	t.Run("full function name", func(t *testing.T) {
		assert := assert.New(t)
		pool := &ErrorPool{}
		result := pool.extractFuncName("github.com/fueripe-desu/parseme.(*CustomError).CustomFunc")
		assert.Equal(result, "(*CustomError).CustomFunc")
	})
}

func Test_getLineContents(t *testing.T) {
	t.Run("line contents", func(t *testing.T) {
		assert := assert.New(t)
		pool := &ErrorPool{}
		result := pool.getLineContents("../test_data/example1.html", 4)
		assert.Equal(result, "<title>Example 1</title>")
	})
}

func Test_GetErrorCount(t *testing.T) {
	t.Run("error count", func(t *testing.T) {
		assert := assert.New(t)
		pool := &ErrorPool{}
		data := &errorData{
			name:    "Not found error",
			message: "Could not find object.",
			code:    "ABC",
			fix:     "Try searching for the object.",
		}
		pool.addError(data, nil)
		pool.addError(data, nil)
		pool.addError(data, nil)
		result := pool.GetErrorCount()
		assert.Equal(result, 3)
	})
}

func Test_checkName(t *testing.T) {
	testcases := []struct {
		name         string
		input        string
		shouldPanic  bool
		panicMessage string
	}{
		{
			"valid name",
			"Not found error",
			false,
			"",
		},
		{
			"empty name",
			"",
			true,
			"Error name must not be empty.",
		},
		{
			"invalid name chars",
			"This 1s an inv@lid name",
			true,
			"Error name must only contain letters and spaces.",
		},
		{
			"name with leading spaces",
			"      This is an invalid name",
			true,
			"Error name must not contain leading or trailing spaces.",
		},
		{
			"name with trailing spaces",
			"This is an invalid name     ",
			true,
			"Error name must not contain leading or trailing spaces.",
		},
		{
			"name with consecutive spaces",
			"This is      an invalid name",
			true,
			"Error name must not contain consecutive spaces.",
		},
		{
			"no sentence case",
			"this is an invalid name",
			true,
			"Error name must follow the sentence case convention (first letter uppercase and all others lowercase).",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			pool := &ErrorPool{}

			if tc.shouldPanic {
				assert.PanicsWithError(tc.panicMessage, func() {
					pool.checkName(tc.input)
				})
			} else {
				pool.checkName(tc.input)
			}
		})
	}
}

func Test_checkMessage(t *testing.T) {
	testcases := []struct {
		name         string
		input        string
		shouldPanic  bool
		panicMessage string
	}{
		{
			"valid message",
			"This is a valid error message.",
			false,
			"",
		},
		{
			"empty message",
			"",
			true,
			"Error message must not be empty.",
		},
		{
			"without letters",
			"@!#@#",
			true,
			"Error message must contain at least one letter.",
		},
		{
			"message with control chars",
			"This is an invalid message.\n",
			true,
			"Error message must not contain control characters.",
		},
		{
			"first char not capitalized",
			"this is an invalid message.",
			true,
			"Error message must start with uppercase letter.",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			pool := &ErrorPool{}

			if tc.shouldPanic {
				assert.PanicsWithError(tc.panicMessage, func() {
					pool.checkMessage(tc.input)
				})
			} else {
				pool.checkMessage(tc.input)
			}
		})
	}
}

func Test_checkCode(t *testing.T) {
	testcases := []struct {
		name         string
		input        string
		shouldPanic  bool
		panicMessage string
	}{
		{
			"valid code",
			"C3A",
			false,
			"",
		},
		{
			"empty code",
			"",
			true,
			"Error code must not be empty.",
		},
		{
			"code with more than 3 chars",
			"ABCD",
			true,
			"Error code must have 3 characters of length.",
		},
		{
			"code with less than 3 chars",
			"AB",
			true,
			"Error code must have 3 characters of length.",
		},
		{
			"code with non alphanum chars",
			"A#2",
			true,
			"Error code must only contain uppercase letters and numbers.",
		},
		{
			"code not capitalized",
			"c3A",
			true,
			"Error code must contain only uppercase letters.",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			pool := &ErrorPool{}

			if tc.shouldPanic {
				assert.PanicsWithError(tc.panicMessage, func() {
					pool.checkCode(tc.input)
				})
			} else {
				pool.checkCode(tc.input)
			}
		})
	}
}

func Test_checkModule(t *testing.T) {
	testcases := []struct {
		name         string
		input        string
		shouldPanic  bool
		panicMessage string
	}{
		{
			"valid module",
			"Not Found Error",
			false,
			"",
		},
		{
			"empty module",
			"",
			true,
			"Error module must not be empty.",
		},
		{
			"invalid module chars",
			"This 1s An Inv@lid Module",
			true,
			"Error module must only contain letters and spaces.",
		},
		{
			"module with leading spaces",
			"      This Is An Invalid Module",
			true,
			"Error module must not contain leading or trailing spaces.",
		},
		{
			"module with trailing spaces",
			"This Is An Invalid Module     ",
			true,
			"Error module must not contain leading or trailing spaces.",
		},
		{
			"module with consecutive spaces",
			"This Is      An Invalid Module",
			true,
			"Error module must not contain consecutive spaces.",
		},
		{
			"no title case",
			"this is an invalid Module",
			true,
			"Error module must follow the title case convention (first letter of each word in uppercase separated by a space).",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			pool := &ErrorPool{}

			if tc.shouldPanic {
				assert.PanicsWithError(tc.panicMessage, func() {
					pool.checkModule(tc.input)
				})
			} else {
				pool.checkModule(tc.input)
			}
		})
	}
}

func Test_checkFix(t *testing.T) {
	testcases := []struct {
		name         string
		input        string
		shouldPanic  bool
		panicMessage string
	}{
		{
			"valid fix",
			"This is a valid error fix.",
			false,
			"",
		},
		{
			"empty fix",
			"",
			true,
			"Error fix must not be empty.",
		},
		{
			"without letters",
			"@!#@#",
			true,
			"Error fix must contain at least one letter.",
		},
		{
			"fix with control chars",
			"This is an invalid fix.\n",
			true,
			"Error fix must not contain control characters.",
		},
		{
			"first char not capitalized",
			"this is an invalid fix.",
			true,
			"Error fix must start with uppercase letter.",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			pool := &ErrorPool{}

			if tc.shouldPanic {
				assert.PanicsWithError(tc.panicMessage, func() {
					pool.checkFix(tc.input)
				})
			} else {
				pool.checkFix(tc.input)
			}
		})
	}
}

func Test_validComplete(t *testing.T) {
	testcases := []struct {
		name         string
		input        errorData
		shouldPanic  bool
		panicMessage string
	}{
		{
			"valid error",
			errorData{
				name:    "Not found error",
				message: "Could not find object.",
				code:    "ABC",
				fix:     "Try searching for the object.",
			},
			false,
			"",
		},
		{
			"invalid name",
			errorData{
				name:    "",
				message: "Could not find object.",
				code:    "ABC",
				fix:     "Try searching for the object.",
			},
			true,
			"Error name must not be empty.",
		},
		{
			"invalid message",
			errorData{
				name:    "Not found error",
				message: "",
				code:    "ABC",
				fix:     "Try searching for the object.",
			},
			true,
			"Error message must not be empty.",
		},
		{
			"invalid code",
			errorData{
				name:    "Not found error",
				message: "Could not find object.",
				code:    "",
				fix:     "Try searching for the object.",
			},
			true,
			"Error code must not be empty.",
		},
		{
			"invalid fix",
			errorData{
				name:    "Not found error",
				message: "Could not find object.",
				code:    "ABC",
				fix:     "",
			},
			true,
			"Error fix must not be empty.",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			pool := &ErrorPool{}

			if tc.shouldPanic {
				assert.PanicsWithError(tc.panicMessage, func() {
					pool.validateComplete(tc.input)
				})
			} else {
				pool.validateComplete(tc.input)
			}
		})
	}
}

func Test_validateIncomplete(t *testing.T) {
	testcases := []struct {
		name         string
		input        errorData
		shouldPanic  bool
		panicMessage string
	}{
		{
			"valid error",
			errorData{
				name:    "Not found error",
				message: "Could not find object.",
				code:    "ABC",
				fix:     "Try searching for the object.",
			},
			false,
			"",
		},
		{
			"invalid name",
			errorData{
				name:    "",
				message: "Could not find object.",
				code:    "ABC",
				fix:     "Try searching for the object.",
			},
			false,
			"",
		},
		{
			"invalid message",
			errorData{
				name:    "Not found error",
				message: "",
				code:    "ABC",
				fix:     "Try searching for the object.",
			},
			true,
			"Error message must not be empty.",
		},
		{
			"invalid code",
			errorData{
				name:    "Not found error",
				message: "Could not find object.",
				code:    "",
				fix:     "Try searching for the object.",
			},
			false,
			"",
		},
		{
			"invalid fix",
			errorData{
				name:    "Not found error",
				message: "Could not find object.",
				code:    "ABC",
				fix:     "",
			},
			false,
			"",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			pool := &ErrorPool{}

			if tc.shouldPanic {
				assert.PanicsWithError(tc.panicMessage, func() {
					pool.validateIncomplete(tc.input)
				})
			} else {
				pool.validateIncomplete(tc.input)
			}
		})
	}
}

func Test_validateData(t *testing.T) {
	testcases := []struct {
		name     string
		complete bool
		level    ErrorLevel
	}{
		{
			"validate fatal",
			true,
			Fatal,
		},
		{
			"validate error",
			true,
			Error,
		},
		{
			"validate warning",
			true,
			Warning,
		},
		{
			"validate info",
			false,
			Info,
		},
		{
			"validate debug",
			false,
			Debug,
		},
		{
			"validate trace",
			false,
			Trace,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			pool := &ErrorPool{}
			complete := errorData{
				name:    "Not found error",
				message: "Could not find object.",
				code:    "ABC",
				fix:     "Try searching for the object.",
			}
			incomplete := errorData{
				message: "Could not find object.",
			}

			if tc.complete {
				pool.validateData(tc.level, complete)
			} else {
				pool.validateData(tc.level, incomplete)
			}
		})
	}
}
