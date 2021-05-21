

echo $random_num

SET=$(seq 0 100)

for i in $SET

do

    echo "Running loop seq "$i
	random_num=$(cat /dev/urandom | tr -dc '0-2' | fold -w 1 | sed 1q)
	tx_num="$(($RANDOM% 50+30))"
	echo "TX_NUM ${tx_num}"
	node index.js $tx_num

	# for i in $tx_num

	# do
	# curl -w "\n" -d '{"id":"CAR0"}' -H "Content-Type: application/json" -X POST http://localhost:8000/api/buycar
	# done
	sleep $random_num
    # some instructions

done


