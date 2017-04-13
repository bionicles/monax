package chains

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/monax/cli/chains/maker"
	"github.com/monax/cli/definitions"
	"github.com/monax/cli/keys"
	"github.com/monax/cli/log"
	"github.com/monax/cli/version"
)

// MakeChain runs the chain maker
// It returns an error. Note that if do.Known, do.AccountTypes
// or do.ChainType are not set the command will run via interactive
// shell.
//
//  do.Name          - name of the chain to be created (required)
//  do.Known         - will parse csv's and create a genesis.json (requires do.ChainMakeVals and do.ChainMakeActs) (optional)
//  do.Output  	     - takes various options for the type of output wanted (tar/zip)
//  do.ChainMakeVals - csv file to use for validators (optional)
//  do.ChainMakeActs - csv file to use for accounts (optional)
//  do.AccountTypes  - use the account-types paradigm (example: Root:1,Participants:25,...) (optional)
//  do.ChainType     - use the chain-types paradigm (example: sprawlchain) (optional)
//  do.Verbose       - verbose output (optional)
//  do.Debug         - debug output (optional)
//

// remove!
//  do.Tarball       - instead of outputing raw files in directories, output packages of tarbals (optional)
//  do.ZipFile       - similar to do.Tarball except uses zipfiles (optional)
func MakeChain(do *definitions.Do) error {

	keys.InitKeyClient()

	// announce.
	log.Info("Hello! I'm the marmot who makes monax chains.")

	// set infos
	// do.Name; already set
	// do.Accounts ...?
	do.ChainImageName = path.Join(version.DefaultRegistry, version.ImageDB)
	do.ExportedPorts = []string{"1337", "46656", "46657"}
	do.UseDataContainer = true
	do.ContainerEntrypoint = ""

	// make it
	if err := maker.MakeChain(do); err != nil {
		return err
	}

	// cm currently is not opinionated about its writers.
	switch do.Output {
	case "tar":
		if err := maker.Tarball(do); err != nil {
			return err
		}
	case "zip":
		if err := maker.Zip(do); err != nil {
			return err
		}
	case "kubernetes":
		return fmt.Errorf("Not yet implemented, see issue #1272")
	default:
		return fmt.Errorf("Output must be one of [tar,zip,kubernetes]")
	}

	// put at end so users see it after any verbose/debug logs
	if len(do.AccountTypes) > 0 {
		numberOfValidators, err := checkNumberValidators(do.AccountTypes)
		if err != nil {
			return err
		}
		if numberOfValidators == 0 {
			log.Warn("WARNING: The chain made did not contain account types (Full/Validator) with validator permissions and will require further modification to run. The marmots recommend making a chain with Full/Validator account types")
		}
	}

	return nil
}

func checkNumberValidators(accountTypes []string) (int, error) {
	var num int = 0
	var err error
	for _, accT := range accountTypes {
		accounts := strings.Split(accT, ":")
		if accounts[0] == "Full" || accounts[0] == "Validator" {
			num, err = strconv.Atoi(accounts[1])
			if err != nil {
				return -1, err
			}
			num += num
		}
	}
	return num, nil
}
