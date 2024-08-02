package main

import (
	"flag"
	"fmt"
	"log/slog"

	"github.com/NYTimes/gziphandler"
	"github.com/banansys/httpserver"
	"github.com/rs/cors"
)

func main() {
	slog.SetDefault(httpserver.DefaultLoggerProduction)
	var (
		mode     = flag.String("mode", "production", "Server mode (default: development) [development|production]")
		port     = flag.Int("port", 3000, "HTTP server port")
		serveDir = flag.String("root", "/var/www/html", "which directory to serve")
	)
	flag.Parse()
	addr := fmt.Sprintf(":%d", *port)

	serverMode := setLoggerAndServerMode(*mode)

	server := setupServer(serverMode, *serveDir, addr)

	slog.Info("", "www-root", *serveDir)

	if err := server.Run(); err != nil {
		slog.Error("server error", "err", err)
	}
}

func setLoggerAndServerMode(modeFlag string) httpserver.Mode {
	serverMode := httpserver.ModeDevelopment
	if modeFlag == "production" {
		slog.SetDefault(httpserver.DefaultLoggerProduction)
		serverMode = httpserver.ModeProduction
	} else {
		slog.SetDefault(httpserver.DefaultLoggerDevelopment)
	}

	return serverMode
}

func setupServer(serverMode httpserver.Mode, rootDir, listenAddr string) *httpserver.Server {
	mux := spaHandler(rootDir)
	// middlewares ..
	muxWithMiddlewares := requestLogging(slog.Default())(mux)
	muxWithMiddlewares = cors.Default().Handler(muxWithMiddlewares)
	muxWithMiddlewares = gziphandler.GzipHandler(muxWithMiddlewares)

	return httpserver.New(muxWithMiddlewares,
		httpserver.WithMode(serverMode),
		httpserver.WithLogger(slog.Default()),
		httpserver.ListenOn(listenAddr),
	)
}
