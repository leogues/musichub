package musichub

import "context"

var ReportError = func(ctx context.Context, err error, args ...interface{}) {}

var ReportPanic = func(err interface{}) {}
