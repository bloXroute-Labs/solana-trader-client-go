# benchmark/txcompare

Simultaneously submits transactions to each of the provided endpoints, then determines which transactions were accepted
into blocks first and which were lost. Judges speed based on block number and position within block. Note that there is
no correction for RPC endpoint distance, so you may find it informative to ping each endpoint beforehand to determine
the fairness of this analysis

## Usage

Go:
```
$ PRIVATE_KEY=... go run ./benchmark/txcompare --iterations 4 --output result.csv --endpoint [rpc-endpoint-1] --endpoint [rpc-endpoint-2] --query-endpoint [rpc-endpoint-with-tx-indexing]
```

You can specify as many `--endpoint` arguments as you would like. `PRIVATE_KEY` must be provided to sign transactions. 
By default, this is a simple transaction with a [Memo Program](https://spl.solana.com/memo) instruction.

For `--query-endpoint` you need an RPC node with transaction
history enabled, as this script calls the `getBlock` and `getTransaction` endpoints. You can test this yourself:

```
$ curl https://api.mainnet-beta.solana.com -X POST -H "Content-Type: application/json" -d '
  {
    "jsonrpc": "2.0",
    "id": 1,
    "method": "getTransaction",
    "params": [
      "daMbnUbiMDFFGCwsLmLTrDRn2Jpx5YnFdnoxynR4ob2jax7M9TiAtnCpBTXbYdCfBVTg8FzpJU3hwcBJgDs8heB",
      "json"
    ]
  }
'
{"jsonrpc":"2.0","result":{...},"id":1}
```

A node without this indexing would look like this:

```
$ curl [endpoint] -X POST -H "Content-Type: application/json" -d '...'
{"jsonrpc":"2.0","error":{"code":-32011,"message":"Transaction history is not available from this node"},"id":1}
```

## Result

Logs should look this:

```
    2022-07-28T11:29:28.154-0500    DEBUG   transaction/submit.go:94        submitted transaction   {"signature": "2UiE9gyTdYH4nFWP1ir92WXqqkX7WSRSAfN6NyCYgVfR29jmsqs1K1ZkK8KAMwWYp4mMR5tj5sfGQN9UgQxiQkKe"}
2022-07-28T11:29:28.154-0500    DEBUG   transaction/submit.go:94        submitted transaction   {"signature": "5SPCcCYgdsiXSdQeTmibc265YeiiQeeGe6WbdGPDNjo1gjVEf6D9DSmBesV1py6ipPoRtMMUUeah7pK38q9e7w8Y"}
2022-07-28T11:29:28.154-0500    DEBUG   transaction/submit.go:65        submitted iteration of transactions     {"iteration": 0, "count": 2}
2022-07-28T11:29:30.251-0500    DEBUG   transaction/submit.go:94        submitted transaction   {"signature": "2Y4hVUejyNM9mymrseuFrvmmj1RFwpV59uxRXgn1LCB9VYnVjn2sFoGnjSmKA4vRxBT13ZkNe6h6SHhGd2jfgUnQ"}
2022-07-28T11:29:30.252-0500    DEBUG   transaction/submit.go:94        submitted transaction   {"signature": "2BLCkuL2U6iKTpGCxNZUCn8zyBYx8254kJibu39mtve43UstBkAQYZXYRjtiwMVybihDaKwq4MJUUy3K18YeqEUi"}
2022-07-28T11:29:30.252-0500    DEBUG   transaction/submit.go:65        submitted iteration of transactions     {"iteration": 1, "count": 2}
2022-07-28T11:29:32.346-0500    DEBUG   transaction/submit.go:94        submitted transaction   {"signature": "5eyRvD6KtQ7qLs28aafUS3DXZzhnpLEhzByQVJscYFd6QuUFakvEnKLUVmUVJfHD5aJvCTbjpqSaM2rYQjninje7"}
2022-07-28T11:29:32.347-0500    DEBUG   transaction/submit.go:94        submitted transaction   {"signature": "3PJdb5STnCnJ4aBBYgxhq1RaayTaF8P1Dvub8GnDPYC8kd7Dui94VBKqDuAznwy4LJLKBgwuWg3S5vNuPsjv9yHt"}
2022-07-28T11:29:32.347-0500    DEBUG   transaction/submit.go:65        submitted iteration of transactions     {"iteration": 2, "count": 2}
2022-07-28T11:29:34.439-0500    DEBUG   transaction/submit.go:94        submitted transaction   {"signature": "5ofkDQsaEm4PrpG1v9u8t7JhDcXx2pAPtZscPrGLUfUdCDymSGNCfyrc85tAzjhM6KkSRh8kiXucP79A8HZSzMYu"}
2022-07-28T11:29:34.440-0500    DEBUG   transaction/submit.go:94        submitted transaction   {"signature": "2gej52a1hEF715UyqDD5Q8CSuvQ8tBhY3twNoV7jCqoJXqxLtCiayVVZAMgW8iSF2DcA2b3TMztysEXnKnigiNFq"}
2022-07-28T11:29:34.440-0500    DEBUG   transaction/submit.go:65        submitted iteration of transactions     {"iteration": 3, "count": 2}
2022-07-28T11:30:47.107-0500    DEBUG   txcompare/main.go:76    iteration results found {"iteration": 0, "winner": "https://nd-223-967-158.p2pify.com/92b9f51421b09d9b68ce6e8cd8d50ebf"}
2022-07-28T11:30:47.107-0500    DEBUG   txcompare/main.go:98    iteration transaction result    {"iteration": 0, "endpoint": "https://api.mainnet-beta.solana.com", "slot": 143533719, "position": 894, "signature": "2UiE9gyTdYH4nFWP1ir92WXqqkX7WSRSAfN6NyCYgVfR29jmsqs1K1ZkK8KAMwWYp4mMR5tj5sfGQN9UgQxiQkKe"}
2022-07-28T11:30:47.107-0500    DEBUG   txcompare/main.go:98    iteration transaction result    {"iteration": 0, "endpoint": "https://nd-223-967-158.p2pify.com/92b9f51421b09d9b68ce6e8cd8d50ebf", "slot": 143533718, "position": 877, "signature": "5SPCcCYgdsiXSdQeTmibc265YeiiQeeGe6WbdGPDNjo1gjVEf6D9DSmBesV1py6ipPoRtMMUUeah7pK38q9e7w8Y"}
2022-07-28T11:30:47.151-0500    DEBUG   txcompare/main.go:76    iteration results found {"iteration": 1, "winner": "https://api.mainnet-beta.solana.com"}
2022-07-28T11:30:47.151-0500    DEBUG   txcompare/main.go:98    iteration transaction result    {"iteration": 1, "endpoint": "https://api.mainnet-beta.solana.com", "slot": 143533718, "position": 1057, "signature": "2Y4hVUejyNM9mymrseuFrvmmj1RFwpV59uxRXgn1LCB9VYnVjn2sFoGnjSmKA4vRxBT13ZkNe6h6SHhGd2jfgUnQ"}
2022-07-28T11:30:47.151-0500    DEBUG   txcompare/main.go:98    iteration transaction result    {"iteration": 1, "endpoint": "https://nd-223-967-158.p2pify.com/92b9f51421b09d9b68ce6e8cd8d50ebf", "slot": 143533718, "position": 1354, "signature": "2BLCkuL2U6iKTpGCxNZUCn8zyBYx8254kJibu39mtve43UstBkAQYZXYRjtiwMVybihDaKwq4MJUUy3K18YeqEUi"}
2022-07-28T11:30:47.392-0500    DEBUG   txcompare/main.go:76    iteration results found {"iteration": 2, "winner": "https://api.mainnet-beta.solana.com"}
2022-07-28T11:30:47.392-0500    DEBUG   txcompare/main.go:98    iteration transaction result    {"iteration": 2, "endpoint": "https://api.mainnet-beta.solana.com", "slot": 143533720, "position": 286, "signature": "5eyRvD6KtQ7qLs28aafUS3DXZzhnpLEhzByQVJscYFd6QuUFakvEnKLUVmUVJfHD5aJvCTbjpqSaM2rYQjninje7"}
2022-07-28T11:30:47.393-0500    DEBUG   txcompare/main.go:98    iteration transaction result    {"iteration": 2, "endpoint": "https://nd-223-967-158.p2pify.com/92b9f51421b09d9b68ce6e8cd8d50ebf", "slot": 143533720, "position": 569, "signature": "3PJdb5STnCnJ4aBBYgxhq1RaayTaF8P1Dvub8GnDPYC8kd7Dui94VBKqDuAznwy4LJLKBgwuWg3S5vNuPsjv9yHt"}
2022-07-28T11:30:57.717-0500    DEBUG   txcompare/main.go:76    iteration results found {"iteration": 3, "winner": "https://nd-223-967-158.p2pify.com/92b9f51421b09d9b68ce6e8cd8d50ebf"}
2022-07-28T11:30:57.718-0500    DEBUG   txcompare/main.go:98    iteration transaction result    {"iteration": 3, "endpoint": "https://api.mainnet-beta.solana.com", "slot": 143533759, "position": 127, "signature": "5ofkDQsaEm4PrpG1v9u8t7JhDcXx2pAPtZscPrGLUfUdCDymSGNCfyrc85tAzjhM6KkSRh8kiXucP79A8HZSzMYu"}
2022-07-28T11:30:57.718-0500    DEBUG   txcompare/main.go:98    iteration transaction result    {"iteration": 3, "endpoint": "https://nd-223-967-158.p2pify.com/92b9f51421b09d9b68ce6e8cd8d50ebf", "slot": 143533733, "position": 2085, "signature": "2gej52a1hEF715UyqDD5Q8CSuvQ8tBhY3twNoV7jCqoJXqxLtCiayVVZAMgW8iSF2DcA2b3TMztysEXnKnigiNFq"}
Iterations:  4
Endpoints:
     https://api.mainnet-beta.solana.com
     https://nd-223-967-158.p2pify.com/92b9f51421b09d9b68ce6e8cd8d50ebf

Win counts: 
    2    https://api.mainnet-beta.solana.com
    2    https://nd-223-967-158.p2pify.com/92b9f51421b09d9b68ce6e8cd8d50ebf

Lost transactions: 
    0    https://api.mainnet-beta.solana.com
    0    https://nd-223-967-158.p2pify.com/92b9f51421b09d9b68ce6e8cd8d50ebf
```

A CSV file will also be generated at the `--output` location, with details of each transaction.