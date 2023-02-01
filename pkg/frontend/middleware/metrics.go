package middleware

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/Azure/ARO-RP/pkg/api"
	"github.com/Azure/ARO-RP/pkg/metrics"
)

// Metric records request metrics for tracking
func Metrics(m metrics.Emitter) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiVersion := r.URL.Query().Get(api.APIVersionKey)
			t := time.Now()
			var routeName string
			if route := mux.CurrentRoute(r); route != nil {
				routeName = route.GetName()
			} else {
				routeName = "unknown"
			}

			w = &logResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			defer func() {
				m.EmitGauge("frontend.count", 1, map[string]string{
					"verb":        r.Method,
					"api-version": apiVersion,
					"code":        strconv.Itoa(w.(*logResponseWriter).statusCode),
					"route":       routeName,
				})

				m.EmitGauge("frontend.duration", time.Since(t).Milliseconds(), map[string]string{
					"verb":        r.Method,
					"api-version": apiVersion,
					"code":        strconv.Itoa(w.(*logResponseWriter).statusCode),
					"route":       routeName,
				})
			}()

			h.ServeHTTP(w, r)
		})
	}
}
