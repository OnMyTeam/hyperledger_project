

echo $random_num

SET=$(seq 0 100)

for i in $SET

do

    echo "Running loop seq "$i
	random_num=$(cat /dev/urandom | tr -dc '0-2' | fold -w 1 | sed 1q)
	tx_num=$(cat /dev/urandom | tr -dc '30-100' | fold -w 1 | sed 1q)
	# node index.js $tx_num
	node invoke.js $tx_num
	sleep $random_num
    # some instructions

done


