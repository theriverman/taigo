# Wait for automatic migrations to apply after executing `docker-compose up -d`

while :
do
	curl -X GET -s --fail -o /dev/null http://localhost:9000/api/v1/stats/discover
	if [[ "$?" != '0' ]]; then
		sleep 1
		continue
	fi
	break
done

echo "Taiga is up!"
exit 0
