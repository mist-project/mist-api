package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/render"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ----- GRPC -----

// ----- ERROR HANDLERS -----

type DataResponse struct {
	Meta interface{} `json:"meta,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Detail string `json:"detail,omitempty"`
}

func HandleGrpcError(w http.ResponseWriter, r *http.Request, err error) {
	s, _ := status.FromError(err)

	// Log the error
	log.Printf("Error from service: %v\n", err)

	// Map gRPC status code to HTTP status and error message
	httpStatus, message := mapGrpcStatusToHTTP(s.Code(), s.Message())

	// Set the HTTP status and send the error response
	render.Status(r, httpStatus)
	render.JSON(w, r, &ErrorResponse{Detail: message})
}

func mapGrpcStatusToHTTP(code codes.Code, grpcMessage string) (int, string) {
	switch code {
	case codes.Unavailable:
		return http.StatusBadGateway, "Server is unresponsive."
	case codes.DeadlineExceeded:
		return http.StatusBadGateway, "Server timed out."
	case codes.Canceled:
		return http.StatusBadGateway, "Server error."
	case codes.Unauthenticated:
		return http.StatusUnauthorized, "Unauthorized request."
	case codes.NotFound:
		return http.StatusNotFound, "Not found."
	case codes.AlreadyExists:
		return http.StatusConflict, "Resource already exists."
	case codes.InvalidArgument:
		return http.StatusBadRequest, grpcMessage
	default:
		// Log unhandled gRPC error codes
		log.Printf("Unhandled gRPC error code: %v\n", code)
		return http.StatusInternalServerError, "Internal Server Error."
	}
}

func DecodeRequestBody(w http.ResponseWriter, r *http.Request, bind interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&bind)
	if err != nil {
		// TODO: use better logging solution
		log.Printf("Error while decoding: %v\n", err)

		// If there is an error in decoding, return 400 Bad Request
		render.Status(r, http.StatusUnprocessableEntity)
		render.JSON(w, r, CreateErrorResponse("Invalid attributes provided."))

		return err
	}
	return nil
}

func CreateErrorResponse(detail string) *ErrorResponse {
	return &ErrorResponse{
		Detail: detail,
	}
}

func CreateResponse(data interface{}) *DataResponse {
	return &DataResponse{
		Meta: nil,
		Data: data,
	}
}
