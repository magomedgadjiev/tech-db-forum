package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/validate"
)

// Forum Информация о форуме.
//
// swagger:model Forum
type Forum struct {

	// Общее кол-во сообщений в данном форуме.
	//
	// Read Only: true
	Posts int64 `json:"posts,omitempty"`

	// Человекопонятный URL (https://ru.wikipedia.org/wiki/%D0%A1%D0%B5%D0%BC%D0%B0%D0%BD%D1%82%D0%B8%D1%87%D0%B5%D1%81%D0%BA%D0%B8%D0%B9_URL).
	// Required: true
	// Pattern: ^(\d|\w|-|_)*(\w|-|_)(\d|\w|-|_)*$
	Slug string `json:"slug"`

	// Общее кол-во ветвей обсуждения в данном форуме.
	//
	// Read Only: true
	Threads int32 `json:"threads,omitempty"`

	// Название форума.
	// Required: true
	Title string `json:"title"`

	// Nickname пользователя, который отвечает за форум (уникальное поле).
	// Required: true
	User string `json:"user"`
}

// Validate validates this forum
func (m *Forum) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSlug(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateTitle(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateUser(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Forum) validateSlug(formats strfmt.Registry) error {

	if err := validate.RequiredString("slug", "body", string(m.Slug)); err != nil {
		return err
	}

	if err := validate.Pattern("slug", "body", string(m.Slug), `^(\d|\w|-|_)*(\w|-|_)(\d|\w|-|_)*$`); err != nil {
		return err
	}

	return nil
}

func (m *Forum) validateTitle(formats strfmt.Registry) error {

	if err := validate.RequiredString("title", "body", string(m.Title)); err != nil {
		return err
	}

	return nil
}

func (m *Forum) validateUser(formats strfmt.Registry) error {

	if err := validate.RequiredString("user", "body", string(m.User)); err != nil {
		return err
	}

	return nil
}
