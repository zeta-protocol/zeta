# Chain Replay

This is a guide to replay a chain using the backup of an existing chain (e.g. Testnet)

## How it works

A Tendermint Core and Zeta Core node store their configuration and data to disk by default at `$HOME/.cometbft` and `$HOME/.zeta`. When you start new instances of those nodes using a copy of these directories as their home, Tendermint re-submits (replays) historical blocks/transactions from the genesis height to Zeta Core.

## Prerequisites

- [Google Cloud SDK ][gcloud]
- Zeta Core Node
- Zeta Wallet
- [Tendermint][tendermint]

## Chain backups

Note you need to first authenticate `gcloud`.

You can find backups for the Zeta networks stored in Google Cloud Storage, e.g. For Testnet Node 01

```
$ gsutil ls gs://zeta-chaindata-n01-testnet/chain_stores
```

## Steps

- Copy backups locally to `<path>`

- Overwrite Zeta node wallet with your own development [node wallet][wallet]. 

```
$ cp -rp ~/.zeta/node_wallets_dev <path>/.zeta
$ cp ~/.zeta/nodewalletstore <path>/.zeta
```

- Update Zeta node configuration

```
$ sed -i 's/\/home\/zeta/<path>' <path>/.zeta/config.toml
```

- Start Zeta and Tendermint using backups

```
$ zeta node --root-path=<path>/.zeta --stores-enabled=false
$ tendermint node --home=<path>/.tendermint
```


## Tips

The Zeta nodes adheres to the Tendermint ABCI contract, therefore breakpoints in the following methods are useful:

```
blockchain/abci/abci.go#BeginBlock
```

## Alternatives

Instead of a backup, which effectively replays the full chain from genesis, you can also use a snapshot of the chain at a given height to bootstrap the Tendermint node. Which only replays blocks/transactions from the given height. This however requires extra tooling.

## References

- https://github.com/tendermint/tendermint/blob/master/docs/introduction/quick-start.md
- https://docs.tendermint.com/master/spec/abci/apps.html
- https://github.com/tendermint/spec/blob/master/spec/abci/README.md
- https://docs.tendermint.com/master/spec/abci/apps.html#state-sync

[wallet]: https://github.com/zetaprotocol/zeta#configuration
[gcloud]: https://cloud.google.com/sdk/docs/install
[tendermint]: https://github.com/tendermint/tendermint/blob/master/docs/introduction/install.md