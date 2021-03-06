package lib

import (
	"fmt"
	"time"

	"github.com/scorredoira/dune"
	"github.com/scorredoira/dune/lib/localization"
)

var defaultLocalizer *localizer

func init() {
	defaultLocalizer = &localizer{
		culture:    localization.DefaultCulture,
		translator: localization.DefaultTranslator,
	}

	dune.RegisterLib(Locale, `

// translate a value.
declare function T(key: string, ...params: any[]): string

declare namespace locale {
	export const defaultLocalizer: Localizer
	export function setDefaultLocalizer(c: Localizer): void

	export function format(format: string, v: any): string
	export function parseNumber(v: string): number
	export function parseDate(v: string, format?: string): time.Time

	export function newCulture(name: string): Culture

	export interface Culture {
		name: string
		numberOfDecimals: number
		decimalSeparator: string
		thousandSeparator: string
		currencySymbol: string
		currency: string
		currencyPattern: string
		dateTimePattern: string
		longDateTimePattern: string
		shortDatePattern: string
		longDatePattern: string
		shortTimePattern: string
		firstDayOfWeek: number
		clone(): Culture
	}

	export function newTranslator(): Translator

	export interface Translator {
		add(language: string, key: string, translation: string): void
		languages(): string[]
	}

	export function newLocalizer(): Localizer

	export interface Localizer {
		culture: Culture
		translator: Translator
		
		translate(language: string, template: string, ...values: any[]): string
		format(format: string, v: any, language?: string): string
		parseNumber(v: string): number
		parseDate(value: string, format?: string): time.Time
	}

}

`)
}

var Locale = []dune.NativeFunction{
	{
		Name:      "->locale.defaultLocalizer",
		Arguments: 0,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			c := dune.NewObject(defaultLocalizer)
			return c, nil
		},
	},
	{
		Name:      "locale.setDefaultLocalizer",
		Arguments: 1,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			loc, ok := args[0].ToObjectOrNil().(*localizer)
			if !ok {
				return dune.NullValue, ErrInvalidType
			}
			defaultLocalizer = loc
			return dune.NullValue, nil
		},
	},
	{
		Name:      "locale.newCulture",
		Arguments: 1,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			if err := ValidateArgs(args, dune.String); err != nil {
				return dune.NullValue, err
			}
			name := args[0].ToString()
			c := dune.NewObject(&culture{localization.NewCulture(name)})
			return c, nil
		},
	},
	{
		Name:      "locale.newTranslator",
		Arguments: 0,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			t := &localization.Translator{}
			lt := dune.NewObject(&translator{t})
			return lt, nil
		},
	},
	{
		Name:      "locale.newLocalizer",
		Arguments: 0,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			loc := &localizer{}
			lt := dune.NewObject(loc)
			return lt, nil
		},
	},
	{
		Name:      "locale.parseNumber",
		Arguments: 1,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			if err := ValidateArgs(args, dune.String); err != nil {
				return dune.NullValue, err
			}

			loc := vm.Localizer
			if loc == nil {
				loc = defaultLocalizer
			}

			v, err := loc.ParseNumber(args[0].ToString())
			if err != nil {
				return dune.NullValue, err
			}

			if v == float64(int64(v)) {
				return dune.NewInt(int(v)), nil
			}
			return dune.NewFloat(v), nil
		},
	},
	{
		Name:      "locale.parseDate",
		Arguments: -1,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			if err := ValidateOptionalArgs(args, dune.String, dune.String); err != nil {
				return dune.NullValue, err
			}

			if len(args) == 0 {
				return dune.NullValue, fmt.Errorf("expected at least 1 argument, got %d", len(args))
			}

			var format string
			if len(args) == 2 {
				format = args[1].ToString()
			} else {
				format = ""
			}

			loc := vm.Localizer
			if loc == nil {
				loc = defaultLocalizer
			}

			v, err := loc.ParseDate(args[0].ToString(), format, vm.Location)
			if err != nil {
				return dune.NullValue, err
			}

			return dune.NewObject(TimeObj(v)), nil
		},
	},
	{
		Name:      "locale.format",
		Arguments: 2,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			if len(args) != 2 {
				return dune.NullValue, fmt.Errorf("expected 2 arguments, got %d", len(args))
			}

			a := args[0]
			if a.Type != dune.String {
				return dune.NullValue, fmt.Errorf("expected argument 1 to be a string, got %v", a.TypeName())
			}

			loc := vm.Localizer
			if loc == nil {
				loc = defaultLocalizer
			}

			b := args[1].Export(0)
			s := loc.Format(vm.Language, a.ToString(), b)
			return dune.NewString(s), nil
		},
	},
	{
		Name:      "T",
		Arguments: -1,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			ln := len(args)
			if ln == 0 {
				return dune.NullValue, fmt.Errorf("expected at least 1 argument, got %d", len(args))
			}

			a := args[0]
			if a.Type == dune.Null {
				return dune.NullValue, nil
			}

			var params []interface{}

			if ln > 1 {
				params = make([]interface{}, ln-1)
				for i, vp := range args[1:] {
					params[i] = vp.Export(0)
				}
			}

			loc := vm.Localizer
			if loc == nil {
				loc = defaultLocalizer
			}

			s := loc.Translate(vm.Language, a.ToString(), params...)
			return dune.NewString(s), nil
		},
	},
}

type translator struct {
	translator *localization.Translator
}

func (t *translator) Type() string {
	return "locale.Translator"
}

func (t *translator) GetMethod(name string) dune.NativeMethod {
	switch name {
	case "add":
		return t.add
	case "languages":
		return t.languages
	}
	return nil
}

func (t *translator) add(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	if err := ValidateArgs(args, dune.String, dune.String, dune.String); err != nil {
		return dune.NullValue, err
	}

	language := args[0].ToString()
	key := args[1].ToString()
	value := args[2].ToString()

	t.translator.Add(language, key, value)
	return dune.NullValue, nil
}

func (t *translator) languages(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	if err := ValidateArgs(args); err != nil {
		return dune.NullValue, err
	}

	keys := t.translator.Languages()

	values := make([]dune.Value, len(keys))

	for i, k := range keys {
		values[i] = dune.NewString(k)
	}

	return dune.NewArrayValues(values), nil
}

type culture struct {
	culture *localization.Culture
}

func (c *culture) Type() string {
	return "locale.Culture"
}

func (c *culture) GetProperty(name string, vm *dune.VM) (dune.Value, error) {
	switch name {
	case "name":
		return dune.NewString(c.culture.Name), nil
	case "currencySymbol":
		return dune.NewString(c.culture.CurrencySymbol), nil
	case "currency":
		return dune.NewString(c.culture.Currency), nil
	case "currencyPattern":
		return dune.NewString(c.culture.CurrencyPattern), nil
	case "numberOfDecimals":
		return dune.NewInt(c.culture.NumberOfDecimals), nil
	case "decimalSeparator":
		return dune.NewString(string(c.culture.DecimalSeparator)), nil
	case "thousandSeparator":
		return dune.NewString(string(c.culture.ThousandSeparator)), nil
	case "shortDatePattern":
		return dune.NewString(c.culture.ShortDatePattern), nil
	case "longDatePattern":
		return dune.NewString(c.culture.LongDatePattern), nil
	case "longDateTimePattern":
		return dune.NewString(c.culture.LongDateTimePattern), nil
	case "dateTimePattern":
		return dune.NewString(c.culture.DateTimePattern), nil
	case "shortTimePattern":
		return dune.NewString(c.culture.ShortTimePattern), nil
	case "firstDayOfWeek":
		return dune.NewInt(int(c.culture.FirstDayOfWeek)), nil
	default:
		return dune.UndefinedValue, nil
	}
}

func (c *culture) SetProperty(name string, v dune.Value, vm *dune.VM) error {
	switch name {
	case "name":
		if v.Type != dune.String {
			return ErrInvalidType
		}
		c.culture.Name = v.ToString()
	case "currencySymbol":
		if v.Type != dune.String {
			return ErrInvalidType
		}
		c.culture.CurrencySymbol = v.ToString()
	case "currency":
		if v.Type != dune.String {
			return ErrInvalidType
		}
		c.culture.Currency = v.ToString()
	case "currencyPattern":
		if v.Type != dune.String {
			return ErrInvalidType
		}
		c.culture.CurrencyPattern = v.ToString()
	case "numberOfDecimals":
		if v.Type != dune.Int {
			return ErrInvalidType
		}
		c.culture.NumberOfDecimals = int(v.ToInt())
	case "decimalSeparator":
		if v.Type != dune.Rune && v.Type != dune.String {
			return ErrInvalidType
		}
		c.culture.DecimalSeparator = v.ToRune()
	case "thousandSeparator":
		if v.Type != dune.Rune && v.Type != dune.String {
			return ErrInvalidType
		}
		c.culture.ThousandSeparator = v.ToRune()
	case "shortDatePattern":
		if v.Type != dune.String {
			return ErrInvalidType
		}
		c.culture.ShortDatePattern = v.ToString()
	case "longDatePattern":
		if v.Type != dune.String {
			return ErrInvalidType
		}
		c.culture.LongDatePattern = v.ToString()
	case "dateTimePattern":
		if v.Type != dune.String {
			return ErrInvalidType
		}
		c.culture.DateTimePattern = v.ToString()
	case "longDateTimePattern":
		if v.Type != dune.String {
			return ErrInvalidType
		}
		c.culture.LongDateTimePattern = v.ToString()
	case "shortTimePattern":
		if v.Type != dune.String {
			return ErrInvalidType
		}
		c.culture.ShortTimePattern = v.ToString()
	case "firstDayOfWeek":
		if v.Type != dune.Int {
			return ErrInvalidType
		}
		c.culture.FirstDayOfWeek = time.Weekday(v.ToInt())
	default:
		return ErrUndefined
	}

	return nil
}

func (c *culture) GetMethod(name string) dune.NativeMethod {
	switch name {
	case "clone":
		return c.clone
	}
	return nil
}

func (c *culture) clone(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	if len(args) != 0 {
		return dune.NullValue, fmt.Errorf("expected 0 arguments, got %d", len(args))
	}

	return dune.NewObject(&culture{culture: c.culture.Clone()}), nil
}

type localizer struct {
	culture    *localization.Culture
	translator *localization.Translator
}

func (l *localizer) Type() string {
	return "locale.Localizer"
}

func (l *localizer) GetProperty(name string, vm *dune.VM) (dune.Value, error) {
	switch name {
	case "culture":
		return dune.NewObject(&culture{l.culture}), nil
	case "translator":
		return dune.NewObject(&translator{l.translator}), nil
	default:
		return dune.UndefinedValue, nil
	}
}

func (l *localizer) SetProperty(name string, v dune.Value, vm *dune.VM) error {
	switch name {
	case "culture":
		c, ok := v.ToObjectOrNil().(*culture)
		if !ok {
			return ErrInvalidType
		}
		l.culture = c.culture
	case "translator":
		t, ok := v.ToObjectOrNil().(*translator)
		if !ok {
			return ErrInvalidType
		}
		l.translator = t.translator
	default:
		return ErrReadOnlyOrUndefined
	}

	return nil
}

func (l *localizer) GetMethod(name string) dune.NativeMethod {
	switch name {
	case "translate":
		return l.translate
	case "format":
		return l.format
	case "parseNumber":
		return l.parseNumber
	case "parseDate":
		return l.parseDate
	}
	return nil
}

func (l *localizer) translate(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	ln := len(args)
	if ln < 2 {
		return dune.NullValue, fmt.Errorf("expected at least 2 arguments, got %d", len(args))
	}

	lang := args[0]
	if lang.Type != dune.String {
		return dune.NullValue, fmt.Errorf("expected argument 1 to be a string, got %v", lang.TypeName())
	}

	a := args[1]
	if a.Type != dune.String {
		return dune.NullValue, fmt.Errorf("expected argument 1 to be a string, got %v", a.TypeName())
	}

	var params []interface{}

	if ln > 1 {
		params = make([]interface{}, ln-2)
		for i, vp := range args[2:] {
			params[i] = vp.Export(0)
		}
	}

	s := l.Translate(lang.ToString(), a.ToString(), params...)
	return dune.NewString(s), nil
}

func (l *localizer) format(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return dune.NullValue, fmt.Errorf("expected 2 or 3 arguments, got %d", len(args))
	}

	a := args[0]
	if a.Type != dune.String {
		return dune.NullValue, fmt.Errorf("expected argument 2 to be a string, got %v", a.TypeName())
	}

	b := args[1].Export(0)

	var lang string
	if len(args) == 3 {
		c := args[2]
		if c.Type != dune.String {
			return dune.NullValue, fmt.Errorf("expected argument 3 to be a string, got %v", c.TypeName())
		}
		lang = c.ToString()
	}

	s := l.Format(lang, a.ToString(), b)
	return dune.NewString(s), nil
}

func (l *localizer) parseNumber(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	if err := ValidateArgs(args, dune.String); err != nil {
		return dune.NullValue, err
	}

	v, err := l.ParseNumber(args[0].ToString())
	if err != nil {
		return dune.NullValue, err
	}

	return dune.NewFloat(v), nil
}

func (l *localizer) parseDate(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	if err := ValidateOptionalArgs(args, dune.String, dune.String); err != nil {
		return dune.NullValue, err
	}

	if len(args) == 0 {
		return dune.NullValue, fmt.Errorf("expected at least 1 argument, got %d", len(args))
	}

	var format string
	if len(args) == 2 {
		format = args[1].ToString()
	} else {
		format = ""
	}

	v, err := l.ParseDate(args[0].ToString(), format, vm.Location)
	if err != nil {
		return dune.NullValue, err
	}

	return dune.NewObject(TimeObj(v)), nil
}

func (l *localizer) Translate(language, template string, values ...interface{}) string {
	template, _ = l.translator.Translate(language, template)
	return FormatTemplate(template, values...)
}

func (l *localizer) Format(language, format string, v interface{}) string {
	return l.culture.Format(format, v, l.translator, language)
}

func (l *localizer) ParseNumber(v string) (float64, error) {
	return l.culture.ParseNumber(v)
}

func (l *localizer) ParseDate(value, format string, loc *time.Location) (time.Time, error) {
	return l.culture.ParseDate(value, format, loc)
}
