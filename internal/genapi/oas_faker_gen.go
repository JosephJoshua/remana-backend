// Code generated by ogen, DO NOT EDIT.

package genapi

import (
	"net/url"

	"github.com/google/uuid"
)

// SetFake set fake values.
func (s *CreateRepairOrderRequest) SetFake() {
	{
		{
			s.CustomerName = "string"
		}
	}
	{
		{
			s.ContactPhoneNumber = "string"
		}
	}
	{
		{
			s.PhoneType = "string"
		}
	}
	{
		{
			s.Imei.SetFake()
		}
	}
	{
		{
			s.PartsNotCheckedYet.SetFake()
		}
	}
	{
		{
			s.Passcode.SetFake()
		}
	}
	{
		{
			s.Color = "string"
		}
	}
	{
		{
			s.InitialCost = int(0)
		}
	}
	{
		{
			s.DownPayment.SetFake()
		}
	}
	{
		{
			s.SalesPersonID = uuid.New()
		}
	}
	{
		{
			s.TechnicianID = uuid.New()
		}
	}
	{
		{
			s.PhoneConditions = nil
			for i := 0; i < 0; i++ {
				var elem uuid.UUID
				{
					elem = uuid.New()
				}
				s.PhoneConditions = append(s.PhoneConditions, elem)
			}
		}
	}
	{
		{
			s.DamageTypes = nil
			for i := 0; i < 1; i++ {
				var elem uuid.UUID
				{
					elem = uuid.New()
				}
				s.DamageTypes = append(s.DamageTypes, elem)
			}
		}
	}
	{
		{
			s.PhoneEquipments = nil
			for i := 0; i < 0; i++ {
				var elem uuid.UUID
				{
					elem = uuid.New()
				}
				s.PhoneEquipments = append(s.PhoneEquipments, elem)
			}
		}
	}
	{
		{
			s.Photos = nil
			for i := 0; i < 1; i++ {
				var elem url.URL
				{
					elem = url.URL{Scheme: "https", Host: "github.com", Path: "/ogen-go/ogen"}
				}
				s.Photos = append(s.Photos, elem)
			}
		}
	}
}

// SetFake set fake values.
func (s *CreateRepairOrderRequestDownPayment) SetFake() {
	{
		{
			s.Amount = int(0)
		}
	}
	{
		{
			s.Method = uuid.New()
		}
	}
}

// SetFake set fake values.
func (s *CreateRepairOrderRequestPasscode) SetFake() {
	{
		{
			s.IsPatternLocked = true
		}
	}
	{
		{
			s.Value = "string"
		}
	}
}

// SetFake set fake values.
func (s *CreateTechnicianRequest) SetFake() {
	{
		{
			s.Name = "string"
		}
	}
}

// SetFake set fake values.
func (s *Error) SetFake() {
	{
		{
			s.Message = "string"
		}
	}
}

// SetFake set fake values.
func (s *LoginCodePrompt) SetFake() {
	{
		{
			s.LoginCode = "string"
		}
	}
}

// SetFake set fake values.
func (s *LoginCredentials) SetFake() {
	{
		{
			s.Username = "string"
		}
	}
	{
		{
			s.Password = "string"
		}
	}
	{
		{
			s.StoreCode = "string"
		}
	}
}

// SetFake set fake values.
func (s *LoginResponse) SetFake() {
	{
		{
			s.Type.SetFake()
		}
	}
}

// SetFake set fake values.
func (s *LoginResponseType) SetFake() {
	*s = LoginResponseTypeAdmin
}

// SetFake set fake values.
func (s *OptCreateRepairOrderRequestDownPayment) SetFake() {
	var elem CreateRepairOrderRequestDownPayment
	{
		elem.SetFake()
	}
	s.SetTo(elem)
}

// SetFake set fake values.
func (s *OptCreateRepairOrderRequestPasscode) SetFake() {
	var elem CreateRepairOrderRequestPasscode
	{
		elem.SetFake()
	}
	s.SetTo(elem)
}

// SetFake set fake values.
func (s *OptString) SetFake() {
	var elem string
	{
		elem = "string"
	}
	s.SetTo(elem)
}

// SetFake set fake values.
func (s *UserDetails) SetFake() {
	{
		{
			s.ID = uuid.New()
		}
	}
	{
		{
			s.Username = "string"
		}
	}
	{
		{
			s.Role.SetFake()
		}
	}
	{
		{
			s.Store.SetFake()
		}
	}
}

// SetFake set fake values.
func (s *UserDetailsRole) SetFake() {
	{
		{
			s.ID = uuid.New()
		}
	}
	{
		{
			s.Name = "string"
		}
	}
	{
		{
			s.IsStoreAdmin = true
		}
	}
}

// SetFake set fake values.
func (s *UserDetailsStore) SetFake() {
	{
		{
			s.ID = uuid.New()
		}
	}
	{
		{
			s.Name = "string"
		}
	}
	{
		{
			s.Code = "string"
		}
	}
}
