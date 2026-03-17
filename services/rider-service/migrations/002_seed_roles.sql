INSERT INTO roles (code, name, description) VALUES
('SUPER_ADMIN', 'Super Admin', 'Global platform administrator'),
('RESTAURANT_OWNER', 'Restaurant Owner', 'Owner-level restaurant access'),
('RESTAURANT_ADMIN', 'Restaurant Admin', 'Restaurant operational admin'),
('RESTAURANT_STAFF', 'Restaurant Staff', 'Restaurant staff member'),
('DISPATCHER_MANAGER', 'Dispatcher / Manager', 'Dispatch and assignment operator'),
('RIDER', 'Rider', 'Delivery rider account')
ON CONFLICT (code) DO NOTHING;
