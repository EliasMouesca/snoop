EXEC=snoop

$(EXEC): *.go
	go build -o $@ $^

