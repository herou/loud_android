package loudfixture

import (
	"flag"
	"testing"

	fixturetestSDK "github.com/Pylons-tech/pylons_sdk/cmd/fixtures_test"
	inttestSDK "github.com/Pylons-tech/pylons_sdk/cmd/test"
)

var runSerialMode bool = false
var connectLocalDaemon bool = false
var useKnownCookbook bool = false

func init() {
	flag.BoolVar(&runSerialMode, "runserial", false, "true/false value to check if test will be running in parallel")
	flag.BoolVar(&connectLocalDaemon, "locald", false, "true/false value to check if test will be connecting to local daemon")
	flag.BoolVar(&useKnownCookbook, "use-known-cookbook", false, "use existing cookbook or not")
}

func TestFixturesViaCLI(t *testing.T) {
	flag.Parse()
	if connectLocalDaemon {
		inttestSDK.CLIOpts.CustomNode = "tcp://localhost:26657"
	} else {
		inttestSDK.CLIOpts.CustomNode = "tcp://35.223.7.2:26657"
	}
	fixturetestSDK.FixtureTestOpts.CreateNewCookbook = !useKnownCookbook
	fixturetestSDK.FixtureTestOpts.IsParallel = !runSerialMode
	fixturetestSDK.RegisterDefaultActionRunners()
	fixturetestSDK.RunTestScenarios("scenarios", t)
}
