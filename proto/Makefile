shortener:
	shortener-server
	shortener-telegram

shortener-server:
	protoc shortener/*.proto --go_out=../shortener-service/proto --go_opt=paths=source_relative --go-grpc_out=../shortener-service/proto --go-grpc_opt=paths=source_relative

shortener-telegram:
	protoc shortener/*.proto --go_out=../telegram-bot-service/proto --go_opt=paths=source_relative --go-grpc_out=../telegram-bot-service/proto --go-grpc_opt=paths=source_relative