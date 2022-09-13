# benchmark/traderapi

Compares Solana Trader API orderbook stream to a direct connection with Solana. Identifies updates with the same slot number,
and indicates how fast updates on each stream were relative to each other. Note that all raw data is collected
immediately, and all processing and deserializing happens at the end to avoid noise from deserialization overhead.
Returns some unused data about messages that were seen only one connection vs. the other for future debugging of
reliability.

## Usage

Go:
```
$ AUTH_HEADER=... go run ./benchmark/traderapi --run-time 10s --output result.csv
```

Docker:
```
$ docker run -e AUTH_HEADER=... --name cmp --rm 033969152235.dkr.ecr.us-east-1.amazonaws.com/serum-api:bm-serumapi
```

## Result

```
2022-09-13T14:23:55.672-0500	DEBUG	arrival/traderws.go:53	connection established	{"source": "traderapi", "address": "ws://54.163.206.248:1809/ws", "market": "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT"}
2022-09-13T14:23:55.873-0500	DEBUG	arrival/solanaws.go:83	connection established	{"source": "solanaws", "address": "ws://185.209.178.55", "market": "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT"}
2022-09-13T14:23:55.874-0500	DEBUG	arrival/traderws.go:81	subscription created	{"source": "traderapi", "address": "ws://54.163.206.248:1809/ws", "market": "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT"}
2022-09-13T14:23:55.878-0500	DEBUG	arrival/solanaws.go:110	subscription created	{"source": "solanaws", "address": "ws://185.209.178.55", "market": "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT"}
2022-09-13T14:24:05.875-0500	INFO	traderapi/main.go:120	waited 10s out of 10s...
2022-09-13T14:24:05.876-0500	DEBUG	arrival/traderws.go:92	closing connection	{"source": "traderapi", "address": "ws://54.163.206.248:1809/ws", "market": "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT", "err": "read tcp 192.168.1.214:60182->54.163.206.248:1809: use of closed network connection"}
2022-09-13T14:24:05.876-0500	INFO	traderapi/main.go:125	finished collecting data points	{"tradercount": 44, "solanacount": 26}
2022-09-13T14:24:05.886-0500	DEBUG	arrival/solanaws.go:123	closing Asks subscription	{"source": "solanaws", "address": "ws://185.209.178.55", "market": "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT", "err": "read tcp 192.168.1.214:60184->185.209.178.55:80: use of closed network connection"}
2022-09-13T14:24:05.886-0500	DEBUG	arrival/solanaws.go:142	closing Bids subscription	{"source": "solanaws", "address": "ws://185.209.178.55", "market": "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT", "err": "read tcp 192.168.1.214:60184->185.209.178.55:80: use of closed network connection"}
2022-09-13T14:24:05.936-0500	DEBUG	traderapi/main.go:132	processed trader API results	{"range": "150523364-150523376", "count": 13, "duplicaterange": "150523366-150523372", "duplicatecount": 2}
2022-09-13T14:24:06.033-0500	DEBUG	traderapi/main.go:139	processed solana results	{"range": "150523364-150523376", "count": 12, "duplicaterange": "150523365-150523377", "duplicatecount": 6}
2022-09-13T14:24:06.033-0500	DEBUG	traderapi/main.go:142	finished processing data points	{"startSlot": 150523364, "endSlot": 150523376, "count": 13}
2022-09-13T14:24:06.033-0500	INFO	traderapi/main.go:149	completed merging: outputting data...
Run time:  10s
Endpoints:
     ws://54.163.206.248:1809/ws  [serum]
     ws://185.209.178.55  [solana]

Total updates:  43

Faster counts:
    19      ws://54.163.206.248:1809/ws
    0       ws://185.209.178.55
Average difference( ms):
    352ms   ws://54.163.206.248:1809/ws
    0ms     ws://185.209.178.55
Unmatched updates:
(updates from each stream without a corresponding result on the other)
    20      ws://54.163.206.248:1809/ws
    0       ws://185.209.178.55
```
