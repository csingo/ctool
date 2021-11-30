
func (s *##SERVICE##HttpClient) ##RPC##(ctx context.Context, req *##REQ##) (rsp *##RSP##, err error) {
	err = call(ctx, s.host, req, rsp)
	if err != nil {
		return
	}

	return
}
