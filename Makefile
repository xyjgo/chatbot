

all:
	@cd cmd/chatbot && go build && cd - &> /dev/null
	@cd cmd/simserver && go build && cd - &> /dev/null

clean:
	@cd cmd/chatbot && rm chatbot && cd - &> /dev/null
	@cd cmd/simserver && rm simserver && cd - &> /dev/null