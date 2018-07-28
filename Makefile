linux_amd64 := GOOS=linux GOARCH=amd64
gobuild := go build

BIN := dist/bot

$(BIN):
	@$(gobuild) -o $@ .
	@chmod +x $@

clean:
	@rm -f $(BIN)
