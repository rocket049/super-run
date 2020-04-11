// main.go
package main

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/therecipe/qt/gui"

	"github.com/rocket049/gettext-go/gettext"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func init() {
	exe1, _ := os.Executable()
	dir1 := path.Dir(exe1)
	locale1 := path.Join(dir1, "locale")
	gettext.BindTextdomain("super-run", locale1, nil)
	gettext.Textdomain("super-run")
}

var T = gettext.T

func getAppPath() string {
	exe1, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(exe1)
}

type MyApp struct {
	app     *widgets.QApplication
	window  *widgets.QMainWindow
	list    *widgets.QListWidget
	console *widgets.QTextEdit
}

func (a *MyApp) Run() {
	a.app = widgets.NewQApplication(len(os.Args), os.Args)
	a.window = widgets.NewQMainWindow(nil, core.Qt__Window)
	a.window.SetWindowTitle(T("Super Command"))
	a.window.SetFixedSize2(800, 600)
	a.createGui()
	a.app.SetActiveWindow(a.window)
	a.window.Show()
	a.app.Exec()
}

func (a *MyApp) setIcon() {
	icon := gui.NewQIcon5(filepath.Join(getAppPath(), "icon.png"))
	a.window.SetWindowIcon(icon)
}

func (a *MyApp) createGui() {
	parent := a.window
	a.setIcon()
	spliter1 := widgets.NewQSplitter2(core.Qt__Horizontal, parent)
	spliter1.SetSizePolicy2(widgets.QSizePolicy__Expanding, widgets.QSizePolicy__Expanding)

	a.list = widgets.NewQListWidget(spliter1)
	a.list.SetSizePolicy2(widgets.QSizePolicy__Preferred, widgets.QSizePolicy__Expanding)
	spliter1.AddWidget(a.list)
	a.list.SetSelectionMode(widgets.QAbstractItemView__SingleSelection)
	a.fillList()

	a.console = widgets.NewQTextEdit(spliter1)
	a.console.SetMinimumWidth(600)
	a.console.SetSizePolicy2(widgets.QSizePolicy__Expanding, widgets.QSizePolicy__Expanding)
	spliter1.AddWidget(a.console)
	a.console.SetReadOnly(true)

	parent.SetCentralWidget(spliter1)
}

func (a *MyApp) fillList() {
	list1 := a.list
	list1.SetToolTip(T("Double Click To Run, Single Click To Read JSON"))
	cfgDir, err := os.Open(filepath.Join(getAppPath(), "conf.d"))
	if err != nil {
		panic(err)
	}
	defer cfgDir.Close()
	names, err := cfgDir.Readdirnames(0)
	if err != nil {
		panic(err)
	}
	sort.Strings(names)
	list1.AddItems(names)
	list1.ConnectSelectionChanged(func(sel *core.QItemSelection, desel *core.QItemSelection) {
		item1 := list1.Item(sel.Indexes()[0].Row())
		fn := item1.Data(0).ToString()
		cfg := filepath.Join(getAppPath(), "conf.d", fn)
		data, err := ioutil.ReadFile(cfg)
		if err != nil {
			a.console.SetText(fn + "\n" + err.Error())
		} else {
			a.console.SetText(string(data))
		}
	})
	list1.ConnectItemDoubleClicked(func(item *widgets.QListWidgetItem) {
		fn := item.Data(0).ToString()
		cfg := filepath.Join(getAppPath(), "conf.d", fn)
		cfgWin, err := readJsonCmd(cfg)
		if err != nil {
			panic(err)
		}
		a.showCmdWin(cfgWin, fn)
	})
}

func (a *MyApp) showCmdWin(cfg *JsonCmd, filename string) {
	//cmd for run
	var jcmd JsonCmd
	jcmd.Command = cfg.Command
	jcmd.Envs = cfg.Envs
	jcmd.PreArgs = cfg.PreArgs

	jcmd.Dirs = []string{}
	jcmd.Files = []string{}
	jcmd.OptDirs = [][]string{}
	jcmd.OptFiles = [][]string{}
	jcmd.Texts = []string{}

	optMap := make(map[string]string)
	argArray := []*widgets.QLineEdit{}

	const entryWidth = 500

	savedCfg, err := a.loadSavedConf(filename)
	if err == nil {
		cfg.WorkDir = savedCfg.WorkDir
	}

	dialog := widgets.NewQDialog(a.window, core.Qt__Dialog)
	dialog.SetWindowTitle(cfg.Title)
	layout := widgets.NewQVBoxLayout()

	//提交运行时赋值给 jcmd.WorkDir
	wdEntry := widgets.NewQLineEdit2(cfg.WorkDir, dialog)
	wdEntry.SetMinimumWidth(entryWidth)
	wd := a.createLine(T("Work Dir"), wdEntry, dialog)
	wdEntry.SetReadOnly(true)
	wdEntry.SetToolTip(T("Double Click To Select Path"))
	wdEntry.SetPlaceholderText(T("Double Click To Select Path"))
	wdEntry.ConnectMouseDoubleClickEvent(func(e *gui.QMouseEvent) {
		home, _ := os.UserHomeDir()
		dir1 := widgets.QFileDialog_GetExistingDirectory(dialog, T("Work Dir"), home, widgets.QFileDialog__ShowDirsOnly)
		wdEntry.SetText(dir1)
	})
	layout.AddWidget(wd, 1, 0)

	for _, v := range cfg.OptDirs {
		name := v[0]
		opt := v[1]
		entry := widgets.NewQLineEdit(dialog)
		entry.SetPlaceholderText(T("Double Click To Select Path"))
		entry.SetMinimumWidth(entryWidth)
		entry.SetReadOnly(true)
		line := a.createLine(name, entry, dialog)
		entry.ConnectMouseDoubleClickEvent(func(e *gui.QMouseEvent) {
			home, _ := os.UserHomeDir()
			dir1 := widgets.QFileDialog_GetExistingDirectory(dialog, name, home, widgets.QFileDialog__ShowDirsOnly)
			entry.SetText(dir1)
			optMap[opt] = dir1
		})
		layout.AddWidget(line, 1, 0)
	}

	for _, v := range cfg.OptFiles {
		name := v[0]
		opt := v[1]
		entry := widgets.NewQLineEdit(dialog)
		entry.SetPlaceholderText(T("Double Click To Select Path"))
		entry.SetMinimumWidth(entryWidth)
		entry.SetReadOnly(true)
		line := a.createLine(name, entry, dialog)
		entry.ConnectMouseDoubleClickEvent(func(e *gui.QMouseEvent) {
			home, _ := os.UserHomeDir()
			path1 := widgets.QFileDialog_GetOpenFileName(dialog, name, home, "*", "*", widgets.QFileDialog__ReadOnly)
			entry.SetText(path1)
			optMap[opt] = path1
		})
		layout.AddWidget(line, 1, 0)
	}
	for _, v := range cfg.OptTexts {
		name := v[0]
		opt := v[1]
		entry := widgets.NewQLineEdit(dialog)
		entry.SetMinimumWidth(entryWidth)
		line := a.createLine(name, entry, dialog)
		entry.ConnectEditingFinished(func() {
			optMap[opt] = entry.Text()
		})
		layout.AddWidget(line, 1, 0)
	}

	for _, v := range cfg.Dirs {
		name := v
		entry := widgets.NewQLineEdit(dialog)
		entry.SetPlaceholderText(T("Double Click To Select Path"))
		entry.SetMinimumWidth(entryWidth)
		entry.SetReadOnly(true)
		line := a.createLine(name, entry, dialog)
		entry.ConnectMouseDoubleClickEvent(func(e *gui.QMouseEvent) {
			home, _ := os.UserHomeDir()
			path1 := widgets.QFileDialog_GetExistingDirectory(dialog, name, home, widgets.QFileDialog__ShowDirsOnly)
			entry.SetText(path1)
		})
		layout.AddWidget(line, 1, 0)
		argArray = append(argArray, entry)
	}

	for _, v := range cfg.Files {
		name := v
		entry := widgets.NewQLineEdit(dialog)
		entry.SetPlaceholderText(T("Double Click To Select Path"))
		entry.SetMinimumWidth(entryWidth)
		entry.SetReadOnly(true)
		line := a.createLine(name, entry, dialog)
		entry.ConnectMouseDoubleClickEvent(func(e *gui.QMouseEvent) {
			home, _ := os.UserHomeDir()
			path1 := widgets.QFileDialog_GetOpenFileName(dialog, name, home, "*", "*", widgets.QFileDialog__ReadOnly)
			entry.SetText(path1)
		})
		layout.AddWidget(line, 1, 0)
		argArray = append(argArray, entry)
	}

	for _, v := range cfg.Texts {
		name := v
		entry := widgets.NewQLineEdit(dialog)
		entry.SetMinimumWidth(entryWidth)
		entry.SetReadOnly(false)
		line := a.createLine(name, entry, dialog)
		layout.AddWidget(line, 1, 0)
		argArray = append(argArray, entry)
	}

	btRun := widgets.NewQPushButton2(T("Run"), dialog)
	layout.AddWidget(btRun, 1, 0)

	output := widgets.NewQTextEdit(dialog)
	layout.AddWidget(output, 1, 0)
	output.SetReadOnly(true)
	output.SetMinimumHeight(200)

	input := widgets.NewQLineEdit(dialog)
	layout.AddWidget(input, 1, 0)
	input.SetReadOnly(false)

	btRun.ConnectClicked(func(b bool) {
		jcmd.OptDirs = [][]string{}
		for k, v := range optMap {
			jcmd.OptDirs = append(jcmd.OptDirs, []string{k, v})
		}
		jcmd.Dirs = []string{}
		for _, v := range argArray {
			jcmd.Dirs = append(jcmd.Dirs, v.Text())
		}
		jcmd.WorkDir = wdEntry.Text()
		pIn, pOut, pErr, err := runJsonCmd(&jcmd)
		if err != nil {
			output.SetText(err.Error())
			return
		}
		a.saveConf(filename, &jcmd)
		go a.controlDialog(pIn, pOut, pErr, output, input)
	})

	dialog.SetLayout(layout)
	dialog.SetModal(true)
	dialog.Show()
}

func (a *MyApp) saveConf(cfgName string, jcmd *JsonCmd) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	cfgDir := filepath.Join(home, ".config", "super-run")
	os.MkdirAll(cfgDir, os.ModePerm)
	cfgPath := filepath.Join(cfgDir, cfgName)
	data, err := json.Marshal(jcmd)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cfgPath, data, 0644)
}

func (a *MyApp) loadSavedConf(cfgName string) (*JsonCmd, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	cfgPath := filepath.Join(home, ".config", "super-run", cfgName)
	data, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}
	res := new(JsonCmd)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *MyApp) controlDialog(pIn io.ReadCloser, pOut io.WriteCloser, pErr io.ReadCloser,
	output *widgets.QTextEdit, input *widgets.QLineEdit) {
	defer pIn.Close()
	defer pOut.Close()
	defer pErr.Close()
	input.ConnectEditingFinished(func() {
		pOut.Write([]byte(input.Text() + "\n"))
	})

	mrd := io.MultiReader(pIn, pErr)
	brd := bufio.NewReader(mrd)
	for {
		line, _, err := brd.ReadLine()
		if err != nil {
			break
		}
		output.Append(string(line))
	}
}

func (a *MyApp) createLine(name string, entry *widgets.QLineEdit, parent widgets.QWidget_ITF) *widgets.QWidget {
	res := widgets.NewQWidget(parent, core.Qt__Widget)
	res.SetContentsMargins(0, 0, 0, 0)
	layout := widgets.NewQHBoxLayout()
	layout.SetContentsMargins(0, 0, 0, 0)
	layout.AddWidget(widgets.NewQLabel2(name, res, core.Qt__Widget), 1, 0)
	layout.AddWidget(entry, 1, 0)

	res.SetLayout(layout)
	return res
}

func main() {
	app := new(MyApp)
	app.Run()
}
