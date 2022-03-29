package mocks

import (
	"github.com/stretchr/testify/mock"
)

type Printer struct {
	mock.Mock
}

func (p *Printer) Print(args ...interface{}) {
	p.Called(args)
}

func (p *Printer) Printf(format string, args ...interface{}) {
	p.Called(format, args)
}
