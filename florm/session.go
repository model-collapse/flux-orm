package florm

import (
	"fmt"
	"strings"

	"github.com/influxdata/influxdb-client-go/v2/api"
)

type SessionOutput struct {
	stream FluxStream
	output interface{}
}

type SessionUpdate struct {
	stream FluxStream
	src    interface{}
}

type FluxSession struct {
	vars     map[string][]string
	varseq   []string
	master   []string
	outputs  []SessionOutput
	updates  []SessionUpdate
	c        int
	writeAPI api.WriteAPI
	queryAPI api.QueryAPI
}

func NewFluxSession() *FluxSession {
	return &FluxSession{
		vars:     make(map[string][]string),
		writeAPI: defaultWriteAPI,
		queryAPI: defaultQueryAPI,
	}
}

func NewFluxSessionCustomAPI(w api.WriteAPI, q api.QueryAPI) *FluxSession {
	return &FluxSession{
		vars:     make(map[string][]string),
		writeAPI: w,
		queryAPI: q,
	}
}

func (f *FluxSession) popVarName() string {
	defer func() {
		f.c++
	}()

	return fmt.Sprintf("v%d", f.c)
}

func (f *FluxSession) registerVarName(n string) {
	f.varseq = append(f.varseq, n)
}

func (f *FluxSession) commitMasterTo(vn string) {
	if _, suc := f.vars[vn]; !suc {
		f.vars[vn] = f.master
	}

	f.master = nil
}

func (f *FluxSession) addToMaster(s string) {
	f.master = append(f.master, s)
}

func (f *FluxSession) registerOutput(s FluxStream, dst interface{}) {
	// TODO add check

	f.outputs = append(f.outputs, SessionOutput{stream: s, output: dst})
}

func (f *FluxSession) registerUpdate(s FluxStream, src interface{}) {
	// TODO add check

	f.updates = append(f.updates, SessionUpdate{stream: s, src: src})
}

func (f *FluxSession) Vars() (map[string][]string, []string) {
	for _, o := range f.outputs {
		initializeVars(f, o.stream)
		buildFluxVars(f, o.stream)

		if len(f.master) > 0 {
			f.commitMasterTo(o.stream.Name())
		}
	}

	for _, o := range f.updates {
		initializeVars(f, o.stream)
		buildFluxVars(f, o.stream)

		if len(f.master) > 0 {
			f.commitMasterTo(o.stream.Name())
		}
	}

	return f.vars, f.varseq
}

func (f *FluxSession) QueryString() string {
	vars, sq := f.Vars()
	return varsToFluxLang(vars, sq)
}

func initializeVars(ss *FluxSession, s FluxStream) {
	multi := false
	splitted := false

	if _, suc := ss.vars[s.Name()]; suc {
		return
	}

	if sf, suc := s.(*FluxSinglePipe); suc {
		if sf.in == nil {
			s.AssignName(ss.popVarName())
			ss.registerVarName(s.Name())
		} else {
			initializeVars(ss, sf.in)

			if sf.in.Ref() > 1 {
				splitted = true
			}

			s.AssignName(sf.in.Name())
		}
	} else if sff, suc := s.(*FluxMultiplePipe); suc {
		for _, i := range sff.ins {
			initializeVars(ss, i)
		}

		multi = true
	}

	if splitted || multi {
		s.AssignName(ss.popVarName())
		ss.registerVarName(s.Name())
	}

	return
}

func buildFluxVars(ss *FluxSession, s FluxStream) {
	if _, suc := ss.vars[s.Name()]; suc {
		return
	}

	if sf, suc := s.(*FluxSinglePipe); suc {
		if sf.in != nil {
			buildFluxVars(ss, sf.in)

			if sf.in.Ref() > 1 {
				ss.commitMasterTo(sf.in.Name())
				ss.addToMaster(sf.in.Name())
			}
		}

		ss.addToMaster(s.Statement())

	} else if sff, suc := s.(*FluxMultiplePipe); suc {
		for _, ii := range sff.ins {
			if ii != nil {
				buildFluxVars(ss, ii)
				ss.commitMasterTo(ii.Name())
			}
		}

		ss.addToMaster(s.Statement())
	}
}

func varsToFluxLang(vars map[string][]string, varseq []string) string {
	var blocks []string
	for _, k := range varseq {
		v := vars[k]
		block := fmt.Sprintf("%s = %s", k, strings.Join(v, "\n|> "))
		blocks = append(blocks, block)
	}

	return strings.Join(blocks, "\n\n")
}
