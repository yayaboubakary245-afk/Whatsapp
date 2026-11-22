import React, { useEffect, useState } from 'react';
import { CheckCircle2, XCircle, RefreshCw } from 'lucide-react';

interface ComplianceCheck {
  id: string;
  title: string;
  criterion: string;
  passed: boolean;
  evidence: string;
}

export const ComplianceNetworkLab: React.FC = () => {
  const [checks, setChecks] = useState<ComplianceCheck[]>([]);
  const [loading, setLoading] = useState(true);

  const [msgType, setMsgType] = useState('SEND_MESSAGE');
  const [payloadText, setPayloadText] = useState('{"conversation_id":"conv-01","content":"Bonjour SUPPTIC"}');
  const [frameResult, setFrameResult] = useState<{
    byteLength: number;
    totalFrameBytes: number;
    hexHeader: string;
    hexPreview: string;
  } | null>(null);

  const [udpCmd, setUdpCmd] = useState('PING');
  const [udpResult, setUdpResult] = useState<{
    sent: string;
    received: string;
    rttMs: number;
    port: number;
  } | null>(null);

  const runCompliance = () => {
    setLoading(true);
    fetch('/api/project/compliance')
      .then((r) => r.json())
      .then((data) => {
        setChecks(data.checks || []);
        setLoading(false);
      })
      .catch(() => setLoading(false));
  };

  useEffect(() => {
    runCompliance();
  }, []);

  const encodeFrame = async (e: React.FormEvent) => {
    e.preventDefault();
    let parsed: Record<string, unknown> = {};
    try {
      parsed = JSON.parse(payloadText);
    } catch {
      parsed = { raw: payloadText };
    }
    const res = await fetch('/api/protocol/frame', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ type: msgType, user_id: 'u-yaya-01', payload: parsed }),
    });
    const data = await res.json();
    setFrameResult(data);
  };

  const sendUdp = async (cmd: string) => {
    setUdpCmd(cmd);
    const res = await fetch('/api/protocol/udp-ping', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ command: cmd }),
    });
    const data = await res.json();
    setUdpResult(data);
  };

  return (
    <div className="space-y-8">
      {/* Section 1 : Audit C-T-07 à C-T-12 */}
      <div className="bg-white border border-slate-200 rounded-xl p-6">
        <div className="flex flex-wrap items-center justify-between gap-4">
          <div>
            <h2 className="text-base font-semibold text-slate-900">
              Vérification des critères d’acceptation (C-T-07 à C-T-12)
            </h2>
            <p className="text-sm text-slate-600 mt-1">
              Analyse statique réelle exécutée sur les répertoires server-go/, client-android/ et database/.
            </p>
          </div>
          <button
            type="button"
            onClick={runCompliance}
            className="px-4 py-2 text-xs font-medium text-slate-700 border border-slate-200 rounded-lg hover:bg-slate-50 flex items-center gap-2 whitespace-nowrap"
          >
            <RefreshCw className={`w-3.5 h-3.5 ${loading ? 'animate-spin' : ''}`} />
            Relancer l’audit
          </button>
        </div>

        <div className="mt-6 overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="border-b border-slate-200 text-xs font-semibold text-slate-600">
                <th className="py-2.5 pr-4">Critère</th>
                <th className="py-2.5 pr-4">Exigence technique</th>
                <th className="py-2.5 pr-4">Preuve vérifiée</th>
                <th className="py-2.5 text-right">Statut</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100 text-xs">
              {checks.map((c) => (
                <tr key={c.id} className="hover:bg-slate-50">
                  <td className="py-3 pr-4 font-mono font-semibold text-slate-900 whitespace-nowrap">{c.id}</td>
                  <td className="py-3 pr-4">
                    <div className="font-medium text-slate-900">{c.title}</div>
                    <div className="text-slate-500 font-mono text-[11px] mt-0.5">{c.criterion}</div>
                  </td>
                  <td className="py-3 pr-4 font-mono text-slate-600">{c.evidence}</td>
                  <td className="py-3 text-right whitespace-nowrap">
                    {c.passed ? (
                      <span className="inline-flex items-center gap-1.5 text-[#16A34A] font-medium">
                        <CheckCircle2 className="w-4 h-4" /> Conforme
                      </span>
                    ) : (
                      <span className="inline-flex items-center gap-1.5 text-[#DC2626] font-medium">
                        <XCircle className="w-4 h-4" /> Échec
                      </span>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Section 2 : Banc de test TCP Length-Prefixed & UDP PING/PONG */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white border border-slate-200 rounded-xl p-6">
          <h3 className="text-sm font-semibold text-slate-900">
            Protocole TCP Port 5000 · Framing Length-Prefixed
          </h3>
          <p className="text-xs text-slate-600 mt-1">
            Format binaire Go : [4 octets big-endian : longueur N][N octets : payload JSON]
          </p>

          <form onSubmit={encodeFrame} className="mt-4 space-y-3">
            <div>
              <label className="block text-xs font-medium text-slate-700 mb-1">Type de paquet (protocol/types.go)</label>
              <select
                value={msgType}
                onChange={(e) => setMsgType(e.target.value)}
                className="w-full px-3 py-2 text-xs border border-slate-200 rounded-lg bg-white font-mono"
              >
                <option value="SEND_MESSAGE">SEND_MESSAGE</option>
                <option value="REGISTER">REGISTER</option>
                <option value="LOGIN">LOGIN</option>
                <option value="CREATE_GROUP">CREATE_GROUP</option>
                <option value="CALL_SIGNAL">CALL_SIGNAL</option>
                <option value="GET_ANNOUNCEMENTS">GET_ANNOUNCEMENTS</option>
              </select>
            </div>
            <div>
              <label className="block text-xs font-medium text-slate-700 mb-1">Payload JSON</label>
              <textarea
                rows={3}
                value={payloadText}
                onChange={(e) => setPayloadText(e.target.value)}
                className="w-full px-3 py-2 text-xs font-mono border border-slate-200 rounded-lg"
              />
            </div>
            <button
              type="submit"
              className="px-4 py-2 bg-[#0B6E4F] text-white text-xs font-medium rounded-lg hover:bg-[#08563D]"
            >
              Encoder la trame TCP
            </button>
          </form>

          {frameResult && (
            <div className="mt-4 p-4 bg-slate-50 border border-slate-200 rounded-lg space-y-2 text-xs font-mono tabular-nums">
              <div className="flex justify-between">
                <span className="text-slate-500">En-tête 4 octets big-endian :</span>
                <span className="font-semibold text-slate-900">{frameResult.hexHeader}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-slate-500">Longueur JSON (N) :</span>
                <span className="font-semibold text-slate-900">{frameResult.byteLength} octets</span>
              </div>
              <div className="flex justify-between">
                <span className="text-slate-500">Taille totale trame :</span>
                <span className="font-semibold text-slate-900">{frameResult.totalFrameBytes} octets</span>
              </div>
              <div className="pt-2 border-t border-slate-200 text-[11px] text-slate-600 break-all">
                Aperçu hex : {frameResult.hexPreview}
              </div>
            </div>
          )}
        </div>

        <div className="bg-white border border-slate-200 rounded-xl p-6">
          <h3 className="text-sm font-semibold text-slate-900">
            Supervision UDP Port 5001 · PING / PONG & Présence
          </h3>
          <p className="text-xs text-slate-600 mt-1">
            Implémenté en Go dans internal/network/udp/server.go via net.ListenPacket et ReadFromUDP.
          </p>

          <div className="mt-4 flex flex-wrap gap-2">
            {['PING', 'PRESENCE:yaya.boubakary', 'PRESENCE:christian.ndue'].map((cmd) => (
              <button
                key={cmd}
                type="button"
                onClick={() => sendUdp(cmd)}
                className="px-3 py-2 text-xs font-mono border border-slate-200 rounded-lg hover:border-[#0B6E4F] hover:text-[#0B6E4F]"
              >
                Envoyer {cmd}
              </button>
            ))}
          </div>

          {udpResult ? (
            <div className="mt-4 p-4 bg-slate-50 border border-slate-200 rounded-lg space-y-2 text-xs font-mono tabular-nums">
              <div className="flex justify-between">
                <span className="text-slate-500">Datagramme envoyé :</span>
                <span className="text-slate-900 font-semibold">{udpResult.sent}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-slate-500">Réponse WriteToUDP :</span>
                <span className="text-[#16A34A] font-semibold">{udpResult.received}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-slate-500">Port / Latence estimée :</span>
                <span className="text-slate-900">:{udpResult.port} · {udpResult.rttMs} ms</span>
              </div>
            </div>
          ) : (
            <div className="mt-4 p-4 bg-slate-50 border border-slate-200 rounded-lg text-xs text-slate-500">
              Cliquez sur une commande ci-dessus pour tester l’échange UDP Go.
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
