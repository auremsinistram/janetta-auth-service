package account

import (
	"auth-service/pkg/econst"
	"auth-service/pkg/tools"
	"regexp"

	"github.com/auremsinistram/go-errors"
)

// TODO: use config
var phoneRegexp = regexp.MustCompile(`^\+7\s{1}\(\d{3}\)\s\d{3}-\d{2}-\d{2}$`)

func (r *SignUpReq) Validate() error {
	if !phoneRegexp.MatchString(r.Phone) {
		return errors.NewWithCode(econst.InvalidPhoneFormat, "is not match")
	}

	if !tools.PasswordStrong(r.Password) {
		return errors.NewWithCode(econst.WeakPassword, "is not strong")
	}

	return nil
}

func (r *SignUpConfirmReq) Validate() error {
	if !phoneRegexp.MatchString(r.Phone) {
		return errors.NewWithCode(econst.InvalidPhoneFormat, "is not match")
	}

	return nil
}

func (r *SignInReq) Validate() error {
	if !phoneRegexp.MatchString(r.Phone) {
		return errors.NewWithCode(econst.InvalidPhoneFormat, "is not match")
	}

	if !tools.PasswordStrong(r.Password) {
		return errors.NewWithCode(econst.WeakPassword, "is not strong")
	}

	return nil
}

func (r *SignInConfirmReq) Validate() error {
	if !phoneRegexp.MatchString(r.Phone) {
		return errors.NewWithCode(econst.InvalidPhoneFormat, "is not match")
	}

	return nil
}

func (r *PasswordForgotReq) Validate() error {
	if !phoneRegexp.MatchString(r.Phone) {
		return errors.NewWithCode(econst.InvalidPhoneFormat, "is not match")
	}

	return nil
}

func (r *PasswordResetReq) Validate() error {
	if !phoneRegexp.MatchString(r.Phone) {
		return errors.NewWithCode(econst.InvalidPhoneFormat, "is not match")
	}

	if !tools.PasswordStrong(r.Password) {
		return errors.NewWithCode(econst.WeakPassword, "is not strong")
	}

	return nil
}

func (r *PasswordChangeReq) Validate() error {
	if !tools.PasswordStrong(r.Password) {
		return errors.NewWithCode(econst.WeakPassword, "is not strong")
	}

	return nil
}
