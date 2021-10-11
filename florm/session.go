package florm

import (
	"errors"
	"fmt"
	"strings"
)

var ErrDoubleModel = errors.New("the Model clause can only be used in a single branch update script, multiple model clause will be invalid")

type SessionOutput struct {
	stream FluxStream
	output interface{}
}

type SessionUpdate struct {
	stream FluxStream
	src    interface{}
	sel    []string
	vals   map[string]interface{}
}

type FluxSession struct {
	buckets []string
	vars    map[string][]string
	varseq  []string
	master  []string
	outputs []SessionOutput
	update  *SessionUpdate
	c       int
	mgr     APIManager
	model   interface{}

	dbg bool
}

func NewFluxSession() *FluxSession {
	return &FluxSession{
		vars: make(map[string][]string),
		mgr:  defaultAPIManager,
	}
}

func NewFluxSessionCustomAPI(mgr APIManager) *FluxSession {
	return &FluxSession{
		vars: make(map[string][]string),
		mgr:  mgr,
	}
}

func (f *FluxSession) Model(m InfluxModel) FluxStream {
	if f.model != nil {
		panic(ErrDoubleModel)
	}

	f.model = m
	return f.From(m.Bucket())
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

func (f *FluxSession) registerUpdate(s FluxStream, sel []string, src interface{}) {
	if f.update != nil {
		panic(errors.New("in fluxorm, update clause can only present once in the stream"))
	}

	f.update = &SessionUpdate{
		stream: s,
		src:    src,
		sel:    sel,
	}
}

func (f *FluxSession) registerUpdateValues(s FluxStream, values map[string]interface{}) {
	if f.update != nil {
		panic(errors.New("in fluxorm, update clause can only present once in the stream"))
	}

	f.update = &SessionUpdate{
		stream: s,
		vals:   values,
	}
}

func (f *FluxSession) Vars() (map[string][]string, []string) {
	defer func() {
		f.buckets = strDedup(f.buckets)
	}()

	for _, o := range f.outputs {
		initializeVars(f, o.stream)
		buildFluxVars(f, o.stream)

		if len(f.master) > 0 {
			f.commitMasterTo(o.stream.Name())
		}
	}

	if f.update != nil {
		o := f.update
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
		if s.GetOp() == OpFrom {
			_, bkt, _ := parseParam(sf.params[0])
			ss.buckets = append(ss.buckets, bkt)
		}
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
