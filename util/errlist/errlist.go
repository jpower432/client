package errlist

import (
	"errors"
	"strings"
)

// TODO(jpower432): Determine if we want to have this errlist implementation which is
// very similar to apimachinery.Aggregate or import another kubernetes lib. However, it
// uses strings.Builder to concatenate strings in the event there is a large slice.

type ErrList interface {
	error
	Errors() []error
	Is(error) bool
}

func NewErrList(list []error) ErrList {
	if len(list) == 0 {
		return nil
	}
	// In case of input error list contains nil
	var errs []error
	for _, e := range list {
		if e != nil {
			errs = append(errs, e)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errlist(errs)
}

type errlist []error

func (el errlist) Error() string {
	if len(el) == 0 {
		// This should never happen, really.
		return ""
	}
	if len(el) == 1 {
		return el[0].Error()
	}
	seenErrs := map[string]struct{}{}
	result := strings.Builder{}
	el.visit(func(err error) bool {
		msg := err.Error()
		if _, has := seenErrs[msg]; has {
			return false
		}
		seenErrs[msg] = struct{}{}
		if len(seenErrs) > 1 {
			result.WriteString(", ")
		}
		result.WriteString(msg)
		return false
	})
	if len(seenErrs) == 1 {
		return result.String()
	}
	return "[" + result.String() + "]"
}

func (el errlist) Is(target error) bool {
	return el.visit(func(err error) bool {
		return errors.Is(err, target)
	})
}

func (el errlist) visit(f func(err error) bool) bool {
	for _, err := range el {
		switch err := err.(type) {
		case errlist:
			if match := err.visit(f); match {
				return match
			}
		case ErrList:
			for _, nestedErr := range err.Errors() {
				if match := f(nestedErr); match {
					return match
				}
			}
		default:
			if match := f(err); match {
				return match
			}
		}
	}

	return false
}

func (el errlist) Errors() []error {
	return []error(el)
}
