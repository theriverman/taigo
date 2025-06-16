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

# Dump the running compose services
docker compose ps

# Copy initial_test_data.json into /taiga-back/media through the `taiga-docker-stable-taiga-back-1` container
docker cp initial_test_data.json taiga-docker-stable-taiga-back-1:/taiga-back/media/initial_test_data.json  || exit 1

# Move into the taiga-docker submodule's folder
cd taiga-docker || exit 1

# Load the default user + default project from the `initial_test_data` fixture
docker compose -f docker-compose.yml -f docker-compose-inits.yml run --rm taiga-manage loaddata /taiga-back/media/initial_test_data.json || exit 1
