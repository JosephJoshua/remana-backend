package permission

const (
	groupNameRepairOrder    = "repair_order"
	groupNameDamageType     = "damage_type"
	groupNamePhoneCondition = "phone_condition"
	groupNamePhoneEquipment = "phone_equipment"
	groupNameTechnician     = "technician"
	groupNameSalesPerson    = "sales_person"
	groupNamePaymentMethod  = "payment_method"
	groupNameRole           = "role"
)

type Permission interface {
	GroupName() string
	Name() string
}

type permission struct {
	groupName string
	name      string
}

func (p permission) GroupName() string {
	return p.groupName
}

func (p permission) Name() string {
	return p.name
}

func CreateRepairOrder() Permission {
	return permission{
		groupName: groupNameRepairOrder,
		name:      "create",
	}
}

func CreateDamageType() Permission {
	return permission{
		groupName: groupNameDamageType,
		name:      "create",
	}
}

func CreatePhoneCondition() Permission {
	return permission{
		groupName: groupNamePhoneCondition,
		name:      "create",
	}
}

func CreatePhoneEquipment() Permission {
	return permission{
		groupName: groupNamePhoneEquipment,
		name:      "create",
	}
}

func CreateTechnician() Permission {
	return permission{
		groupName: groupNameTechnician,
		name:      "create",
	}
}

func CreateSalesPerson() Permission {
	return permission{
		groupName: groupNameSalesPerson,
		name:      "create",
	}
}

func CreatePaymentMethod() Permission {
	return permission{
		groupName: groupNamePaymentMethod,
		name:      "create",
	}
}

func CreateRole() Permission {
	return permission{
		groupName: groupNameRole,
		name:      "create",
	}
}

func AssignPermissionsToRole() Permission {
	return permission{
		groupName: groupNameRole,
		name:      "assign_permissions",
	}
}
