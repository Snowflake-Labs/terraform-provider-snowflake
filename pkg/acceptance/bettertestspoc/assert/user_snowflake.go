package assert

import (
	"fmt"
	"testing"
	"time"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type UserAssert struct {
	*SnowflakeObjectAssert[sdk.User, sdk.AccountObjectIdentifier]
}

func User(t *testing.T, id sdk.AccountObjectIdentifier) *UserAssert {
	t.Helper()
	return &UserAssert{
		NewSnowflakeObjectAssertWithProvider(sdk.ObjectTypeUser, id, acc.TestClient().User.Show),
	}
}

func UserFromObject(t *testing.T, user *sdk.User) *UserAssert {
	t.Helper()
	return &UserAssert{
		NewSnowflakeObjectAssertWithObject(sdk.ObjectTypeUser, user.ID(), user),
	}
}

func (w *UserAssert) HasName(expected string) *UserAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.Name != expected {
			return fmt.Errorf("expected name: %v; got: %v", expected, o.Name)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasCreatedOn(expected time.Time) *UserAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.CreatedOn != expected {
			return fmt.Errorf("expected created on: %v; got: %v", expected, o.CreatedOn)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasLoginName(expected string) *UserAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.LoginName != expected {
			return fmt.Errorf("expected login name: %v; got: %v", expected, o.LoginName)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasDisplayName(expected string) *UserAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.DisplayName != expected {
			return fmt.Errorf("expected display name: %v; got: %v", expected, o.DisplayName)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasFirstName(expected string) *UserAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.FirstName != expected {
			return fmt.Errorf("expected first name: %v; got: %v", expected, o.FirstName)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasLastName(expected string) *UserAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.LastName != expected {
			return fmt.Errorf("expected last name: %v; got: %v", expected, o.LastName)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasEmail(expected string) *UserAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.Email != expected {
			return fmt.Errorf("expected email: %v; got: %v", expected, o.Email)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasMinsToUnlock(expected string) *UserAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.MinsToUnlock != expected {
			return fmt.Errorf("expected mins to unlock: %v; got: %v", expected, o.MinsToUnlock)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasDaysToExpiry(expected string) *UserAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.DaysToExpiry != expected {
			return fmt.Errorf("expected days to expiry: %v; got: %v", expected, o.DaysToExpiry)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasComment(expected string) *UserAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.Comment != expected {
			return fmt.Errorf("expected comment: %v; got: %v", expected, o.Comment)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasDisabled(expected bool) *UserAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.Disabled != expected {
			return fmt.Errorf("expected disabled: %v; got: %v", expected, o.Disabled)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasMustChangePassword(expected bool) *UserAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.MustChangePassword != expected {
			return fmt.Errorf("expected must change password: %v; got: %v", expected, o.MustChangePassword)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasSnowflakeLock(expected bool) *UserAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.SnowflakeLock != expected {
			return fmt.Errorf("expected snowflake lock: %v; got: %v", expected, o.SnowflakeLock)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasDefaultWarehouse(expected string) *UserAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.DefaultWarehouse != expected {
			return fmt.Errorf("expected default warehouse: %v; got: %v", expected, o.DefaultWarehouse)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasDefaultNamespace(expected string) *UserAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.DefaultNamespace != expected {
			return fmt.Errorf("expected default namespace: %v; got: %v", expected, o.DefaultNamespace)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasDefaultRole(expected string) *UserAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.DefaultRole != expected {
			return fmt.Errorf("expected default role: %v; got: %v", expected, o.DefaultRole)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasDefaultSecondaryRoles(expected string) *UserAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.DefaultSecondaryRoles != expected {
			return fmt.Errorf("expected default secondary roles: %v; got: %v", expected, o.DefaultSecondaryRoles)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasExtAuthnDuo(expected bool) *UserAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.ExtAuthnDuo != expected {
			return fmt.Errorf("expected ext auth duo: %v; got: %v", expected, o.ExtAuthnDuo)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasExtAuthnUid(expected string) *UserAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.ExtAuthnUid != expected {
			return fmt.Errorf("expected ext authn uid: %v; got: %v", expected, o.ExtAuthnUid)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasMinsToBypassMfa(expected string) *UserAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.MinsToBypassMfa != expected {
			return fmt.Errorf("expected mins to bypass mfa: %v; got: %v", expected, o.MinsToBypassMfa)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasOwner(expected string) *UserAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.Owner != expected {
			return fmt.Errorf("expected owner: %v; got: %v", expected, o.Owner)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasLastSuccessLogin(expected time.Time) *UserAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.LastSuccessLogin != expected {
			return fmt.Errorf("expected last success login: %v; got: %v", expected, o.LastSuccessLogin)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasExpiresAtTime(expected time.Time) *UserAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.ExpiresAtTime != expected {
			return fmt.Errorf("expected expires at time: %v; got: %v", expected, o.ExpiresAtTime)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasLockedUntilTime(expected time.Time) *UserAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.LockedUntilTime != expected {
			return fmt.Errorf("expected locked until time: %v; got: %v", expected, o.LockedUntilTime)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasHasPassword(expected bool) *UserAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.HasPassword != expected {
			return fmt.Errorf("expected has password: %v; got: %v", expected, o.HasPassword)
		}
		return nil
	})
	return w
}

func (w *UserAssert) HasHasRsaPublicKey(expected bool) *UserAssert {
	w.assertions = append(w.assertions, func(t *testing.T, o *sdk.User) error {
		t.Helper()
		if o.HasRsaPublicKey != expected {
			return fmt.Errorf("expected has rsa public key: %v; got: %v", expected, o.HasRsaPublicKey)
		}
		return nil
	})
	return w
}
