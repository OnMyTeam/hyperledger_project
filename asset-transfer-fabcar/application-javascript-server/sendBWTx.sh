


total_tx_num=1000
tx_num=0
max_num=50
min_num=30
while [ $tx_num -le $total_tx_num ]
do
	
    
	random_num=$(cat /dev/urandom | tr -dc '1-2' | fold -w 1 | sed 1q)
	send_tx_num=$(($RANDOM% $max_num+$min_num))
	remain=$((total_tx_num-tx_num))
	if [ ${remain} -le $((max_num+min_num)) ] && [ ${remain} -ge $min_num ]; then
		send_tx_num=$remain
	fi
	tx_num=$((tx_num+send_tx_num))
	
	# # sendBWTransaction
	node index.js $send_tx_num
	
	# curl -w "\n" \
	# 	-d '{"id":"CAR0", "txnum":'${send_tx_num}'}' \
	# 	-H "Content-Type: application/json" \
	# 	-X POST http://localhost:8000/api/buycar
	sleep $random_num
	echo "send Tx_num ${send_tx_num}"
	echo "Total Tx_num ${tx_num}"
	if [ $tx_num -eq $total_tx_num ]; then
		tx_num=$((tx_num+1))
	fi
	
    # some instructions

done

