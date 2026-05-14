package rest

// Response is a unified REST API response wrapper.
type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

// Success creates a success response.
func Success[T any](data T) Response[T] {
	return Response[T]{Code: 0, Message: "success", Data: data}
}

// Error creates an error response.
func Error[T any](code int, message string) Response[T] {
	return Response[T]{Code: code, Message: message}
}
