import { useState, useEffect } from 'react';
import { InitUser, GetIdentity, IsLocked, VerifyPassword, SetPassword, GetSettings, SaveSettings } from '../wailsjs/go/main/App';
import { main } from '../wailsjs/go/models';
import { AetherClient } from 'aether-sdk';
import './index.css';

interface Capsule {
  capsule_id: string;
  created_at: string;
  publisher: string;
  visibility: string;
  owner_key_id: string;
  signature: string;
  version: number;
}

const client = new AetherClient('http://localhost:5000');

function App() {
  const [authMode, setAuthMode] = useState<'loading' | 'setup' | 'login' | 'unlocked'>('loading');
  const [password, setPassword] = useState('');
  const [authError, setAuthError] = useState('');

  const [activeTab, setActiveTab] = useState('identity');
  const [username, setUsername] = useState('');
  const [identity, setIdentity] = useState<string | null>(null);
  const [error, setError] = useState('');
  const [publishPayload, setPublishPayload] = useState('');
  const [publishStatus, setPublishStatus] = useState('');
  
  const [fetchUsername, setFetchUsername] = useState('');
  const [capsules, setCapsules] = useState<Capsule[]>([]);
  const [traceError, setTraceError] = useState('');
  const [decodedPayloads, setDecodedPayloads] = useState<Record<string, string>>({});
  const [decodingIds, setDecodingIds] = useState<Set<string>>(new Set());
  const [swarmHealth, setSwarmHealth] = useState<Record<string, number>>({});

  const [reputation, setReputation] = useState<Record<string, number>>({});

  // Settings
  const [settings, setSettings] = useState<main.GatewaySettings | null>(null);

  useEffect(() => {
    checkAuth();
  }, []);

  const checkAuth = async () => {
    try {
      const locked = await IsLocked();
      if (locked) {
        setAuthMode('login');
      } else {
        setAuthMode('setup');
      }
    } catch (err) {
      console.error(err);
    }
  };

  const loadSettings = async () => {
    const s = await GetSettings();
    setSettings(s);
  };

  const handleAuthSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setAuthError('');
    if (!password) {
      setAuthError('ERROR: PASSWORD CANNOT BE EMPTY');
      return;
    }

    try {
      if (authMode === 'setup') {
        await SetPassword(password);
        setAuthMode('unlocked');
        loadSettings();
      } else if (authMode === 'login') {
        const ok = await VerifyPassword(password);
        if (ok) {
          setAuthMode('unlocked');
          loadSettings();
        } else {
          setAuthError('ACCESS DENIED: INVALID CREDENTIALS');
        }
      }
    } catch (err: any) {
      setAuthError(`ERROR: ${err.message}`);
    }
  };

  const handleCreateIdentity = async () => {
    try {
      setError('');
      const keyId = await InitUser(username);
      setIdentity(keyId);
    } catch (err: any) {
      setError(`[ERR] ${err.message || String(err)}`);
    }
  };

  const handleGetIdentity = async () => {
    try {
      setError('');
      const keyId = await GetIdentity(username);
      setIdentity(keyId);
    } catch (err: any) {
      setError(`[ERR] ${err.message || String(err)}`);
    }
  };

  const handleClaimUsername = async () => {
    try {
      if (!identity) throw new Error("LOAD KEY FIRST");
      setError('> CLAIMING ALIAS IN NETWORK...');
      const res = await client.claimUsername(username, identity);
      if (res.ok) setError(`> ALIAS [${username}] SUCCESSFULLY CLAIMED`);
    } catch (err: any) {
      setError(`[ERR] ${err.message}`);
    }
  };

  const handlePublish = async () => {
    setPublishStatus('> INITIALIZING DATA TRANSFER...');
    try {
      const res = await client.publish(username, publishPayload, { persistent: true });
      setPublishStatus(`> TRANSFER SUCCESSFUL\n> CAPSULE_ID: ${res.capsule_id}`);
    } catch (err: any) {
      setPublishStatus(`> [ERR] TRANSFER FAILED: ${err.message}`);
    }
  };

  const handleTrace = async () => {
    setCapsules([]);
    setTraceError('');
    setDecodedPayloads({});
    setSwarmHealth({});
    try {
      const res = await client.resolve(fetchUsername);
      if (!res.ok) throw new Error(res.error || 'Failed to resolve');
      setCapsules(res.capsules || []);
      
      // Check availability for all capsules
      res.capsules?.forEach(async (c: any) => {
        try {
          const avail = await client.checkAvailability(c.capsule_id);
          setSwarmHealth(prev => ({ ...prev, [c.capsule_id]: avail.seeds }));
        } catch (e) {
          console.error(e);
        }
      });
    } catch (err: any) {
      setTraceError(`[ERR] TRACE FAILED: ${err.message}`);
    }
  };

  const handleDecode = async (capsuleId: string) => {
    setDecodingIds(prev => new Set(prev).add(capsuleId));
    try {
      const res = await client.fetch(capsuleId);
      if (!res.ok) throw new Error(res.error || 'Failed to decode');
      setDecodedPayloads(prev => ({ ...prev, [capsuleId]: res.data || '' }));
    } catch (err: any) {
      setDecodedPayloads(prev => ({ ...prev, [capsuleId]: `[ERR] ${err.message}` }));
    } finally {
      setDecodingIds(prev => {
        const next = new Set(prev);
        next.delete(capsuleId);
        return next;
      });
    }
  };

  const fetchReputation = async () => {
    try {
      const res = await client.getReputation();
      setReputation(res.reputation || {});
    } catch (err) {
      console.error(err);
    }
  };

  const saveSettings = async (updates: Partial<main.GatewaySettings>) => {
    if (!settings) return;
    const next = { ...settings, ...updates } as main.GatewaySettings;
    setSettings(next);
    await SaveSettings(next);
  };

  const formatTime = (iso: string) => {
    try {
      return new Date(iso).toLocaleString();
    } catch {
      return iso;
    }
  };

  if (authMode === 'loading') {
    return <div className="lock-screen">LOADING...</div>;
  }

  if (authMode === 'setup' || authMode === 'login') {
    return (
      <div className="lock-screen">
        <div className="lock-container">
          <h1 style={{ color: 'var(--ctos-green)', marginBottom: '10px' }}>// AETHER GATEWAY</h1>
          <p style={{ color: 'var(--ctos-text-dim)', marginBottom: '30px' }}>
            {authMode === 'setup' ? 'INITIALIZE NEW MASTER PASSWORD TO ENCRYPT LOCAL KEYS' : 'ENTER MASTER PASSWORD TO DECRYPT SYSTEM'}
          </p>
          <form onSubmit={handleAuthSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '15px' }}>
            <input 
              type="password" 
              placeholder={authMode === 'setup' ? "NEW_PASSWORD" : "MASTER_PASSWORD"} 
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              autoFocus
              style={{ textAlign: 'center', letterSpacing: '3px' }}
            />
            <button type="submit" style={{ width: '100%' }}>
              {authMode === 'setup' ? 'INITIALIZE_SYSTEM' : 'UNLOCK_SYSTEM'}
            </button>
          </form>
          {authError && <div style={{ color: 'var(--ctos-red)', marginTop: '20px', fontSize: '14px' }}>{authError}</div>}
        </div>
      </div>
    );
  }

  return (
    <div className={`app-container theme-${settings?.theme || 'ctos'}`}>
      <div className="sidebar">
        <div className="logo">
          <span>//</span> AETHER
        </div>
        
        <div style={{ padding: '0 15px', color: '#555', fontSize: '11px', marginBottom: '8px' }}>
          SYSTEM_ACCESS_GRANTED
        </div>

        <div 
          className={`nav-item ${activeTab === 'identity' ? 'active' : ''}`}
          onClick={() => setActiveTab('identity')}
        >
          [01] Identity.Mgr
        </div>
        <div 
          className={`nav-item ${activeTab === 'publish' ? 'active' : ''}`}
          onClick={() => setActiveTab('publish')}
        >
          [02] Data.Uplink
        </div>
        <div 
          className={`nav-item ${activeTab === 'explore' ? 'active' : ''}`}
          onClick={() => setActiveTab('explore')}
        >
          [03] Net.Trace
        </div>
        <div 
          className={`nav-item ${activeTab === 'network' ? 'active' : ''}`}
          onClick={() => { setActiveTab('network'); fetchReputation(); }}
        >
          [04] Net.Swarm
        </div>
        <div 
          className={`nav-item ${activeTab === 'control' ? 'active' : ''}`}
          onClick={() => setActiveTab('control')}
        >
          [05] Sys.Control
        </div>

        <div className="sys-info">
          OS: AETHER_OS v1.0.5<br/>
          NET: {settings?.mode === 'offline' ? 'OFFLINE' : 'DHT_KADEMLIA_ACTIVE'}<br/>
          SEC: ED25519_AES_GCM<br/>
          STS: UNLOCKED
        </div>
      </div>

      <div className="main-content">
        {activeTab === 'identity' && (
          <div className="ctos-panel">
            <h2>Identity_Manager.exe</h2>
            <p style={{ color: 'var(--ctos-text-dim)', marginBottom: '20px' }}>
              // GENERATE NEW DECENTRALIZED IDENTIFIER OR LOAD EXISTING KEYRING
            </p>
            
            <div style={{ display: 'flex', gap: '10px', marginBottom: '15px' }}>
              <input 
                type="text" 
                placeholder="TARGET_ALIAS (e.g. alice)" 
                value={username}
                onChange={(e) => setUsername(e.target.value)}
              />
              <button onClick={handleCreateIdentity}>Generate</button>
              <button onClick={handleGetIdentity} style={{ borderColor: '#555', color: '#888' }}>Load_Key</button>
            </div>
            
            <div style={{ display: 'flex', gap: '10px', marginBottom: '15px' }}>
              <button onClick={handleClaimUsername} style={{ borderColor: 'var(--ctos-cyan)', color: 'var(--ctos-cyan)' }}>Broadcast_Claim</button>
            </div>

            {error && <div style={{ color: 'var(--ctos-red)', margin: '15px 0', fontFamily: 'monospace' }}>{error}</div>}
            
            {identity && (
              <div className="data-output">
                <div style={{ color: '#555', marginBottom: '8px' }}>// KEY GENERATION SUCCESSFUL</div>
                <div style={{ color: '#fff' }}>ALIAS : {username}</div>
                <div>DID   : {identity}</div>
              </div>
            )}
          </div>
        )}

        {activeTab === 'publish' && (
          <div className="ctos-panel">
            <h2>Data_Uplink.exe</h2>
            <p style={{ color: 'var(--ctos-text-dim)', marginBottom: '20px' }}>
              // INJECT RAW JSON PAYLOAD INTO GLOBAL DHT NETWORK
            </p>

            <textarea 
              rows={8}
              placeholder="{ 'message': 'HELLO_WORLD' }"
              value={publishPayload}
              onChange={(e) => setPublishPayload(e.target.value)}
              style={{ marginBottom: '15px' }}
            />
            
            <button onClick={handlePublish}>Execute_Uplink</button>

            {publishStatus && (
              <div className="data-output">
                {publishStatus}
              </div>
            )}
          </div>
        )}

        {activeTab === 'explore' && (
          <div className="ctos-panel">
            <h2>Net_Trace.exe</h2>
            <p style={{ color: 'var(--ctos-text-dim)', marginBottom: '20px' }}>
              // QUERY GLOBAL DHT FOR TARGET SIGNATURES
            </p>

            <div style={{ display: 'flex', gap: '10px', marginBottom: '20px' }}>
              <input 
                type="text" 
                placeholder="TARGET_ALIAS" 
                value={fetchUsername}
                onChange={(e) => setFetchUsername(e.target.value)}
              />
              <button onClick={handleTrace}>Trace_Target</button>
            </div>

            {traceError && (
              <div className="data-output" style={{ color: 'var(--ctos-red)' }}>
                {traceError}
              </div>
            )}

            {capsules.length > 0 && (
              <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
                <div style={{ color: 'var(--ctos-cyan)', fontSize: '12px' }}>
                  // FOUND {capsules.length} CAPSULE{capsules.length > 1 ? 'S' : ''}
                </div>
                {capsules.map((cap) => (
                  <div key={cap.capsule_id} className="capsule-card">
                    <div className="capsule-header">
                      <div className="capsule-row">
                        <span className="capsule-label">CAPSULE</span>
                        <span className="capsule-value">{cap.capsule_id}</span>
                      </div>
                      <div className="capsule-row">
                        <span className="capsule-label">PUBLISHER</span>
                        <span className="capsule-value">{cap.publisher}</span>
                      </div>
                      <div className="capsule-row">
                        <span className="capsule-label">CREATED</span>
                        <span className="capsule-value">{formatTime(cap.created_at)}</span>
                      </div>
                      <div className="capsule-row">
                        <span className="capsule-label">SWARM</span>
                        <span className="capsule-value" style={{ color: swarmHealth[cap.capsule_id] > 0 ? 'var(--ctos-green)' : 'var(--ctos-red)' }}>
                          {swarmHealth[cap.capsule_id] !== undefined ? `${swarmHealth[cap.capsule_id]} SEEDS` : 'SCANNING...'}
                        </span>
                      </div>
                    </div>

                    {!decodedPayloads[cap.capsule_id] && (
                      <button 
                        className="decode-btn"
                        onClick={() => handleDecode(cap.capsule_id)}
                        disabled={decodingIds.has(cap.capsule_id)}
                      >
                        {decodingIds.has(cap.capsule_id) ? 'Decoding...' : 'Decode_Payload'}
                      </button>
                    )}

                    {decodedPayloads[cap.capsule_id] && (
                      <div className="decoded-payload">
                        <div style={{ color: 'var(--ctos-cyan)', fontSize: '11px', marginBottom: '6px' }}>// DECODED PAYLOAD</div>
                        <div>{decodedPayloads[cap.capsule_id]}</div>
                      </div>
                    )}
                  </div>
                ))}
              </div>
            )}
          </div>
        )}

        {activeTab === 'network' && (
          <div className="ctos-panel">
            <h2>Net_Swarm.exe</h2>
            <p style={{ color: 'var(--ctos-text-dim)', marginBottom: '20px' }}>
              // GLOBAL SEEDING REPUTATION SCOREBOARD
            </p>
            <button onClick={fetchReputation} style={{ marginBottom: '20px' }}>Refresh_Scores</button>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
              {Object.keys(reputation).length === 0 ? (
                <div style={{ color: 'var(--ctos-text-dim)' }}>// NO REPUTATION DATA FOUND</div>
              ) : (
                Object.entries(reputation)
                  .sort(([,a], [,b]) => b - a)
                  .map(([peerId, score], i) => (
                    <div key={peerId} className="capsule-card" style={{ flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center' }}>
                      <div style={{ display: 'flex', alignItems: 'center', gap: '15px' }}>
                        <span style={{ color: i === 0 ? 'var(--ctos-cyan)' : 'var(--ctos-text-dim)' }}>[{i + 1}]</span>
                        <span style={{ fontFamily: 'monospace' }}>{peerId}</span>
                      </div>
                      <div style={{ color: 'var(--ctos-green)' }}>{score.toLocaleString()} BYTES</div>
                    </div>
                  ))
              )}
            </div>
          </div>
        )}

        {activeTab === 'control' && settings && (
          <div className="ctos-panel">
            <h2>Sys_Control.exe</h2>
            <p style={{ color: 'var(--ctos-text-dim)', marginBottom: '30px' }}>
              // CONFIGURE LOCAL NODE PARAMETERS
            </p>

            <div className="settings-group">
              <label>HARD DISK ALLOCATION (GB)</label>
              <div style={{ display: 'flex', alignItems: 'center', gap: '15px' }}>
                <input 
                  type="range" 
                  min="1" max="100" 
                  value={settings.disk_limit_gb}
                  onChange={(e) => saveSettings({ disk_limit_gb: parseInt(e.target.value) })}
                  style={{ flexGrow: 1 }}
                />
                <span style={{ width: '50px', textAlign: 'right', fontFamily: 'monospace' }}>
                  {settings.disk_limit_gb} GB
                </span>
              </div>
            </div>

            <div className="settings-group">
              <label>NETWORK MODE</label>
              <select 
                value={settings.mode}
                onChange={(e) => saveSettings({ mode: e.target.value })}
              >
                <option value="online">ONLINE (P2P ACTIVE)</option>
                <option value="relay">RELAY ONLY</option>
                <option value="offline">OFFLINE (LOCAL CACHE)</option>
              </select>
            </div>

            <div className="settings-group">
              <label>SYSTEM THEME</label>
              <select 
                value={settings.theme}
                onChange={(e) => saveSettings({ theme: e.target.value })}
              >
                <option value="ctos">HACKER (PURE BLACK / GREEN)</option>
                <option value="dark">DARK (GREY / BLUE)</option>
                <option value="light">LIGHT (WHITE / BLACK)</option>
              </select>
            </div>
            
            <div className="settings-group" style={{ marginTop: '40px', borderTop: '1px solid #222', paddingTop: '20px' }}>
              <label style={{ color: 'var(--ctos-red)' }}>SECURITY</label>
              <button 
                style={{ borderColor: 'var(--ctos-red)', color: 'var(--ctos-red)' }}
                onClick={() => setAuthMode('login')}
              >
                LOCK_SESSION
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

export default App;
