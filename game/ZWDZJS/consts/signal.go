package consts

import "context"

var ClosingSignal, ClosingCommand = context.WithCancel(context.Background())
