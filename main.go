package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Flags
var portFlag = flag.Int("port", 9009, "The port to run server on")
var backupDirFlag = flag.String("backupDir", "", "The directory for backups")
var saveBackupsFlag = flag.Bool("saveBackups", true, "Save backups to disk")

func main() {
	flag.Parse()

	logger := log.New(os.Stdout, "[main]: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile|log.Lmsgprefix)

	logger.Println("Flags:")
	flag.VisitAll(func(f *flag.Flag) {
		logger.Printf("Flag %s: %s [def:%s]\n", f.Name, f.Value, f.DefValue)
	})

	logger.Println("Starting up")

	// Initialize app
	serverLogger := log.New(os.Stdout, "[server]: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile|log.Lmsgprefix)
	app := &EtlApp{
		logger: serverLogger,
		engine: &Engine{},
	}
	app.engine.Init(app)

	// Start http server
	app.SetupHttpServer()
}

func helperFn(handler string, fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return promhttp.InstrumentHandlerDuration(
		promHttpLatency.MustCurryWith(prometheus.Labels{"handler": handler}),
		promhttp.InstrumentHandlerCounter(
			promHttpCalls.MustCurryWith(prometheus.Labels{"handler": handler}),
			promhttp.InstrumentHandlerResponseSize(
				promResponseSize,
				http.HandlerFunc(fn),
			),
		),
	)
}

type EtlApp struct {
	logger    *log.Logger
	engine    *Engine
	appId     string
	backupDir string
}

type SolveHandlerReq struct {
	Script []string `json:"script"`
}

type JupyterCmdHandlerReq struct {
	Cmd                  string `json:"cmd"`
	AppId                string `json:"appId"`
	ReturnLastCmdDetails bool   `json:"returnLastCmdDetails"`
}

type SolveHandlerRsp struct {
	Success  bool   `json:"success"`
	ErrorMsg string `json:"error"`
	Output   string `json:"output"`
}

func sendResponse(w http.ResponseWriter, code int, format string, a ...any) {
	w.Header().Set("X-SG-Server", "easyetl")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	io.WriteString(w, fmt.Sprintf(format, a...))
}

func (app *EtlApp) serializeSendHandlerRsp(success bool, errorMsg, output string) string {
	logger := app.logger
	rsp := SolveHandlerRsp{Success: success, ErrorMsg: errorMsg, Output: output}
	bytes, err := json.Marshal(rsp)
	if err != nil {
		logger.Printf("Error serializing rsp: %v %s\n", rsp, err)
		return fmt.Sprintf("{'success':false,'error':'Error serializing rsp: %s'}", err)
	}
	return string(bytes)
}

func (app *EtlApp) healthHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "ok")
}

func (app *EtlApp) jupyterCmdHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		payload := &JupyterCmdHandlerReq{}
		err := json.NewDecoder(r.Body).Decode(payload)
		app.logger.Printf("jupyterCmdHandler payload: %v %s\n", payload, err)
		if err != nil {
			sendResponse(w, http.StatusBadRequest, app.serializeSendHandlerRsp(false, fmt.Sprintf("Bad request: %s", err), ""))
			return
		}

		app2, found := savedApps[payload.AppId]
		if !found {
			serverLogger := log.New(os.Stdout, fmt.Sprintf("[%s]: ", payload.AppId), log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile|log.Lmsgprefix)
			app2 = &EtlApp{
				logger:    serverLogger,
				engine:    &Engine{},
				appId:     payload.AppId,
				backupDir: *backupDirFlag,
			}
			app2.engine.Init(app2)
			savedApps[payload.AppId] = app2
		}
		code, lastLhs, err := app2.solveScript(&SolveHandlerReq{Script: strings.Split(payload.Cmd, "\n")})
		if err != nil {
			sendResponse(w, http.StatusOK, app.serializeSendHandlerRsp(false, err.Error(), ""))
			return
		}

		if payload.ReturnLastCmdDetails {
			output, err := app2.engine.getTableOrDb(lastLhs, NUM_ROWS)
			if err != nil {
				sendResponse(w, http.StatusOK, app.serializeSendHandlerRsp(false, err.Error(), ""))
				return
			}
			sendResponse(w, http.StatusOK, app.serializeSendHandlerRsp(true, "", output))
			return
		}
		sendResponse(w, http.StatusOK, app.serializeSendHandlerRsp(true, "", fmt.Sprintf("code: %d", code)))
	}
}

func (app *EtlApp) solveHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		payload := &SolveHandlerReq{}
		err := json.NewDecoder(r.Body).Decode(payload)
		app.logger.Printf("solveHandler payload: %s\n", payload)
		if err != nil {
			sendResponse(w, http.StatusBadRequest, app.serializeSendHandlerRsp(false, fmt.Sprintf("Bad request: %s", err), ""))
			return
		}

		app.logger.Printf("script: %s\n", payload.Script)

		res, _, err := app.solveScript(payload)
		if err != nil {
			sendResponse(w, http.StatusBadRequest, app.serializeSendHandlerRsp(false, err.Error(), ""))
			return
		}
		sendResponse(w, http.StatusOK, app.serializeSendHandlerRsp(true, "", fmt.Sprintf("%d", res)))
	default:
		sendResponse(w, http.StatusBadRequest, app.serializeSendHandlerRsp(false, "unsupported", ""))
	}
}

func (app *EtlApp) solveScript(payload *SolveHandlerReq) (int, string, error) {
	commands := make([]string, 0)
	var cmd string
	for _, l := range payload.Script {
		line := strings.TrimSpace(l)
		if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "#") || line == "" {
			// Comments
			continue
		}

		if strings.HasPrefix(line, NEW_COMMAND_START) {
			// New command
			if cmd != "" {
				commands = append(commands, cmd)
			}
			cmd = strings.TrimPrefix(line, "-> ")
		} else {
			// Continuation of previous command
			cmd = cmd + "\n" + line
		}
	}
	if cmd != "" {
		commands = append(commands, cmd)
	}

	var lastLhs string
	for idx, line := range commands {
		startTime := time.Now()
		_code, _lhs, err := app.engine.EvaluateCmd(line)
		noop(_code)
		lastLhs = _lhs

		timeTakenMs := time.Since(startTime).Milliseconds()
		if timeTakenMs > 100 {
			app.logger.Printf("Executing command (%s) took %d ms\n", line, timeTakenMs)
		}

		if err != nil {
			out := fmt.Errorf("error executing line %d %s : %s", idx, line, err)
			app.logger.Println(out)
			return GENERAL_ERROR, lastLhs, out
		}

		// Last line
		if idx == len(payload.Script)-1 {
			return CODE_OK, lastLhs, nil
		}
	}
	return CODE_OK, lastLhs, nil
}

func (app *EtlApp) handleAll(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello!")
}

func (app *EtlApp) SetupHttpServer() {
	logger := app.logger
	mux := http.NewServeMux()

	// Register all of the metrics in the standard registry.
	prometheus.MustRegister(promHttpCalls, promHttpLatency, promResponseSize)

	mux.HandleFunc("/healthz", helperFn("healthz", app.healthHandler))
	mux.HandleFunc("/solve", helperFn("solve", app.solveHandler))
	mux.HandleFunc("/jupyterCmd", helperFn("jupyterCmd", app.jupyterCmdHandler))

	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/", helperFn("/", app.handleAll))

	logger.Println("Server starting")
	ctx, cancelCtx := context.WithCancel(context.Background())
	serverOne := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", *portFlag),
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, ctxKey, l.Addr().String())
			return ctx
		},
	}

	// Start the server
	go func() {
		err := serverOne.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			logger.Println("Error: server one closed")
		} else if err != nil {
			logger.Printf("Error listening for server one: %s\n", err)
		}
		cancelCtx()
	}()

	WaitForKillSignal(ctx, logger)
}

func WaitForKillSignal(ctx context.Context, logger *log.Logger) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	if ctx != nil {
		select {
		case <-ctx.Done():
			logger.Println("Context done, exiting")
			return
		case <-c:
			logger.Println("Interrupt received, exiting")
			return
		}
	} else {
		for range c {
			logger.Println("Interrupt received, exiting")
			return
		}
	}
}

const NEW_COMMAND_START = "-> "
const NUM_ROWS = 100

type key int

const ctxKey key = iota

var savedApps = make(map[string]*EtlApp)
