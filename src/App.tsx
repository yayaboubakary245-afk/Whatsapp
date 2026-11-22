/**
 * @license
 * SPDX-License-Identifier: Apache-2.0
 */
import React, { useState } from 'react';
import { Download } from 'lucide-react';
import { AndroidSimulator } from './components/AndroidSimulator';
import { CodeExplorer } from './components/CodeExplorer';
import { SchemaExplorer } from './components/SchemaExplorer';
import { ComplianceNetworkLab } from './components/ComplianceNetworkLab';

type TabId = 'simulator' | 'code' | 'schema' | 'compliance';

export default function App() {
  const [activeTab, setActiveTab] = useState<TabId>('simulator');

  return (
    <div className="min-h-screen bg-[#F8FAFC] text-[#0F172A] flex flex-col">
      {/* Top Bar Contract : 1 ligne, 3 zones */}
      <header className="bg-white border-b border-slate-200 sticky top-0 z-30">
        <div className="max-w-[1280px] mx-auto px-6 h-16 flex items-center justify-between gap-8">
          {/* Zone 1 : Wordmark 1 ligne */}
          <span className="text-lg font-display font-bold tracking-tight text-slate-900 whitespace-nowrap shrink-0">
            WhatsApp Campus
          </span>

          {/* Zone 2 : 4 liens de navigation single-line */}
          <nav className="hidden md:flex items-center gap-6 text-sm font-medium">
            {(
              [
                ['simulator', 'Simulateur Android'],
                ['code', 'Explorateur Code'],
                ['schema', 'Schéma PostgreSQL'],
                ['compliance', 'Conformité & Réseau'],
              ] as [TabId, string][]
            ).map(([id, label]) => (
              <button
                key={id}
                type="button"
                onClick={() => setActiveTab(id)}
                className={`py-5 border-b-2 transition-colors whitespace-nowrap shrink-0 ${
                  activeTab === id
                    ? 'border-[#0B6E4F] text-[#0B6E4F]'
                    : 'border-transparent text-slate-600 hover:text-slate-900'
                }`}
              >
                {label}
              </button>
            ))}
          </nav>

          {/* Zone 3 : 1 action principale */}
          <div className="shrink-0">
            <a
              href="/api/project/download-zip"
              download="whatsapp-campus-groupe12.zip"
              className="inline-flex items-center gap-2 px-4 py-2 text-xs font-medium text-white bg-[#0B6E4F] hover:bg-[#08563D] rounded-lg transition-colors whitespace-nowrap shrink-0"
            >
              <Download className="w-3.5 h-3.5" />
              Télécharger ZIP
            </a>
          </div>
        </div>

        {/* Navigation mobile */}
        <div className="md:hidden border-t border-slate-100 px-4 py-2 flex items-center gap-2 overflow-x-auto">
          {(
            [
              ['simulator', 'Simulateur'],
              ['code', 'Code'],
              ['schema', 'PostgreSQL'],
              ['compliance', 'Conformité'],
            ] as [TabId, string][]
          ).map(([id, label]) => (
            <button
              key={id}
              type="button"
              onClick={() => setActiveTab(id)}
              className={`px-3 py-1.5 text-xs font-medium rounded-md whitespace-nowrap shrink-0 ${
                activeTab === id ? 'bg-[#0B6E4F] text-white' : 'text-slate-600 bg-slate-100'
              }`}
            >
              {label}
            </button>
          ))}
        </div>
      </header>

      {/* En-tête contextuel épuré */}
      <main className="flex-1 max-w-[1280px] w-full mx-auto px-6 py-8">
        <div className="mb-8 flex flex-col md:flex-row md:items-end justify-between gap-4 border-b border-slate-200 pb-6">
          <div>
            <div className="text-xs text-slate-500">
              SUP’PTIC · Groupe 12 · Yaya, Ndue, Mengue
            </div>
            <h1 className="text-2xl sm:text-3xl font-display font-bold text-slate-900 mt-1 text-balance">
              Architecture Android Studio, Serveur Golang & PostgreSQL
            </h1>
          </div>
          <div className="text-xs text-slate-600 tabular-nums">
            REST :8080 · TCP :5000 · UDP :5001 · 22 Tables SQL · 16+ Écrans Java
          </div>
        </div>

        {activeTab === 'simulator' && <AndroidSimulator />}
        {activeTab === 'code' && <CodeExplorer />}
        {activeTab === 'schema' && <SchemaExplorer />}
        {activeTab === 'compliance' && <ComplianceNetworkLab />}
      </main>
    </div>
  );
}
