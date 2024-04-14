package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"

	"github.com/benderr/keypass/internal/client/app"
	"github.com/benderr/keypass/internal/client/logic"
	"github.com/benderr/keypass/internal/client/queries"
	"github.com/benderr/keypass/internal/client/repository"
	"github.com/benderr/keypass/pkg/logger"
	"github.com/benderr/keypass/pkg/sender"
)

//go:embed localhost.crt
var public []byte

func main() {
	log, close := logger.New()
	defer close()
	home, _ := os.UserHomeDir()
	path := fmt.Sprintf("%s%s%s", home, string(os.PathSeparator), "keypass")

	httpClient := sender.New("https://localhost:8080", string(public), log)
	qClient := queries.New(httpClient)
	repo := repository.New(path, log)
	appLogin := logic.New(qClient, repo)
	a := app.New(appLogin)

	ctx := context.Background()

	if err := a.Run(ctx); err != nil {
		fmt.Println(err)
		panic(err)
	}
}
