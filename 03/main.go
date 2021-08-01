package main

import (
	"context"
	_ "errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	group, errCtx := errgroup.WithContext(ctx)
	srv := &http.Server{Addr: ":8090"}

	group.Go(func() error {
		return StartHttpServer(srv)
	})

	group.Go(func() error {
		<-errCtx.Done()
		fmt.Println("http server stop")
		return srv.Shutdown(errCtx)
	})
	chanel := make(chan os.Signal, 1)
	signal.Notify(chanel)

	group.Go(func() error {
		for {
			select {
			case <-errCtx.Done():
				return errCtx.Err()
			case <-chanel:
				cancel()
			}
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		fmt.Println("group error: ", err)
	}
	fmt.Println("OK")

}

func StartHttpServer(srv *http.Server) error {
	http.HandleFunc("/hello", HelloServer2)
	fmt.Println("http server start")
	err := srv.ListenAndServe()
	return err
}

func HelloServer2(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "let's go!\n")
}
