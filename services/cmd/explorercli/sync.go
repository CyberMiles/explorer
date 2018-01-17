package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	_ "github.com/spf13/viper"


	_ "github.com/cosmos/cosmos-sdk/modules/auth"
	_ "github.com/cosmos/cosmos-sdk/modules/base"
	_ "github.com/cosmos/cosmos-sdk/modules/coin"
	_ "github.com/cosmos/cosmos-sdk/modules/nonce"
	_ "github.com/cosmos/cosmos-sdk/modules/fee"
	"github.com/cosmos/cosmos-sdk/client/commands"

	_ "github.com/cybermiles/explorer/services/modules/stake"
)

var (
	syncCmd = &cobra.Command{
		Use:   "sync",
		Long:  `sync`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmdSync(cmd, args)
		},
	}
)

const (
	BatchSize = 20
)

func cmdSync(cmd *cobra.Command, args []string) error {
//   fmt.Println(viper.GetString("current-block-height"))
// viper.Set("current-block-height",11)
// viper.WriteConfig()
//   fmt.Println(viper.GetString("current-block-height"))
	
  c := commands.GetNode()
  scanned_block_height := int64(500000)
  latest := int64(0)
  tx := []int64{}

  for ok := true; ok; ok = (scanned_block_height < latest) {
    end := scanned_block_height + BatchSize
    blocks, err := c.BlockchainInfo(0, end)
    if err != nil {
    	log.Fatal(err)
      return err
    }
    latest = blocks.LastHeight

    for _, block := range blocks.BlockMetas {
      current := block.Header.Height
      if (current <= scanned_block_height){
        break
      }
      if (block.Header.NumTxs > 0){
        tx = append([]int64{current}, tx...)
      }
    }
    scanned_block_height = blocks.BlockMetas[0].Header.Height
    fmt.Println(scanned_block_height)
  }
  fmt.Println(tx)
  return nil
}
