/*
 * Copyright (C) 2018 The dad-go Authors
 * This file is part of The dad-go library.
 *
 * The dad-go is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The dad-go is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The dad-go.  If not, see <http://www.gnu.org/licenses/>.
 */
package cmd

import "fmt"

func showAssetHelp() {
	var assetHelp = `
   Name:
      dad-go asset                       asset operation

   Usage:
      dad-go asset [command options] [args]

   Description:
      With this command, you can control assert through transaction.

   Command:
      transfer
         --caddr     value                 smart contract address
         --from      value                 wallet address base58, which will transfer from
         --to        value                 wallet address base58, which will transfer to
         --value     value                 how much asset will be transfered
         --password  value                 use password who transfer from

      status
         --hash     value                  transfer transaction hash
`
	fmt.Println(assetHelp)
}

func showQueryAssetTransferHelp() {
	var queryAssetTransferHelp = `
   Name:
      dad-go asset query              asset transfer resule query

   Usage:
      dad-go asset query [command options] [args]

   Description:
      With this command, you can query transfer assert status.

   Command:
      --hash     value                    transfer transaction hash
`
	fmt.Println(queryAssetTransferHelp)
}

func showAssetTransferHelp() {
	var assetTransferHelp = `
   Name:
      dad-go asset transfer              asset transfer

   Usage:
      dad-go asset transfer [command options] [args]

   Description:
      With this command, you can transfer assert through transaction.

   Command:
      --caddr     value                    smart contract address
      --from      value                    wallet address base58, which will transfer from
      --to        value                    wallet address base58, which will transfer to
      --value     value                    how much asset will be transfered
      --password  value                    use password who transfer from
`
	fmt.Println(assetTransferHelp)
}

func showContractHelp() {
	var contractUsingHelp = `
   Name:
      dad-go contract      deploy or invoke a smart contract by this command
   Usage:
      dad-go contract [command options] [args]

   Description:
      With this command, you can invoke a smart contract

   Command:
     invoke
       --caddr      value               smart contract address that will be invoke
       --params     value               params will be  
			
     deploy
       --type       value               contract type ,value: 1 (NEOVM) | 2 (WASM)
       --store      value               does this contract will be stored, value: true or false
       --code       value               directory of smart contract that will be deployed
       --cname      value               contract name that will be deployed
       --cversion   value               contract version which will be deployed
       --author     value               owner of deployed smart contract
       --email      value               owner email who deploy the smart contract
       --desc       value               contract description when deploy one
`
	fmt.Println(contractUsingHelp)
}

func showDeployHelp() {
	var deployHelp = `
   Name:
      dad-go contract deploy        deploy a smart contract by this command
   Usage:
      dad-go contract deploy [command options] [args]

   Description:
      With this command, you can deploy a smart contract

   Command:
      --type       value              contract type ,value: 1 (NEOVM) | 2 (WASM)
      --store      value              does this contract will be stored, value: true or false
      --code       value              directory of smart contract that will be deployed
      --cname      value              contract name that will be deployed
      --cversion   value              contract version which will be deployed
      --author     value              owner of deployed smart contract
      --email      value              owner email who deploy the smart contract
      --desc       value              contract description when deploy one
`
	fmt.Println(deployHelp)
}
func showInvokeHelp() {
	var invokeHelp = `
   Name:
      dad-go contract invoke          invoke a smart contract by this command
   Usage:
      dad-go contract invoke [command options] [args]

   Description:
      With this command, you can invoke a smart contract

   Command:
      --caddr      value                smart contract address that will be invoke
      --params     value                params will be
`
	fmt.Println(invokeHelp)
}

func showInfoHelp() {
	var infoHelp = `
   Name:
      dad-go info                    Show blockchain information

   Usage:
      dad-go info [command options] [args]

   Description:
      With dad-go info, you can look up blocks, transactions, etc.

   Command:
      version

      block
         --hash value                  block hash value
         --height value                block height value

      tx
         --hash value                  transaction hash value

`
	fmt.Println(infoHelp)
}

func showVersionInfoHelp() {
	var versionInfoHelp = `
   Name:
      dad-go info version            Show dad-go node version

   Usage:
      dad-go info version

   Description:
      With this command, you can look up the dad-go node version.

`
	fmt.Println(versionInfoHelp)
}

func showBlockInfoHelp() {
	var blockInfoHelp = `
   Name:
      dad-go info block             Show blockchain information

   Usage:
      dad-go info block [command options] [args]

   Description:
      With this command, you can look up block information.

   Options:
      --hash value                    block hash value
      --height value                  block height value
`
	fmt.Println(blockInfoHelp)
}

func showTxInfoHelp() {
	var txInfoHelp = `
   Name:
      dad-go info tx               Show transaction information

   Usage:
      dad-go info tx [command options] [args]

   Description:
      With this command, you can look up transaction information.

   Options:
      --hash value                   transaction hash value

`
	fmt.Println(txInfoHelp)
}

func showSettingHelp() {
	var settingHelp = `
   Name:
      dad-go set                       Show blockchain information

   Usage:
      dad-go set [command options] [args]

   Description:
      With dad-go set, you can configure the node.

   Command:
      --debuglevel value                 debug level(0~6) will be set
      --consensus value                  [ on / off ]
`
	fmt.Println(settingHelp)
}

func showWalletHelp() {
	var walletHelp = `
   Name:
      dad-go wallet                  User wallet operation

   Usage:
      dad-go wallet [command options] [args]

   Description:
      With dad-go wallet, you could control your account.

   Command:
      create
      --name value                     wallet name
      show
      --name value                     wallet name (default: wallet.dat)
      balance
      --name value                     wallet name (default: wallet.dat)
`
	fmt.Println(walletHelp)
}
