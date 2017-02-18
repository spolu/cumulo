package endpoint

import (
	"context"
	"regexp"

	"github.com/spolu/cumulo/lib/errors"
)

// Possible email: von.neumann+foo@ias.edu
var emailRegexp = regexp.MustCompile(
	"^([a-zA-Z0-9-_.]{1,256})(\\+[a-zA-Z0-9-_.]+){0,1}@" +
		"([a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+)$")

// Possible phone: +3314152165700
var phoneRegexp = regexp.MustCompile(
	"^\\+[0-9]{1,3}[0-9]+$")

// Possible username: von.neuman-23_86
var usernameRegexp = regexp.MustCompile("^([a-zA-Z0-9-_.]{1,256})$")

// ValidateEmail validates an email address.
func ValidateEmail(
	ctx context.Context,
	email string,
) (*string, error) {

	if !emailRegexp.MatchString(email) {
		return nil, errors.Trace(errors.NewUserErrorf(nil,
			400, "email_invalid",
			"The email you provided is invalid: %s.",
			email,
		))
	}

	return &email, nil
}

// ValidateUsername validates an email address.
func ValidateUsername(
	ctx context.Context,
	username string,
) (*string, error) {

	if !usernameRegexp.MatchString(username) {
		return nil, errors.Trace(errors.NewUserErrorf(nil,
			400, "username_invalid",
			"The username you provided is invalid: %s.",
			username,
		))
	}

	return &username, nil
}

// ValidatePhone validates (structurally) a phone number.
func ValidatePhone(
	ctx context.Context,
	phone string,
) (*string, error) {

	if !phoneRegexp.MatchString(phone) {
		return nil, errors.Trace(errors.NewUserErrorf(nil,
			400, "phone_invalid",
			"The phone number you provided is invalid: %s.",
			phone,
		))
	}

	return &phone, nil
}
