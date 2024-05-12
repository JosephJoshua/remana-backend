-- +migrate Up
CREATE TABLE permission_groups (
  permission_group_id UUID NOT NULL PRIMARY KEY,
  permission_group_name TEXT NOT NULL UNIQUE
);

CREATE TABLE permissions (
  permission_id UUID NOT NULL PRIMARY KEY,
  permission_group_id UUID NOT NULL REFERENCES permission_groups (permission_group_id),
  permission_name TEXT NOT NULL UNIQUE,
  permission_display_name TEXT NOT NULL UNIQUE
);

CREATE TABLE role_permissions (
  role_id UUID NOT NULL REFERENCES roles (role_id) ON DELETE CASCADE,
  permission_id UUID NOT NULL REFERENCES permissions (permission_id) ON DELETE CASCADE,
  PRIMARY KEY (role_id, permission_id)
);

-- +migrate Down
DROP TABLE permission_groups;
DROP TABLE permissions;
DROP TABLE role_permissions;
