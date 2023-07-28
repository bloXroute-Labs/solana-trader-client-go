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
$ AUTH_HEADER=... PRIVATE_KEY=... go run ./benchmark/quotes --iterations 10 --runtime 5m --swap-amount 0.1 --swap-interval 30s --trigger-activity --output updates.csv
```

Single swap test (1 swap for 0.1 USDC => GOSU, useful for debugging):
```
$ AUTH_HEADER=... PRIVATE_KEY=... go run ./benchmark/quotes --iterations 1 --runtime 1m --swap-amount 0.1 --swap-interval 5s --trigger-activity --output updates.csv
```

No swaps (useful for debugging, though WS might not produce any results):
```
$ AUTH_HEADER=... go run ./benchmark/quotes --runtime 10s --output updates.csv
```

## Result



```
2023-06-02T13:10:29.125-0500    INFO    quotes/main.go:99       trader API clients connected    {"env": "mainnet"}
2023-06-02T13:10:29.125-0500    INFO    quotes/main.go:133      starting all routines   {"duration": "1m0s"}
2023-06-02T13:10:39.126-0500    INFO    actor/jupiter_swap.go:82        starting swap submission        {"source": "jupiterActor", "total": 1}
2023-06-02T13:10:44.127-0500    INFO    actor/jupiter_swap.go:88        submitting swap {"source": "jupiterActor", "count": 0}
2023-06-02T13:10:45.613-0500    INFO    actor/jupiter_swap.go:97        completed swap  {"source": "jupiterActor", "transactions": [{"signature":"FkQABUjAXbqNroR5Y4xnnhS5RhtDe7abuyYusumJCpVpbrTGjvDtNdo4QdCWmTbhZVbiHFuHyDxkyfYqTzSupf6","submitted":true}]}
2023-06-02T13:10:55.618-0500    INFO    quotes/main.go:195      ignoring jupiter duplicates     {"count": 19}
2023-06-02T13:10:55.618-0500    INFO    quotes/main.go:198      ignoring tradeWS duplicates     {"count": 0}
2023-06-02T13:10:55.618-0500    INFO    quotes/main.go:201      ignoring tradeWS duplicates     {"count": 1}

Trader API vs. Jupiter API Benchmark

Swaps placed:  1
Jun  2 13:10:45.613: FkQABUjAXbqNroR5Y4xnnhS5RhtDe7abuyYusumJCpVpbrTGjvDtNdo4QdCWmTbhZVbiHFuHyDxkyfYqTzSupf6

Jupiter:  26  samples
Start time: Jun  2 13:10:30.126
End time: Jun  2 13:10:55.128
Slot range: 197398959 => 197399010
Price change: 0.00047147519075931205 => 0.0004718409190889233 
Distinct prices: 3

Trader WS:  2  samples
Start time: Jun  2 13:10:29.335
End time: Jun  2 13:10:47.571
Slot range: 197398958 => 197398999
Buy change: 0.000467 => 0.000467
Sell change: 0.0004714926053447434 => 0.00047185835217295884
Distinct buy prices: 1
Distinct sell prices: 2

Trader HTTP:  26  samples
Start time: Jun  2 13:10:29.626
End time: Jun  2 13:10:54.626
Buy change: 0.000467 => 0.000467
Sell change: 0.0004714926053447434 => 0.00047185835217295884
Distinct buy prices: 2
Distinct sell prices: 3

jupiter API
[197398959] 2023-06-02 13:10:30.126294 -0500 CDT m=+1.282104543 [790ms]: B: 0.00047147519075931205 | S: 0.00047147519075931205
[197398979] 2023-06-02 13:10:40.128294 -0500 CDT m=+11.284147626 [129ms]: B: 0.000471492615933144 | S: 0.000471492615933144
[197398981] 2023-06-02 13:10:41.126732 -0500 CDT m=+12.282590376 [322ms]: B: 0.00047147519075931205 | S: 0.00047147519075931205
[197398987] 2023-06-02 13:10:45.126687 -0500 CDT m=+16.282562418 [420ms]: B: 0.00047147519075931205 | S: 0.00047147519075931205
[197398988] 2023-06-02 13:10:44.126723 -0500 CDT m=+15.282594084 [268ms]: B: 0.000471492615933144 | S: 0.000471492615933144
[197398996] 2023-06-02 13:10:49.127456 -0500 CDT m=+20.283349126 [297ms]: B: 0.0004718409190889233 | S: 0.0004718409190889233

traderWS
[197398958] 2023-06-02 13:10:29.335606 -0500 CDT m=+0.491412793 [0ms]: B: 0.000467 | S: 0.0004714926053447434
[197398999] 2023-06-02 13:10:47.571456 -0500 CDT m=+18.727342418 [0ms]: B: 0.000467 | S: 0.00047185835217295884

traderHTTP
[-1] 2023-06-02 13:10:29.860785 -0500 CDT m=+1.016593918 [233ms]: B: 0.000467 | S: 0.0004714926053447434
[-1] 2023-06-02 13:10:45.471033 -0500 CDT m=+16.626910209 [844ms]: B: 0.000466 | S: 0.00047147894730418404
[-1] 2023-06-02 13:10:46.486188 -0500 CDT m=+17.642068876 [859ms]: B: 0.000467 | S: 0.0004714926053447434
[-1] 2023-06-02 13:10:47.683645 -0500 CDT m=+18.839531959 [56ms]: B: 0.000467 | S: 0.00047185835217295884
```

A CSV file will also be generated with a time-sorted event list.

```
timestamp,source,slot,processingTime,buy,sell
2023-06-02T13:10:29.335606-05:00,traderWS,197398958,0s,0.000467,0.0004714926053447434
2023-06-02T13:10:29.860785-05:00,traderHTTP,-1,233.970375ms,0.000467,0.0004714926053447434
2023-06-02T13:10:30.126294-05:00,jupiter,197398959,790.252125ms,0.00047147519075931205,0.00047147519075931205
2023-06-02T13:10:40.128294-05:00,jupiter,197398979,129.040042ms,0.000471492615933144,0.000471492615933144
2023-06-02T13:10:41.126732-05:00,jupiter,197398981,322.407708ms,0.00047147519075931205,0.00047147519075931205
2023-06-02T13:10:44.126723-05:00,jupiter,197398988,268.150959ms,0.000471492615933144,0.000471492615933144
2023-06-02T13:10:45.126687-05:00,jupiter,197398987,420.006041ms,0.00047147519075931205,0.00047147519075931205
2023-06-02T13:10:45.471033-05:00,traderHTTP,-1,844.338583ms,0.000466,0.00047147894730418404
2023-06-02T13:10:46-05:00,transaction-FkQABUjAXbqNroR5Y4xnnhS5RhtDe7abuyYusumJCpVpbrTGjvDtNdo4QdCWmTbhZVbiHFuHyDxkyfYqTzSupf6,197398996,386.508ms,0,0
2023-06-02T13:10:46.486188-05:00,traderHTTP,-1,859.4435ms,0.000467,0.0004714926053447434
2023-06-02T13:10:47.571456-05:00,traderWS,197398999,0s,0.000467,0.00047185835217295884
2023-06-02T13:10:47.683645-05:00,traderHTTP,-1,56.971583ms,0.000467,0.00047185835217295884
2023-06-02T13:10:49.127456-05:00,jupiter,197398996,297.499125ms,0.0004718409190889233,0.0004718409190889233
```