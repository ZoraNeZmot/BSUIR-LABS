package ui

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"oop/lab4/internal/plugins"
	storagepkg "oop/lab4/internal/storage"
	"oop/lab4/internal/vehicle"
)

// State holds the in-memory list shown by the master view.
type State struct {
	items []vehicle.Vehicle
}

// PluginDir is the relative directory scanned for *.json plugin files.
// It is exported so the main module can override it on the command line.
var PluginDir = "plugins"

// Run boots the Fyne application and blocks until the window is closed.
func Run() {
	a := app.New()
	w := a.NewWindow("Lab 4 — Vehicles + plugins (text serialization)")

	state := &State{}
	listLabel := widget.NewLabel("Vehicles: 0")
	pluginLabel := widget.NewLabel("")

	updatePluginLabel := func() {
		pluginLabel.SetText(fmt.Sprintf("Registered types: %s", strings.Join(vehicle.Names(), ", ")))
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

	saveBtn := widget.NewButton("Save…", func() {
		dlg := dialog.NewFileSave(func(f fyne.URIWriteCloser, err error) {
			if err != nil || f == nil {
				return
			}
			defer f.Close()
			if err := storagepkg.Marshal(f, state.items); err != nil {
				dialog.ShowError(err, w)
				return
			}
			dialog.ShowInformation("Save", "Saved "+itoa(len(state.items))+" object(s).", w)
		}, w)
		dlg.SetFileName("vehicles.txt")
		dlg.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
		dlg.Show()
	})

	loadBtn := widget.NewButton("Load…", func() {
		dlg := dialog.NewFileOpen(func(f fyne.URIReadCloser, err error) {
			if err != nil || f == nil {
				return
			}
			defer f.Close()
			items, err := storagepkg.Unmarshal(f)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			state.items = items
			selected = -1
			list.UnselectAll()
			refresh()
		}, w)
		dlg.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
		dlg.Show()
	})

	loadDefaultBtn := widget.NewButton("Load default", func() {
		const path = "vehicles.txt"
		f, err := os.Open(path)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				dialog.ShowError(err, w)
			}
			return
		}
		defer f.Close()
		items, err := storagepkg.Unmarshal(f)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		state.items = items
		selected = -1
		list.UnselectAll()
		refresh()
	})

	reloadBtn := widget.NewButton("Reload plugins", func() {
		n, errs := plugins.LoadDir(PluginDir)
		msg := fmt.Sprintf("Loaded %d plugin(s) from %s.", n, PluginDir)
		if len(errs) > 0 {
			parts := []string{msg, "", "Errors:"}
			for _, e := range errs {
				parts = append(parts, "• "+e.Error())
			}
			dialog.ShowInformation("Plugins", strings.Join(parts, "\n"), w)
		} else {
			dialog.ShowInformation("Plugins", msg, w)
		}
		refresh()
	})

	loadPluginBtn := widget.NewButton("Add plugin file…", func() {
		dlg := dialog.NewFileOpen(func(f fyne.URIReadCloser, err error) {
			if err != nil || f == nil {
				return
			}
			path := f.URI().Path()
			f.Close()
			if err := plugins.LoadFile(path); err != nil {
				dialog.ShowError(err, w)
				return
			}
			dialog.ShowInformation("Plugins", "Loaded plugin: "+path, w)
			refresh()
		}, w)
		dlg.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
		dlg.Show()
	})

	toolbar := container.NewHBox(addBtn, editBtn, delBtn, saveBtn, loadBtn, loadDefaultBtn)
	pluginBar := container.NewHBox(reloadBtn, loadPluginBtn)
	root := container.NewBorder(
		container.NewVBox(toolbar, pluginBar, listLabel, pluginLabel),
		nil, nil, nil,
		container.NewScroll(list),
	)
	updatePluginLabel()
	w.SetContent(root)
	w.Resize(fyne.NewSize(720, 520))
	w.ShowAndRun()
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
