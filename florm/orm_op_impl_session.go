package florm

func (ss *FluxSession) From(bucket string) FluxStream {
	return From(bucket, ss)
}

func (ss *FluxSession) FromC(bucket string, host string, org string, token string) FluxStream {
	return FromC(bucket, host, org, token, ss)
}

func (ss *FluxSession) Buckets() FluxStream {
	return Buckets(ss)
}
