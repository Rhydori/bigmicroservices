gate:
	@cd cmd/gateway && go run .

login:
	@cd cmd/loginserver && go run .

chat: 
	@cd cmd/chatserver && go run .

gc:
	@cd cmd/gateway && go run . &
	@cd cmd/chatserver && go run .

all:
	@cd cmd/gateway && go run . &
	@cd cmd/loginserver && go run . &
	@cd cmd/chatserver && go run .