export interface TableSpec {
  name: string;
  sprint: 'S1' | 'S2' | 'S3' | 'S4' | 'S5' | 'S6';
  owner: 'Mengue' | 'Yaya';
  description: string;
  columns: { name: string; type: string; constraint?: string }[];
}

export const POSTGRES_TABLES: TableSpec[] = [
  {
    name: 'roles',
    sprint: 'S1',
    owner: 'Mengue',
    description: 'Rôles et permissions RBAC (etudiant, enseignant, admin)',
    columns: [
      { name: 'id', type: 'SERIAL', constraint: 'PK' },
      { name: 'name', type: 'VARCHAR(32)', constraint: 'UNIQUE' },
      { name: 'permissions', type: 'JSONB' },
    ],
  },
  {
    name: 'departments',
    sprint: 'S1',
    owner: 'Mengue',
    description: "Départements académiques de SUP'PTIC",
    columns: [
      { name: 'id', type: 'SERIAL', constraint: 'PK' },
      { name: 'code', type: 'VARCHAR(16)', constraint: 'UNIQUE' },
      { name: 'name', type: 'VARCHAR(120)' },
    ],
  },
  {
    name: 'programs',
    sprint: 'S1',
    owner: 'Mengue',
    description: 'Filières rattachées aux départements (ITT, IT, IGL, CYB, IP)',
    columns: [
      { name: 'id', type: 'SERIAL', constraint: 'PK' },
      { name: 'department_id', type: 'INT', constraint: 'FK -> departments' },
      { name: 'code', type: 'VARCHAR(16)', constraint: 'UNIQUE' },
    ],
  },
  {
    name: 'academic_levels',
    sprint: 'S1',
    owner: 'Mengue',
    description: "Niveaux d'étude (L1 à M2)",
    columns: [
      { name: 'id', type: 'SERIAL', constraint: 'PK' },
      { name: 'code', type: 'VARCHAR(16)', constraint: 'UNIQUE' },
      { name: 'label', type: 'VARCHAR(64)' },
    ],
  },
  {
    name: 'users',
    sprint: 'S2',
    owner: 'Mengue',
    description: 'Comptes authentifiés par e-mail universitaire @ensptic.cm',
    columns: [
      { name: 'id', type: 'UUID', constraint: 'PK' },
      { name: 'email', type: 'VARCHAR(160)', constraint: 'UNIQUE @ensptic.cm' },
      { name: 'password_hash', type: 'TEXT' },
      { name: 'role_id', type: 'INT', constraint: 'FK -> roles' },
      { name: 'is_verified', type: 'BOOLEAN' },
      { name: 'is_suspended', type: 'BOOLEAN' },
    ],
  },
  {
    name: 'profiles',
    sprint: 'S2',
    owner: 'Yaya',
    description: 'Profils personnels, filière, bio et avatar',
    columns: [
      { name: 'id', type: 'UUID', constraint: 'PK' },
      { name: 'user_id', type: 'UUID', constraint: 'UNIQUE FK -> users' },
      { name: 'first_name', type: 'VARCHAR(80)' },
      { name: 'last_name', type: 'VARCHAR(80)' },
      { name: 'department_id', type: 'INT', constraint: 'FK -> departments' },
    ],
  },
  {
    name: 'contacts',
    sprint: 'S3',
    owner: 'Mengue',
    description: 'Répertoire de contacts et état de blocage',
    columns: [
      { name: 'id', type: 'UUID', constraint: 'PK' },
      { name: 'owner_id', type: 'UUID', constraint: 'FK -> users' },
      { name: 'contact_user_id', type: 'UUID', constraint: 'FK -> users' },
      { name: 'is_blocked', type: 'BOOLEAN' },
    ],
  },
  {
    name: 'conversations',
    sprint: 'S3',
    owner: 'Mengue',
    description: 'Fils de discussion privés (1-1) et collectifs (group)',
    columns: [
      { name: 'id', type: 'UUID', constraint: 'PK' },
      { name: 'type', type: 'VARCHAR(16)', constraint: 'private | group' },
      { name: 'last_message_at', type: 'TIMESTAMPTZ' },
    ],
  },
  {
    name: 'conversation_members',
    sprint: 'S3',
    owner: 'Mengue',
    description: 'Participants et horodatage de lecture',
    columns: [
      { name: 'id', type: 'UUID', constraint: 'PK' },
      { name: 'conversation_id', type: 'UUID', constraint: 'FK -> conversations' },
      { name: 'user_id', type: 'UUID', constraint: 'FK -> users' },
    ],
  },
  {
    name: 'messages',
    sprint: 'S3',
    owner: 'Mengue',
    description: 'Messages persistés avec citations et suppression logique',
    columns: [
      { name: 'id', type: 'UUID', constraint: 'PK' },
      { name: 'conversation_id', type: 'UUID', constraint: 'FK -> conversations' },
      { name: 'sender_id', type: 'UUID', constraint: 'FK -> users' },
      { name: 'reply_to_message_id', type: 'UUID', constraint: 'FK -> messages' },
      { name: 'content', type: 'TEXT' },
    ],
  },
  {
    name: 'message_status',
    sprint: 'S3',
    owner: 'Mengue',
    description: 'Accusés envoyé / reçu / lu par destinataire',
    columns: [
      { name: 'id', type: 'UUID', constraint: 'PK' },
      { name: 'message_id', type: 'UUID', constraint: 'FK -> messages' },
      { name: 'user_id', type: 'UUID', constraint: 'FK -> users' },
      { name: 'status', type: 'VARCHAR(16)', constraint: 'sent | delivered | read' },
    ],
  },
  {
    name: 'groups',
    sprint: 'S4',
    owner: 'Mengue',
    description: 'Groupes universitaires (7 types : classe, filière, niveau, département, TD, projet, pédagogique)',
    columns: [
      { name: 'id', type: 'UUID', constraint: 'PK' },
      { name: 'conversation_id', type: 'UUID', constraint: 'UNIQUE FK' },
      { name: 'name', type: 'VARCHAR(120)' },
      { name: 'group_type', type: 'VARCHAR(24)', constraint: '7 types' },
    ],
  },
  {
    name: 'group_members',
    sprint: 'S4',
    owner: 'Mengue',
    description: 'Membres des groupes, rôles et permissions de publication',
    columns: [
      { name: 'id', type: 'UUID', constraint: 'PK' },
      { name: 'group_id', type: 'UUID', constraint: 'FK -> groups' },
      { name: 'user_id', type: 'UUID', constraint: 'FK -> users' },
      { name: 'member_role', type: 'VARCHAR(16)', constraint: 'member | moderator | admin' },
    ],
  },
  {
    name: 'attachments',
    sprint: 'S4',
    owner: 'Mengue',
    description: 'Références des documents pédagogiques (PDF, Word, PPT, images, vidéos)',
    columns: [
      { name: 'id', type: 'UUID', constraint: 'PK' },
      { name: 'message_id', type: 'UUID', constraint: 'FK -> messages' },
      { name: 'file_category', type: 'VARCHAR(24)' },
      { name: 'checksum_sha256', type: 'VARCHAR(64)' },
    ],
  },
  {
    name: 'announcements',
    sprint: 'S5',
    owner: 'Mengue',
    description: 'Annonces officielles administratives, examens et emplois du temps',
    columns: [
      { name: 'id', type: 'UUID', constraint: 'PK' },
      { name: 'author_id', type: 'UUID', constraint: 'FK -> users' },
      { name: 'category', type: 'VARCHAR(32)' },
      { name: 'is_pinned', type: 'BOOLEAN' },
    ],
  },
  {
    name: 'notifications',
    sprint: 'S5',
    owner: 'Mengue',
    description: 'Notifications push (messages, mentions, annonces, appels)',
    columns: [
      { name: 'id', type: 'UUID', constraint: 'PK' },
      { name: 'user_id', type: 'UUID', constraint: 'FK -> users' },
      { name: 'type', type: 'VARCHAR(32)' },
      { name: 'is_read', type: 'BOOLEAN' },
    ],
  },
  {
    name: 'reports',
    sprint: 'S6',
    owner: 'Mengue',
    description: 'Signalements de contenus ou comportements',
    columns: [
      { name: 'id', type: 'UUID', constraint: 'PK' },
      { name: 'reporter_id', type: 'UUID', constraint: 'FK -> users' },
      { name: 'status', type: 'VARCHAR(24)' },
    ],
  },
  {
    name: 'audit_logs',
    sprint: 'S6',
    owner: 'Mengue',
    description: "Journalisation des opérations sensibles d'administration",
    columns: [
      { name: 'id', type: 'UUID', constraint: 'PK' },
      { name: 'actor_id', type: 'UUID', constraint: 'FK -> users' },
      { name: 'action', type: 'VARCHAR(80)' },
      { name: 'metadata', type: 'JSONB' },
    ],
  },
  {
    name: 'calls',
    sprint: 'S6',
    owner: 'Mengue',
    description: 'Métadonnées des appels audio/vidéo WebRTC',
    columns: [
      { name: 'id', type: 'UUID', constraint: 'PK' },
      { name: 'caller_id', type: 'UUID', constraint: 'FK -> users' },
      { name: 'call_type', type: 'VARCHAR(16)', constraint: 'audio | video' },
      { name: 'duration_seconds', type: 'INT' },
    ],
  },
  {
    name: 'call_history',
    sprint: 'S6',
    owner: 'Mengue',
    description: 'Historique entrant, sortant et manqué par utilisateur',
    columns: [
      { name: 'id', type: 'UUID', constraint: 'PK' },
      { name: 'call_id', type: 'UUID', constraint: 'FK -> calls' },
      { name: 'direction', type: 'VARCHAR(16)' },
    ],
  },
  {
    name: 'call_participants',
    sprint: 'S6',
    owner: 'Mengue',
    description: "Participants d'un appel et état micro/vidéo",
    columns: [
      { name: 'id', type: 'UUID', constraint: 'PK' },
      { name: 'call_id', type: 'UUID', constraint: 'FK -> calls' },
      { name: 'user_id', type: 'UUID', constraint: 'FK -> users' },
    ],
  },
  {
    name: 'ice_servers',
    sprint: 'S6',
    owner: 'Mengue',
    description: 'Configuration STUN/TURN pour WebRTC',
    columns: [
      { name: 'id', type: 'SERIAL', constraint: 'PK' },
      { name: 'urls', type: 'TEXT' },
      { name: 'server_type', type: 'VARCHAR(16)', constraint: 'stun | turn' },
    ],
  },
];

export const GROUP_TYPES = [
  { id: 'classe', label: 'Classe' },
  { id: 'filiere', label: 'Filière' },
  { id: 'niveau', label: 'Niveau' },
  { id: 'departement', label: 'Département' },
  { id: 'td', label: 'TD / TP' },
  { id: 'projet', label: 'Projet' },
  { id: 'pedagogique', label: 'Pédagogique' },
] as const;
