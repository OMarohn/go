package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"kom.com/m/v2/src/kom.com/coaster/coaster"
)

func TestCoasterEchoServer(t *testing.T) {
	var assert = assert.New(t)
	var require = require.New(t)
	var testCoaster = coaster.Coaster{ID: "id123", Name: "TestCoaster", Manufacture: "TestManufature", Height: 123}

	server := CreateEchoServer()
	require.NotNil(server)

	go func() {
		if err := server.Start("localhost:8080"); err != nil && err != http.ErrServerClosed {
			server.Logger.Fatal("shutting down the server")
		}
	}()

	t.Run("lesen initial leerer mem coasters", func(t *testing.T) {
		res, err := http.Get("http://localhost:8080/mem/coasters")
		require.NoError(err)
		assert.NotNil(res)
		require.Equal(http.StatusOK, res.StatusCode)

		data, err := ioutil.ReadAll(res.Body)
		require.NoError(err)
		assert.NotNil(data)
		require.Equal(http.StatusOK, res.StatusCode)

		var coasters []coaster.Coaster
		err = json.Unmarshal(data, &coasters)
		require.NoError(err)
		assert.Len(coasters, 0)

	})

	t.Run("anlegen eines mem coasters", func(t *testing.T) {
		payloadBuf := new(bytes.Buffer)
		err := json.NewEncoder(payloadBuf).Encode(testCoaster)
		require.NoError(err)

		res, err := http.Post("http://localhost:8080/mem/coasters", "application/json", payloadBuf)
		require.NoError(err)
		assert.NotNil(res)
		require.Equal(http.StatusCreated, res.StatusCode)
	})

	t.Run("anlegen eines vorhandenen mem coasters", func(t *testing.T) {
		payloadBuf := new(bytes.Buffer)
		err := json.NewEncoder(payloadBuf).Encode(testCoaster)
		require.NoError(err)

		res, err := http.Post("http://localhost:8080/mem/coasters", "application/json", payloadBuf)
		require.NoError(err)
		assert.NotNil(res)
		require.Equal(http.StatusBadRequest, res.StatusCode)
	})

	t.Run("lesen mem coasters", func(t *testing.T) {
		res, err := http.Get("http://localhost:8080/mem/coasters")
		require.NoError(err)
		assert.NotNil(res)
		require.Equal(http.StatusOK, res.StatusCode)

		data, err := ioutil.ReadAll(res.Body)
		require.NoError(err)
		assert.NotNil(data)
		require.Equal(http.StatusOK, res.StatusCode)

		var coasters []coaster.Coaster
		err = json.Unmarshal(data, &coasters)
		require.NoError(err)
		assert.Len(coasters, 1)
		assert.Equal(testCoaster, coasters[0])
	})

	t.Run("lesen mem coaster by unknown id", func(t *testing.T) {
		res, err := http.Get("http://localhost:8080/mem/coasters/id9999")
		require.NoError(err)
		assert.NotNil(res)
		require.Equal(http.StatusNotFound, res.StatusCode)
	})

	t.Run("lesen mem coaster by id", func(t *testing.T) {
		res, err := http.Get("http://localhost:8080/mem/coasters/id123")
		require.NoError(err)
		assert.NotNil(res)

		data, err := ioutil.ReadAll(res.Body)
		require.NoError(err)
		assert.NotNil(data)
		require.Equal(http.StatusOK, res.StatusCode)

		var coaster coaster.Coaster
		err = json.Unmarshal(data, &coaster)
		require.NoError(err)
		assert.Equal(testCoaster, coaster)
	})

	t.Run("löschen mem coaster by id", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/mem/coasters/id123", nil)
		require.NoError(err)
		res, err := http.DefaultClient.Do(req)
		require.NoError(err)
		require.NotNil(res)

		require.Equal(http.StatusOK, res.StatusCode)
	})

	t.Run("löschen mem coaster by unknown id", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/mem/coasters/id99999", nil)
		require.NoError(err)
		res, err := http.DefaultClient.Do(req)
		require.NoError(err)
		require.NotNil(res)

		require.Equal(http.StatusNotFound, res.StatusCode)
	})

	/**
	Das ist noch ggf. mit dem Aufruf an dem 'Port' zu versehen -- bin mir aber noch nicht klar wie genau das funktioniert und ob das der bessere Weg ist
	*/
	t.Run("lesen (recorder) mem coaster by id", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := server.NewContext(req, rec)
		c.SetPath("/mem/coasters/:id")
		c.SetParamNames("id")
		c.SetParamValues("id123")

		res := rec.Result()
		fmt.Println(ioutil.ReadAll(res.Body))
		defer res.Body.Close()

	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		server.Logger.Fatal(err)
	}

}
