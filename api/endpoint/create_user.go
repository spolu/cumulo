package endpoint

import (
	"context"
	"net/http"

	"github.com/spolu/cumulo/api/model"
	"github.com/spolu/cumulo/lib/db"
	"github.com/spolu/cumulo/lib/errors"
	"github.com/spolu/cumulo/lib/format"
	"github.com/spolu/cumulo/lib/logging"
	"github.com/spolu/cumulo/lib/ptr"
	"github.com/spolu/cumulo/lib/svc"
)

const (
	// EndPtCreateUser creates a new user.
	EndPtCreateUser EndPtName = "CreateUser"
)

func init() {
	registrar[EndPtCreateUser] = NewCreateUser
}

// CreateUser a new user by username and email and send its secret over eail.
type CreateUser struct {
	Phone string
}

// NewCreateUser constructs and initializes the endpoint.
func NewCreateUser(
	r *http.Request,
) (Endpoint, error) {
	return &CreateUser{}, nil
}

// Validate validates the input parameters.
func (e *CreateUser) Validate(
	r *http.Request,
) error {
	ctx := r.Context()

	// Validate username.
	phone, err := ValidatePhone(ctx, r.PostFormValue("phone"))
	if err != nil {
		return errors.Trace(err) // 400
	}
	e.Phone = *phone

	return nil
}

// Execute executes the endpoint.
func (e *CreateUser) Execute(
	ctx context.Context,
) (*int, *svc.Resp, error) {
	ctx = db.Begin(ctx, "api")
	defer db.LoggedRollback(ctx)

	user, err := model.CreateUser(ctx,
		e.Phone,
	)
	if err != nil {
		switch err := errors.Cause(err).(type) {
		case model.ErrUniqueConstraintViolation:
			return nil, nil, errors.Trace(errors.NewUserErrorf(err,
				400, "phone_already_registered",
				"The provided phone number is already registered: %s",
				e.Phone,
			))
		default:
			return nil, nil, errors.Trace(err) // 500
		}
	}

	// TODO(stan): validate phone with twilio
	// TODO(stan): send SMS or call Login endpoint
	//             see https://www.twilio.com/lookup

	db.Commit(ctx)

	logging.Logf(ctx,
		"Created user: token=%s created=%q phone=%s",
		user.Token, user.Created, user.Phone)

	return ptr.Int(http.StatusCreated), &svc.Resp{
		"user": format.JSONPtr(model.NewUserResource(ctx, user)),
	}, nil
}
