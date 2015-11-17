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

	if hs.useHttps {
		// secure files
		if *confCAFile == "" {
			return errorf(errMissingFlag, FlagCAFile)
		}
		if *confKeyFile == "" {
			return errorf(errMissingFlag, FlagKeyFile)
		}
	}

	Log.Infof("HTTP Server Listen on %s, use https: %v", hs.addr, hs.useHttps)
	go func() {
		if hs.useHttps == false {
			http.ListenAndServe(hs.addr, hs.handler)
		} else {
			err := http.ListenAndServeTLS(hs.addr, *confCAFile, *confKeyFile, hs.handler)
			if err != nil {
				Log.Fatal(err.Error())
			}
		}
	}()

	return nil
}
