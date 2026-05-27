package ui

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"oop/lab5/internal/funcplugins"
	"oop/lab5/internal/plugins"
	storagepkg "oop/lab5/internal/storage"
	"oop/lab5/internal/vehicle"
)

// State holds the in-memory list shown by the master view.
type State struct {
	items []vehicle.Vehicle
}

// PluginDir is the directory scanned for hierarchy plugins (*.json).
var PluginDir = "plugins"

// FuncPluginDir is the directory scanned for functional plugins.
var FuncPluginDir = "funcplugins"

// Run boots the Fyne application and blocks until the window is closed.
func Run() {
	a := app.New()
	w := a.NewWindow("Lab 5 - Vehicles + plugins + encryption")

	state := &State{}
	listLabel := widget.NewLabel("Vehicles: 0")
	pluginLabel := widget.NewLabel("")

	updatePluginLabel := func() {
		enabled := []string{}
		for _, p := range funcplugins.EnabledPlugins() {
			enabled = append(enabled, p.Name)
		}
		funcLine := "Active functional plugins: " + strings.Join(enabled, " -> ")
		if len(enabled) == 0 {
			funcLine = "Active functional plugins: (none)"
		}
		pluginLabel.SetText(fmt.Sprintf(
			"Registered types: %s\n%s",
			strings.Join(vehicle.Names(), ", "),
			funcLine,
		))
	}

	list := widget.NewList(
		func() int { return len(state.items) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(state.items[id].Summary())
		},
	)

	selected := -1
	list.OnSelected = func(id widget.ListItemID) { selected = id }
	list.OnUnselected = func(id widget.ListItemID) {
		if selected == id {
			selected = -1
		}
	}

	refresh := func() {
		listLabel.SetText("Vehicles: " + itoa(len(state.items)))
		list.Refresh()
		updatePluginLabel()
	}

	addBtn := widget.NewButton("Add", func() {
		pickType(w, func(v vehicle.Vehicle) {
			editVehicle(w, "New "+v.TypeName(), v, func() {
				state.items = append(state.items, v)
				refresh()
			})
		})
	})

	editBtn := widget.NewButton("Edit", func() {
		if selected < 0 || selected >= len(state.items) {
			dialog.ShowInformation("Edit", "Select an item first.", w)
			return
		}
		v := state.items[selected]
		editVehicle(w, "Edit "+v.TypeName(), v, refresh)
	})

	delBtn := widget.NewButton("Delete", func() {
		if selected < 0 || selected >= len(state.items) {
			dialog.ShowInformation("Delete", "Select an item first.", w)
			return
		}
		state.items = append(state.items[:selected], state.items[selected+1:]...)
		selected = -1
		list.UnselectAll()
		refresh()
	})

	saveBtn := widget.NewButton("Save...", func() {
		dlg := dialog.NewFileSave(func(f fyne.URIWriteCloser, err error) {
			if err != nil || f == nil {
				return
			}
			defer f.Close()
			var buf bytes.Buffer
			if err := storagepkg.Marshal(&buf, state.items); err != nil {
				dialog.ShowError(err, w)
				return
			}
			processed, err := funcplugins.Encode(buf.Bytes())
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if _, err := f.Write(processed); err != nil {
				dialog.ShowError(err, w)
				return
			}
			dialog.ShowInformation("Save", fmt.Sprintf(
				"Saved %d object(s) through %d active plugin(s).",
				len(state.items), len(funcplugins.EnabledPlugins()),
			), w)
		}, w)
		dlg.SetFileName("vehicles.txt")
		dlg.SetFilter(storage.NewExtensionFileFilter([]string{".txt", ".bin", ".enc"}))
		dlg.Show()
	})

	loadBtn := widget.NewButton("Load...", func() {
		dlg := dialog.NewFileOpen(func(f fyne.URIReadCloser, err error) {
			if err != nil || f == nil {
				return
			}
			defer f.Close()
			raw, err := io.ReadAll(f)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			plain, err := funcplugins.Decode(raw)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			items, err := storagepkg.Unmarshal(bytes.NewReader(plain))
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			state.items = items
			selected = -1
			list.UnselectAll()
			refresh()
		}, w)
		dlg.SetFilter(storage.NewExtensionFileFilter([]string{".txt", ".bin", ".enc"}))
		dlg.Show()
	})

	loadDefaultBtn := widget.NewButton("Load default", func() {
		const path = "vehicles.txt"
		raw, err := os.ReadFile(path)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				dialog.ShowError(err, w)
			}
			return
		}
		plain, err := funcplugins.Decode(raw)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		items, err := storagepkg.Unmarshal(bytes.NewReader(plain))
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		state.items = items
		selected = -1
		list.UnselectAll()
		refresh()
	})

	reloadHierarchyBtn := widget.NewButton("Reload hierarchy plugins", func() {
		n, errs := plugins.LoadDir(PluginDir)
		showLoadResult(w, "Hierarchy plugins", PluginDir, n, errs)
		refresh()
	})

	reloadFuncBtn := widget.NewButton("Reload functional plugins", func() {
		n, errs := funcplugins.LoadDir(FuncPluginDir)
		showLoadResult(w, "Functional plugins", FuncPluginDir, n, errs)
		refresh()
	})

	settingsBtn := widget.NewButton("Plugin settings...", func() {
		showFuncPluginSettings(w, refresh)
	})

	loadPluginBtn := widget.NewButton("Add plugin file...", func() {
		dlg := dialog.NewFileOpen(func(f fyne.URIReadCloser, err error) {
			if err != nil || f == nil {
				return
			}
			path := f.URI().Path()
			f.Close()
			if err := plugins.LoadFile(path); err != nil {
				if err2 := funcplugins.LoadFile(path); err2 != nil {
					dialog.ShowError(fmt.Errorf("not a hierarchy plugin (%v) and not a functional plugin (%v)", err, err2), w)
					return
				}
			}
			dialog.ShowInformation("Plugins", "Loaded: "+path, w)
			refresh()
		}, w)
		dlg.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
		dlg.Show()
	})

	row1 := container.NewHBox(addBtn, editBtn, delBtn, saveBtn, loadBtn, loadDefaultBtn)
	row2 := container.NewHBox(reloadHierarchyBtn, reloadFuncBtn, settingsBtn, loadPluginBtn)
	root := container.NewBorder(
		container.NewVBox(row1, row2, listLabel, pluginLabel),
		nil, nil, nil,
		container.NewScroll(list),
	)
	updatePluginLabel()
	w.SetContent(root)
	w.Resize(fyne.NewSize(820, 560))
	w.ShowAndRun()
}

// showLoadResult prints a dialog summarising a (re)load operation.
func showLoadResult(w fyne.Window, title, dir string, n int, errs []error) {
	if len(errs) == 0 {
		dialog.ShowInformation(title, fmt.Sprintf("Loaded %d plugin(s) from %s.", n, dir), w)
		return
	}
	parts := []string{fmt.Sprintf("Loaded %d plugin(s) from %s.", n, dir), "", "Errors:"}
	for _, e := range errs {
		parts = append(parts, "* "+e.Error())
	}
	dialog.ShowInformation(title, strings.Join(parts, "\n"), w)
}

// itoa is a tiny helper so we don't pull in strconv just for a couple
// of formatting calls.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
