package main

import (
	"github.com/KhaiHust/email-notification-service/internal/bootstrap"
	"go.uber.org/fx"
)

func main() {
	fx.New(bootstrap.All()).Run()
}
