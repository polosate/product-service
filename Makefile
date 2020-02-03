build:
	echo "Build proto file"
	protoc -I. --go_out=plugins=micro:. \
		  proto/product/product.proto

	echo "Build docker image"
	docker build -t product-service .

run:
	echo "Run docker image"
	docker run -p 50051:50051 -e MICRO_SERVER_ADDRESS=:50051 product-service
