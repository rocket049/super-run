g++ -c qtermwidget.cpp -fPIC `pkg-config --cflags qtermwidget5`
go build -ldflags '-s -w'

