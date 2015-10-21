// http server library.
package server

import (
	"net/http"
)

type HTTPServer struct {
	addr     string
	handler  http.Handler
	useHttps bool
}

func (hs *HTTPServer) Start() error {
	// field check
	if hs.handler == nil {
		return errorf("Start HTTP Server error : http handler not registered!")
	}

	Log.Infof("HTTP Server Listen on %s, use https: %v", hs.addr, hs.useHttps)
	go func() {
		if hs.useHttps == false {
			http.ListenAndServe(hs.addr, hs.handler)
		} else {
			// secure files
			if *confCAFile == "" {
				errorf(errMissingFlag, flagCAFile)
				return
			}
			if *confKeyFile == "" {
				errorf(errMissingFlag, flagKeyFile)
				return
			}
			err := http.ListenAndServeTLS(hs.addr, *confCAFile, *confKeyFile, hs.handler)
			if err != nil {
				Log.Fatal(err.Error())
			}
		}
	}()

	return nil
}
