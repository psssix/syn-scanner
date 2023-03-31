package reporters

import (
	"github.com/stretchr/testify/mock"
)

type printerMock struct { //nolint:forbidigo // linter false-positive
	mock.Mock
}

func (p *printerMock) Print(args ...interface{}) { //nolint:forbidigo // linter false-positive
	p.Called(args)
}

func (p *printerMock) Printf(format string, args ...interface{}) { //nolint:forbidigo // linter false-positive
	p.Called(format, args)
}
