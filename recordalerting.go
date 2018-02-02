package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/brotherlogic/goserver"
	"google.golang.org/grpc"

	pbgh "github.com/brotherlogic/githubcard/proto"
	pbg "github.com/brotherlogic/goserver/proto"
	"github.com/brotherlogic/goserver/utils"
	pbrc "github.com/brotherlogic/recordcollection/proto"
)

type gh interface {
	alert(r *pbrc.Record, text string) error
}

type prodGh struct{}

func (gh *prodGh) recordMissing(r *pbrc.Record, text string) error {
	host, port, err := utils.Resolve("githubcard")

	if err != nil {
		return err
	}

	conn, err := grpc.Dial(host+":"+strconv.Itoa(int(port)), grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	client := pbgh.NewGithubClient(conn)
	_, err = client.AddIssue(ctx, &pbgh.Issue{Title: "Problematic Record", Body: fmt.Sprintf("%v - %v", text, r.GetRelease().Title), Service: "recordcollection"})
	return err
}

//Server main server type
type Server struct {
	*goserver.GoServer
}

// Init builds the server
func Init() *Server {
	s := &Server{GoServer: &goserver.GoServer{}}
	return s
}

// DoRegister does RPC registration
func (s *Server) DoRegister(server *grpc.Server) {
	// Do nothing
}

// ReportHealth alerts if we're not healthy
func (s *Server) ReportHealth() bool {
	return true
}

// Mote promotes/demotes this server
func (s *Server) Mote(master bool) error {
	return nil
}

// GetState gets the state of the server
func (s *Server) GetState() []*pbg.State {
	return []*pbg.State{}
}

func main() {
	var quiet = flag.Bool("quiet", false, "Show all output")
	flag.Parse()

	//Turn off logging
	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	server := Init()
	server.PrepServer()
	server.Register = server

	server.RegisterServer("recordalerting", false)
	server.Log("Starting!")
	server.Serve()
}
