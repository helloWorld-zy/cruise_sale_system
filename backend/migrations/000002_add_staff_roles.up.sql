-- HI-04: Add staff_roles join table (many-to-many between staffs and roles)
CREATE TABLE staff_roles (
    staff_id BIGINT NOT NULL REFERENCES staffs(id) ON DELETE CASCADE,
    role_id  BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (staff_id, role_id)
);
CREATE INDEX idx_staff_roles_staff_id ON staff_roles(staff_id);
CREATE INDEX idx_staff_roles_role_id ON staff_roles(role_id);
