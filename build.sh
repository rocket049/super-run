g++ -c qtermwidget.cpp -fPIC `pkg-config --cflags qtermwidget5`
g++ -c qtermwidget.moc.cpp -fPIC `pkg-config --cflags qtermwidget5`

#go build -ldflags '-s -w'
~/go/bin/qtdeploy build

