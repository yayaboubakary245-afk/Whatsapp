import React, { useEffect, useState } from 'react';
import { Copy, Check, Search, FileCode2 } from 'lucide-react';

interface FileItem {
  path: string;
  name: string;
  category: 'go' | 'android' | 'sql' | 'docs';
  size: number;
  lines: number;
}

export const CodeExplorer: React.FC = () => {
  const [files, setFiles] = useState<FileItem[]>([]);
  const [counts, setCounts] = useState({ total: 0, go: 0, android: 0, sql: 0 });
  const [filter, setFilter] = useState<'all' | 'go' | 'android' | 'sql'>('all');
  const [query, setQuery] = useState('');
  const [selectedPath, setSelectedPath] = useState('server-go/cmd/server/main.go');
  const [content, setContent] = useState('');
  const [loading, setLoading] = useState(true);
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    fetch('/api/project/tree')
      .then((r) => r.json())
      .then((data) => {
        setFiles(data.files || []);
        setCounts(data.counts || { total: 0, go: 0, android: 0, sql: 0 });
        setLoading(false);
      })
      .catch(() => setLoading(false));
  }, []);

  useEffect(() => {
    if (!selectedPath) return;
    fetch(`/api/project/file?path=${encodeURIComponent(selectedPath)}`)
      .then((r) => r.json())
      .then((data) => {
        setContent(data.content || '');
      })
      .catch(() => setContent('Impossible de charger ce fichier.'));
  }, [selectedPath]);

  const filtered = files.filter((f) => {
    const matchCat = filter === 'all' || f.category === filter;
    const matchQuery = !query.trim() || f.path.toLowerCase().includes(query.toLowerCase());
    return matchCat && matchQuery;
  });

  const handleCopy = () => {
    navigator.clipboard.writeText(content);
    setCopied(true);
    setTimeout(() => setCopied(false), 1500);
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-4 bg-white border border-slate-200 rounded-xl p-4">
        <div className="flex items-center gap-1 p-1 bg-slate-100 rounded-lg">
          {(
            [
              ['all', `Tous (${counts.total})`],
              ['go', `Serveur Go (${counts.go})`],
              ['android', `Android Java (${counts.android})`],
              ['sql', `PostgreSQL (${counts.sql})`],
            ] as const
          ).map(([id, label]) => (
            <button
              key={id}
              type="button"
              onClick={() => setFilter(id)}
              className={`px-3 py-1.5 text-xs font-medium rounded-md transition-colors whitespace-nowrap tabular-nums ${
                filter === id ? 'bg-white text-slate-900 shadow-sm' : 'text-slate-600 hover:text-slate-900'
              }`}
            >
              {label}
            </button>
          ))}
        </div>

        <div className="relative w-full sm:w-72">
          <Search className="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
          <input
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Filtrer par chemin de fichier..."
            className="w-full pl-9 pr-3 py-2 text-xs border border-slate-200 rounded-lg focus:outline-none focus:border-[#0B6E4F]"
          />
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 items-stretch">
        <div className="lg:col-span-4 bg-white border border-slate-200 rounded-xl overflow-hidden flex flex-col max-h-[620px]">
          <div className="px-4 py-3 border-b border-slate-200 text-xs font-semibold text-slate-700">
            Arborescence du projet ({filtered.length} fichiers)
          </div>
          <div className="overflow-y-auto divide-y divide-slate-100">
            {loading ? (
              <div className="p-4 text-xs text-slate-500">Chargement de l’arborescence...</div>
            ) : filtered.length === 0 ? (
              <div className="p-4 text-xs text-slate-500">Aucun fichier ne correspond à votre recherche.</div>
            ) : (
              filtered.map((f) => (
                <button
                  key={f.path}
                  type="button"
                  onClick={() => setSelectedPath(f.path)}
                  className={`w-full text-left px-4 py-2.5 flex items-center justify-between gap-2 text-xs transition-colors ${
                    selectedPath === f.path ? 'bg-[#0B6E4F]/8 text-[#0B6E4F] font-medium' : 'text-slate-700 hover:bg-slate-50'
                  }`}
                >
                  <span className="truncate font-mono">{f.path}</span>
                  <span className="text-[11px] text-slate-400 tabular-nums shrink-0">{f.lines} L</span>
                </button>
              ))
            )}
          </div>
        </div>

        <div className="lg:col-span-8 bg-slate-900 border border-slate-800 rounded-xl overflow-hidden flex flex-col max-h-[620px]">
          <div className="px-4 py-3 border-b border-slate-800 flex items-center justify-between gap-4">
            <div className="flex items-center gap-2 text-xs font-mono text-slate-200 truncate">
              <FileCode2 className="w-4 h-4 text-emerald-400 shrink-0" />
              <span className="truncate">{selectedPath}</span>
            </div>
            <button
              type="button"
              onClick={handleCopy}
              className="px-3 py-1.5 text-xs font-medium text-slate-200 bg-slate-800 hover:bg-slate-700 rounded-md flex items-center gap-1.5 whitespace-nowrap shrink-0"
            >
              {copied ? <Check className="w-3.5 h-3.5 text-emerald-400" /> : <Copy className="w-3.5 h-3.5" />}
              {copied ? 'Copié' : 'Copier'}
            </button>
          </div>
          <pre className="p-4 overflow-auto text-xs font-mono text-slate-200 leading-relaxed flex-1">
            <code>{content}</code>
          </pre>
        </div>
      </div>
    </div>
  );
};
