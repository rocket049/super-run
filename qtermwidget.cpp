#include "qtermwidget.h"
#include <QString>
#include <QList>
#include <QStringList>
#include <QObject>
#include <stdio.h>

extern "C" {
	void *createTermWidget(int startnow, void *parent);
    void termChangeDir(void *p, char *wd);
    void termSendText(void *p,char *s);
    void termSetMinimumHeight(void *p,int minh);
    char *termSelectedText(void *p);
    void termConnectFinish2Close(void *p);
    void termSetTermFont(void *p, void *f);
    void termSendKeyEvent(void *p, void *e);
}

void *createTermWidget(int startnow, void *parent)
{
    return (void*) new QTermWidget(startnow, (QWidget*)parent);
}

void termChangeDir(void *p, char *wd){
    QTermWidget *t=(QTermWidget*)p;
    t->setWorkingDirectory(*(new QString(wd)));
}

void termSendText(void *p,char *s){
    QTermWidget *t=(QTermWidget*)p;
    t->sendText(*(new QString(s)));
}

void termSetMinimumHeight(void *p,int minh){
    QTermWidget *t=(QTermWidget*)p;
    t->setMinimumHeight(minh);
}

char *termSelectedText(void *p){
    QTermWidget *t=(QTermWidget*)p;
    QString s=t->selectedText(true);
    return s.toUtf8().data();
}

void termConnectFinish2Close(void *p){
    QTermWidget *t=(QTermWidget*)p;
    QObject::connect(t,SIGNAL(finished()),t->parentWidget(),SLOT(close()));
}

void termSetTermFont(void *p, void *f) {
    QTermWidget *t=(QTermWidget*)p;
    QFont *ft1=(QFont *)f;
    t->setTerminalFont(*ft1);
}

void termSendKeyEvent(void *p, void *e) {
    QTermWidget *t=(QTermWidget*)p;
    QKeyEvent *evt=(QKeyEvent*)e;
    t->sendKeyEvent(evt);
}