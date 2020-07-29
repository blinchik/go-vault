package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	secret "github.com/blinchik/go-aws/lib/secrets"
	sys "github.com/blinchik/go-vault/sys"
)

var consulVault string

const (
	vaultPort = "8200"
)

func init() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

}

func main() {

	init := flag.Bool("init", false, "init vault")
	enable := flag.Bool("enable", false, "enable secret engine")

	flag.Parse()

	if *init {

		vaultAddress := os.Args[2]
		vaultPort := os.Args[3]

		output := sys.VaultInit(vaultAddress, vaultPort, 5, 3)

		secret.CreateSecret("vault root token (should be not saved in future)", "vault_root_token", output.RootToken)

		for idx, key := range output.Recoverykeys {

			secret.CreateSecret("vault backup key", fmt.Sprintf("vault_backup_key_%d", idx), key)

		}

		if *enable {

			secretEngines := strings.Split(os.Args[2], ",")

			for _, se := range secretEngines {

				sys.EnableSecretEngine(consulVault, vaultPort, &output.RootToken, se)

			}

		}

	}

}
