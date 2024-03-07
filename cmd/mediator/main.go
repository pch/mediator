package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/pch/mediator/internal"
)

var Version string

func main() {
	setLogger()
	handleVersionFlag()

	vips.LoggingSettings(func(msgDomain string, level vips.LogLevel, msg string) {
		slog.Debug("vips", "msgDomain", msgDomain, "level", level, "msg", msg)
	}, vips.LogLevelDebug)
	vips.Startup(nil)
	defer vips.Shutdown()

	config, err := internal.NewConfig()
	if err != nil {
		panic(err)
	}

	handler := internal.NewHandler(config)
	server := internal.NewHttpServer(config, handler)
	server.Start()
	defer server.Stop()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
}

func setLogger() {
	level := slog.LevelInfo

	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		level = slog.LevelDebug
	case "error":
		level = slog.LevelError
	case "warn":
		level = slog.LevelWarn
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})))
}

func handleVersionFlag() {
	versionFlag := flag.Bool("version", false, "Print the version of the app")
	flag.Parse()

	if *versionFlag {
		fmt.Println("Version:", Version)
		os.Exit(0)
	}
}
