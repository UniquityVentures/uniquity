package p_uniquity_accounting

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// journalTransferPickerIncludeClosestForm adds hx-include so the account picker GET includes the enclosing form values (e.g. the other account id).
func journalTransferPickerIncludeClosestForm(context.Context) (Node, error) {
	return Attr("hx-include", "closest form"), nil
}

// journalAccountTransferForeignKey is [components.InputForeignKey] for Account plus optional getter-rendered nodes on the picker trigger (same idea as InputForeignKey.Attr on the hidden input).
// It avoids extending lamu for journal transfer–specific htmx behavior.
// [components.Page] is embedded at this level so [components.GetRequiredRoles] can find field "Page" (see lamu components/page.go).
type journalAccountTransferForeignKey struct {
	components.Page
	FK components.InputForeignKey[Account]
	// Attr is merged onto the picker div (hx-get target), e.g. [journalTransferPickerIncludeClosestForm].
	Attr getters.Getter[Node]
}

func (e journalAccountTransferForeignKey) GetKey() string {
	return e.FK.GetKey()
}

func (e journalAccountTransferForeignKey) GetRoles() []string {
	return e.FK.GetRoles()
}

func (e journalAccountTransferForeignKey) GetName() string {
	return e.FK.GetName()
}

func (e journalAccountTransferForeignKey) Parse(v any, ctx context.Context) (any, error) {
	return e.FK.Parse(v, ctx)
}

func (e journalAccountTransferForeignKey) Build(ctx context.Context) Node {
	if e.FK.Hidden {
		return e.FK.Build(ctx)
	}

	inner := e.FK
	valuePk := ""
	displayValue := ""

	if inner.Getter != nil {
		value, err := inner.Getter(ctx)
		if err != nil {
			slog.Error("journalAccountTransferForeignKey getter failed", "error", err, "key", inner.Key)
		} else {
			valueMap := getters.MapFromStruct(value)
			if len(valueMap) > 0 {
				haveSelectedID := false
				if idVal, exists := valueMap["ID"]; exists {
					if rv := reflect.ValueOf(idVal); rv.IsValid() && !rv.IsZero() {
						valuePk = fmt.Sprintf("%v", idVal)
						haveSelectedID = true
					}
				} else if idVal, exists := valueMap["id"]; exists {
					if rv := reflect.ValueOf(idVal); rv.IsValid() && !rv.IsZero() {
						valuePk = fmt.Sprintf("%v", idVal)
						haveSelectedID = true
					}
				}
				if inner.Display != nil && haveSelectedID {
					displayStr, err := inner.Display(context.WithValue(ctx, "$in", valueMap))
					if err != nil {
						slog.Error("journalAccountTransferForeignKey display getter failed", "error", err, "key", inner.Key)
					} else {
						displayValue = displayStr
					}
				}
			}
		}
	}

	placeholder := inner.Placeholder
	if placeholder == "" {
		placeholder = "Select..."
	}

	urlStr := ""
	if inner.Url != nil {
		var err error
		urlStr, err = inner.Url(ctx)
		if err != nil {
			slog.Error("journalAccountTransferForeignKey url getter failed", "error", err, "key", inner.Key)
			urlStr = ""
		}
	}

	alpinePayload, errAlpine := json.Marshal(map[string]string{
		"value":       valuePk,
		"display":     displayValue,
		"placeholder": placeholder,
	})
	if errAlpine != nil {
		alpinePayload = []byte(`{"value":"","display":"","placeholder":""}`)
	}
	alpineData := string(alpinePayload)
	eventHandler := fmt.Sprintf("if ($event.detail.name === '%s') { value = $event.detail.value; display = $event.detail.display }", inner.Name)

	return Div(
		Class(fmt.Sprintf("my-1 relative %s", inner.Classes)),
		Attr("x-data", alpineData),
		Attr("@fk-select.window", eventHandler),
		Label(Class("label text-sm font-bold flex flex-col items-start gap-1"),
			Text(inner.Label),
			Input(Type("hidden"), Name(inner.Name), Attr(":value", "value"),
				If(inner.Required, Required()),
				Iff(inner.Attr != nil, func() (out Node) {
					out = Raw("")
					defer func() {
						if r := recover(); r != nil {
							slog.Error("journalAccountTransferForeignKey attr getter panicked", "panic", r, "key", inner.Key)
						}
					}()
					n, err := inner.Attr(ctx)
					if err != nil {
						slog.Error("journalAccountTransferForeignKey attr getter failed", "error", err, "key", inner.Key)
						return out
					}
					if n == nil {
						return out
					}
					v := reflect.ValueOf(n)
					if (v.Kind() == reflect.Pointer || v.Kind() == reflect.Map || v.Kind() == reflect.Slice || v.Kind() == reflect.Interface || v.Kind() == reflect.Func) && v.IsNil() {
						return out
					}
					return n
				}),
			),
			Div(Class("flex w-full items-stretch gap-1"),
				Div(Class("input input-bordered flex-1 flex items-center cursor-pointer"),
					Attr(":class", "display ? '' : 'opacity-50'"),
					Attr("hx-get", urlStr),
					Iff(e.Attr != nil, func() (out Node) {
						out = Raw("")
						defer func() {
							if r := recover(); r != nil {
								slog.Error("journalAccountTransferForeignKey picker attr panicked", "panic", r, "key", inner.Key)
							}
						}()
						n, err := e.Attr(ctx)
						if err != nil {
							slog.Error("journalAccountTransferForeignKey picker attr failed", "error", err, "key", inner.Key)
							return out
						}
						if n == nil {
							return out
						}
						v := reflect.ValueOf(n)
						if (v.Kind() == reflect.Pointer || v.Kind() == reflect.Map || v.Kind() == reflect.Slice || v.Kind() == reflect.Interface || v.Kind() == reflect.Func) && v.IsNil() {
							return out
						}
						return n
					}),
					Attr("hx-target", components.HTMXTargetBodyModal),
					Attr("hx-swap", components.HTMXSwapBodyModal),
					Attr("hx-push-url", "false"),
					El("span", Attr("x-text", "display || placeholder")),
				),
				If(!inner.Required,
					Button(
						Type("button"),
						Class("btn btn-ghost btn-square shrink-0"),
						Attr("@click.stop", "value = ''; display = ''"),
						Attr("x-show", "value"),
						Attr("aria-label", "Clear selection"),
						components.Render(components.Icon{Name: "x-mark"}, ctx),
					),
				),
			),
		),
	)
}
