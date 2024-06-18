package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testLoggerObserver struct {
}

func (o *testLoggerObserver) OnUpdate(info ErrorInfo) {
}

func Test_InitLogger(t *testing.T) {
	t.Run("init logger", func(t *testing.T) {
		t.Cleanup(func() {
			loggerInstance = nil
		})
		assert := assert.New(t)
		pool := &ErrorPool{}
		InitLogger(pool)
		assert.Equal(loggerInstance != nil, true)
	})

	t.Run("nil pool", func(t *testing.T) {
		t.Cleanup(func() {
			loggerInstance = nil
		})
		assert := assert.New(t)
		assert.PanicsWithError("Cannot init logger if pool is nil.", func() {
			InitLogger(nil)
		})
	})
}

func Test_GetLogger(t *testing.T) {
	t.Run("initialized logger", func(t *testing.T) {
		t.Cleanup(func() {
			loggerInstance = nil
		})
		assert := assert.New(t)
		pool := &ErrorPool{}
		data := &errorData{
			name:    "Not found error",
			message: "Could not find object.",
			code:    "ABC",
			fix:     "Try searching for the object.",
		}

		// Init logger
		InitLogger(pool)

		// Add Error
		logger := GetLogger()
		logger.pool.error(Error, "Testing", data, nil)

		// Recover instance
		logger2 := GetLogger()
		size := logger2.pool.GetErrorCount()

		assert.Equal(size, 1)
	})

	t.Run("uninitialized logger", func(t *testing.T) {
		t.Cleanup(func() {
			loggerInstance = nil
		})
		assert := assert.New(t)
		assert.PanicsWithError("Must initialize logger before accessing it.", func() {
			GetLogger()
		})
	})
}

func Test_IsLoggerInitialized(t *testing.T) {
	t.Run("initialized logger", func(t *testing.T) {
		t.Cleanup(func() {
			loggerInstance = nil
		})
		assert := assert.New(t)
		pool := &ErrorPool{}

		InitLogger(pool)
		assert.Equal(IsLoggerInitialized(), true)
	})

	t.Run("uninitialized logger", func(t *testing.T) {
		t.Cleanup(func() {
			loggerInstance = nil
		})
		assert := assert.New(t)
		assert.Equal(IsLoggerInitialized(), false)
	})
}

func Test_CloseLogger(t *testing.T) {
	t.Run("initialized logger", func(t *testing.T) {
		t.Cleanup(func() {
			loggerInstance = nil
		})
		assert := assert.New(t)
		pool := &ErrorPool{}
		observer := &testLoggerObserver{}

		data := &errorData{name: "Some error", message: "Some message", code: "ABC", fix: "Some fix"}
		pool.addError(data, nil)
		pool.Subscribe(observer)

		InitLogger(pool)
		assert.Equal(loggerInstance != nil, true)
		assert.Equal(loggerInstance.pool.HasErrors(), true)
		assert.Equal(len(loggerInstance.pool.observers) != 0, true)

		CloseLogger()
		assert.Equal(loggerInstance != nil, true)
		assert.Equal(loggerInstance.pool.HasErrors(), false)
		assert.Equal(len(loggerInstance.pool.observers) == 0, true)
	})

	t.Run("uninitialized logger", func(t *testing.T) {
		t.Cleanup(func() {
			loggerInstance = nil
		})
		assert := assert.New(t)
		assert.Equal(loggerInstance == nil, true)

		CloseLogger()
		assert.Equal(loggerInstance == nil, true)
	})
}

func Test_SetLoggerModule(t *testing.T) {
	t.Run("set module", func(t *testing.T) {
		t.Cleanup(func() {
			loggerInstance = nil
		})
		assert := assert.New(t)
		pool := &ErrorPool{}
		InitLogger(pool)

		assert.Equal(loggerInstance.module, "")

		SetLoggerModule("Testing Module")
		assert.Equal(loggerInstance.module, "Testing Module")
	})

	t.Run("uninitiliazed logger", func(t *testing.T) {
		t.Cleanup(func() {
			loggerInstance = nil
		})
		assert := assert.New(t)

		assert.PanicsWithError("Cannot set module if logger is not initialized.", func() {
			SetLoggerModule("Testing Module")
		})
	})
}
