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

docker cp initial_test_data.json taiga-docker_taiga-back_1:/taiga-back/media/initial_test_data.json  || exit 1
docker-compose -f docker-compose.yml -f docker-compose-inits.yml run --rm taiga-manage loaddata /taiga-back/media/initial_test_data.json || exit 1
