package gateway

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/Nick-Spencer-SumUp/test-router/internal/config"
	"github.com/Nick-Spencer-SumUp/test-router/internal/config/mappings"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

var (
	BadRequestError          = errors.New("bad request")
	InternalServerError      = errors.New("internal server error")
	CountryNotSupportedError = errors.New("country not supported")
	PathNotSupportedError    = errors.New("path not supported")
)

func New() *Handler {
	return &Handler{}
}

// HandleRequest is a generic handler that can route any request to the appropriate upstream service
func (h *Handler) HandleRequest(c echo.Context) error {
	// Extract country from request header
	countryString := c.Request().Header.Get("country")
	country, err := config.GetCountryFromConfig(countryString)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":   CountryNotSupportedError.Error(),
			"details": err.Error(),
		})
	}

	// Get the request path
	requestPath := c.Request().URL.Path

	// Get configuration for this country and path
	countryConfig, err := config.SelectConfigByPath(country, requestPath)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":   PathNotSupportedError.Error(),
			"details": err.Error(),
		})
	}

	// Get endpoint configuration
	endpointConfig, err := countryConfig.GetEndpointConfigByPath(requestPath)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":   PathNotSupportedError.Error(),
			"details": err.Error(),
		})
	}

	// Read request body
	var requestBody []byte
	if c.Request().Body != nil {
		requestBody, err = io.ReadAll(c.Request().Body)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error":   BadRequestError.Error(),
				"details": "Failed to read request body",
			})
		}
	}

	// Forward request to upstream service
	response, err := h.forwardRequest(c.Request().Method, endpointConfig, requestBody, c.Request().Header)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   InternalServerError.Error(),
			"details": err.Error(),
		})
	}

	return h.streamResponse(c, response)
}

// forwardRequest sends the request to the upstream service
func (h *Handler) forwardRequest(method string, endpointConfig mappings.EndpointConfig, requestBody []byte, headers http.Header) (*http.Response, error) {
	// Build the upstream URL
	upstreamURL := endpointConfig.BaseURL + endpointConfig.Endpoint

	request, err := http.NewRequest(method, upstreamURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	// Copy headers (excluding hop-by-hop headers)
	for key, values := range headers {
		if isHopByHopHeader(key) {
			continue
		}
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	return http.DefaultClient.Do(request)
}

// streamResponse streams the upstream response back to the client
func (h *Handler) streamResponse(c echo.Context, response *http.Response) error {
	defer response.Body.Close()

	// Copy response headers
	for key, values := range response.Header {
		if isHopByHopHeader(key) {
			continue
		}
		for _, value := range values {
			c.Response().Header().Add(key, value)
		}
	}

	// Set status code
	c.Response().WriteHeader(response.StatusCode)

	// Stream response body
	_, err := io.Copy(c.Response().Writer, response.Body)
	return err
}

// isHopByHopHeader checks if a header is hop-by-hop and should not be forwarded
func isHopByHopHeader(header string) bool {
	hopByHopHeaders := []string{
		"Connection",
		"Keep-Alive",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"Te", // canonicalized version of "TE"
		"Trailers",
		"Transfer-Encoding",
		"Upgrade",
	}

	for _, hopHeader := range hopByHopHeaders {
		if header == hopHeader {
			return true
		}
	}
	return false
}
