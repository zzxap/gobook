package main

import (
	"fmt"
	"time"
)

type Server struct {
	*http.Server
	ln           net.Listener
	SignalHooks  map[int]map[os.Signal][]func()
	sigChan      chan os.Signal
	isChild      bool
	state        uint8
	Network      string
	terminalChan chan error
}

func main() {
	myHandler := mux.NewRouter()
	//myHandler := http.NewServeMux()
	myHandler.HandleFunc("/"+apiprefix+"/test", WelcomeTask)
	public.Log("start http server")
	errr := grace.ListenAndServe(":8080", myHandler)
	if errr != nil {
		public.Log("ListenAndServe  error: %v"+public.GetCurDateTime(), errr)
		//panic("http server stop exit at" + public.GetCurDateTime())
	} else {
		public.Log("ListenAndServe success")
	}
}

func (srv *Server) ListenAndServe() (err error) {

	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}
	//倾听热更新信号
	go srv.handleSignals()

	srv.ln, err = srv.getListener(addr)
	if err != nil {
		log.Println(err)
		return err
	}

	if srv.isChild {
		process, err := os.FindProcess(os.Getppid())
		if err != nil {
			log.Println(err)
			return err
		}
		err = process.Signal(syscall.SIGTERM)
		if err != nil {
			return err
		}
	}

	log.Println(os.Getpid(), srv.Addr)
	return srv.Serve()
}

func (srv *Server) handleSignals() {

	var sig os.Signal

	signal.Notify(
		srv.sigChan,
		hookableSignals...,
	)

	pid := syscall.Getpid()
	for {
		sig = <-srv.sigChan
		srv.signalHooks(PreSignal, sig)
		switch sig {
		case syscall.SIGHUP:
			log.Println(pid, "接收到热更新信号")
			//开启子进程
			err := srv.fork()
			if err != nil {
				log.Println("Fork err:", err)
			}
		case syscall.SIGINT:
			srv.shutdown()
		case syscall.SIGTERM:
			srv.shutdown()
		default:
			log.Printf("Received %v: nothing i care about...\n", sig)
		}
		srv.signalHooks(PostSignal, sig)
	}
}
