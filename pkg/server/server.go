package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/studio-b12/parrot/pkg/models"
)

const maxBodySize = 4 * 1024 * 1024 // 4 MiB

type Server struct {
	upstream string
	server   *http.Server
}

func New(bindAddress string, upstream string) (*Server, error) {
	_, err := url.Parse(upstream)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	server := http.Server{Addr: bindAddress, Handler: mux}

	t := &Server{
		upstream: upstream,
		server:   &server,
	}

	mux.HandleFunc("POST /{topic}", t.handlePush)
	mux.HandleFunc("PUT /{topic}", t.handlePush)

	return t, nil
}

func (t *Server) ListenAndServe() error {
	return t.server.ListenAndServe()
}

func (t *Server) handlePush(w http.ResponseWriter, r *http.Request) {
	topic := r.PathValue("topic")

	slog.Info("req", "method", r.Method, "topic", topic)

	var payload models.WebhookPayload
	limitReader := io.LimitReader(r.Body, maxBodySize)
	err := json.NewDecoder(limitReader).Decode(&payload)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid body format", err)
		return
	}

	notifications := payload.ToNotifications(topic)

	var errs error
	for _, notification := range notifications {
		var body bytes.Buffer
		err := json.NewEncoder(&body).Encode(notification)
		if err != nil {
			slog.Error("failed encoding notification payload", "topic", topic, "err", err)
			errs = errors.Join(errs, err)
			continue
		}

		slog.Debug("sending notification", "topic", topic, "content", body.String())

		req, _ := http.NewRequest(r.Method, t.upstream, &body)
		if authValue := r.Header.Get("Authorization"); authValue != "" {
			req.Header.Set("Authorization", authValue)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			slog.Error("failed sending upstream request", "topic", topic, "err", err)
			errs = errors.Join(errs, err)
			continue
		}
		if resp.StatusCode > 299 {
			slog.Error("upstream request had non-ok status code", "topic", topic, "status", resp.StatusCode)
			errs = errors.Join(errs, err)
			continue
		}
		slog.Info("sent notification", "topic", topic, "status", resp.StatusCode)
	}

	if errs != nil {
		respondError(w, http.StatusInternalServerError,
			"failed sending some notifications to upstream", errs)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func respondError(w http.ResponseWriter, statusCode int, message string, err ...error) {
	w.WriteHeader(statusCode)
	if len(err) != 0 {
		fmt.Fprintf(w, "%s: %s", message, err[0].Error())
	} else {
		w.Write([]byte(message))
	}
}
