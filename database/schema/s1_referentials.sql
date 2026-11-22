-- Référence Schéma Sprint 1 — 4 tables initiales (Mengue)
-- tables: roles, departments, programs, academic_levels

-- Vérification rapide des tables S1 :
SELECT table_name
FROM information_schema.tables
WHERE table_schema = 'public'
  AND table_name IN ('roles', 'departments', 'programs', 'academic_levels')
ORDER BY table_name;
