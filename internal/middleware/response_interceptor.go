package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ResponseWrapper struct {
	Version string      `json:"version"`
	Data    interface{} `json:"data"`
}

func ResponseInterceptor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		rec := context.Response().Writer
		buf := new(bytes.Buffer)
		context.Response().Writer = &bodyWriter{ResponseWriter: rec, body: buf}

		if err := next(context); err != nil {
			context.Error(err)
		}

		respBody := buf.Bytes()
		var originalResponse interface{}

		if len(respBody) > 0 {
			if err := json.Unmarshal(respBody, &originalResponse); err != nil {
				return err
			}
		}

		wrappedResponse := ResponseWrapper{
			Version: "1.0.0", // Set your API version here
			Data:    originalResponse,
		}

		finalResponse, err := json.Marshal(wrappedResponse)
		if err != nil {
			return err
		}

		context.Response().Writer = rec
		context.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		context.Response().WriteHeader(http.StatusOK)
		context.Response().Write(finalResponse)

		return nil
	}
}

type bodyWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (writer *bodyWriter) Write(bytes []byte) (int, error) {
	return writer.body.Write(bytes)
}
