// Copyright 2022-2023 The Connect Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"

	entitiesv1 "github.com/astronlabltd/protos-pkg/gen/go/entities/demo/v1"
	servicesv1 "github.com/astronlabltd/protos-pkg/gen/go/services/demo/v1"
	"github.com/astronlabltd/protos-pkg/gen/go/services/demo/v1/demov1connect"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type service1Server struct {
}

// NewDemoServer returns a new Eliza implementation which sleeps for the  provided duration between streaming responses.
func NewDemoServer() demov1connect.DemoServiceHandler {
	return &service1Server{}
}

func (s *service1Server) ListPrompts(ctx context.Context, req *connect.Request[servicesv1.ListPromptsRequest]) (*connect.Response[servicesv1.ListPromptsResponse], error) {
	reply := servicesv1.ListPromptsResponse{
		Prompts: []*entitiesv1.Prompt{
			{
				Id:   120,
				Text: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
			},
		},
		Pagination: &entitiesv1.Pagination{
			TotalPages: 10,
			Offset:     1,
		},
	}

	return &connect.Response[servicesv1.ListPromptsResponse]{
		Msg: &reply,
	}, nil

}

func (s *service1Server) SavePrompt(ctx context.Context, req *connect.Request[servicesv1.SavePromptRequest]) (*connect.Response[servicesv1.SavePromptResponse], error) {
	log.Println(req)
	log.Println("SavePrompt calleddd")
	reply := servicesv1.SavePromptResponse{
		Message: "Prompt saved successfully",
	}

	return &connect.Response[servicesv1.SavePromptResponse]{
		Msg: &reply,
	}, nil
}

// StreamQuotes implements demov1connect.DemoServiceHandler.
func (s *service1Server) StreamQuotes(ctx context.Context, stream *connect.ClientStream[servicesv1.StreamQuotesRequest]) (*connect.Response[servicesv1.StreamQuotesResponse], error) {
	log.Println("StreamQuotes called")

	var quotes strings.Builder
	for stream.Receive() {
		g := fmt.Sprintf("Hello, %s!\n", stream.Msg().NumberOfQuotes)
		if _, err := quotes.WriteString(g); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	if err := stream.Err(); err != nil {
		return nil, connect.NewError(connect.CodeUnknown, err)
	}
	res := connect.NewResponse(&servicesv1.StreamQuotesResponse{
		Quote: quotes.String(),
	})
	res.Header().Set("Quote-version", "1.0")
	return res, nil
}

// StreamMovieNames implements demov1connect.DemoServiceHandler.
func (s *service1Server) StreamMovieNames(ctx context.Context, req *connect.Request[servicesv1.StreamMovieNamesRequest], stream *connect.ServerStream[servicesv1.StreamMovieNamesResponse]) error {
	log.Println("StreamMovieNames called")

	var moviesName strings.Builder
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 10000; i++ {
		movieReq := req.Any()
		if movieReq == false {
			break
		}

		// Generate random movie name
		randomMovieName := generateRandomMovieName()

		// Generate timestamp
		timestamp := time.Now().Format(time.RFC3339)

		// Append movie name with timestamp to the string builder
		moviesName.WriteString(fmt.Sprintf("%s -. %s\n", timestamp, randomMovieName))
		stream.Send(&servicesv1.StreamMovieNamesResponse{
			MovieName: randomMovieName,
		})

		// send each msg after 5 sec
		time.Sleep(1 * time.Second)
	}

	// Example of how to log the constructed movie names
	log.Println(moviesName.String())

	return nil
}

// Function to generate a random movie name
func generateRandomMovieName() string {
	movieNames := []string{"The Matrix", "Inception", "Interstellar", "The Shawshank Redemption", "Pulp Fiction", "The Godfather", "Fight Club", "Forrest Gump", "The Dark Knight", "Goodfellas"}

	return movieNames[rand.Intn(len(movieNames))]
}

func newCORS() *cors.Cors {
	// To let web developers play with the demo service from browsers, we need a
	// very permissive CORS setup.
	return cors.New(cors.Options{
		Debug: true,
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowOriginFunc: func(origin string) bool {
			// Allow all origins, which effectively disables CORS.
			return true
		},
		// AllowedHeaders: []string{"connect-protocol-version"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{
			// Content-Type is in the default safelist.
			"Accept",
			"Accept-Encoding",
			"Accept-Post",
			"Connect-Accept-Encoding",
			"Connect-Content-Encoding",
			"Content-Encoding",
			"Grpc-Accept-Encoding",
			"Grpc-Encoding",
			"Grpc-Message",
			"Grpc-Status",
			"Grpc-Status-Details-Bin",
			"Access-Control-Allow-Origin",
			// "connect-protocol-version",
		},
		// Let browsers cache CORS information for longer, which reduces the number
		// of preflight requests. Any changes to ExposedHeaders won't take effect
		// until the cached data expires. FF caps this value at 24h, and modern
		// Chrome caps it at 2h.
		MaxAge: int(2 * time.Hour / time.Second),
	})
}

func main() {
	log.SetOutput(os.Stdout)

	mux := http.NewServeMux()
	compress1KB := connect.WithCompressMinBytes(1024)
	mux.Handle(demov1connect.NewDemoServiceHandler(
		NewDemoServer(),
		compress1KB,
	))
	mux.Handle(grpchealth.NewHandler(
		grpchealth.NewStaticChecker(demov1connect.DemoServiceName),
		compress1KB,
	))
	mux.Handle(grpcreflect.NewHandlerV1(
		grpcreflect.NewStaticReflector(demov1connect.DemoServiceName),
		compress1KB,
	))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(
		grpcreflect.NewStaticReflector(demov1connect.DemoServiceName),
		compress1KB,
	))

	addr := "127.0.0.1:8000"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}
	srv := &http.Server{
		Addr: addr,
		Handler: h2c.NewHandler(
			newCORS().Handler(mux),
			&http2.Server{},
		),
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		MaxHeaderBytes:    8 * 1024, // 8KiB
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		log.Printf("Listening on http://%s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP listen and serve: %v", err)
		}
	}()

	<-signals
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP shutdown: %v", err) //nolint:gocritic
	}
}
