// snake v1.0.0 23da2b4947ef517e519e579f85d4618609d9bd42
// --
// Code generated by webrpc-gen@v0.14.0-dev with ../../gen-golang generator. DO NOT EDIT.
//
// webrpc-gen -schema=proto/snake.ridl -target=../../gen-golang -pkg=proto -server -client -fmt=false -out=proto/snake.gen.go
package proto

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// WebRPC description and code-gen version
func WebRPCVersion() string {
	return "v1"
}

// Schema version of your RIDL schema
func WebRPCSchemaVersion() string {
	return "v1.0.0"
}

// Schema hash generated from your RIDL schema
func WebRPCSchemaHash() string {
	return "23da2b4947ef517e519e579f85d4618609d9bd42"
}


//
// Types
//

type Direction uint

const (
	Direction_left Direction = 0
	Direction_right Direction = 1
	Direction_up Direction = 2
	Direction_down Direction = 3
)

var Direction_name = map[uint]string{
	0: "left",
	1: "right",
	2: "up",
	3: "down",
}

var Direction_value = map[string]uint{
	"left": 0,
	"right": 1,
	"up": 2,
	"down": 3,
}

func (x Direction) String() string {
	return Direction_name[uint(x)]
}

func (x Direction) MarshalText() ([]byte, error) {
	return []byte(Direction_name[uint(x)]), nil
}

func (x *Direction) UnmarshalText(b []byte) error {
	*x = Direction(Direction_value[string(b)])
	return nil
}

func (x *Direction) Is(values ...Direction) bool {
	if x == nil {
		return false
	}
	for _, v := range values {
		if *x == v {
			return true
		}
	}
	return false
}

type ItemType uint

const (
	ItemType_bite ItemType = 0
)

var ItemType_name = map[uint]string{
	0: "bite",
}

var ItemType_value = map[string]uint{
	"bite": 0,
}

func (x ItemType) String() string {
	return ItemType_name[uint(x)]
}

func (x ItemType) MarshalText() ([]byte, error) {
	return []byte(ItemType_name[uint(x)]), nil
}

func (x *ItemType) UnmarshalText(b []byte) error {
	*x = ItemType(ItemType_value[string(b)])
	return nil
}

func (x *ItemType) Is(values ...ItemType) bool {
	if x == nil {
		return false
	}
	for _, v := range values {
		if *x == v {
			return true
		}
	}
	return false
}

type State struct {
	Width uint `json:"width"`
	Height uint `json:"height"`
	Snakes map[uint64]*Snake `json:"snakes"`
	Items map[uint64]*Item `json:"items"`
}

type Snake struct {
	Id uint64 `json:"id"`
	Name string `json:"name"`
	Color string `json:"color"`
	Body []*Square `json:"body"`
	Direction *Direction `json:"direction"`
	NextDirections []*Direction `json:"nextDirections"`
	Length int `json:"length"`
	BornAt time.Time `json:"bornAt"`
	DiedAt time.Time `json:"diedAt"`
}

type Item struct {
	Id uint64 `json:"id"`
	Color string `json:"color"`
	Type *ItemType `json:"type"`
	Body *Square `json:"body"`
}

type Event struct {
}

type Square struct {
	X uint `json:"x"`
	Y uint `json:"y"`
}

var WebRPCServices = map[string][]string{
	"SnakeGame": {
		"JoinGame",
		"CreateSnake",
		"TurnSnake",
	},
}

//
// Server types
//

type SnakeGame interface {
	JoinGame(ctx context.Context, stream JoinGameStreamWriter) error
	CreateSnake(ctx context.Context, username string) (uint64, error)
	TurnSnake(ctx context.Context, snakeId uint64, direction *Direction) error
}
type JoinGameStreamWriter interface {
	Write(state *State, event *Event) error
}



type joinGameStreamWriter struct {
	streamWriter
}

func (w *joinGameStreamWriter) Write(state *State, event *Event) error {
	out := struct {
		Ret0 *State `json:"state"`
		Ret1 *Event `json:"event"`
	}{
		Ret0: state,
		Ret1: event,
	}

	return w.streamWriter.write(out)
}

type streamWriter struct {
	mu sync.Mutex // Guards concurrent writes to w.
	w  http.ResponseWriter
	f  http.Flusher
	e  *json.Encoder

	sendError func(w http.ResponseWriter, r *http.Request, rpcErr WebRPCError)
}

const StreamKeepAliveInterval = 10*time.Second

func (w *streamWriter) keepAlive(ctx context.Context) {
	for {
		select {
		case <-time.After(StreamKeepAliveInterval):
			err := w.ping()
			if err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (w *streamWriter) ping() error {
	defer w.f.Flush()

	w.mu.Lock()
	defer w.mu.Unlock()

	_, err := w.w.Write([]byte("\n"))
	return err
}

func (w *streamWriter) write(respPayload interface{}) error {
	defer w.f.Flush()

	w.mu.Lock()
	defer w.mu.Unlock()

	return w.e.Encode(respPayload)
}

//
// Client types
//

type SnakeGameClient interface {
	JoinGame(ctx context.Context) (JoinGameStreamReader, error)
	CreateSnake(ctx context.Context, username string) (uint64, error)
	TurnSnake(ctx context.Context, snakeId uint64, direction *Direction) error
}
type JoinGameStreamReader interface {
	Read() (state *State, event *Event, err error)
}




//
// Server
//

type WebRPCServer interface {
	http.Handler
}

type snakeGameServer struct {
	SnakeGame
	OnError func(r *http.Request, rpcErr *WebRPCError)
}

func NewSnakeGameServer(svc SnakeGame) *snakeGameServer {
	return &snakeGameServer{
		SnakeGame: svc,
	}
}

func (s *snakeGameServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		// In case of a panic, serve a HTTP 500 error and then panic.
		if rr := recover(); rr != nil {
			s.sendErrorJSON(w, r, ErrWebrpcServerPanic.WithCause(fmt.Errorf("%v", rr)))
			panic(rr)
		}
	}()

	ctx := r.Context()
	ctx = context.WithValue(ctx, HTTPResponseWriterCtxKey, w)
	ctx = context.WithValue(ctx, HTTPRequestCtxKey, r)
	ctx = context.WithValue(ctx, ServiceNameCtxKey, "SnakeGame")

	var handler func(ctx context.Context, w http.ResponseWriter, r *http.Request)
	switch r.URL.Path {
	case "/rpc/SnakeGame/JoinGame":
		handler = s.serveJoinGameJSONStream
	case "/rpc/SnakeGame/CreateSnake":
		handler = s.serveCreateSnakeJSON
	case "/rpc/SnakeGame/TurnSnake":
		handler = s.serveTurnSnakeJSON
	default:
		err := ErrWebrpcBadRoute.WithCause(fmt.Errorf("no handler for path %q", r.URL.Path))
		s.sendErrorJSON(w, r, err)
		return
	}

	if r.Method != "POST" {
		w.Header().Add("Allow", "POST") // RFC 9110.
		err := ErrWebrpcBadMethod.WithCause(fmt.Errorf("unsupported method %q (only POST is allowed)", r.Method))
		s.sendErrorJSON(w, r, err)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if i := strings.Index(contentType, ";"); i >= 0 {
		contentType = contentType[:i]
	}
	contentType = strings.TrimSpace(strings.ToLower(contentType))

	switch contentType {
	case "application/json":
		handler(ctx, w, r)
	default:
		err := ErrWebrpcBadRequest.WithCause(fmt.Errorf("unexpected Content-Type: %q", r.Header.Get("Content-Type")))
		s.sendErrorJSON(w, r, err)
	}
}

func (s *snakeGameServer) serveJoinGameJSONStream(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	ctx = context.WithValue(ctx, MethodNameCtxKey, "JoinGame")

	

	f, ok := w.(http.Flusher)
	if !ok {
		s.sendErrorJSON(w, r, ErrWebrpcInternalError.WithCause(fmt.Errorf("server http.ResponseWriter doesn't support .Flush() method")))
		return
	}

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "application/x-ndjson")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)

	streamWriter := &joinGameStreamWriter{streamWriter{w: w, f: f, e: json.NewEncoder(w), sendError: s.sendErrorJSON}}
	if err := streamWriter.ping(); err != nil {
		s.sendErrorJSON(w, r, ErrWebrpcStreamLost.WithCause(fmt.Errorf("failed to establish SSE stream: %w", err)))
		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	go streamWriter.keepAlive(ctx)

	// Call service method implementation.
	if err := s.SnakeGame.JoinGame(ctx, streamWriter); err != nil {
		rpcErr, ok := err.(WebRPCError)
		if !ok {
			rpcErr = ErrWebrpcEndpoint.WithCause(err)
		}
		streamWriter.sendError(w, r, rpcErr)
		return
	}
}
func (s *snakeGameServer) serveCreateSnakeJSON(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	ctx = context.WithValue(ctx, MethodNameCtxKey, "CreateSnake")

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		s.sendErrorJSON(w, r, ErrWebrpcBadRequest.WithCause(fmt.Errorf("failed to read request data: %w", err)))
		return
	}
	defer r.Body.Close()

	reqPayload := struct {
		Arg0 string `json:"username"`
	}{}
	if err := json.Unmarshal(reqBody, &reqPayload); err != nil {
		s.sendErrorJSON(w, r, ErrWebrpcBadRequest.WithCause(fmt.Errorf("failed to unmarshal request data: %w", err)))
		return
	}

	// Call service method implementation.
	ret0, err := s.SnakeGame.CreateSnake(ctx, reqPayload.Arg0)
	if err != nil {
		rpcErr, ok := err.(WebRPCError)
		if !ok {
			rpcErr = ErrWebrpcEndpoint.WithCause(err)
		}
		s.sendErrorJSON(w, r, rpcErr)
		return
	}

	respPayload := struct {
		Ret0 uint64 `json:"snakeId"`
	}{ret0}
	respBody, err := json.Marshal(respPayload)
	if err != nil {
		s.sendErrorJSON(w, r, ErrWebrpcBadResponse.WithCause(fmt.Errorf("failed to marshal json response: %w", err)))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func (s *snakeGameServer) serveTurnSnakeJSON(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	ctx = context.WithValue(ctx, MethodNameCtxKey, "TurnSnake")

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		s.sendErrorJSON(w, r, ErrWebrpcBadRequest.WithCause(fmt.Errorf("failed to read request data: %w", err)))
		return
	}
	defer r.Body.Close()

	reqPayload := struct {
		Arg0 uint64 `json:"snakeId"`
		Arg1 *Direction `json:"direction"`
	}{}
	if err := json.Unmarshal(reqBody, &reqPayload); err != nil {
		s.sendErrorJSON(w, r, ErrWebrpcBadRequest.WithCause(fmt.Errorf("failed to unmarshal request data: %w", err)))
		return
	}

	// Call service method implementation.
	err = s.SnakeGame.TurnSnake(ctx, reqPayload.Arg0, reqPayload.Arg1)
	if err != nil {
		rpcErr, ok := err.(WebRPCError)
		if !ok {
			rpcErr = ErrWebrpcEndpoint.WithCause(err)
		}
		s.sendErrorJSON(w, r, rpcErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}


func (s *snakeGameServer) sendErrorJSON(w http.ResponseWriter, r *http.Request, rpcErr WebRPCError) {
	if s.OnError != nil {
		s.OnError(r, &rpcErr)
	}

	if w.Header().Get("Content-Type") == "application/x-ndjson" {
		out := struct {
			WebRPCError WebRPCError `json:"webrpcError"`
		}{ WebRPCError: rpcErr }
		json.NewEncoder(w).Encode(out)
		return	
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rpcErr.HTTPStatus)

	respBody, _ := json.Marshal(rpcErr)
	w.Write(respBody)
}

func RespondWithError(w http.ResponseWriter, err error) {
	rpcErr, ok := err.(WebRPCError)
	if !ok {
		rpcErr = ErrWebrpcEndpoint.WithCause(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rpcErr.HTTPStatus)

	respBody, _ := json.Marshal(rpcErr)
	w.Write(respBody)
}



//
// Client
//

const SnakeGamePathPrefix = "/rpc/SnakeGame/"

type snakeGameClient struct {
	client HTTPClient
	urls	 [3]string
}

func NewSnakeGameClient(addr string, client HTTPClient) SnakeGameClient {
	prefix := urlBase(addr) + SnakeGamePathPrefix
	urls := [3]string{
		prefix + "JoinGame",
		prefix + "CreateSnake",
		prefix + "TurnSnake",
	}
	return &snakeGameClient{
		client: client,
		urls:	 urls,
	}
}

func (c *snakeGameClient) JoinGame(ctx context.Context) (JoinGameStreamReader, error) {

	resp, err := doHTTPRequest(ctx, c.client, c.urls[0], nil, nil)
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return nil, err
	}

	buf := bufio.NewReader(resp.Body)
	return &joinGameStreamReader{streamReader{ctx: ctx, c: resp.Body, r: buf, d: json.NewDecoder(buf)}}, nil
}

type joinGameStreamReader struct {
	streamReader
}

func (r *joinGameStreamReader) Read() (*State, *Event, error) {
	out := struct {
		Ret0 *State `json:"state"`
		Ret1 *Event `json:"event"`
		WebRPCError *WebRPCError `json:"webrpcError"`
	}{}

	err := r.streamReader.read(&out)
	if err != nil {
		return out.Ret0, out.Ret1, err
	}

	if out.WebRPCError != nil {
		return out.Ret0, out.Ret1, out.WebRPCError
	}

	return out.Ret0, out.Ret1, nil
}

func (c *snakeGameClient) CreateSnake(ctx context.Context, username string) (uint64, error) {
	in := struct {
		Arg0 string `json:"username"`
	}{username}
	out := struct {
		Ret0 uint64 `json:"snakeId"`
	}{}

	resp, err := doHTTPRequest(ctx, c.client, c.urls[1], in, &out)
	if resp != nil {
		cerr := resp.Body.Close()
		if err == nil && cerr != nil {
			err = ErrWebrpcRequestFailed.WithCause(fmt.Errorf("failed to close response body: %w", cerr))
		}
	}

	return out.Ret0, err
}

func (c *snakeGameClient) TurnSnake(ctx context.Context, snakeId uint64, direction *Direction) error {
	in := struct {
		Arg0 uint64 `json:"snakeId"`
		Arg1 *Direction `json:"direction"`
	}{snakeId, direction}

	resp, err := doHTTPRequest(ctx, c.client, c.urls[2], in, nil)
	if resp != nil {
		cerr := resp.Body.Close()
		if err == nil && cerr != nil {
			err = ErrWebrpcRequestFailed.WithCause(fmt.Errorf("failed to close response body: %w", cerr))
		}
	}

	return err
}

type streamReader struct {
	ctx context.Context
	c   io.Closer
	r   *bufio.Reader
	d   *json.Decoder
}

func (r *streamReader) read(v interface {}) error {
	for {
		// Read newlines (keep-alive pings) and unblock decoder on ctx timeout.
		select {
		case <-r.ctx.Done():
			r.c.Close()
			return ErrWebrpcClientDisconnected.WithCause(r.ctx.Err())
		default:
		}

		b, err := r.r.ReadByte()
		if err != nil {
			return r.handleReadError(err)
		}
		if b != '\n' {
			r.r.UnreadByte()
			break
		}
	}

	if err := r.d.Decode(&v); err != nil {
		return r.handleReadError(err)
	}

	return nil
}

func (r *streamReader) handleReadError(err error) error {
	defer r.c.Close()
	if errors.Is(err, io.EOF) {
		return ErrWebrpcStreamFinished.WithCause(err)
	}
	if errors.Is(err, io.ErrUnexpectedEOF) {
		return ErrWebrpcStreamLost.WithCause(err)
	}
	return ErrWebrpcBadResponse.WithCause(fmt.Errorf("reading stream: %w", err))
}

// HTTPClient is the interface used by generated clients to send HTTP requests.
// It is fulfilled by *(net/http).Client, which is sufficient for most users.
// Users can provide their own implementation for special retry policies.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// urlBase helps ensure that addr specifies a scheme. If it is unparsable
// as a URL, it returns addr unchanged.
func urlBase(addr string) string {
	// If the addr specifies a scheme, use it. If not, default to
	// http. If url.Parse fails on it, return it unchanged.
	url, err := url.Parse(addr)
	if err != nil {
		return addr
	}
	if url.Scheme == "" {
		url.Scheme = "http"
	}
	return url.String()
}

// newRequest makes an http.Request from a client, adding common headers.
func newRequest(ctx context.Context, url string, reqBody io.Reader, contentType string) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", contentType)
	req.Header.Set("Content-Type", contentType)
	if headers, ok := HTTPRequestHeaders(ctx); ok {
		for k := range headers {
			for _, v := range headers[k] {
				req.Header.Add(k, v)
			}
		}
	}
	return req, nil
}

// doHTTPRequest is common code to make a request to the remote service.
func doHTTPRequest(ctx context.Context, client HTTPClient, url string, in, out interface{}) (*http.Response, error) {
	reqBody, err := json.Marshal(in)
	if err != nil {
		return nil, ErrWebrpcRequestFailed.WithCause(fmt.Errorf("failed to marshal JSON body: %w", err))
	}
	if err = ctx.Err(); err != nil {
		return nil, ErrWebrpcRequestFailed.WithCause(fmt.Errorf("aborted because context was done: %w", err))
	}

	req, err := newRequest(ctx, url, bytes.NewBuffer(reqBody), "application/json")
	if err != nil {
		return nil, ErrWebrpcRequestFailed.WithCause(fmt.Errorf("could not build request: %w", err))
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, ErrWebrpcRequestFailed.WithCause(err)
	}

	if resp.StatusCode != 200 {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, ErrWebrpcBadResponse.WithCause(fmt.Errorf("failed to read server error response body: %w", err))
		}

		var rpcErr WebRPCError
		if err := json.Unmarshal(respBody, &rpcErr); err != nil {
			return nil, ErrWebrpcBadResponse.WithCause(fmt.Errorf("failed to unmarshal server error: %w", err))
		}
		if rpcErr.Cause != "" {
			rpcErr.cause = errors.New(rpcErr.Cause)
		}
		return nil, rpcErr
	}

	if out != nil {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, ErrWebrpcBadResponse.WithCause(fmt.Errorf("failed to read response body: %w", err))
		}

		err = json.Unmarshal(respBody, &out)
		if err != nil {
			return nil, ErrWebrpcBadResponse.WithCause(fmt.Errorf("failed to unmarshal JSON response body: %w", err))
		}
	}

	return resp, nil
}

func WithHTTPRequestHeaders(ctx context.Context, h http.Header) (context.Context, error) {
	if _, ok := h["Accept"]; ok {
		return nil, errors.New("provided header cannot set Accept")
	}
	if _, ok := h["Content-Type"]; ok {
		return nil, errors.New("provided header cannot set Content-Type")
	}

	copied := make(http.Header, len(h))
	for k, vv := range h {
		if vv == nil {
			copied[k] = nil
			continue
		}
		copied[k] = make([]string, len(vv))
		copy(copied[k], vv)
	}

	return context.WithValue(ctx, HTTPClientRequestHeadersCtxKey, copied), nil
}

func HTTPRequestHeaders(ctx context.Context) (http.Header, bool) {
	h, ok := ctx.Value(HTTPClientRequestHeadersCtxKey).(http.Header)
	return h, ok
}

//
// Helpers
//

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "webrpc context value " + k.name
}

var (
	HTTPClientRequestHeadersCtxKey = &contextKey{"HTTPClientRequestHeaders"}
	HTTPResponseWriterCtxKey = &contextKey{"HTTPResponseWriter"}

	HTTPRequestCtxKey = &contextKey{"HTTPRequest"}

	ServiceNameCtxKey = &contextKey{"ServiceName"}

	MethodNameCtxKey = &contextKey{"MethodName"}
)

func ServiceNameFromContext(ctx context.Context) string {
	service, _ := ctx.Value(ServiceNameCtxKey).(string)
	return service
}

func MethodNameFromContext(ctx context.Context) string {
	method, _ := ctx.Value(MethodNameCtxKey).(string)
	return method
}

func RequestFromContext(ctx context.Context) *http.Request {
	r, _ := ctx.Value(HTTPRequestCtxKey).(*http.Request)
	return r
}
func ResponseWriterFromContext(ctx context.Context) http.ResponseWriter {
	w, _ := ctx.Value(HTTPResponseWriterCtxKey).(http.ResponseWriter)
	return w
}

//
// Errors
//

type WebRPCError struct {
	Name       string `json:"error"`
	Code       int    `json:"code"`
	Message    string `json:"msg"`
	Cause      string `json:"cause,omitempty"`
	HTTPStatus int    `json:"status"`
	cause      error
}

var _ error = WebRPCError{}

func (e WebRPCError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s %d: %s: %v", e.Name, e.Code, e.Message, e.cause)
	}
	return fmt.Sprintf("%s %d: %s", e.Name, e.Code, e.Message)
}

func (e WebRPCError) Is(target error) bool {
	if rpcErr, ok := target.(WebRPCError); ok {
		return rpcErr.Code == e.Code
	}
	return errors.Is(e.cause, target)
}

func (e WebRPCError) Unwrap() error {
	return e.cause
}

func (e WebRPCError) WithCause(cause error) WebRPCError {
	err := e
	err.cause = cause
	err.Cause = cause.Error()
	return err
}

// Deprecated: Use .WithCause() method on WebRPCError.
func ErrorWithCause(rpcErr WebRPCError, cause error) WebRPCError {
	return rpcErr.WithCause(cause)
}

// Webrpc errors
var (
	ErrWebrpcEndpoint = WebRPCError{Code: 0, Name: "WebrpcEndpoint", Message: "endpoint error", HTTPStatus: 400}
	ErrWebrpcRequestFailed = WebRPCError{Code: -1, Name: "WebrpcRequestFailed", Message: "request failed", HTTPStatus: 400}
	ErrWebrpcBadRoute = WebRPCError{Code: -2, Name: "WebrpcBadRoute", Message: "bad route", HTTPStatus: 404}
	ErrWebrpcBadMethod = WebRPCError{Code: -3, Name: "WebrpcBadMethod", Message: "bad method", HTTPStatus: 405}
	ErrWebrpcBadRequest = WebRPCError{Code: -4, Name: "WebrpcBadRequest", Message: "bad request", HTTPStatus: 400}
	ErrWebrpcBadResponse = WebRPCError{Code: -5, Name: "WebrpcBadResponse", Message: "bad response", HTTPStatus: 500}
	ErrWebrpcServerPanic = WebRPCError{Code: -6, Name: "WebrpcServerPanic", Message: "server panic", HTTPStatus: 500}
	ErrWebrpcInternalError = WebRPCError{Code: -7, Name: "WebrpcInternalError", Message: "internal error", HTTPStatus: 500}
	ErrWebrpcClientDisconnected = WebRPCError{Code: -8, Name: "WebrpcClientDisconnected", Message: "client disconnected", HTTPStatus: 400}
	ErrWebrpcStreamLost = WebRPCError{Code: -9, Name: "WebrpcStreamLost", Message: "stream lost", HTTPStatus: 400}
	ErrWebrpcStreamFinished = WebRPCError{Code: -10, Name: "WebrpcStreamFinished", Message: "stream finished", HTTPStatus: 200}
)

// Schema errors
var (
	ErrSnakeNotFound = WebRPCError{Code: 100, Name: "SnakeNotFound", Message: "Snake not found.", HTTPStatus: 400}
	ErrInvalidTurn = WebRPCError{Code: 101, Name: "InvalidTurn", Message: "Can't reverse direction.", HTTPStatus: 400}
)
