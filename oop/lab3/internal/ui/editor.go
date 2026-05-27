// Package ui contains the Fyne front-end. Both the master list and the
// editor dialog are built generically from FieldDescriptor data so that
// adding a new vehicle class never requires UI changes.
package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"oop/lab3/internal/vehicle"
)

// buildForm constructs a *widget.Form that mirrors the FieldDescriptors
// of the supplied vehicle. Each entry's submit handler updates the
// underlying field through the descriptor's Set closure, so the dialog
// stays type-agnostic.
func buildForm(v vehicle.Vehicle) (*widget.Form, []func() error) {
	form := widget.NewForm()
	commits := []func() error{}
	for _, f := range v.Fields() {
		// `f` is captured by closure; copy it explicitly to avoid
		// referencing the loop variable after the iteration ends.
		field := f
		switch field.Kind {
		case vehicle.KindBool:
			chk := widget.NewCheck("", nil)
			chk.SetChecked(field.Get() == "true")
			form.Append(field.Label, chk)
			commits = append(commits, func() error {
				if chk.Checked {
					return field.Set("true")
				}
				return field.Set("false")
			})
		default:
			entry := widget.NewEntry()
			entry.SetText(field.Get())
			form.Append(field.Label, entry)
			commits = append(commits, func() error {
				return field.Set(entry.Text)
			})
		}
	}
	return form, commits
}

// editVehicle pops up a modal dialog allowing the user to edit every
// field of v. onSave is invoked only after every commit closure returns
// without error.
func editVehicle(parent fyne.Window, title string, v vehicle.Vehicle, onSave func()) {
	form, commits := buildForm(v)
	d := dialog.NewCustomConfirm(title, "Save", "Cancel", form, func(ok bool) {
		if !ok {
			return
		}
		for _, c := range commits {
			if err := c(); err != nil {
				dialog.ShowError(err, parent)
				return
			}
		}
		onSave()
	}, parent)
	d.Resize(fyne.NewSize(420, 480))
	d.Show()
}

// pickType asks the user to choose a registered vehicle type, then
// hands the freshly-created instance to onPicked.
func pickType(parent fyne.Window, onPicked func(vehicle.Vehicle)) {
	names := vehicle.Names()
	if len(names) == 0 {
		dialog.ShowInformation("No types", "No vehicle types registered.", parent)
		return
	}
	sel := widget.NewSelect(names, nil)
	sel.SetSelected(names[0])
	dialog.NewCustomConfirm("New vehicle", "Create", "Cancel",
		container.NewVBox(widget.NewLabel("Select type:"), sel),
		func(ok bool) {
			if !ok {
				return
			}
			v, err := vehicle.Create(sel.Selected)
			if err != nil {
				dialog.ShowError(err, parent)
				return
			}
			onPicked(v)
		}, parent).Show()
}
