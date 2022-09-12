# benchmark/serumapi

Compares Serum API orderbook stream to a direct connection with Solana. Identifies updates with the same slot number,
and indicates how fast updates on each stream were relative to each other. Note that all raw data is collected
immediately, and all processing and deserializing happens at the end to avoid noise from deserialization overhead.
Returns some unused data about messages that were seen only one connection vs. the other for future debugging of
reliability.

## Usage

Go:
```
$ go run ./benchmark/serumapi --run-time 10s --output result.csv
```

Docker:
```
$ docker run --name cmp --rm 033969152235.dkr.ecr.us-east-1.amazonaws.com/serum-api:bm-serumapi
```

## Result

```
2022-06-30T16:54:20.280-0500	DEBUG	serumapi/serumws.go:43	connection established	{"source": "serum", "address": "ws://serum-nlb-53baf45ef9775263.elb.us-east-1.amazonaws.com/ws", "market": "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT"}
2022-06-30T16:54:20.476-0500	DEBUG	serumapi/solanaws.go:64	connection established	{"source": "solanaws", "address": "ws://34.203.186.197:8900/ws", "market": "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT"}
2022-06-30T16:54:20.476-0500	DEBUG	serumapi/serumws.go:67	subscription created	{"source": "serum", "address": "ws://serum-nlb-53baf45ef9775263.elb.us-east-1.amazonaws.com/ws", "market": "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT"}
2022-06-30T16:54:20.483-0500	DEBUG	serumapi/solanaws.go:87	subscription created	{"source": "solanaws", "address": "ws://34.203.186.197:8900/ws", "market": "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT"}
2022-06-30T16:54:30.477-0500	INFO	serumapi/main.go:132	waited 10s out of 10s...
2022-06-30T16:54:30.477-0500	DEBUG	serumapi/solanaws.go:100	closing asks subscription	{"source": "solanaws", "address": "ws://34.203.186.197:8900/ws", "market": "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT", "err": "read tcp 10.0.41.87:49403->34.203.186.197:8900: use of closed network connection"}
2022-06-30T16:54:30.477-0500	DEBUG	serumapi/serumws.go:78	closing connection	{"source": "serum", "address": "ws://serum-nlb-53baf45ef9775263.elb.us-east-1.amazonaws.com/ws", "market": "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT", "err": "read tcp 10.0.41.87:49401->44.205.114.20:80: use of closed network connection"}
2022-06-30T16:54:30.477-0500	INFO	serumapi/main.go:137	finished collecting data points	{"serumcount": 10, "solanacount": 24}
2022-06-30T16:54:30.477-0500	DEBUG	serumapi/solanaws.go:119	closing bids subscription	{"source": "solanaws", "address": "ws://34.203.186.197:8900/ws", "market": "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT", "err": "read tcp 10.0.41.87:49403->34.203.186.197:8900: use of closed network connection"}
2022-06-30T16:54:30.527-0500	INFO	serumapi/main.go:150	finished processing data points	{"startSlot": 139696869, "endSlot": 139696884}
2022-06-30T16:54:30.527-0500	INFO	serumapi/main.go:157	completed merging: outputting data...
comparison points:
slot 139696869 (bid): serum n/a, solana 16:54:20.677327, diff n/a
slot 139696869 (ask): serum n/a, solana 16:54:20.713047, diff n/a
slot 139696870 (bid): serum 16:54:20.783392, solana 16:54:21.096202, diff -312.800953ms
slot 139696870 (ask): serum 16:54:20.859227, solana 16:54:21.13843, diff -279.195527ms
slot 139696870 (n/a): serum 16:54:20.925673, solana n/a, diff n/a
slot 139696870 (n/a): serum 16:54:20.960147, solana n/a, diff n/a
slot 139696871 (ask): serum 16:54:21.148357, solana 16:54:21.44604, diff -297.675057ms
slot 139696871 (bid): serum n/a, solana 16:54:21.480208, diff n/a
slot 139696872 (ask): serum 16:54:23.073898, solana 16:54:23.205255, diff -131.353613ms
slot 139696872 (bid): serum 16:54:23.088174, solana 16:54:23.240584, diff -152.406286ms
slot 139696873 (bid): serum n/a, solana 16:54:23.863466, diff n/a
slot 139696873 (ask): serum n/a, solana 16:54:23.897648, diff n/a
slot 139696874 (ask): serum 16:54:24.146858, solana 16:54:24.663496, diff -516.623513ms
slot 139696874 (bid): serum n/a, solana 16:54:24.694405, diff n/a
slot 139696875 (bid): serum 16:54:24.79706, solana 16:54:25.141467, diff -344.397383ms
slot 139696875 (ask): serum n/a, solana 16:54:25.175605, diff n/a
slot 139696880 (ask): serum 16:54:27.731076, solana 16:54:28.17689, diff -445.801797ms
slot 139696880 (bid): serum n/a, solana 16:54:28.21031, diff n/a
slot 139696881 (ask): serum n/a, solana 16:54:28.592455, diff n/a
slot 139696881 (bid): serum n/a, solana 16:54:28.627385, diff n/a
slot 139696882 (ask): serum n/a, solana 16:54:29.269987, diff n/a
slot 139696882 (bid): serum n/a, solana 16:54:29.272968, diff n/a
slot 139696883 (bid): serum n/a, solana 16:54:29.630217, diff n/a
slot 139696883 (ask): serum n/a, solana 16:54:29.66079, diff n/a
slot 139696884 (ask): serum n/a, solana 16:54:30.340842, diff n/a
slot 139696884 (bid): serum n/a, solana 16:54:30.342766, diff n/a
```
