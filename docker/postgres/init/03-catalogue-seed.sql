INSERT INTO catalogue.category (name, description) VALUES
('Torte', 'Kriške torti'),
('Kolači', 'Pojedinačni kolači i deserti');

INSERT INTO catalogue.article (category_id, code, name, description, price, available_quantity) VALUES
((SELECT category_id FROM catalogue.category WHERE name = 'Torte'),  'S-001', 'Čokoladna torta', 'Kriška, tamna čokolada i vrhnje', 4.50, 8),
((SELECT category_id FROM catalogue.category WHERE name = 'Torte'),  'S-002', 'Sacher kocka', 'Marelica i čokoladna glazura', 3.80, 12),
((SELECT category_id FROM catalogue.category WHERE name = 'Kolači'), 'S-003', 'Krempita', 'Domaća, s vanilija kremom', 3.20, 0),
((SELECT category_id FROM catalogue.category WHERE name = 'Kolači'), 'S-004', 'Palačinka s Nutellom', '2 komada, posipana lješnjacima', 3.00, 2),
((SELECT category_id FROM catalogue.category WHERE name = 'Kolači'), 'S-005', 'Tiramisu', 'Mascarpone, kava i piškote', 3.50, 6),
((SELECT category_id FROM catalogue.category WHERE name = 'Torte'),  'S-006', 'Cheesecake od jagoda', 'Kriška, preljev od svježih jagoda', 4.20, 1);