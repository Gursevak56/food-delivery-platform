INSERT INTO system_configurations (config_key, config_value) VALUES
('order_accept_timeout_seconds', '45'::jsonb),
('pickup_otp_required', 'true'::jsonb),
('delivery_otp_required', 'true'::jsonb),
('otp_max_retries', '5'::jsonb),
('otp_max_resends', '3'::jsonb),
('rider_max_active_orders', '1'::jsonb),
('surge_multiplier_default', '1.0'::jsonb),
('payout_minimum_threshold', '500'::jsonb),
('break_rules', '{"max_break_minutes":30}'::jsonb),
('shift_rules', '{"auto_start":true}'::jsonb),
('incentive_rules', '{"distance_bonus_threshold_km":5,"distance_bonus_amount":12}'::jsonb)
ON CONFLICT (config_key) DO UPDATE SET config_value = EXCLUDED.config_value, updated_at = NOW();
