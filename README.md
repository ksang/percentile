#percentile
A command line tool to fetch x% percentile sample value from prometheus time series query

###Usage:

	  -duration string
	    	duration, e.g 2s, 3m, 4h.
	  -f string
	    	full URI to prometheus series api or result file output
	  -p int
	    	percentile to get, a int number in range 0-100 (default -1)
	  -prom string
	    	promethues host url, e.g http://1.1.1.1:9090
	  -query string
	    	query used for query_range, without time series
	  -start string
	    	start timestamp, int RFC3339 or unix timestamp
	  -step string
	    	step, e.g 10s, 15s, 1m (default "30s")
	  -t	print table result, including max/avg/min
	  -v	verbose output



###Example:

	percentile -p 95 -prom http://127.0.0.1:9090 \
		-query "sum(rate(node_network_receive_bytes[2m]))by(node)" \
		-start 2016-09-26T00:00:00+08:00 -duration 10h -step 5m -t

###Output:
	+-------------+-------------+-------------+-------------+----------+
	|     95%     |     MAX     |     MIN     |     AVG     |  LABELS  |
	+-------------+-------------+-------------+-------------+----------+
	|    49630.28 |    51362.50 |    44744.68 |    47127.35 | {node3}  |
	|   159774.87 |   166482.98 |   154748.27 |   158583.75 | {node1}  |
	|   172777.18 |   176971.40 |   167431.70 |   170969.79 | {node5}  |
	|   186383.87 |   251204.15 |   180375.62 |   185220.67 | {node2}  |
	|   439206.48 |   480525.60 |   293597.97 |   367577.75 | {node4}  |
	|    60933.70 |    61621.82 |    56347.10 |    59338.97 | {node0}  |
	+-------------+-------------+-------------+-------------+----------+