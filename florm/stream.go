package florm

import (
	"errors"
	"fmt"
	"strings"
)

// FluxStream interface
type FluxStream interface {
	GetOp() FluxOp
	Name() string
	Session() *FluxSession
	AssignName(n string) FluxStream
	Statement() string
	Yield(dst interface{}) FluxStream
	Update(sel []string, rc interface{}) FluxStream
	AddRef()
	Ref() int

	FluxOperatable
}

type FluxSinglePipe struct {
	in      FluxStream
	session *FluxSession
	name    string
	op      FluxOp
	params  []string
	ref     int
}

func (f *FluxSinglePipe) GetOp() FluxOp {
	return f.op
}

func (f *FluxSinglePipe) Yield(dst interface{}) FluxStream {
	if err := checkYieldReceiver(dst); err != nil {
		panic(err)
	}

	f.AddRef()
	f.session.registerOutput(f, dst)
	return f
}

func (f *FluxSinglePipe) Update(sel []string, src interface{}) FluxStream {
	if !isInfluxModel(src) {
		panic(errors.New("the update model should be a InfluxModel"))
	}

	f.AddRef()
	f.session.registerUpdate(f, sel, src)
	return f
}

func (f *FluxSinglePipe) Session() *FluxSession {
	return f.session
}

func (f *FluxSinglePipe) Name() string {
	return f.name
}

func (f *FluxSinglePipe) AddRef() {
	f.ref++
}

func (f *FluxSinglePipe) Ref() int {
	return f.ref
}

func (f *FluxSinglePipe) AssignName(n string) FluxStream {
	ret := f
	ret.name = n

	return ret
}

func (f *FluxSinglePipe) Statement() string {
	return fmt.Sprintf("%s(%s)", opToStr[f.op], strings.Join(f.params, ", "))
}

type FluxMultiplePipe struct {
	ins     []FluxStream
	name    string
	session *FluxSession
	op      FluxOp
	params  []string
	ref     int
}

func (f *FluxMultiplePipe) GetOp() FluxOp {
	return f.op
}

func (f *FluxMultiplePipe) Yield(dst interface{}) FluxStream {
	if err := checkYieldReceiver(dst); err != nil {
		panic(err)
	}

	f.AddRef()
	f.session.registerOutput(f, dst)
	return f
}

func (f *FluxMultiplePipe) Update(sel []string, src interface{}) FluxStream {
	f.AddRef()
	f.session.registerUpdate(f, sel, src)
	return f
}

func (f *FluxMultiplePipe) Name() string {
	return f.name
}

func (f *FluxMultiplePipe) AssignName(n string) FluxStream {
	ret := f
	ret.name = n

	return ret
}

func (f *FluxMultiplePipe) Session() *FluxSession {
	return f.session
}

func (f *FluxMultiplePipe) AddRef() {
	f.ref++
}

func (f *FluxMultiplePipe) Ref() int {
	return f.ref
}

// Do not try this without a var name
func (f *FluxMultiplePipe) Statement() string {
	var ps []string
	for _, i := range f.ins {
		ps = append(ps, i.Name())
	}

	ps = append(ps, f.params...)

	return fmt.Sprintf("%s(%s)", opToStr[f.op], strings.Join(ps, ", "))
}
