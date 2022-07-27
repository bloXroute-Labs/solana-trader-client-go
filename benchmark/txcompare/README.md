# benchmark/txcompare

Simultaneously submits transactions to each of the provided endpoints, then determines which transactions were accepted
into blocks first and which were lost. Judges speed based on block number and position within block.

## Usage

Go:
```
$ PRIVATE_KEY=... go run ./benchmark/txcompare --iterations 4 --output result.csv --endpoint [rpc-endpoint-1] --endpoint [rpc-endpoint-2] --query-endpoint [rpc-endpoint-with-tx-indexing]
```

You can specify as many `--endpoint` arguments as you would like. For `--query-endpoint` you need an RPC node with transaction
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
2022-07-27T14:41:39.399-0500    DEBUG   transaction/submit.go:94        submitted transaction   {"signature": "24DHDPwTaCz3UJsdkXM2of1oFVoS15KUXpNn8upfWkuTVa7L2SiS11NxvmAEQ9Wpor73CxpKW3GcEaG5ySVqxQGH"}
2022-07-27T14:41:39.400-0500    DEBUG   transaction/submit.go:94        submitted transaction   {"signature": "4efHFGHSJ9qQbtf35L4uZdf6Nt6nHnFwkFkbPURnjubGWvMi6sbXqxbSot523Ty7FxXgTqrNSg8e5kkff8iSxqMu"}
2022-07-27T14:41:39.400-0500    DEBUG   transaction/submit.go:65        submitted iteration of transactions     {"iteration": 0, "count": 2}
2022-07-27T14:41:41.443-0500    DEBUG   transaction/submit.go:94        submitted transaction   {"signature": "3DhCLYeENXAuMnwrTyew7hNprofJShSGB4GGGfpPx58YzVswZwEX7wC1EGjezLrNfdBxY8pXsdFAzVvKMC97Wt6J"}
2022-07-27T14:41:41.443-0500    DEBUG   transaction/submit.go:94        submitted transaction   {"signature": "5weRaUwR3m7NiA3s8kdzWjJQinBT58XVr3BcjBn9QmE5EFXGuteKq1d3Ms3mD6eHc7jTo7aG7ybAye57tskTfgEr"}
2022-07-27T14:41:41.443-0500    DEBUG   transaction/submit.go:65        submitted iteration of transactions     {"iteration": 1, "count": 2}
2022-07-27T14:41:43.487-0500    DEBUG   transaction/submit.go:94        submitted transaction   {"signature": "2M6KGkLoof371zNmZfnYeNnfyDh1beePzbRnZGpXhugM1zqgjouSZiVVZDhomej4HDV1hPhNwVLk1MDKHNn1nig5"}
2022-07-27T14:41:43.487-0500    DEBUG   transaction/submit.go:94        submitted transaction   {"signature": "5odrbxxoPFsUHsnssm7W26YDgvhM3cfgToxFj9dTDBiMtEyLhLdVC3u17h3CNWB9KNbxR6X6KHsWwEhCDFM4Leaz"}
2022-07-27T14:41:43.487-0500    DEBUG   transaction/submit.go:65        submitted iteration of transactions     {"iteration": 2, "count": 2}
2022-07-27T14:41:45.531-0500    DEBUG   transaction/submit.go:94        submitted transaction   {"signature": "2a9Rs95ikWwM7PtXQLxk5fWadKcKemc1dJyrRs2Q2L95ep61zGv8gKqaCZ2tC5FoEsrQMdpMpuxaJxQ57GW8ZCF8"}
2022-07-27T14:41:45.531-0500    DEBUG   transaction/submit.go:94        submitted transaction   {"signature": "5P1nLzay7DMnFThKobhz18G5P6Gbjnoai5ia2Fk8xeurajgLjvZnB2uiQF5u2WXkxjk6xTUGYF7gtCKMCn6v3MD8"}
2022-07-27T14:41:45.531-0500    DEBUG   transaction/submit.go:65        submitted iteration of transactions     {"iteration": 3, "count": 2}
2022-07-27T14:43:28.051-0500    DEBUG   transaction/status.go:122       transaction failed execution    {"signature": "24DHDPwTaCz3UJsdkXM2of1oFVoS15KUXpNn8upfWkuTVa7L2SiS11NxvmAEQ9Wpor73CxpKW3GcEaG5ySVqxQGH"}
2022-07-27T14:43:28.055-0500    DEBUG   transaction/status.go:122       transaction failed execution    {"signature": "4efHFGHSJ9qQbtf35L4uZdf6Nt6nHnFwkFkbPURnjubGWvMi6sbXqxbSot523Ty7FxXgTqrNSg8e5kkff8iSxqMu"}
2022-07-27T14:43:28.055-0500    DEBUG   txcompare/main.go:85    iteration no transactions confirmed     {"iteration": 0}
2022-07-27T14:43:28.055-0500    DEBUG   txcompare/main.go:104   iteration transaction result    {"iteration": 0, "endpoint": "https://api.mainnet-beta.solana.com", "slot": 0, "position": -1, "signature": "24DHDPwTaCz3UJsdkXM2of1oFVoS15KUXpNn8upfWkuTVa7L2SiS11NxvmAEQ9Wpor73CxpKW3GcEaG5ySVqxQGH"}
2022-07-27T14:43:28.055-0500    DEBUG   txcompare/main.go:104   iteration transaction result    {"iteration": 0, "endpoint": "https://my-rpc-endpoint.com", "slot": 0, "position": -1, "signature": "4efHFGHSJ9qQbtf35L4uZdf6Nt6nHnFwkFkbPURnjubGWvMi6sbXqxbSot523Ty7FxXgTqrNSg8e5kkff8iSxqMu"}
2022-07-27T14:45:08.515-0500    DEBUG   transaction/status.go:122       transaction failed execution    {"signature": "5weRaUwR3m7NiA3s8kdzWjJQinBT58XVr3BcjBn9QmE5EFXGuteKq1d3Ms3mD6eHc7jTo7aG7ybAye57tskTfgEr"}
2022-07-27T14:45:08.520-0500    DEBUG   transaction/status.go:122       transaction failed execution    {"signature": "3DhCLYeENXAuMnwrTyew7hNprofJShSGB4GGGfpPx58YzVswZwEX7wC1EGjezLrNfdBxY8pXsdFAzVvKMC97Wt6J"}
2022-07-27T14:45:08.521-0500    DEBUG   txcompare/main.go:85    iteration no transactions confirmed     {"iteration": 1}
2022-07-27T14:45:08.521-0500    DEBUG   txcompare/main.go:104   iteration transaction result    {"iteration": 1, "endpoint": "https://api.mainnet-beta.solana.com", "slot": 0, "position": -1, "signature": "3DhCLYeENXAuMnwrTyew7hNprofJShSGB4GGGfpPx58YzVswZwEX7wC1EGjezLrNfdBxY8pXsdFAzVvKMC97Wt6J"}
2022-07-27T14:45:08.521-0500    DEBUG   txcompare/main.go:104   iteration transaction result    {"iteration": 1, "endpoint": "https://my-rpc-endpoint.com", "slot": 0, "position": -1, "signature": "5weRaUwR3m7NiA3s8kdzWjJQinBT58XVr3BcjBn9QmE5EFXGuteKq1d3Ms3mD6eHc7jTo7aG7ybAye57tskTfgEr"}
2022-07-27T14:46:48.966-0500    DEBUG   transaction/status.go:122       transaction failed execution    {"signature": "5odrbxxoPFsUHsnssm7W26YDgvhM3cfgToxFj9dTDBiMtEyLhLdVC3u17h3CNWB9KNbxR6X6KHsWwEhCDFM4Leaz"}
2022-07-27T14:46:48.966-0500    DEBUG   transaction/status.go:122       transaction failed execution    {"signature": "2M6KGkLoof371zNmZfnYeNnfyDh1beePzbRnZGpXhugM1zqgjouSZiVVZDhomej4HDV1hPhNwVLk1MDKHNn1nig5"}
2022-07-27T14:46:48.966-0500    DEBUG   txcompare/main.go:85    iteration no transactions confirmed     {"iteration": 2}
2022-07-27T14:46:48.966-0500    DEBUG   txcompare/main.go:104   iteration transaction result    {"iteration": 2, "endpoint": "https://api.mainnet-beta.solana.com", "slot": 0, "position": -1, "signature": "2M6KGkLoof371zNmZfnYeNnfyDh1beePzbRnZGpXhugM1zqgjouSZiVVZDhomej4HDV1hPhNwVLk1MDKHNn1nig5"}
2022-07-27T14:46:48.966-0500    DEBUG   txcompare/main.go:104   iteration transaction result    {"iteration": 2, "endpoint": "https://my-rpc-endpoint.com", "slot": 0, "position": -1, "signature": "5odrbxxoPFsUHsnssm7W26YDgvhM3cfgToxFj9dTDBiMtEyLhLdVC3u17h3CNWB9KNbxR6X6KHsWwEhCDFM4Leaz"}
2022-07-27T14:48:29.517-0500    DEBUG   transaction/status.go:122       transaction failed execution    {"signature": "5P1nLzay7DMnFThKobhz18G5P6Gbjnoai5ia2Fk8xeurajgLjvZnB2uiQF5u2WXkxjk6xTUGYF7gtCKMCn6v3MD8"}
2022-07-27T14:48:29.517-0500    DEBUG   transaction/status.go:122       transaction failed execution    {"signature": "2a9Rs95ikWwM7PtXQLxk5fWadKcKemc1dJyrRs2Q2L95ep61zGv8gKqaCZ2tC5FoEsrQMdpMpuxaJxQ57GW8ZCF8"}
2022-07-27T14:48:29.517-0500    DEBUG   txcompare/main.go:85    iteration no transactions confirmed     {"iteration": 3}
2022-07-27T14:48:29.517-0500    DEBUG   txcompare/main.go:104   iteration transaction result    {"iteration": 3, "endpoint": "https://api.mainnet-beta.solana.com", "slot": 0, "position": -1, "signature": "2a9Rs95ikWwM7PtXQLxk5fWadKcKemc1dJyrRs2Q2L95ep61zGv8gKqaCZ2tC5FoEsrQMdpMpuxaJxQ57GW8ZCF8"}
2022-07-27T14:48:29.517-0500    DEBUG   txcompare/main.go:104   iteration transaction result    {"iteration": 3, "endpoint": "https://my-rpc-endpoint.com", "slot": 0, "position": -1, "signature": "5P1nLzay7DMnFThKobhz18G5P6Gbjnoai5ia2Fk8xeurajgLjvZnB2uiQF5u2WXkxjk6xTUGYF7gtCKMCn6v3MD8"}
Iterations:  4
Endpoints:
     https://api.mainnet-beta.solana.com
     https://my-rpc-endpoint.com

Win counts: 
    0    https://api.mainnet-beta.solana.com
    0    https://my-rpc-endpoint.com

Lost transactions: 
    4    https://api.mainnet-beta.solana.com
    4    https://my-rpc-endpoint.com
```

A CSV file for also be generated at the `--output` location, with details of each transaction.