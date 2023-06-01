# benchmark/quotes

Compares Solana Trader API prices stream to Jupiter's API. This comparison is most optimally done by choosing a low
liquidity market (GOSU by default) then using this script to generate swaps for that market so you can identify when 
each source responds to the update in the chain.

Note that this comparison is not entirely fair: Jupiter's API does not support streaming, so a poll-based approach
must be taken. You can configure a poll timeout to your preference, but will have to deal with rate limits. To help
with this somewhat, an HTTP polling based approach is also included for Trader API, though using websockets is
expected to be significantly more helpful.

You will need a bloXroute `AUTH_HEADER` to access Trader API, and a Solana `PRIVATE_KEY` to generate swaps.

## Usage

Full test (10 swaps for 0.1 USDC => GOSU):
```
$ AUTH_HEADER=... PRIVATE_KEY=... go run ./benchmark/quotes --iterations 10 --runtime 5m --swap-amount 0.1 --trigger-activity --output updates.csv
```

Single swap test (1 swap for 0.1 USDC => GOSU, useful for debugging):
```
$ AUTH_HEADER=... PRIVATE_KEY=... go run ./benchmark/quotes --iterations 1 --runtime 1m --swap-amount 0.1 --trigger-activity --output updates.csv
```

No swaps (useful for debugging, though WS might not produce any results):
```
$ AUTH_HEADER=... go run ./benchmark/quotes --runtime 10s --output updates.csv
```

## Result



```
2023-06-01T17:08:54.987-0500    INFO    quotes/main.go:93       trader API clients connected    {"env": "mainnet"}
2023-06-01T17:08:54.990-0500    INFO    quotes/main.go:127      starting all routines   {"duration": "5m0s"}
2023-06-01T17:09:04.990-0500    INFO    actor/jupiter_swap.go:82        starting swap submission        {"source": "jupiterActor", "total": 10}
2023-06-01T17:09:07.991-0500    INFO    actor/jupiter_swap.go:88        submitting swap {"source": "jupiterActor", "count": 0}
2023-06-01T17:09:10.991-0500    INFO    actor/jupiter_swap.go:88        submitting swap {"source": "jupiterActor", "count": 1}
2023-06-01T17:09:11.923-0500    INFO    actor/jupiter_swap.go:97        completed swap  {"source": "jupiterActor", "transactions": [{"signature":"2XXH3GJhWcmZX2o9qqUp29apsue6cckcAanGg3um5gWA2NGX4iKPzgTMBc8QJaCo5VuA5CHMh3vpfcKaS1YWtMaJ","submitted":true}]}
2023-06-01T17:09:13.992-0500    INFO    actor/jupiter_swap.go:88        submitting swap {"source": "jupiterActor", "count": 2}
2023-06-01T17:09:16.087-0500    INFO    actor/jupiter_swap.go:97        completed swap  {"source": "jupiterActor", "transactions": [{"signature":"AwdiZeSzPTuVCUi2CSLAP1xD6FTJHNzkwPBwqRU2LqMjcmVU8gWXKWQQwHLoiqixzkC6mns9AMd5oaFaXmRGcu4","submitted":true}]}
2023-06-01T17:09:16.567-0500    INFO    actor/jupiter_swap.go:97        completed swap  {"source": "jupiterActor", "transactions": [{"signature":"y7zv1QFJjz2eZVRQksn1nHEvLVtorKuCmv2A5Ud4pCb6jpGvsYkBXFTnTTDghRvhZwmGnCMiwRoBo5b6do6cGJ6","submitted":true}]}
2023-06-01T17:09:16.991-0500    INFO    actor/jupiter_swap.go:88        submitting swap {"source": "jupiterActor", "count": 3}
2023-06-01T17:09:18.191-0500    INFO    actor/jupiter_swap.go:97        completed swap  {"source": "jupiterActor", "transactions": [{"signature":"2V6Xydunh6Z5Mzmz2Z8LWogT8tq4BQgxTFbCLY5wYW7ogHLAvZEjeeSJG1qMBWvfdA9tji1ozDGrchmQA2ijKU2n","submitted":true}]}
2023-06-01T17:09:19.991-0500    INFO    actor/jupiter_swap.go:88        submitting swap {"source": "jupiterActor", "count": 4}
2023-06-01T17:09:21.860-0500    INFO    actor/jupiter_swap.go:97        completed swap  {"source": "jupiterActor", "transactions": [{"signature":"5xTvn3LMyUqwrtv5MyAwb2GrUn6D38kKnZ1oVmkTJv2LjSqQEjbBLVef7FKBR5do2oVAqyFR2JZSUxxTqb4t5t3w","submitted":true}]}
2023-06-01T17:09:22.991-0500    INFO    actor/jupiter_swap.go:88        submitting swap {"source": "jupiterActor", "count": 5}
2023-06-01T17:09:24.497-0500    INFO    actor/jupiter_swap.go:97        completed swap  {"source": "jupiterActor", "transactions": [{"signature":"5yKJEQxNHRw6g9s1YbsjRbZwWfZvhYgi4z32eHFzeu58WiPgMgfxbyquQRqovvmJCSfesEmkTWgFHP449Bfdh5aw","submitted":true}]}
2023-06-01T17:09:25.991-0500    INFO    actor/jupiter_swap.go:88        submitting swap {"source": "jupiterActor", "count": 6}
2023-06-01T17:09:27.964-0500    INFO    actor/jupiter_swap.go:97        completed swap  {"source": "jupiterActor", "transactions": [{"signature":"4BbUk4vox6rHNNui9zmuBd5yvLHPw8SVdHPAzQvwQ6BA8Ww2vHUzMddg9k35sZiZLZ6ohKzXtx4XsepXmPvpXz8k","submitted":true}]}
2023-06-01T17:09:28.991-0500    INFO    actor/jupiter_swap.go:88        submitting swap {"source": "jupiterActor", "count": 7}
2023-06-01T17:09:30.080-0500    INFO    actor/jupiter_swap.go:97        completed swap  {"source": "jupiterActor", "transactions": [{"signature":"221QKF7UKMJLADPx9MCySQCCS26wseS64W2qigswFHN1ttQoK7qA8KrraV3jD9kyqC4DjEiK6DAfbK3crX6Jregk","submitted":true}]}
2023-06-01T17:09:31.991-0500    INFO    actor/jupiter_swap.go:88        submitting swap {"source": "jupiterActor", "count": 8}
2023-06-01T17:09:33.920-0500    INFO    actor/jupiter_swap.go:97        completed swap  {"source": "jupiterActor", "transactions": [{"signature":"57ptnTz5RA6SRSasJGq1PGUgY28mfaG6cRHvxmETrxC5q6etDCp1QYRHs93ZXH63AWsRLfqRy4C9aLPfESR7g4iP","submitted":true}]}
2023-06-01T17:09:34.991-0500    INFO    actor/jupiter_swap.go:88        submitting swap {"source": "jupiterActor", "count": 9}
2023-06-01T17:09:35.775-0500    INFO    actor/jupiter_swap.go:97        completed swap  {"source": "jupiterActor", "transactions": [{"signature":"4wcMSuxV6omcFNJe8Qtc5cnfhoScNVAQjjrcj3BDQPpzw1QRKe9FLEWwFmwDe4eWcHMPmtEoaegnMPnf5zUf8xHQ","submitted":true}]}
2023-06-01T17:09:45.777-0500    ERROR   stream/traderhttp_price.go:59   could not fetch price   {"source": "traderapi/http", "err": "Get \"https://virginia.solana.dex.blxrbdn.com/api/v1/market/price?tokens=6D7nXHAhsRbwj8KFZR2agB6GEjMLg4BM7MAqZzRT8F1j\": context canceled"}
github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/stream.traderHTTPPriceStream.Run.func2
        /Users/aspin/workspace/solana-trader-client-go/benchmark/internal/stream/traderhttp_price.go:59
github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/stream.collectOrderedUpdates[...].func1
        /Users/aspin/workspace/solana-trader-client-go/benchmark/internal/stream/source.go:70
2023-06-01T17:09:45.777-0500    INFO    quotes/main.go:189      ignoring jupiter duplicates     {"count": 18}
2023-06-01T17:09:45.777-0500    INFO    quotes/main.go:192      ignoring tradeWS duplicates     {"count": 0}
2023-06-01T17:09:45.777-0500    INFO    quotes/main.go:195      ignoring tradeWS duplicates     {"count": 1}

Trader API vs. Jupiter API Benchmark

Swaps placed:  10
Jun  1 17:09:11.923: 2XXH3GJhWcmZX2o9qqUp29apsue6cckcAanGg3um5gWA2NGX4iKPzgTMBc8QJaCo5VuA5CHMh3vpfcKaS1YWtMaJ
Jun  1 17:09:16.087: AwdiZeSzPTuVCUi2CSLAP1xD6FTJHNzkwPBwqRU2LqMjcmVU8gWXKWQQwHLoiqixzkC6mns9AMd5oaFaXmRGcu4
Jun  1 17:09:16.567: y7zv1QFJjz2eZVRQksn1nHEvLVtorKuCmv2A5Ud4pCb6jpGvsYkBXFTnTTDghRvhZwmGnCMiwRoBo5b6do6cGJ6
Jun  1 17:09:18.191: 2V6Xydunh6Z5Mzmz2Z8LWogT8tq4BQgxTFbCLY5wYW7ogHLAvZEjeeSJG1qMBWvfdA9tji1ozDGrchmQA2ijKU2n
Jun  1 17:09:21.860: 5xTvn3LMyUqwrtv5MyAwb2GrUn6D38kKnZ1oVmkTJv2LjSqQEjbBLVef7FKBR5do2oVAqyFR2JZSUxxTqb4t5t3w
Jun  1 17:09:24.502: 5yKJEQxNHRw6g9s1YbsjRbZwWfZvhYgi4z32eHFzeu58WiPgMgfxbyquQRqovvmJCSfesEmkTWgFHP449Bfdh5aw
Jun  1 17:09:27.964: 4BbUk4vox6rHNNui9zmuBd5yvLHPw8SVdHPAzQvwQ6BA8Ww2vHUzMddg9k35sZiZLZ6ohKzXtx4XsepXmPvpXz8k
Jun  1 17:09:30.080: 221QKF7UKMJLADPx9MCySQCCS26wseS64W2qigswFHN1ttQoK7qA8KrraV3jD9kyqC4DjEiK6DAfbK3crX6Jregk
Jun  1 17:09:33.920: 57ptnTz5RA6SRSasJGq1PGUgY28mfaG6cRHvxmETrxC5q6etDCp1QYRHs93ZXH63AWsRLfqRy4C9aLPfESR7g4iP
Jun  1 17:09:35.775: 4wcMSuxV6omcFNJe8Qtc5cnfhoScNVAQjjrcj3BDQPpzw1QRKe9FLEWwFmwDe4eWcHMPmtEoaegnMPnf5zUf8xHQ

Jupiter:  50  samples
Start time: Jun  1 17:08:55.990
End time: Jun  1 17:09:44.989
Slot range: 197243271 => 197243373
Price change: 0.0004576648625989163 => 0.0004612747283519241 
Distinct prices: 15

Trader WS:  10  samples
Start time: Jun  1 17:08:55.429
End time: Jun  1 17:09:37.647
Slot range: 197243269 => 197243361
Buy change: 0.000452 => 0.000456
Sell change: 0.00045766530815864395 => 0.0004612747220814633
Distinct buy prices: 5
Distinct sell prices: 10

Trader HTTP:  50  samples
Start time: Jun  1 17:08:55.490
End time: Jun  1 17:09:45.489
Buy change: 0.000452 => 0.000456
Sell change: 0.00045766530815864395 => 0.0004612747220814633
Distinct buy prices: 6
Distinct sell prices: 16

jupiter API
[197243271] 2023-06-01 17:08:55.990535 -0500 CDT m=+1.444809126 [473ms]: B: 0.0004576648625989163 | S: 0.0004576648625989163
[197243271] 2023-06-01 17:08:56.988544 -0500 CDT m=+2.442817959 [161ms]: B: 0.00045765754919974524 | S: 0.00045765754919974524
[197243273] 2023-06-01 17:08:57.992009 -0500 CDT m=+3.446284251 [107ms]: B: 0.0004576648625989163 | S: 0.0004576648625989163
[197243277] 2023-06-01 17:08:59.989702 -0500 CDT m=+5.443977459 [94ms]: B: 0.0004576648625989163 | S: 0.0004576648625989163
[197243278] 2023-06-01 17:08:58.988636 -0500 CDT m=+4.442911251 [274ms]: B: 0.00045765754919974524 | S: 0.00045765754919974524
[197243281] 2023-06-01 17:09:00.989456 -0500 CDT m=+6.443731751 [153ms]: B: 0.00045765754919974524 | S: 0.00045765754919974524
[197243283] 2023-06-01 17:09:02.988901 -0500 CDT m=+8.443178459 [313ms]: B: 0.0004576648625989163 | S: 0.0004576648625989163
[197243291] 2023-06-01 17:09:06.993928 -0500 CDT m=+12.448207751 [164ms]: B: 0.00045765754919974524 | S: 0.00045765754919974524
[197243297] 2023-06-01 17:09:07.989974 -0500 CDT m=+13.444253542 [95ms]: B: 0.0004576648625989163 | S: 0.0004576648625989163
[197243299] 2023-06-01 17:09:08.991257 -0500 CDT m=+14.445537626 [110ms]: B: 0.00045765754919974524 | S: 0.00045765754919974524
[197243299] 2023-06-01 17:09:10.989275 -0500 CDT m=+16.443556709 [223ms]: B: 0.0004576648625989163 | S: 0.0004576648625989163
[197243303] 2023-06-01 17:09:11.989282 -0500 CDT m=+17.443564292 [245ms]: B: 0.00045765754919974524 | S: 0.00045765754919974524
[197243305] 2023-06-01 17:09:12.994719 -0500 CDT m=+18.449001126 [241ms]: B: 0.0004576648625989163 | S: 0.0004576648625989163
[197243305] 2023-06-01 17:09:13.988555 -0500 CDT m=+19.442837376 [287ms]: B: 0.00045765754919974524 | S: 0.00045765754919974524
[197243307] 2023-06-01 17:09:15.989258 -0500 CDT m=+21.443542126 [120ms]: B: 0.0004580178764348867 | S: 0.0004580178764348867
[197243316] 2023-06-01 17:09:18.989353 -0500 CDT m=+24.443637959 [265ms]: B: 0.00045874631535347717 | S: 0.00045874631535347717
[197243318] 2023-06-01 17:09:20.989754 -0500 CDT m=+26.444040459 [415ms]: B: 0.00045873898464079065 | S: 0.00045873898464079065
[197243319] 2023-06-01 17:09:19.989843 -0500 CDT m=+25.444128334 [329ms]: B: 0.00045909976202499695 | S: 0.00045909976202499695
[197243325] 2023-06-01 17:09:21.989743 -0500 CDT m=+27.444030251 [868ms]: B: 0.0004591070985143641 | S: 0.0004591070985143641
[197243332] 2023-06-01 17:09:24.991684 -0500 CDT m=+30.445972709 [196ms]: B: 0.0004594679980028692 | S: 0.0004594679980028692
[197243335] 2023-06-01 17:09:26.989321 -0500 CDT m=+32.443610459 [116ms]: B: 0.0004598290644239391 | S: 0.0004598290644239391
[197243341] 2023-06-01 17:09:28.989229 -0500 CDT m=+34.443519251 [136ms]: B: 0.00045982171637443613 | S: 0.00045982171637443613
[197243342] 2023-06-01 17:09:29.989253 -0500 CDT m=+35.443544209 [86ms]: B: 0.0004598290644239391 | S: 0.0004598290644239391
[197243345] 2023-06-01 17:09:30.989252 -0500 CDT m=+36.443543792 [129ms]: B: 0.0004601828933233086 | S: 0.0004601828933233086
[197243348] 2023-06-01 17:09:31.989195 -0500 CDT m=+37.443486709 [392ms]: B: 0.00046055159683767413 | S: 0.00046055159683767413
[197243355] 2023-06-01 17:09:35.989234 -0500 CDT m=+41.443528001 [479ms]: B: 0.0004609130917751887 | S: 0.0004609130917751887
[197243359] 2023-06-01 17:09:36.989263 -0500 CDT m=+42.443557251 [358ms]: B: 0.0004612747283519241 | S: 0.0004612747283519241
[197243363] 2023-06-01 17:09:38.989257 -0500 CDT m=+44.443551959 [455ms]: B: 0.0004612673571546972 | S: 0.0004612673571546972
[197243364] 2023-06-01 17:09:39.989279 -0500 CDT m=+45.443575376 [75ms]: B: 0.0004612747283519241 | S: 0.0004612747283519241
[197243368] 2023-06-01 17:09:41.989215 -0500 CDT m=+47.443511501 [95ms]: B: 0.0004612673571546972 | S: 0.0004612673571546972
[197243372] 2023-06-01 17:09:43.989225 -0500 CDT m=+49.443523376 [275ms]: B: 0.0004612747283519241 | S: 0.0004612747283519241

traderWS
[197243269] 2023-06-01 17:08:55.429213 -0500 CDT m=+0.883486417 [0ms]: B: 0.000452 | S: 0.00045766530815864395
[197243274] 2023-06-01 17:08:56.91252 -0500 CDT m=+2.366794376 [0ms]: B: 0.000452 | S: 0.00045766486203128743
[197243317] 2023-06-01 17:09:18.664077 -0500 CDT m=+24.118362251 [0ms]: B: 0.000453 | S: 0.00045874631795672123
[197243321] 2023-06-01 17:09:20.760472 -0500 CDT m=+26.214758001 [0ms]: B: 0.000453 | S: 0.0004591066391381409
[197243326] 2023-06-01 17:09:22.333984 -0500 CDT m=+27.788270876 [0ms]: B: 0.000453 | S: 0.0004591070880807876
[197243329] 2023-06-01 17:09:23.712465 -0500 CDT m=+29.166752292 [0ms]: B: 0.000454 | S: 0.0004594680110284085
[197243345] 2023-06-01 17:09:31.446862 -0500 CDT m=+36.901153376 [0ms]: B: 0.000455 | S: 0.0004601902435243478
[197243349] 2023-06-01 17:09:33.872055 -0500 CDT m=+39.326347376 [0ms]: B: 0.000455 | S: 0.00046055158636167136
[197243358] 2023-06-01 17:09:36.29817 -0500 CDT m=+41.752463626 [0ms]: B: 0.000455 | S: 0.00046091307942387635
[197243361] 2023-06-01 17:09:37.647866 -0500 CDT m=+43.102160876 [0ms]: B: 0.000456 | S: 0.0004612747220814633

traderHTTP
[-1] 2023-06-01 17:08:55.849004 -0500 CDT m=+1.303277917 [358ms]: B: 0.000452 | S: 0.00045766530815864395
[-1] 2023-06-01 17:08:56.985607 -0500 CDT m=+2.439881667 [495ms]: B: 0.000452 | S: 0.00045766486203128743
[-1] 2023-06-01 17:09:09.23887 -0500 CDT m=+14.693150876 [746ms]: B: 0.000453 | S: 0.00045769001810883585
[-1] 2023-06-01 17:09:14.382206 -0500 CDT m=+19.836488917 [2892ms]: B: 0.000452 | S: 0.00045766486203128743
[-1] 2023-06-01 17:09:12.783528 -0500 CDT m=+18.237810209 [294ms]: B: 0.000453 | S: 0.00045769001810883585
[-1] 2023-06-01 17:09:14.531465 -0500 CDT m=+19.985748417 [42ms]: B: 0.000452 | S: 0.00045766486203128743
[-1] 2023-06-01 17:09:17.090663 -0500 CDT m=+22.544947584 [600ms]: B: 0.000453 | S: 0.00045805036559882833
[-1] 2023-06-01 17:09:18.708819 -0500 CDT m=+24.163104417 [220ms]: B: 0.000453 | S: 0.00045874631795672123
[-1] 2023-06-01 17:09:20.760614 -0500 CDT m=+26.214900626 [1271ms]: B: 0.000453 | S: 0.0004591066391381409
[-1] 2023-06-01 17:09:22.334148 -0500 CDT m=+27.788434751 [845ms]: B: 0.000453 | S: 0.0004591070880807876
[-1] 2023-06-01 17:09:24.162181 -0500 CDT m=+29.616468584 [668ms]: B: 0.000455 | S: 0.0004594932273574347
[-1] 2023-06-01 17:09:26.531783 -0500 CDT m=+31.986072459 [42ms]: B: 0.000454 | S: 0.0004594680110284085
[-1] 2023-06-01 17:09:29.03424 -0500 CDT m=+34.488530084 [545ms]: B: 0.000455 | S: 0.00045985434215983697
[-1] 2023-06-01 17:09:30.589175 -0500 CDT m=+36.043465917 [99ms]: B: 0.000455 | S: 0.0004601902435243478
[-1] 2023-06-01 17:09:32.963572 -0500 CDT m=+38.417864376 [474ms]: B: 0.000456 | S: 0.00046057688701531165
[-1] 2023-06-01 17:09:34.726204 -0500 CDT m=+40.180497542 [236ms]: B: 0.000455 | S: 0.00046055158636167136
[-1] 2023-06-01 17:09:35.875643 -0500 CDT m=+41.329936459 [387ms]: B: 0.000456 | S: 0.00046093841981183545
[-1] 2023-06-01 17:09:38.06379 -0500 CDT m=+43.518084834 [574ms]: B: 0.000457 | S: 0.00046130006738659575
[-1] 2023-06-01 17:09:42.628165 -0500 CDT m=+48.082462042 [139ms]: B: 0.000456 | S: 0.0004612747220814633
```

A CSV file will also be generated with a time-sorted event list. Transaction info is only written to stdout.