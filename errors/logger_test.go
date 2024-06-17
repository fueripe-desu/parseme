package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_InitLogger(t *testing.T) {
	t.Run("init logger", func(t *testing.T) {
		t.Cleanup(func() {
			CloseLogger()
		})
		assert := assert.New(t)
		pool := &ErrorPool{}
		InitLogger(pool)
		assert.Equal(loggerInstance != nil, true)
	})

	t.Run("nil pool", func(t *testing.T) {
		t.Cleanup(func() {
			CloseLogger()
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
			CloseLogger()
		})
		assert := assert.New(t)
		pool := &ErrorPool{}

		// Init logger
		InitLogger(pool)

		// Add Error
		logger := GetLogger()
		logger.pool.error(Error, nilErrorDataError, nil)

		// Recover instance
		logger2 := GetLogger()
		size := logger2.pool.GetErrorCount()

		assert.Equal(size, 1)
	})

	t.Run("uninitialized logger", func(t *testing.T) {
		t.Cleanup(func() {
			CloseLogger()
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
			CloseLogger()
		})
		assert := assert.New(t)
		pool := &ErrorPool{}

		InitLogger(pool)
		assert.Equal(IsLoggerInitialized(), true)
	})

	t.Run("uninitialized logger", func(t *testing.T) {
		t.Cleanup(func() {
			CloseLogger()
		})
		assert := assert.New(t)
		assert.Equal(IsLoggerInitialized(), false)
	})
}

func Test_CloseLogger(t *testing.T) {
	t.Run("initialized logger", func(t *testing.T) {
		t.Cleanup(func() {
			CloseLogger()
		})
		assert := assert.New(t)
		pool := &ErrorPool{}

		InitLogger(pool)
		assert.Equal(loggerInstance != nil, true)

		CloseLogger()
		assert.Equal(loggerInstance == nil, true)
	})

	t.Run("uninitialized logger", func(t *testing.T) {
		t.Cleanup(func() {
			CloseLogger()
		})
		assert := assert.New(t)
		assert.Equal(loggerInstance == nil, true)

		CloseLogger()
		assert.Equal(loggerInstance == nil, true)
	})
}
