import React, { useState } from 'react';
import {
  Phone,
  Video,
  Send,
  Paperclip,
  Users,
  Megaphone,
  Shield,
  ArrowLeft,
  CheckCheck,
  Mic,
  MicOff,
  Volume2,
  PhoneOff,
  Plus,
  UserPlus,
  Bell,
} from 'lucide-react';
import { GROUP_TYPES } from '../data/campusSchema';

type Screen =
  | 'LoginActivity'
  | 'ConversationsActivity'
  | 'ChatActivity'
  | 'GroupsActivity'
  | 'CreateGroupActivity'
  | 'AnnouncementsActivity'
  | 'ContactsActivity'
  | 'AdminActivity'
  | 'CallActivity';

type Role = 'etudiant' | 'enseignant' | 'admin';

interface SimMessage {
  id: string;
  sender: string;
  mine: boolean;
  text: string;
  time: string;
  status: 'sent' | 'delivered' | 'read';
  attachment?: string;
}

interface SimGroup {
  id: string;
  name: string;
  type: string;
  members: number;
}

interface SimAnnouncement {
  id: string;
  title: string;
  category: string;
  author: string;
  content: string;
  date: string;
}

export const AndroidSimulator: React.FC = () => {
  const [screen, setScreen] = useState<Screen>('ConversationsActivity');
  const [role, setRole] = useState<Role>('etudiant');
  const [emailInput, setEmailInput] = useState('yaya.boubakary@ensptic.cm');
  const [loginError, setLoginError] = useState('');

  const [activePeer, setActivePeer] = useState('Ndue Christian (Backend Go)');
  const [messages, setMessages] = useState<SimMessage[]>([
    {
      id: 'm1',
      sender: 'Ndue Christian',
      mine: false,
      text: 'Le serveur TCP sur le port 5000 et UDP 5001 est opérationnel.',
      time: '09:12',
      status: 'read',
    },
    {
      id: 'm2',
      sender: 'Yaya Boubakary',
      mine: true,
      text: 'Parfait, ApiClient.java et WebRTCManager.java sont prêts sur Android Studio.',
      time: '09:14',
      status: 'read',
    },
    {
      id: 'm3',
      sender: 'Mengue',
      mine: false,
      text: 'Les 22 tables PostgreSQL sont migrées (S1 à S6).',
      time: '09:16',
      status: 'read',
      attachment: 'schema_22_tables_supptic.pdf',
    },
  ]);
  const [draft, setDraft] = useState('');

  const [groups, setGroups] = useState<SimGroup[]>([
    { id: 'g1', name: 'IGL Niveau 4 — Architecture Réseau', type: 'classe', members: 42 },
    { id: 'g2', name: 'Projet WhatsApp Campus (G12)', type: 'projet', members: 3 },
    { id: 'g3', name: 'TD Systèmes Distribués & Go', type: 'td', members: 28 },
    { id: 'g4', name: 'Département TIC — Officiel', type: 'departement', members: 185 },
  ]);
  const [newGroupName, setNewGroupName] = useState('');
  const [newGroupType, setNewGroupType] = useState('projet');

  const [announcements, setAnnouncements] = useState<SimAnnouncement[]>([
    {
      id: 'a1',
      title: 'Soutenance Projets Génie Logiciel & Réseaux',
      category: 'examen',
      author: 'Direction Académique',
      content: 'Les soutenances du Groupe 12 se tiendront en Amphi B le 24 novembre 2026 à 09h00.',
      date: '22 nov. 2026',
    },
    {
      id: 'a2',
      title: 'Mise à jour Emploi du Temps — Semaine 10',
      category: 'emploi_du_temps',
      author: 'Pr. Mballa (Enseignant)',
      content: 'Le cours de Sécurité des Réseaux est déplacé au jeudi 14h00 en Salle T3.',
      date: '21 nov. 2026',
    },
  ]);
  const [newAnnTitle, setNewAnnTitle] = useState('');
  const [newAnnContent, setNewAnnContent] = useState('');

  const [callActive, setCallActive] = useState(false);
  const [callMuted, setCallMuted] = useState(false);
  const [callSpeaker, setCallSpeaker] = useState(false);
  const [signalLog, setSignalLog] = useState<string[]>([
    'POST /api/login -> JWT HS256 délivré par Go',
    'WS /api/ws -> Connexion gorilla/websocket établie',
    'UDP :5001 -> PING -> PONG (0.4 ms)',
  ]);

  const appendLog = (entry: string) => {
    setSignalLog((prev) => [entry, ...prev].slice(0, 8));
  };

  const handleLogin = (e: React.FormEvent) => {
    e.preventDefault();
    if (!emailInput.trim().toLowerCase().endsWith('@ensptic.cm')) {
      setLoginError('Domaine @ensptic.cm exigé par le serveur Go');
      return;
    }
    setLoginError('');
    appendLog(`POST /api/login (${emailInput}) -> 200 OK`);
    setScreen('ConversationsActivity');
  };

  const handleSendMessage = (e: React.FormEvent) => {
    e.preventDefault();
    if (!draft.trim()) return;
    const msg: SimMessage = {
      id: `m-${Date.now()}`,
      sender: 'Moi',
      mine: true,
      text: draft.trim(),
      time: new Date().toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' }),
      status: 'delivered',
    };
    setMessages((prev) => [...prev, msg]);
    setDraft('');
    appendLog(`POST /api/messages -> Persisté PostgreSQL + Push TCP :5000`);
  };

  const handleAttachPdf = () => {
    const msg: SimMessage = {
      id: `m-${Date.now()}`,
      sender: 'Moi',
      mine: true,
      text: 'Support de cours partagé',
      time: new Date().toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' }),
      status: 'delivered',
      attachment: 'TP_Sockets_Golang_SUPPTIC.pdf',
    };
    setMessages((prev) => [...prev, msg]);
    appendLog('POST /api/attachments/upload -> SHA-256 + table attachments');
  };

  const handleCreateGroup = (e: React.FormEvent) => {
    e.preventDefault();
    if (!newGroupName.trim()) return;
    setGroups((prev) => [
      { id: `g-${Date.now()}`, name: newGroupName.trim(), type: newGroupType, members: 1 },
      ...prev,
    ]);
    appendLog(`POST /api/groups (type=${newGroupType}) -> 201 Created`);
    setNewGroupName('');
    setScreen('GroupsActivity');
  };

  const handlePublishAnnouncement = (e: React.FormEvent) => {
    e.preventDefault();
    if (role === 'etudiant') return;
    if (!newAnnTitle.trim() || !newAnnContent.trim()) return;
    setAnnouncements((prev) => [
      {
        id: `a-${Date.now()}`,
        title: newAnnTitle.trim(),
        category: 'administrative',
        author: role === 'admin' ? 'Administrateur SUP’PTIC' : 'Enseignant SUP’PTIC',
        content: newAnnContent.trim(),
        date: "Aujourd'hui",
      },
      ...prev,
    ]);
    appendLog('POST /api/announcements -> Broadcast WebSocket + table announcements');
    setNewAnnTitle('');
    setNewAnnContent('');
  };

  const startCall = (type: 'audio' | 'video') => {
    setCallActive(true);
    setScreen('CallActivity');
    appendLog(`POST /api/call/start (${type}) -> SDP Offer via gorilla/websocket`);
    appendLog('WebRTCManager.java -> ICE P2P UDP media établi (STUN Google)');
  };

  return (
    <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
      {/* Colonne gauche : Contrôles de rôle & Trace réseau Go */}
      <div className="lg:col-span-5 space-y-6">
        <div className="bg-white border border-slate-200 rounded-xl p-6">
          <h2 className="text-base font-semibold text-slate-900">Contexte utilisateur simulé</h2>
          <p className="text-sm text-slate-600 mt-1">
            Basculez entre les 3 rôles RBAC définis dans PostgreSQL pour tester les permissions (F1.6, F7.5, F9).
          </p>

          <div className="mt-4 flex items-center gap-1 p-1 bg-slate-100 rounded-lg">
            {(['etudiant', 'enseignant', 'admin'] as Role[]).map((r) => (
              <button
                key={r}
                type="button"
                onClick={() => {
                  setRole(r);
                  appendLog(`Rôle actif basculé sur ${r}`);
                }}
                className={`flex-1 py-2 px-3 text-xs font-medium rounded-md transition-colors whitespace-nowrap ${
                  role === r ? 'bg-white text-slate-900 shadow-sm' : 'text-slate-600 hover:text-slate-900'
                }`}
              >
                {r === 'etudiant' ? 'Étudiant' : r === 'enseignant' ? 'Enseignant' : 'Administrateur'}
              </button>
            ))}
          </div>

          <div className="mt-5 pt-5 border-t border-slate-100">
            <h3 className="text-xs font-semibold text-slate-700">Navigation directe (16 écrans Android)</h3>
            <div className="mt-3 grid grid-cols-2 gap-2">
              {(
                [
                  ['LoginActivity', 'Login @ensptic.cm'],
                  ['ConversationsActivity', 'Conversations'],
                  ['ChatActivity', 'Chat & Fichiers'],
                  ['GroupsActivity', 'Groupes (7 types)'],
                  ['CreateGroupActivity', 'Créer un groupe'],
                  ['AnnouncementsActivity', 'Annonces'],
                  ['ContactsActivity', 'Contacts'],
                  ['AdminActivity', 'Administration'],
                ] as [Screen, string][]
              ).map(([scr, label]) => (
                <button
                  key={scr}
                  type="button"
                  onClick={() => setScreen(scr)}
                  className={`text-left px-3 py-2 rounded-lg text-xs font-medium border transition-colors truncate ${
                    screen === scr
                      ? 'border-[#0B6E4F] bg-[#0B6E4F]/5 text-[#0B6E4F]'
                      : 'border-slate-200 text-slate-700 hover:bg-slate-50'
                  }`}
                >
                  {label}
                </button>
              ))}
            </div>
          </div>
        </div>

        <div className="bg-white border border-slate-200 rounded-xl p-6">
          <h3 className="text-sm font-semibold text-slate-900">Journal des échanges Android ↔ Serveur Go</h3>
          <p className="text-xs text-slate-500 mt-1">
            Flux REST (ApiClient.java), Signalisation WebSocket (gorilla/websocket) et Sockets Go (:5000 / :5001)
          </p>
          <ul className="mt-4 space-y-2 font-mono text-xs text-slate-700 tabular-nums">
            {signalLog.map((item, idx) => (
              <li key={idx} className="py-1.5 px-2.5 bg-slate-50 border border-slate-200 rounded-md break-words">
                {item}
              </li>
            ))}
          </ul>
        </div>
      </div>

      {/* Colonne droite : Émulateur Mobile Android Studio */}
      <div className="lg:col-span-7 flex justify-center">
        <div className="w-full max-w-[400px] bg-slate-900 rounded-[36px] p-3 shadow-lg border border-slate-800">
          <div className="bg-white rounded-[28px] overflow-hidden flex flex-col h-[680px]">
            {/* Barre supérieure Android Material */}
            <div className="bg-[#0B6E4F] text-white px-4 py-3.5 flex items-center justify-between shrink-0">
              <div className="flex items-center gap-2.5 min-w-0">
                {screen !== 'ConversationsActivity' && screen !== 'LoginActivity' && (
                  <button
                    type="button"
                    onClick={() => setScreen('ConversationsActivity')}
                    className="p-1 rounded hover:bg-white/10"
                    aria-label="Retour"
                  >
                    <ArrowLeft className="w-4 h-4" />
                  </button>
                )}
                <div className="min-w-0">
                  <div className="text-sm font-semibold truncate">
                    {screen === 'LoginActivity' && 'WhatsApp Campus'}
                    {screen === 'ConversationsActivity' && 'WhatsApp Campus'}
                    {screen === 'ChatActivity' && activePeer}
                    {screen === 'GroupsActivity' && 'Groupes Universitaires'}
                    {screen === 'CreateGroupActivity' && 'Nouveau Groupe'}
                    {screen === 'AnnouncementsActivity' && "Annonces SUP'PTIC"}
                    {screen === 'ContactsActivity' && 'Contacts Campus'}
                    {screen === 'AdminActivity' && 'Espace Administration'}
                    {screen === 'CallActivity' && 'Appel P2P WebRTC'}
                  </div>
                  <div className="text-[11px] text-emerald-100 truncate">
                    {screen} · {role}
                  </div>
                </div>
              </div>

              {screen === 'ChatActivity' && (
                <div className="flex items-center gap-1.5 shrink-0">
                  <button
                    type="button"
                    onClick={() => startCall('audio')}
                    className="p-2 rounded-lg hover:bg-white/15"
                    title="Appel audio WebRTC"
                  >
                    <Phone className="w-4 h-4" />
                  </button>
                  <button
                    type="button"
                    onClick={() => startCall('video')}
                    className="p-2 rounded-lg hover:bg-white/15"
                    title="Appel vidéo WebRTC"
                  >
                    <Video className="w-4 h-4" />
                  </button>
                </div>
              )}
            </div>

            {/* Corps de l'écran actif */}
            <div className="flex-1 overflow-y-auto bg-slate-50">
              {screen === 'LoginActivity' && (
                <form onSubmit={handleLogin} className="p-6 flex flex-col justify-center h-full space-y-4">
                  <div>
                    <h3 className="text-lg font-bold text-slate-900">Connexion SUP’PTIC</h3>
                    <p className="text-xs text-slate-500 mt-1">
                      Authentification JWT via POST /api/login (domaine @ensptic.cm)
                    </p>
                  </div>
                  <div>
                    <label className="block text-xs font-medium text-slate-700 mb-1">E-mail universitaire</label>
                    <input
                      type="email"
                      value={emailInput}
                      onChange={(e) => setEmailInput(e.target.value)}
                      className="w-full px-3 py-2.5 text-sm border border-slate-300 rounded-lg bg-white"
                      placeholder="prenom.nom@ensptic.cm"
                    />
                  </div>
                  {loginError && <p className="text-xs text-red-600 font-medium">{loginError}</p>}
                  <button
                    type="submit"
                    className="w-full py-2.5 bg-[#0B6E4F] text-white text-sm font-medium rounded-lg hover:bg-[#08563D]"
                  >
                    Se connecter
                  </button>
                </form>
              )}

              {screen === 'ConversationsActivity' && (
                <div>
                  <div className="grid grid-cols-4 border-b border-slate-200 bg-white text-[11px] font-medium text-slate-600">
                    <button
                      type="button"
                      onClick={() => setScreen('GroupsActivity')}
                      className="py-2.5 flex flex-col items-center gap-1 hover:text-[#0B6E4F]"
                    >
                      <Users className="w-4 h-4" />
                      Groupes
                    </button>
                    <button
                      type="button"
                      onClick={() => setScreen('AnnouncementsActivity')}
                      className="py-2.5 flex flex-col items-center gap-1 hover:text-[#0B6E4F]"
                    >
                      <Megaphone className="w-4 h-4" />
                      Annonces
                    </button>
                    <button
                      type="button"
                      onClick={() => setScreen('ContactsActivity')}
                      className="py-2.5 flex flex-col items-center gap-1 hover:text-[#0B6E4F]"
                    >
                      <UserPlus className="w-4 h-4" />
                      Contacts
                    </button>
                    <button
                      type="button"
                      onClick={() => setScreen('AdminActivity')}
                      className="py-2.5 flex flex-col items-center gap-1 hover:text-[#0B6E4F]"
                    >
                      <Shield className="w-4 h-4" />
                      Admin
                    </button>
                  </div>

                  <div className="divide-y divide-slate-200 bg-white">
                    {[
                      {
                        name: 'Ndue Christian (Backend Go)',
                        last: 'Serveur TCP :5000 & UDP :5001 actifs',
                        time: '09:16',
                      },
                      {
                        name: 'Mengue (Base PostgreSQL)',
                        last: 'schema_22_tables_supptic.pdf',
                        time: '09:12',
                      },
                      {
                        name: 'IGL Niveau 4 — Architecture Réseau',
                        last: 'Support de cours TCP/UDP disponible',
                        time: 'Hier',
                      },
                    ].map((conv) => (
                      <button
                        key={conv.name}
                        type="button"
                        onClick={() => {
                          setActivePeer(conv.name);
                          setScreen('ChatActivity');
                        }}
                        className="w-full text-left px-4 py-3.5 hover:bg-slate-50 flex items-start justify-between gap-3"
                      >
                        <div className="min-w-0">
                          <div className="text-sm font-semibold text-slate-900 truncate">{conv.name}</div>
                          <div className="text-xs text-slate-500 truncate mt-0.5">{conv.last}</div>
                        </div>
                        <span className="text-[11px] text-slate-400 tabular-nums shrink-0">{conv.time}</span>
                      </button>
                    ))}
                  </div>
                </div>
              )}

              {screen === 'ChatActivity' && (
                <div className="flex flex-col h-full">
                  <div className="flex-1 p-3 space-y-2.5 overflow-y-auto">
                    {messages.map((m) => (
                      <div key={m.id} className={`flex ${m.mine ? 'justify-end' : 'justify-start'}`}>
                        <div
                          className={`max-w-[82%] rounded-xl px-3 py-2 text-xs ${
                            m.mine
                              ? 'bg-[#0B6E4F] text-white'
                              : 'bg-white border border-slate-200 text-slate-800'
                          }`}
                        >
                          {!m.mine && <div className="font-semibold text-[11px] text-[#0B6E4F] mb-0.5">{m.sender}</div>}
                          {m.attachment && (
                            <div
                              className={`mb-1.5 px-2 py-1 rounded text-[11px] font-mono ${
                                m.mine ? 'bg-white/15 text-white' : 'bg-slate-100 text-slate-700'
                              }`}
                            >
                              PDF · {m.attachment}
                            </div>
                          )}
                          <p className="leading-relaxed">{m.text}</p>
                          <div
                            className={`mt-1 flex items-center justify-end gap-1 text-[10px] tabular-nums ${
                              m.mine ? 'text-emerald-100' : 'text-slate-400'
                            }`}
                          >
                            <span>{m.time}</span>
                            {m.mine && <CheckCheck className="w-3 h-3" />}
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>

                  <form onSubmit={handleSendMessage} className="p-2.5 bg-white border-t border-slate-200 flex items-center gap-2">
                    <button
                      type="button"
                      onClick={handleAttachPdf}
                      className="p-2 text-slate-500 hover:text-[#0B6E4F] rounded-lg hover:bg-slate-100"
                      title="Joindre un document PDF"
                    >
                      <Paperclip className="w-4 h-4" />
                    </button>
                    <input
                      type="text"
                      value={draft}
                      onChange={(e) => setDraft(e.target.value)}
                      placeholder="Message (TCP/REST Go)..."
                      className="flex-1 px-3 py-2 text-xs border border-slate-200 rounded-lg focus:outline-none focus:border-[#0B6E4F]"
                    />
                    <button
                      type="submit"
                      className="p-2 bg-[#0B6E4F] text-white rounded-lg hover:bg-[#08563D]"
                      aria-label="Envoyer"
                    >
                      <Send className="w-4 h-4" />
                    </button>
                  </form>
                </div>
              )}

              {screen === 'GroupsActivity' && (
                <div className="p-3 space-y-3">
                  <button
                    type="button"
                    onClick={() => setScreen('CreateGroupActivity')}
                    className="w-full py-2.5 px-3 bg-[#0B6E4F] text-white text-xs font-medium rounded-lg flex items-center justify-center gap-2"
                  >
                    <Plus className="w-4 h-4" />
                    Créer un groupe universitaire (7 types)
                  </button>
                  <div className="space-y-2">
                    {groups.map((g) => (
                      <button
                        key={g.id}
                        type="button"
                        onClick={() => {
                          setActivePeer(g.name);
                          setScreen('ChatActivity');
                        }}
                        className="w-full text-left p-3 bg-white border border-slate-200 rounded-lg hover:border-[#0B6E4F]"
                      >
                        <div className="text-sm font-semibold text-slate-900">{g.name}</div>
                        <div className="text-xs text-slate-500 mt-0.5 tabular-nums">
                          Type : {g.type} · {g.members} membres
                        </div>
                      </button>
                    ))}
                  </div>
                </div>
              )}

              {screen === 'CreateGroupActivity' && (
                <form onSubmit={handleCreateGroup} className="p-4 space-y-4">
                  <div>
                    <label className="block text-xs font-medium text-slate-700 mb-1">Nom du groupe</label>
                    <input
                      type="text"
                      value={newGroupName}
                      onChange={(e) => setNewGroupName(e.target.value)}
                      placeholder="Ex: TD Réseaux TCP/UDP"
                      className="w-full px-3 py-2 text-sm border border-slate-300 rounded-lg bg-white"
                    />
                  </div>
                  <div>
                    <label className="block text-xs font-medium text-slate-700 mb-1">Type universitaire (7 types)</label>
                    <select
                      value={newGroupType}
                      onChange={(e) => setNewGroupType(e.target.value)}
                      className="w-full px-3 py-2 text-sm border border-slate-300 rounded-lg bg-white"
                    >
                      {GROUP_TYPES.map((gt) => (
                        <option key={gt.id} value={gt.id}>
                          {gt.label} ({gt.id})
                        </option>
                      ))}
                    </select>
                  </div>
                  <button
                    type="submit"
                    className="w-full py-2.5 bg-[#0B6E4F] text-white text-sm font-medium rounded-lg"
                  >
                    Créer via POST /api/groups
                  </button>
                </form>
              )}

              {screen === 'AnnouncementsActivity' && (
                <div className="p-3 space-y-3">
                  {role !== 'etudiant' ? (
                    <form onSubmit={handlePublishAnnouncement} className="p-3 bg-white border border-slate-200 rounded-lg space-y-2">
                      <div className="text-xs font-semibold text-slate-800">Publier une annonce officielle</div>
                      <input
                        type="text"
                        value={newAnnTitle}
                        onChange={(e) => setNewAnnTitle(e.target.value)}
                        placeholder="Titre de l'annonce"
                        className="w-full px-2.5 py-1.5 text-xs border border-slate-300 rounded"
                      />
                      <textarea
                        rows={2}
                        value={newAnnContent}
                        onChange={(e) => setNewAnnContent(e.target.value)}
                        placeholder="Contenu..."
                        className="w-full px-2.5 py-1.5 text-xs border border-slate-300 rounded"
                      />
                      <button type="submit" className="w-full py-2 bg-[#0B6E4F] text-white text-xs font-medium rounded">
                        Diffuser (POST /api/announcements)
                      </button>
                    </form>
                  ) : (
                    <div className="p-3 bg-amber-50 border border-amber-200 rounded-lg text-xs text-amber-800">
                      Lecture seule : seuls les rôles Enseignant et Admin peuvent publier (F7.5).
                    </div>
                  )}

                  {announcements.map((a) => (
                    <div key={a.id} className="p-3 bg-white border border-slate-200 rounded-lg">
                      <div className="text-[11px] text-[#0B6E4F] font-medium">
                        {a.category} · {a.date}
                      </div>
                      <div className="text-sm font-semibold text-slate-900 mt-0.5">{a.title}</div>
                      <p className="text-xs text-slate-600 mt-1 leading-relaxed">{a.content}</p>
                      <div className="text-[11px] text-slate-400 mt-1.5">{a.author}</div>
                    </div>
                  ))}
                </div>
              )}

              {screen === 'ContactsActivity' && (
                <div className="divide-y divide-slate-200 bg-white">
                  {[
                    { name: 'Ndue Christian', email: 'christian.ndue@ensptic.cm', role: 'Développeur Go' },
                    { name: 'Mengue', email: 'mengue@ensptic.cm', role: 'Base PostgreSQL' },
                    { name: 'Yaya Boubakary', email: 'yaya.boubakary@ensptic.cm', role: 'Lead Android & WebRTC' },
                  ].map((c) => (
                    <div key={c.email} className="px-4 py-3 flex items-center justify-between">
                      <div>
                        <div className="text-sm font-semibold text-slate-900">{c.name}</div>
                        <div className="text-xs text-slate-500">
                          {c.email} · {c.role}
                        </div>
                      </div>
                      <button
                        type="button"
                        onClick={() => {
                          setActivePeer(c.name);
                          setScreen('ChatActivity');
                        }}
                        className="px-2.5 py-1.5 text-xs font-medium text-[#0B6E4F] border border-[#0B6E4F]/30 rounded-lg hover:bg-[#0B6E4F]/5"
                      >
                        Discuter
                      </button>
                    </div>
                  ))}
                </div>
              )}

              {screen === 'AdminActivity' && (
                <div className="p-3 space-y-3">
                  {role !== 'admin' ? (
                    <div className="p-4 bg-red-50 border border-red-200 rounded-lg text-xs text-red-700">
                      Accès refusé (HTTP 403). Passez le rôle sur « Administrateur » à gauche pour consulter audit_logs et gérer les comptes.
                    </div>
                  ) : (
                    <div className="space-y-2.5">
                      <div className="p-3 bg-white border border-slate-200 rounded-lg">
                        <div className="text-xs font-semibold text-slate-900 flex items-center gap-1.5">
                          <Bell className="w-3.5 h-3.5 text-[#0B6E4F]" />
                          Journal d’audit (audit_logs)
                        </div>
                        <ul className="mt-2 space-y-1.5 font-mono text-[11px] text-slate-600">
                          <li>USER_LOGIN · christian.ndue@ensptic.cm</li>
                          <li>GROUP_CREATE · IGL Niveau 4</li>
                          <li>CALL_START · WebRTC audio P2P</li>
                        </ul>
                      </div>
                    </div>
                  )}
                </div>
              )}

              {screen === 'CallActivity' && (
                <div className="h-full bg-slate-900 text-white p-6 flex flex-col items-center justify-between">
                  <div className="text-center mt-8">
                    <div className="w-16 h-16 rounded-full bg-[#0B6E4F] flex items-center justify-center mx-auto text-xl font-bold">
                      {activePeer.slice(0, 2).toUpperCase()}
                    </div>
                    <h3 className="text-lg font-semibold mt-4">{activePeer}</h3>
                    <p className="text-xs text-emerald-300 mt-1">
                      {callActive ? 'Flux média P2P actif (org.webrtc)' : 'Appel terminé'}
                    </p>
                  </div>

                  <div className="flex items-center gap-4 mb-8">
                    <button
                      type="button"
                      onClick={() => setCallMuted(!callMuted)}
                      className={`p-3.5 rounded-full ${callMuted ? 'bg-amber-600' : 'bg-slate-700'}`}
                      aria-label="Muet"
                    >
                      {callMuted ? <MicOff className="w-5 h-5" /> : <Mic className="w-5 h-5" />}
                    </button>
                    <button
                      type="button"
                      onClick={() => setCallSpeaker(!callSpeaker)}
                      className={`p-3.5 rounded-full ${callSpeaker ? 'bg-[#0B6E4F]' : 'bg-slate-700'}`}
                      aria-label="Haut-parleur"
                    >
                      <Volume2 className="w-5 h-5" />
                    </button>
                    <button
                      type="button"
                      onClick={() => {
                        setCallActive(false);
                        appendLog('POST /api/call/end -> Durée enregistrée dans call_history');
                        setScreen('ChatActivity');
                      }}
                      className="p-3.5 rounded-full bg-red-600 hover:bg-red-700"
                      aria-label="Raccrocher"
                    >
                      <PhoneOff className="w-5 h-5" />
                    </button>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
