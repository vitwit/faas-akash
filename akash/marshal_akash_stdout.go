package akash

import (
	"encoding/json"
	"fmt"
	akashTypes "github.com/vitwit/faas-akash/types"
	"strings"
)

func marshalAkashStdout(stdout string) (*akashTypes.AkashStdout, error) {
	// akash stdout is not single json object
	// this is a workaround for capturing Leases Object from entire stdout for akash deployment
	out := strings.Split(stdout, "}{")
	if len(out) != 2 {
		return nil, fmt.Errorf("%s", "invalid stdout captured")
	}

	var leaseOut akashTypes.AkashStdout
	// split on }{ would remove opening json curly-brace, so we're adding it here manually
	// out[1] is assumed to be the leases object
	validJson := []byte("{" + out[1])
	if err := json.Unmarshal(validJson, &leaseOut); err != nil {
		return nil, err
	}

	return &leaseOut, nil
}
