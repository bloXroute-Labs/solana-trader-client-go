# benchmark/pumpfun_newtoken_compare

Compares Solana Trader API pumpfun new token stream to another thirdparty blocksubscribe stream.

## Usage

Go:
```
$ AUTH_HEADER=... go run ./benchmark/pumpfun_newtoken_compare
```

## Result

```

Mahmoud Taabodi
  5:35 PM
BlockTime         TraderAPIEventTime     ThirdPartyEventTime       Diff(thirdParty)       Diff(Blocktime)
1726263118000     1726263140243          1726263140243,            -21.091899 sec,        1.151549 sec,
1726263118000     1726263140245          1726263140245,            -21.044734 sec,        1.200886 sec,
1726263133000     1726263156966          1726263156966,            -21.959157 sec,        2.006891 sec,
1726263156000     1726263171369          1726263171369,            -12.996627 sec,        2.373115 sec,
1726263165000     1726263181012          1726263181012,            -14.125343 sec,        1.887493 sec,
1726263171000     1726263187884          1726263187884,            -15.026972 sec,        1.857471 sec,
1726263174000     1726263190033          1726263190033,            -15.338945 sec,        0.694350 sec,
1726263215000     1726263235128          1726263235128,            -18.032521 sec,        2.095713 sec,
Run time:  3m0s

Total events:  8

Faster counts: 
 traderAPIFaster   8
 thirdPartyFaster   0
```
