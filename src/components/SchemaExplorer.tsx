import React, { useState } from 'react';
import { POSTGRES_TABLES } from '../data/campusSchema';

export const SchemaExplorer: React.FC = () => {
  const [sprintFilter, setSprintFilter] = useState<'ALL' | 'S1' | 'S2' | 'S3' | 'S4' | 'S5' | 'S6'>('ALL');

  const visibleTables = POSTGRES_TABLES.filter(
    (t) => sprintFilter === 'ALL' || t.sprint === sprintFilter
  );

  return (
    <div className="space-y-6">
      <div className="bg-white border border-slate-200 rounded-xl p-6 flex flex-wrap items-center justify-between gap-4">
        <div>
          <h2 className="text-base font-semibold text-slate-900">
            Modèle Relationnel PostgreSQL (22 Tables · Sprints S1 à S6)
          </h2>
          <p className="text-sm text-slate-600 mt-1">
            Conçu par Mengue et accessible exclusivement par le serveur Go via pgx/v5 et requêtes SQL paramétrées.
          </p>
        </div>

        <div className="flex items-center gap-1 p-1 bg-slate-100 rounded-lg">
          {(['ALL', 'S1', 'S2', 'S3', 'S4', 'S5', 'S6'] as const).map((sp) => (
            <button
              key={sp}
              type="button"
              onClick={() => setSprintFilter(sp)}
              className={`px-3 py-1.5 text-xs font-medium rounded-md transition-colors whitespace-nowrap ${
                sprintFilter === sp ? 'bg-white text-slate-900 shadow-sm' : 'text-slate-600 hover:text-slate-900'
              }`}
            >
              {sp === 'ALL' ? 'Toutes (22)' : sp}
            </button>
          ))}
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {visibleTables.map((table) => (
          <div key={table.name} className="bg-white border border-slate-200 rounded-xl p-5 flex flex-col justify-between">
            <div>
              <div className="flex items-center justify-between gap-2">
                <h3 className="text-sm font-mono font-semibold text-slate-900">{table.name}</h3>
                <span className="text-xs text-slate-500 tabular-nums">
                  {table.sprint} · {table.owner}
                </span>
              </div>
              <p className="text-xs text-slate-600 mt-1.5 leading-relaxed">{table.description}</p>
            </div>

            <div className="mt-4 pt-3 border-t border-slate-100 space-y-1.5">
              {table.columns.map((col) => (
                <div key={col.name} className="flex items-center justify-between text-xs font-mono">
                  <span className="text-slate-800">{col.name}</span>
                  <span className="text-slate-500">
                    {col.type}
                    {col.constraint ? ` · ${col.constraint}` : ''}
                  </span>
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
