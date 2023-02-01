package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/brotherlogic/goserver"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	dsc "github.com/brotherlogic/dstore/client"

	dspb "github.com/brotherlogic/dstore/proto"
	gdpb "github.com/brotherlogic/godiscogs/proto"
	pbg "github.com/brotherlogic/goserver/proto"
	pb "github.com/brotherlogic/recordalerting/proto"
	rcpb "github.com/brotherlogic/recordcollection/proto"
	pbro "github.com/brotherlogic/recordsorganiser/proto"
	google_protobuf "github.com/golang/protobuf/ptypes/any"
)

const (
	CONFIG_KEY = "github.com/brotherlogic/recordalerting/config"
)

var (
	tracked = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "recordalerting_tracked_issues",
		Help: "The size of the print queue",
	})
)

type ro interface {
	getLocation(ctx context.Context, name string) (*pbro.Location, error)
}

type prodRO struct {
	dial func(ctx context.Context, server string) (*grpc.ClientConn, error)
}

func (gh *prodRO) getLocation(ctx context.Context, name string) (*pbro.Location, error) {
	conn, err := gh.dial(ctx, "recordsorganiser")
	if err != nil {
		return &pbro.Location{}, err
	}
	defer conn.Close()

	client := pbro.NewOrganiserServiceClient(conn)
	resp, err := client.GetOrganisation(ctx, &pbro.GetOrganisationRequest{Locations: []*pbro.Location{&pbro.Location{Name: name}}})

	if err != nil {
		return &pbro.Location{}, err
	}

	if len(resp.GetLocations()) != 1 {
		return &pbro.Location{}, fmt.Errorf("Too many locations returned: %v", len(resp.GetLocations()))
	}

	return resp.GetLocations()[0], nil
}

type rc interface {
	getRecord(ctx context.Context, instanceID int32) (*rcpb.Record, error)
	clean(ctx context.Context, instanceID int32) error
	getLibraryRecords(ctx context.Context) ([]*rcpb.Record, error)
	getRecordsInFolder(ctx context.Context, folder int32) ([]int32, error)
}

type prodRC struct {
	dial func(ctx context.Context, server string) (*grpc.ClientConn, error)
}

func (gh *prodRC) getRecord(ctx context.Context, i int32) (*rcpb.Record, error) {
	conn, err := gh.dial(ctx, "recordcollection")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := rcpb.NewRecordCollectionServiceClient(conn)
	recs, err := client.GetRecord(ctx, &rcpb.GetRecordRequest{InstanceId: i})

	if err != nil {
		return nil, err
	}

	return recs.GetRecord(), nil
}

func (gh *prodRC) clean(ctx context.Context, i int32) error {
	conn, err := gh.dial(ctx, "recordcollection")
	if err != nil {
		return err
	}
	defer conn.Close()

	client := rcpb.NewRecordCollectionServiceClient(conn)
	_, err = client.UpdateRecord(ctx, &rcpb.UpdateRecordRequest{Reason: "alert-clean", Update: &rcpb.Record{
		Release:  &gdpb.Release{InstanceId: i},
		Metadata: &rcpb.ReleaseMetadata{MoveFolder: 3386035},
	}})

	if err != nil {
		return err
	}

	return nil
}

func (s *Server) loadConfig(ctx context.Context) (*pb.Config, error) {
	res, err := s.dstoreClient.Read(ctx, &dspb.ReadRequest{Key: CONFIG_KEY})
	if err != nil {
		if status.Convert(err).Code() == codes.NotFound {
			return &pb.Config{}, nil
		}

		return nil, err

	}

	if res.GetConsensus() < 0.5 {
		return nil, fmt.Errorf("could not get read consensus (%v)", res.GetConsensus())
	}

	config := &pb.Config{}
	err = proto.Unmarshal(res.GetValue().GetValue(), config)
	if err != nil {
		return nil, err
	}

	tracked.Set(float64(len(config.GetProblems())))

	return config, nil
}

func (s *Server) saveConfig(ctx context.Context, config *pb.Config) error {
	data, err := proto.Marshal(config)
	if err != nil {
		return err
	}
	res, err := s.dstoreClient.Write(ctx, &dspb.WriteRequest{Key: CONFIG_KEY, Value: &google_protobuf.Any{Value: data}})
	if err != nil {
		return err
	}

	if res.GetConsensus() < 0.5 {
		return fmt.Errorf("could not get write consensus (%v)", res.GetConsensus())
	}

	tracked.Set(float64(len(config.GetProblems())))

	return nil
}

func (gh *prodRC) getLibraryRecords(ctx context.Context) ([]*rcpb.Record, error) {
	conn, err := gh.dial(ctx, "recordsorganiser")
	if err != nil {
		return []*rcpb.Record{}, err
	}
	defer conn.Close()

	client := pbro.NewOrganiserServiceClient(conn)
	resp, err := client.GetOrganisation(ctx, &pbro.GetOrganisationRequest{Locations: []*pbro.Location{&pbro.Location{Name: "Library Records"}}})

	if err != nil {
		return []*rcpb.Record{}, err
	}

	if len(resp.GetLocations()) != 1 {
		return []*rcpb.Record{}, fmt.Errorf("Too many locations returned: %v", len(resp.GetLocations()))
	}

	recs := make([]*rcpb.Record, 0)
	for _, loc := range resp.GetLocations()[0].GetReleasesLocation() {
		rec, err := gh.getRecord(ctx, loc.GetInstanceId())
		if err != nil {
			return []*rcpb.Record{}, err
		}
		recs = append(recs, rec)
	}

	return recs, nil
}

func (gh *prodRC) getRecordsInFolder(ctx context.Context, folder int32) ([]int32, error) {
	conn, err := gh.dial(ctx, "recordcollection")
	if err != nil {
		return []int32{}, err
	}
	defer conn.Close()

	client := rcpb.NewRecordCollectionServiceClient(conn)
	recs, err := client.QueryRecords(ctx, &rcpb.QueryRecordsRequest{Query: &rcpb.QueryRecordsRequest_FolderId{folder}})

	if err != nil {
		return []int32{}, err
	}

	return recs.GetInstanceIds(), nil
}

//Server main server type
type Server struct {
	*goserver.GoServer
	rc             rc
	ro             ro
	invalidRecords int
	alertCount     int
	dstoreClient   *dsc.DStoreClient
}

// Init builds the server
func Init() *Server {
	s := &Server{GoServer: &goserver.GoServer{}}
	s.rc = &prodRC{s.FDialServer}
	s.ro = &prodRO{s.FDialServer}
	s.dstoreClient = &dsc.DStoreClient{Gs: s.GoServer}

	return s
}

// DoRegister does RPC registration
func (s *Server) DoRegister(server *grpc.Server) {
	// Do nothing
	rcpb.RegisterClientUpdateServiceServer(server, s)
}

// ReportHealth alerts if we're not healthy
func (s *Server) ReportHealth() bool {
	return true
}

// Shutdown the server
func (s *Server) Shutdown(ctx context.Context) error {
	return nil
}

// Mote promotes/demotes this server
func (s *Server) Mote(ctx context.Context, master bool) error {
	return nil
}

// GetState gets the state of the server
func (s *Server) GetState() []*pbg.State {
	return []*pbg.State{
		&pbg.State{Key: "testv", Value: int64(12344)},
	}
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
	server.PrepServer("recordalerting")
	server.Register = server

	err := server.RegisterServerV2(false)
	if err != nil {
		return
	}

	fmt.Printf("%v", server.Serve())
}
