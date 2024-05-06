-- +migrate Up
CREATE TABLE technicians (
  technician_id UUID NOT NULL PRIMARY KEY,
  store_id UUID NOT NULL REFERENCES stores (store_id),
  technician_name TEXT NOT NULL
);

CREATE TABLE sales_persons (
  sales_person_id UUID NOT NULL PRIMARY KEY,
  store_id UUID NOT NULL REFERENCES stores (store_id),
  sales_person_name TEXT NOT NULL
);

CREATE TABLE damage_types (
  damage_type_id UUID NOT NULL PRIMARY KEY,
  store_id UUID NOT NULL REFERENCES stores (store_id),
  damage_type_name TEXT NOT NULL
);

CREATE TABLE phone_conditions (
  phone_condition_id UUID NOT NULL PRIMARY KEY,
  store_id UUID NOT NULL REFERENCES stores (store_id),
  phone_condition_name TEXT NOT NULL
);

CREATE TABLE payment_methods (
  payment_method_id UUID NOT NULL PRIMARY KEY,
  store_id UUID NOT NULL REFERENCES stores (store_id),
  payment_method_name TEXT NOT NULL
);

CREATE TABLE phone_equipments (
  phone_equipment_id UUID NOT NULL PRIMARY KEY,
  store_id UUID NOT NULL REFERENCES stores (store_id),
  phone_equipment_name TEXT NOT NULL
);

CREATE TABLE repair_orders (
  repair_order_id UUID NOT NULL PRIMARY KEY,
  creation_time TIMESTAMPTZ NOT NULL,
  slug TEXT NOT NULL UNIQUE,
  store_id UUID NOT NULL REFERENCES stores (store_id),
  customer_name TEXT NOT NULL,
  contact_number TEXT NOT NULL,
  phone_type TEXT NOT NULL,
  imei TEXT,
  parts_not_checked_yet TEXT,
  color TEXT NOT NULL,
  passcode_or_pattern TEXT,
  is_pattern_locked BOOLEAN,
  pick_up_time TIMESTAMPTZ,
  completion_time TIMESTAMPTZ,
  cancellation_time TIMESTAMPTZ,
  cancellation_reason TEXT,
  confirmation_time TIMESTAMPTZ,
  confirmation_content TEXT,
  warranty_days INTEGER,
  down_payment_amount INTEGER,
  down_payment_method_id UUID REFERENCES payment_methods (payment_method_id),
  repayment_amount INTEGER,
  repayment_method_id UUID REFERENCES payment_methods (payment_method_id),
  technician_id UUID REFERENCES technicians (technician_id),
  sales_person_id UUID NOT NULL REFERENCES sales_persons (sales_person_id)
);

CREATE TABLE repair_order_damages (
  repair_order_damage_id UUID NOT NULL PRIMARY KEY,
  repair_order_id UUID NOT NULL REFERENCES repair_orders (repair_order_id) ON DELETE CASCADE,
  damage_name TEXT NOT NULL
);

CREATE TABLE repair_order_phone_conditions (
  repair_order_phone_condition_id UUID NOT NULL PRIMARY KEY,
  repair_order_id UUID NOT NULL REFERENCES repair_orders (repair_order_id) ON DELETE CASCADE,
  phone_condition_name TEXT NOT NULL
);

CREATE TABLE repair_order_phone_equipments (
  repair_order_phone_equipment_id UUID NOT NULL PRIMARY KEY,
  repair_order_id UUID NOT NULL REFERENCES repair_orders (repair_order_id) ON DELETE CASCADE,
  phone_equipment_name TEXT NOT NULL
);

CREATE TABLE repair_order_costs (
  repair_order_cost_id UUID NOT NULL PRIMARY KEY,
  repair_order_id UUID NOT NULL REFERENCES repair_orders (repair_order_id) ON DELETE CASCADE,
  amount INTEGER NOT NULL,
  reason TEXT,
  creation_time TIMESTAMPTZ NOT NULL
);

CREATE TABLE repair_order_photos (
  repair_order_photo_id UUID NOT NULL PRIMARY KEY,
  repair_order_id UUID NOT NULL REFERENCES repair_orders (repair_order_id) ON DELETE CASCADE,
  photo_url TEXT NOT NULL
);

-- +migrate Down
DROP TABLE technicians;
DROP TABLE sales_persons;
DROP TABLE damage_types;
DROP TABLE phone_conditions;
DROP TABLE payment_methods;
DROP TABLE phone_equipments;
DROP TABLE repair_order_confirmations;
DROP TABLE repair_orders;
DROP TABLE repair_order_damages;
DROP TABLE repair_order_phone_conditions;
DROP TABLE repair_order_phone_equipments;
DROP TABLE repair_order_costs;
DROP TABLE repair_order_photos;
