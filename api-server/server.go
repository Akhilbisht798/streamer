package main

import (
	// "encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	listenAddr string
	Router     *mux.Router
}

func NewServer(addr string) *Server {
	return &Server{
		listenAddr: addr,
		Router:     mux.NewRouter(),
	}
}

func (s *Server) Run() {
	s.Router.HandleFunc("/", s.rootHandler)
	s.Router.HandleFunc("/run", s.runStreamer)
	s.Router.HandleFunc("/stop", s.stopStreamer)
	http.ListenAndServe(s.listenAddr, s.Router)
}

func (s *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}

type RunTaskResponseData struct {
	ARN    string `json:"arn"`
	Status string `json:"status"`
	IP     string `json:"ip"`
}

func (s *Server) runStreamer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("starting aws container")
	startTask()
	// resp := RunTaskResponseData{
	// 	ARN:    arn,
	// 	IP:     publicIp,
	// 	Status: "success",
	// }
	// jsonResp, err := json.Marshal(resp)
	// if err != nil {
	// 	http.Error(w, "Failed to create a json Response", http.StatusInternalServerError)
	// 	return
	// }
	// w.Header().Set("Content-Type", "application/json")
	// w.Write(jsonResp)
	w.Write([]byte("hello world"))
}

func (s *Server) stopStreamer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("stopping the container")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error Reading Body")
		return
	}
	stopTask(string(body))
	fmt.Println(body)
	w.Write([]byte("Successfully stopped the container"))
}
