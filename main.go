package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/NYTimes/gziphandler"
	"github.com/banansys/httpserver"
	"github.com/rs/cors"
)

const (
	HealthCheckPath = "/_health"
	IndexFile       = "index.html"
)

func main() {
	defer handlePanic()

	slog.SetDefault(httpserver.DefaultLoggerProduction)
	var (
		mode = flag.String("mode", "production", "Server mode (default: development) [development|production]")
		port = flag.Int("port", 3000, "HTTP server port")
	)
	flag.Parse()
	addr := fmt.Sprintf(":%d", *port)

	serverMode := setLoggerAndServerMode(*mode)

	serveDir := ""
	if len(flag.Args()) == 1 {
		serveDir = flag.Args()[0]
	} else {
		serveDir = Must(os.Getwd())
	}

	serveDir = ResolvePath(serveDir)

	server := setupServer(serverMode, serveDir, addr)

	slog.Info("", "www-root", serveDir)

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
	mux := http.NewServeMux()

	mux.Handle(HealthCheckPath, healthcheckHandler(rootDir))
	mux.Handle("/", spaHandler(rootDir))
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

func handlePanic() {
	if err := recover(); err != nil {
		slog.Error("panic", "err", err)
		os.Exit(1)
	}
}

func Must[T any](res T, err error) T {
	if err != nil {
		panic(err)
	}

	return res
}

// ResolvePath resolves a path (absolute, relative, or home path with ~) to an absolute path,
// and also resolves any symbolic links along the way.
func ResolvePath(inputPath string) string {
	// Handle home directory path starting with `~`
	if strings.HasPrefix(inputPath, "~") {
		// Replace ~ with the user's home directory
		inputPath = filepath.Join(Must(user.Current()).HomeDir, inputPath[1:])
	}

	// Convert relative path to absolute path
	absPath := Must(filepath.Abs(inputPath))

	// Resolve any symbolic links in the path
	resolvedPath := Must(filepath.EvalSymlinks(absPath))

	return resolvedPath
}
