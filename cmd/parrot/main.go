package main

import (
	"log/slog"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/studio-b12/parrot/pkg/server"
)

type Args struct {
	BindAddress  string     `arg:"--bind-address,env:BIND_ADDRESS" default:"0.0.0.0:8080" help:"HTTP bind address"`
	NtfyUpstream string     `arg:"--ntfy-upstream,env:NTFY_UPSTREAM,required" help:"Address of the upstream NTFY server"`
	LogLevel     slog.Level `arg:"--log-level,env:LOG_LEVEL" default:"info" help:"Log level"`
}

func main() {
	var args Args
	arg.MustParse(&args)

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: args.LogLevel}))
	slog.SetDefault(logger)

	s, err := server.New(args.BindAddress, args.NtfyUpstream)
	if err != nil {
		slog.Error("failed initializing server", "err", err)
		os.Exit(1)
	}

	slog.Info("starting server ...", "bindAddress", args.BindAddress, "upstream", args.NtfyUpstream)
	err = s.ListenAndServe()
	if err != nil {
		slog.Error("failed starting server", "err", err)
		os.Exit(1)
	}
}
