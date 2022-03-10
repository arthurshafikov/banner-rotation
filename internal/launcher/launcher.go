package launcher

import (
	"github.com/thewolf27/banner-rotation/internal/server"
)

func Launch() {
	s := server.NewServer()

	s.Serve()
}
