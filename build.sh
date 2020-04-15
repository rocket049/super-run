g++ -c qtermwidget.cpp -fPIC `pkg-config --cflags qtermwidget5`
go build -ldflags '-s -w -r /home/fuhz/src/qterm:/home/fuhz/Qt514/5.14.1/gcc_64/lib'

