package main

/*
#cgo CXXFLAGS: -I/usr/local/include -I/usr/local/include/qtermwidget5 -I/home/fuhz/Qt514/5.14.1/gcc_64/include/QtWidgets -I/home/fuhz/Qt514/5.14.1/gcc_64/include -I/home/fuhz/Qt514/5.14.1/gcc_64/include/QtGui -I/home/fuhz/Qt514/5.14.1/gcc_64/include/QtCore
#cgo LDFLAGS: -L/home/fuhz/src/qterm -lqtermwidget5 -L/home/fuhz/Qt514/5.14.1/gcc_64/lib -lQt5Widgets -lQt5Gui -lQt5Core -lQt5Quick -lQt5Designer -lQt5Qml -lQt5Multimedia -lQt5Network -lQt5Xml -lQt5DBus -lQt5RemoteObjects
#include <stdlib.h>

extern void* createTermWidget(int startnow, void * parent);
extern void termSendText(void *p,char *s);
extern void termSetMinimumHeight(void *p,int minh);
extern char *termSelectedText(void *p);
*/
import "C"

import (
	"strings"
	"unsafe"
)

func getQTermPtr() uintptr {
	t := C.createTermWidget(1, nil)
	return uintptr(t)
}

func termChangeDir(p uintptr, d string) {
	arg := strings.Replace(d, "'", "\\'", 0)
	termSendText(p, "cd '"+arg+"'\n")
}

func termSendText(p uintptr, s string) {
	C.termSendText(unsafe.Pointer(p), C.CString(s))
}

func termSetMiniHeight(p uintptr, h int) {
	C.termSetMinimumHeight(unsafe.Pointer(p), C.int(h))
}

func termSelectedText(p uintptr) string {
	return C.GoString(C.termSelectedText(unsafe.Pointer(p)))
}

func buildCmdLine(prog string, envs, args []string) string {
	var cmd = []string{}
	for _, v := range envs {
		cmd = append(cmd, v)
	}
	cmd = append(cmd, prog)
	for _, v := range args {
		arg := strings.Replace(v, "'", "\\'", 0)
		cmd = append(cmd, "'"+arg+"'")
	}
	return strings.Join(cmd, " \\\n") + "\n"
}
