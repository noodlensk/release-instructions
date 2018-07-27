package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type appHTTPHandler struct {
	app *App
	h   func(app *App, w http.ResponseWriter, r *http.Request) (int, error)
}

func (ah appHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	curLog := log.WithFields(log.Fields{
		"remoteAddr": r.RemoteAddr,
		"method":     r.Method,
		"URL":        r.URL.String(),
	})
	curLog.Info("Recieved request")
	status, err := ah.h(ah.app, w, r)
	if err != nil {
		curLog.WithField("status", status).Errorf("Got error: %v", err)
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(status), status)
		default:
			http.Error(w, http.StatusText(status), status)
		}
	}
}

func indexHandler(a *App, w http.ResponseWriter, r *http.Request) (int, error) {
	http.StripPrefix("", http.FileServer(http.Dir("react-app/dist/"))).ServeHTTP(w, r)
	return http.StatusOK, nil
}

func releaseInstrucionsHandler(a *App, w http.ResponseWriter, r *http.Request) (int, error) {
	var project, fixVersion string
	if project = r.URL.Query().Get("project"); project == "" {
		return http.StatusBadRequest, errors.Errorf("project param cannot be empty")
	}

	if fixVersion = r.URL.Query().Get("fixVersion"); fixVersion == "" {
		return http.StatusBadRequest, errors.Errorf("fixVersion param cannot be empty")
	}

	ri, err := a.GenerateRI(context.Background(), project, fixVersion)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if err := json.NewEncoder(w).Encode(ri); err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}
