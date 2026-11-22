-- Sprint 1 Seed : Rôles, Départements, Filières et Niveaux académiques SUP'PTIC
-- Responsable: Mengue

INSERT INTO roles (name, description, permissions) VALUES
('etudiant', 'Étudiant de SUP''PTIC', '{"can_message": true, "can_create_project_group": true, "can_publish_announcement": false, "can_admin": false}'::jsonb),
('enseignant', 'Enseignant / Chercheur', '{"can_message": true, "can_create_pedagogical_group": true, "can_publish_announcement": true, "can_admin": false}'::jsonb),
('admin', 'Administrateur du campus numérique', '{"can_message": true, "can_create_all_groups": true, "can_publish_announcement": true, "can_admin": true}'::jsonb)
ON CONFLICT (name) DO NOTHING;

INSERT INTO departments (code, name, description) VALUES
('TEL', 'Télécommunications & Réseaux', 'Ingénierie télécoms, réseaux IP et systèmes sans fil'),
('TIC', 'Technologies de l''Information et de la Communication', 'Génie logiciel, cybersécurité et systèmes distribués'),
('POS', 'Management des Services Postaux & Financiers', 'Économie numérique et services postaux')
ON CONFLICT (code) DO NOTHING;

INSERT INTO programs (department_id, code, name)
SELECT d.id, p.code, p.name FROM (
    VALUES
        ('TEL', 'ITT', 'Ingénieur des Travaux des Télécommunications'),
        ('TEL', 'IT', 'Ingénieur des Télécommunications'),
        ('TIC', 'IGL', 'Ingénierie Génie Logiciel & Systèmes'),
        ('TIC', 'CYB', 'Cybersécurité et Réseaux d''Entreprise'),
        ('POS', 'IP', 'Inspecteur des Postes et Services Financiers')
) AS p(dept_code, code, name)
JOIN departments d ON d.code = p.dept_code
ON CONFLICT (code) DO NOTHING;

INSERT INTO academic_levels (code, label, order_index) VALUES
('L1', 'Niveau 1 (Année préparatoire)', 1),
('L2', 'Niveau 2', 2),
('L3', 'Niveau 3 (Licence / Tronc commun)', 3),
('M1', 'Niveau 4 (Master 1 / Cycle Ingénieur)', 4),
('M2', 'Niveau 5 (Master 2 / Fin de cycle)', 5)
ON CONFLICT (code) DO NOTHING;
