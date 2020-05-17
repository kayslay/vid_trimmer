run:
	PORT=9000 go run main.go

docker:
	docker build  -t vid_trimmer .

docker_run:
# run the container in the host network so as to use the same db the local machine is using
	docker run --rm  --env-file .env --network host vid_trimmer


