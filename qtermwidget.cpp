#include "qtermwidget.h"
#include <QString>
#include <QList>
#include <QStringList>
#include <stdio.h>

extern "C" {
	void *createTermWidget(int startnow, void *parent);
    void termChangeDir(void *p, char *wd);
    void termSendText(void *p,char *s);
    void termSetMinimumHeight(void *p,int minh);
    char *termSelectedText(void *p);
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