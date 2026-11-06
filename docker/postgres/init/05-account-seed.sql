-- Testni korisnici, isti kao u Android aplikaciji (SampleUsers.kt).
-- akovac / trgovac123, admin / admin123, ihorvat / trgovac123 (deaktiviran)
INSERT INTO account.users (username, password_hash, role, active) VALUES
('akovac',  '$2a$10$safpWuFTm3yxXLVRPNrt0.00fEhUdr953oI2s9luQsBfmygaho8vy', 'MERCHANT', true),
('admin',   '$2a$10$t8GU39jbIL6z.AIG1HQyk.lvalmQhnOWXY9NIsZNC1e6NR2dW8Zqu', 'ADMIN',    true),
('ihorvat', '$2a$10$qIFqkvYy3IB8IV7JnB21cO4Q00fY7VlGabs9dAHbbvOWVxT2x/GLq', 'MERCHANT', false);