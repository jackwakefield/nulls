package nulls

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strconv"
	"encoding/xml"
)

// Float64 replaces sql.NullFloat64 with an implementation
// that supports proper JSON encoding/decoding.
type Float64 sql.NullFloat64

// Interface implements the nullable interface. It returns nil if
// the float64 is not valid, otherwise it returns the float64 value.
func (ns Float64) Interface() interface{} {
	if !ns.Valid {
		return nil
	}
	return ns.Float64
}

// NewFloat64 returns a new, properly instantiated
// Float64 object.
func NewFloat64(i float64) Float64 {
	return Float64{Float64: i, Valid: true}
}

// Scan implements the Scanner interface.
func (ns *Float64) Scan(value interface{}) error {
	n := sql.NullFloat64{Float64: ns.Float64}
	err := n.Scan(value)
	ns.Float64, ns.Valid = n.Float64, n.Valid
	return err
}

// Value implements the driver Valuer interface.
func (ns Float64) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.Float64, nil
}

// MarshalJSON marshals the underlying value to a
// proper JSON representation.
func (ns Float64) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.Float64)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON will unmarshal a JSON value into
// the propert representation of that value.
func (ns *Float64) UnmarshalJSON(text []byte) error {
	t := string(text)
	ns.Valid = true
	if t == "null" {
		ns.Valid = false
		return nil
	}
	i, err := strconv.ParseFloat(t, 64)
	if err != nil {
		ns.Valid = false
		return err
	}
	ns.Float64 = i
	return nil
}

// UnmarshalText will unmarshal text value into
// the propert representation of that value.
func (ns *Float64) UnmarshalText(text []byte) error {
	return ns.UnmarshalJSON(text)
}

func (ns Float64) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if ns.Valid {
		return e.EncodeElement(ns.Float64, start)
	}
	return nil
}

func (ns *Float64) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var data string
	err := d.DecodeElement(&data, &start)

	if err != nil {
		return err
	}
	if data == "" {
		return nil
	}

	if data == "null" {
		return nil
	}

	val, err := strconv.ParseFloat(data, 64)

	if err != nil {
		return err
	}

	ns.Valid = true
	ns.Float64 = val

	return nil
}

func (ns Float64) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if ns.Valid {
		var value string
		value = strconv.FormatFloat(ns.Float64,'f', -1, 64)

		return xml.Attr{
			Name:  name,
			Value: value,
		}, nil
	}
	return xml.Attr{}, nil
}

func (ns *Float64) UnmarshalXMLAttr(attr xml.Attr) error {
	if attr.Value == "" {
		return nil
	}

	if attr.Value == "null" {
		return nil
	}

	val, err := strconv.ParseFloat(attr.Value, 64)

	if err != nil {
		return err
	}

	ns.Valid = true
	ns.Float64 = val

	return nil
}