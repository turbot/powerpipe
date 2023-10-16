package perr

import (
	"github.com/rs/xid"
)

func reference() string {
	return "fperr_" + xid.New().String()
}
