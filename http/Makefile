.PHONY: all

TARGET=chttp

all: $(TARGET)

$(TARGET): main.go
	go build -o $@ ./

clean:
	rm -fv $(TARGET)
