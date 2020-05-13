package main

/*
#cgo pkg-config: qtermwidget5 Qt5Widgets Qt5Gui Qt5Core
#include <stdlib.h>

extern void* createTermWidget(int startnow, void * parent);
extern void termSendText(void *p,char *s);
extern void termSetMinimumHeight(void *p,int minh);
extern char *termSelectedText(void *p);
extern void termConnectFinish2Close(void *p);
*/
import "C"

import (
	"strings"
	"unsafe"
)

func getQTermPtr(p unsafe.Pointer) uintptr {
	t := C.createTermWidget(1, p)
	return uintptr(t)
}

func termChangeDir(p uintptr, d string) {
	arg := strings.Replace(d, "\"", "\\\"", 0)
	termSendText(p, "cd \""+arg+"\"\n")
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

func termConnectFinish2Close(p uintptr) {
	C.termConnectFinish2Close(unsafe.Pointer(p))
}

func buildCmdLine(prog string, envs, args []string) string {
	var cmd = []string{}
	for _, v := range envs {
		cmd = append(cmd, v)
	}
	cmd = append(cmd, prog)
	for _, v := range args {
		arg := strings.Replace(v, "\"", "\\\"", 0)
		cmd = append(cmd, "\""+arg+"\"")
	}
	return strings.Join(cmd, " \\\n") + "\n"
}
