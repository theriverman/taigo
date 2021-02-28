docker cp initial_test_data.json taiga-docker_taiga-back_1:/taiga-back/media/initial_test_data.json
docker-compose -f docker-compose.yml -f docker-compose-inits.yml run --rm taiga-manage loaddata /taiga-back/media/initial_test_data.json
