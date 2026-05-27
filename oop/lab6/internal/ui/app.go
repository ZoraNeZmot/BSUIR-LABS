package ui

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"oop/lab6/internal/appsettings"
	"oop/lab6/internal/events"
	"oop/lab6/internal/funcplugins"
	"oop/lab6/internal/plugins"
	storagepkg "oop/lab6/internal/storage"
	"oop/lab6/internal/vehicle"
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

	// Singleton-pattern usage: the same Settings instance is consumed
	// by the title bar, by the storage layer (via LastOpenedFile), and
	// by the plugin loader.
	cfg := appsettings.Instance()
	cfg.Update(func(s *appsettings.Settings) {
		s.PluginDir = PluginDir
		s.FuncPluginDir = FuncPluginDir
	})

	// Observer-pattern usage: a single bus is published to whenever
	// the vehicle list mutates. The list widget, the audit log and the
	// title-bar dirty indicator subscribe independently.
	bus := events.NewBus()

	w := a.NewWindow(buildTitle(false))
	dirty := false
	setDirty := func(v bool) {
		if dirty == v {
			return
		}
		dirty = v
		w.SetTitle(buildTitle(dirty))
	}

	state := &State{}
	listLabel := widget.NewLabel("Vehicles: 0")
	pluginLabel := widget.NewLabel("")
	auditLog := widget.NewMultiLineEntry()
	auditLog.SetPlaceHolder("Audit log...")
	auditLog.Wrapping = fyne.TextWrapWord
	auditLog.Disable()

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

	refreshAll := func() {
		listLabel.SetText("Vehicles: " + itoa(len(state.items)))
		list.Refresh()
		updatePluginLabel()
	}

	// Three independent observers attach to the bus.
	bus.Subscribe(func(e events.Event) {
		// Observer #1: the master list & header redraw on every event.
		refreshAll()
	})
	bus.Subscribe(func(e events.Event) {
		// Observer #2: the dirty flag tracks unsaved mutations.
		switch e.Kind {
		case events.VehicleAdded, events.VehicleEdited, events.VehicleRemoved:
			setDirty(true)
		case events.SaveCompleted:
			setDirty(false)
		case events.LoadCompleted, events.ListReplaced:
			setDirty(false)
		}
	})
	bus.Subscribe(func(e events.Event) {
		// Observer #3: an append-only audit log.
		ts := time.Now().Format("15:04:05")
		line := fmt.Sprintf("[%s] %s\n", ts, describeEvent(e))
		auditLog.SetText(auditLog.Text + line)
	})

	addBtn := widget.NewButton("Add", func() {
		pickType(w, func(v vehicle.Vehicle) {
			editVehicle(w, "New "+v.TypeName(), v, func() {
				state.items = append(state.items, v)
				bus.Publish(events.Event{Kind: events.VehicleAdded, Payload: v.Summary()})
			})
		})
	})

	editBtn := widget.NewButton("Edit", func() {
		if selected < 0 || selected >= len(state.items) {
			dialog.ShowInformation("Edit", "Select an item first.", w)
			return
		}
		v := state.items[selected]
		editVehicle(w, "Edit "+v.TypeName(), v, func() {
			bus.Publish(events.Event{Kind: events.VehicleEdited, Payload: v.Summary()})
		})
	})

	delBtn := widget.NewButton("Delete", func() {
		if selected < 0 || selected >= len(state.items) {
			dialog.ShowInformation("Delete", "Select an item first.", w)
			return
		}
		removed := state.items[selected].Summary()
		state.items = append(state.items[:selected], state.items[selected+1:]...)
		selected = -1
		list.UnselectAll()
		bus.Publish(events.Event{Kind: events.VehicleRemoved, Payload: removed})
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
			cfg.Update(func(s *appsettings.Settings) { s.LastOpenedFile = f.URI().Path() })
			bus.Publish(events.Event{
				Kind:    events.SaveCompleted,
				Payload: fmt.Sprintf("%d objects -> %s", len(state.items), f.URI().Path()),
			})
		}, w)
		dlg.SetFileName("vehicles.txt")
		dlg.SetFilter(storage.NewExtensionFileFilter([]string{".txt", ".bin", ".enc"}))
		dlg.Show()
	})

	loadFromBytes := func(raw []byte, source string) {
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
		cfg.Update(func(s *appsettings.Settings) { s.LastOpenedFile = source })
		bus.Publish(events.Event{Kind: events.LoadCompleted, Payload: fmt.Sprintf("%d objects from %s", len(items), source)})
	}

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
			loadFromBytes(raw, f.URI().Path())
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
		loadFromBytes(raw, path)
	})

	reloadHierarchyBtn := widget.NewButton("Reload hierarchy plugins", func() {
		n, errs := plugins.LoadDir(PluginDir)
		showLoadResult(w, "Hierarchy plugins", PluginDir, n, errs)
		bus.Publish(events.Event{Kind: events.PluginsChanged, Payload: fmt.Sprintf("hierarchy: %d", n)})
	})

	reloadFuncBtn := widget.NewButton("Reload functional plugins", func() {
		n, errs := funcplugins.LoadDir(FuncPluginDir)
		showLoadResult(w, "Functional plugins", FuncPluginDir, n, errs)
		bus.Publish(events.Event{Kind: events.PluginsChanged, Payload: fmt.Sprintf("functional: %d", n)})
	})

	settingsBtn := widget.NewButton("Plugin settings...", func() {
		showFuncPluginSettings(w, func() {
			bus.Publish(events.Event{Kind: events.PluginsChanged, Payload: "settings updated"})
		})
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
			bus.Publish(events.Event{Kind: events.PluginsChanged, Payload: "loaded " + path})
		}, w)
		dlg.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
		dlg.Show()
	})

	row1 := container.NewHBox(addBtn, editBtn, delBtn, saveBtn, loadBtn, loadDefaultBtn)
	row2 := container.NewHBox(reloadHierarchyBtn, reloadFuncBtn, settingsBtn, loadPluginBtn)

	auditScroll := container.NewScroll(auditLog)
	auditScroll.SetMinSize(fyne.NewSize(0, 120))

	mainPane := container.NewBorder(nil, nil, nil, nil, container.NewScroll(list))
	root := container.NewBorder(
		container.NewVBox(row1, row2, listLabel, pluginLabel),
		container.NewVBox(widget.NewSeparator(), widget.NewLabel("Audit log (Observer pattern)"), auditScroll),
		nil, nil,
		mainPane,
	)
	updatePluginLabel()
	w.SetContent(root)
	w.Resize(fyne.NewSize(900, 700))
	w.ShowAndRun()
}

// describeEvent renders an Event into a single audit-log line. The
// switch is on the closed Kind enum, so it is fine.
func describeEvent(e events.Event) string {
	switch e.Kind {
	case events.VehicleAdded:
		return "+ Added " + asString(e.Payload)
	case events.VehicleEdited:
		return "~ Edited " + asString(e.Payload)
	case events.VehicleRemoved:
		return "- Removed " + asString(e.Payload)
	case events.ListReplaced:
		return "= List replaced (" + asString(e.Payload) + ")"
	case events.PluginsChanged:
		return "P Plugins: " + asString(e.Payload)
	case events.SaveCompleted:
		return "S Saved " + asString(e.Payload)
	case events.LoadCompleted:
		return "L Loaded " + asString(e.Payload)
	}
	return "? unknown event"
}

func asString(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

// buildTitle uses the singleton's user name and exposes the dirty flag
// in the window title -- a tiny demonstration that the singleton lives
// throughout the lifetime of the program.
func buildTitle(dirty bool) string {
	user := appsettings.Instance().Get().UserName
	mark := ""
	if dirty {
		mark = " *"
	}
	return fmt.Sprintf("Lab 6 - Vehicles + plugins + patterns (%s)%s", user, mark)
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
