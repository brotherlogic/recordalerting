package main

import "golang.org/x/net/context"
import rcpb "github.com/brotherlogic/recordcollection/proto"

//ClientUpdate process new record
func (s *Server) ClientUpdate(ctx context.Context, req *rcpb.ClientUpdateRequest) (*rcpb.ClientUpdateResponse, error) {
	r, err := s.rc.getRecord(ctx, req.GetInstanceId())
	if err != nil {
		return nil, err
	}

	return &rcpb.ClientUpdateResponse{}, s.assessRecord(r)
}
