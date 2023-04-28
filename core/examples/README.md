# Examples

This folder contains examples use-cases which can be built into clients to a local zeta system.

## Nullchain

This is a go-package with a high-level of abstraction that gives the ability to manipulate a Zeta node running with the `null-blockchain`. The idea is that with this package trading-scenarios can be written from a point of view that doesn't require a lot of technical knowledge of how zeta works. For example, the below will connect to zeta, create 3 parties, and then move the blockchain forward:

```

func main() {
  w := nullchain.NewWallet(config.WalletFolder, config.Passphrase)
  conn, _ := nullchain.NewConnection()
  parties, err := w.MakeParties(3)

  nullchain.MoveByDuration(config.BlockDuration)

}
```

Prerequistes:
- A Zeta node in nullchain mode up and running
- Data-node up and running
- The Faucet up and running
- At least 3 users created in a local zeta wallet
- The details in `nullchain/config/config.go` updated to reflect your local environment

The following bash should get you some way there:
```
git clone git@github.com:zetaprotocol/zeta.git
git clone git@github.com:zetaprotocol/data-node.git

# cd into zeta data-node directories and run
go install ./...

# initialise zeta
zeta init -f --home=zetahome
zeta nodewallet generate --chain zeta --home=zetahome
zeta nodewallet generate --chain ethereum --home=zetahome

# initialise the faucet
zeta faucet init -f --home=/zetahome --update-in-place

# initialise TM just so we can auto-generate a genesis file to fill in
zeta tm init --home=zetahome
zeta genesis generate --home=zetahome
zeta genesis update --tm-home=/tenderminthome --home=/zetahome

# initialise the data-node
data-node init -f --home=zetahome

# initialise a zeta wallet and make some parties
zeta wallet init -f --home=zetahome
zeta wallet key generate --wallet=A --home=zetahome
zeta wallet key generate --wallet=B --home=zetahome
zeta wallet key generate --wallet=C --home=zetahome
```

Next you need to fiddle with the zeta config file to switch the blockchain on by changing the `BlockChain` section in `zetahome/config/node/config.toml` to look like this:
```
[Blockchain]
  Level = "Info"
  LogTimeDebug = true
  LogOrderSubmitDebug = true
  LogOrderAmendDebug = false
  LogOrderCancelDebug = false
  ChainProvider = "nullchain"
  [Blockchain.Tendermint]
    Level = "Info"
    LogTimeDebug = true
    ClientAddr = "tcp://0.0.0.0:26657"
    ClientEndpoint = "/websocket"
    ServerPort = 26658
    ServerAddr = "localhost"
    ABCIRecordDir = ""
    ABCIReplayFile = ""
  [Blockchain.Null]
    Level = "Debug"
    BlockDuration = "1s"
    TransactionsPerBlock = 3
    GenesisFile = "zetahome/config/genesis.json"
    IP = "0.0.0.0"
    Port = 3101
```

Now update the genesis file in `zetahome/config/genesis.json` to include the following assets in the appstate:

```
"assets": {
      "VOTE": {
        "name": "VOTE",
        "symbol": "VOTE",
        "total_supply": "0",
        "decimals": 5,
        "min_lp_stake": "1",
        "source": {
          "builtin_asset": {
            "max_faucet_amount_mint": "10000000000"
          }
        }
      },
      "XYZ": {
        "name": "XYZ",
        "symbol": "XYZ",
        "total_supply": "0",
        "decimals": 5,
        "min_lp_stake": "1",
        "source": {
          "builtin_asset": {
            "max_faucet_amount_mint": "10000000000"
          }
        }
      }
    },
```

Now spin up all the services by running the following each in their own terminal:

```
zeta node run --home=zetahome
data-node run --home=zetahome
zeta faucet run --home=zetahome

```

Once all is running, the example app can be run be doing the following in the `zeta` directory
```
go run ./cmd/examples/nullchain/nullchain
```
