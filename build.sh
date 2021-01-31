#!/bin/sh
g++ -c qtermwidget.cpp -fPIC `pkg-config --cflags qtermwidget5`
g++ -c qtermwidget.moc.cpp -fPIC `pkg-config --cflags qtermwidget5`

go build -tags minimal -ldflags -s
#~/go/bin/qtdeploy build

