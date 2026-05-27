package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"oop/lab6/internal/funcplugins"
)

// showFuncPluginSettings displays the settings menu where the user can
// enable/disable functional plugins and edit their parameters. The
// dialog is fully generic: every parameter editor is built from the
// algorithm's ParamSpec so adding new algorithms or new parameter kinds
// requires no changes here.
func showFuncPluginSettings(parent fyne.Window, onChange func()) {
	pls := funcplugins.Plugins()
	if len(pls) == 0 {
		dialog.ShowInformation(
			"Plugin settings",
			"No functional plugins loaded.\nDrop *.json descriptors into the funcplugins/ directory.",
			parent,
		)
		return
	}

	rows := []fyne.CanvasObject{}
	for _, p := range pls {
		plugin := p
		algo, _ := funcplugins.LookupAlgorithm(plugin.AlgorithmID)

		enabled := widget.NewCheck("", func(v bool) { plugin.Enabled = v })
		enabled.SetChecked(plugin.Enabled)

		name := widget.NewLabel(fmt.Sprintf("%s  [%s]", plugin.Name, plugin.AlgorithmID))
		name.TextStyle = fyne.TextStyle{Bold: true}

		desc := plugin.Description
		if desc == "" && algo != nil {
			desc = algo.Description()
		}
		descLabel := widget.NewLabel(desc)
		descLabel.Wrapping = fyne.TextWrapWord

		paramsBtn := widget.NewButton("Parameters...", func() {
			editParameters(parent, plugin, onChange)
		})

		rows = append(rows,
			container.NewBorder(
				nil, nil,
				container.NewHBox(enabled, name),
				paramsBtn,
				descLabel,
			),
			widget.NewSeparator(),
		)
	}

	body := container.NewVBox(rows...)
	d := dialog.NewCustom("Functional plugin settings", "Close",
		container.NewScroll(body), parent)
	d.Resize(fyne.NewSize(640, 480))
	d.SetOnClosed(func() {
		if onChange != nil {
			onChange()
		}
	})
	d.Show()
}

// editParameters opens a small modal dialog with one entry per
// parameter declared by the plugin's algorithm.
func editParameters(parent fyne.Window, p *funcplugins.FuncPlugin, onChange func()) {
	algo, ok := funcplugins.LookupAlgorithm(p.AlgorithmID)
	if !ok {
		dialog.ShowError(fmt.Errorf("algorithm %q is no longer registered", p.AlgorithmID), parent)
		return
	}
	specs := algo.Parameters()
	if len(specs) == 0 {
		dialog.ShowInformation("Parameters",
			fmt.Sprintf("Algorithm %q does not expose any parameters.", algo.DisplayName()), parent)
		return
	}

	form := widget.NewForm()
	commits := []func(){}
	for _, s := range specs {
		spec := s
		entry := paramEntryFor(spec.Kind)
		current := p.Parameters[spec.Name]
		if current == "" {
			current = spec.Default
		}
		entry.SetText(current)
		form.Append(spec.Label, entry)
		commits = append(commits, func() {
			p.Parameters[spec.Name] = entry.Text
		})
	}

	d := dialog.NewCustomConfirm(
		fmt.Sprintf("%s parameters", p.Name),
		"Apply", "Cancel", form,
		func(ok bool) {
			if !ok {
				return
			}
			for _, c := range commits {
				c()
			}
			if onChange != nil {
				onChange()
			}
		}, parent)
	d.Resize(fyne.NewSize(420, 320))
	d.Show()
}

// paramEntryFor picks an Entry widget appropriate for the parameter
// kind. The mapping is closed (only four primitive kinds exist), so
// adding a new functional plugin never requires changing this code.
func paramEntryFor(kind funcplugins.ParamKind) *widget.Entry {
	switch kind {
	case funcplugins.ParamSecret:
		return widget.NewPasswordEntry()
	default:
		return widget.NewEntry()
	}
}
