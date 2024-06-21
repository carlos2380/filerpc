package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	log "filerpc/internal/logger"
	"filerpc/internal/proto"
)

type HTTPFileResponse struct {
	Type    string `json:"type"`
	Version string `json:"version"`
	Hash    string `json:"hash"`
	Content string `json:"content"`
}

type responseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func (rw *responseWriter) Write(p []byte) (int, error) {
	return rw.buf.Write(p)
}

func customResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		crw := &responseWriter{ResponseWriter: w, buf: buf}

		next.ServeHTTP(crw, r)

		if crw.buf.Len() == 0 {
			http.NotFound(w, r)
			return
		}

		var fileResp proto.FileResponse
		err := json.Unmarshal(crw.buf.Bytes(), &fileResp)
		if err != nil {
			grpcErr, _ := status.FromError(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(runtime.HTTPStatusFromCode(grpcErr.Code()))
			err := json.NewEncoder(w).Encode(grpcErr.Proto())
			if err != nil {
				log.Logger.Error("Failed to encode grpc error response: " + err.Error())
			}
			return
		}

		if fileResp.Type == "" && fileResp.Version == "" && fileResp.Hash == "" && len(fileResp.Content) == 0 {
			http.Error(w, "Invalid request parameters", http.StatusBadRequest)
			return
		}

		modifiedResp := &HTTPFileResponse{
			Type:    fileResp.Type,
			Version: fileResp.Version,
			Hash:    fileResp.Hash,
			Content: string(fileResp.Content),
		}

		responseBytes, err := json.Marshal(modifiedResp)
		if err != nil {
			log.Logger.Error("Failed to marshal modified response: " + err.Error())
			http.Error(w, "Failed to marshal modified response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(responseBytes)))
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(responseBytes); err != nil {
			log.Logger.Error("Failed to write modified response: " + err.Error())
		}
	})
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RunGateway(ctx context.Context, host string, grpcPort string, gatewayPort string) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	grpcServerEndpoint := fmt.Sprintf("%s:%s", host, grpcPort)
	if err := proto.RegisterFileServiceHandlerFromEndpoint(ctx, mux, grpcServerEndpoint, opts); err != nil {
		return err
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/v1/file", customResponseMiddleware(mux))
	httpMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	log.Logger.Infof("Serving gRPC-Gateway on http://%s:%s", host, gatewayPort)
	return http.ListenAndServe(":"+gatewayPort, enableCORS(httpMux))
}
