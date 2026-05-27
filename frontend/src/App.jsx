import React, { useState, useEffect } from 'react';

function App() {
  const [activeTab, setActiveTab] = useState('validator'); // 'validator' or 'docs'
  const [email, setEmail] = useState('');
  const [loading, setLoading] = useState(false);
  const [currentStep, setCurrentStep] = useState(0);
  const [error, setError] = useState(null);
  const [result, setResult] = useState(null);

  // Read API base URL from Vite environment, default to localhost for local testing
  const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

  const steps = [
    { id: 0, label: 'RFC 5322 Syntax Assessment' },
    { id: 1, label: 'DNS/MX Record Resolution' },
    { id: 2, label: 'Infrastructure Fingerprinting' },
    { id: 3, label: 'SMTP Connection & Handshake' },
    { id: 4, label: 'Telemetry & Reputation Calculation' }
  ];

  // Micro-animations: step transitions while validating
  useEffect(() => {
    let interval;
    if (loading) {
      setCurrentStep(0);
      interval = setInterval(() => {
        setCurrentStep((prev) => {
          if (prev < steps.length - 1) {
            return prev + 1;
          }
          clearInterval(interval);
          return prev;
        });
      }, 700);
    } else {
      setCurrentStep(0);
    }
    return () => clearInterval(interval);
  }, [loading]);

  const handleValidate = async (e) => {
    e.preventDefault();
    if (!email.trim()) return;

    setLoading(true);
    setError(null);
    setResult(null);

    try {
      const response = await fetch(`${API_URL}/v1/public-validate`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email: email.trim() }),
      });

      if (!response.ok) {
        if (response.status === 429) {
          throw new Error('Rate limit exceeded. Please try again in a minute.');
        }
        throw new Error('Validation request failed. Ensure backend API is online.');
      }

      const data = await response.json();
      
      // Delay finishing slightly so the beautiful state animations complete
      setTimeout(() => {
        setResult(data);
        setLoading(false);
      }, 1000);

    } catch (err) {
      setError(err.message);
      setLoading(false);
    }
  };

  // Helper to get matching color class for scores
  const getScoreClass = (score) => {
    if (score > 85) return 'verified';
    if (score > 50) return 'suspicious';
    return 'invalid';
  };

  const getVerdictClass = (verdict) => {
    switch (verdict?.toLowerCase()) {
      case 'verified': return 'verified';
      case 'suspicious': return 'suspicious';
      case 'invalid': return 'invalid';
      default: return 'disposable';
    }
  };

  return (
    <div className="app-container">
      {/* Brand Header */}
      <header className="app-header">
        <div className="logo-section">
          <span className="logo-shield">🛡️</span>
          <div className="logo-text">
            <h1 className="brand-title">Email Intelligence</h1>
            <span className="brand-subtitle">Platform by Sanjana Mahajan</span>
          </div>
        </div>
        <div className="header-meta">
          <span className="badge-tag">Network Tier 1</span>
          <span className="badge-tag">Autonomous Pump</span>
        </div>
      </header>

      {/* Tab Navigation */}
      <div className="tab-navigation">
        <button 
          className={`tab-btn ${activeTab === 'validator' ? 'active' : ''}`}
          onClick={() => setActiveTab('validator')}
        >
          🔍 Live Engine Validator
        </button>
        <button 
          className={`tab-btn ${activeTab === 'docs' ? 'active' : ''}`}
          onClick={() => setActiveTab('docs')}
        >
          📊 Engine Capabilities & Docs
        </button>
      </div>

      {/* Tab Contents */}
      {activeTab === 'validator' ? (
        <main className="dashboard-grid">
          
          {/* Left Side - Input Panel */}
          <section className="glass-card">
            <h2 className="card-title">
              <span className="spin-slow">⚙️</span> Independent Core Engine
            </h2>
            <p className="card-desc">
              Verify address availability, catch-all filters, domain age integrity, and disposable provider fingerprints in real-time.
            </p>

            <form onSubmit={handleValidate}>
              <div className="input-group">
                <input
                  type="email"
                  className="premium-input"
                  placeholder="Enter email for deep lookup..."
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  disabled={loading}
                  required
                />
                <span className="input-icon">🔍</span>
              </div>

              <button type="submit" className="premium-btn" disabled={loading || !email.trim()}>
                {loading ? 'Executing Real-Time Analysis...' : 'Verify Email Security'}
              </button>
            </form>

            {error && (
              <div style={{ marginTop: '1.5rem', color: 'var(--color-invalid)', fontSize: '0.9rem', display: 'flex', gap: '0.5rem', alignItems: 'center' }}>
                <span>⚠️</span> {error}
              </div>
            )}

            {/* Live Progress Tracker */}
            {loading && (
              <div className="loading-tracker">
                <div className="tracker-title">
                  <span>Active Pipeline Logs</span>
                  <span>{Math.round(((currentStep + 1) / steps.length) * 100)}%</span>
                </div>
                <div className="tracker-steps">
                  {steps.map((step) => {
                    const isActive = currentStep === step.id;
                    const isCompleted = currentStep > step.id;
                    return (
                      <div 
                        key={step.id} 
                        className={`tracker-step ${isActive ? 'active' : ''} ${isCompleted ? 'completed' : ''}`}
                      >
                        <div className="step-indicator">
                          {isCompleted ? '✓' : step.id + 1}
                        </div>
                        <span>{step.label}</span>
                      </div>
                    );
                  })}
                </div>
              </div>
            )}
          </section>

          {/* Right Side - Detailed Results Panel */}
          <section className="glass-card">
            {!result && !loading && (
              <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', height: '100%', minHeight: '300px', color: 'var(--text-muted)' }}>
                <span style={{ fontSize: '3rem', marginBottom: '1rem', opacity: 0.5 }}>📡</span>
                <h3>No Validation Data Loaded</h3>
                <p style={{ fontSize: '0.85rem', textAlign: 'center', marginTop: '0.5rem', maxWidth: '300px' }}>
                  Type an email on the left panel to execute the multi-signal network inspection.
                </p>
              </div>
            )}

            {loading && !result && (
              <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', height: '100%', minHeight: '300px', color: 'var(--primary)' }}>
                <span className="spin-slow" style={{ fontSize: '3rem', marginBottom: '1rem' }}>🛡️</span>
                <h3>Executing Pipeline Inspection</h3>
                <p style={{ fontSize: '0.85rem', color: 'var(--text-muted)', marginTop: '0.5rem' }}>
                  Connecting directly to TCP sockets & MX domains...
                </p>
              </div>
            )}

            {result && !loading && (
              <div>
                {/* Verdict Header */}
                <div className="results-header">
                  <div className="results-verdict">
                    <span className="detail-label">Validation Outcome</span>
                    <span className={`verdict-tag ${getVerdictClass(result.authenticity_status)}`}>
                      {result.authenticity_status}
                    </span>
                  </div>
                  <div className="gauge-container">
                    <span className={`score-value ${getScoreClass(result.reputation_score)}`}>
                      {result.reputation_score}
                    </span>
                    <span className="score-label">Reputation Index</span>
                  </div>
                </div>

                {/* Status Message */}
                <div style={{ background: 'hsla(0, 0%, 100%, 0.03)', border: '1px solid var(--border-glow)', borderRadius: '10px', padding: '1rem', marginBottom: '2rem' }}>
                  <span className="detail-label">System Flag Comment</span>
                  <p style={{ fontSize: '0.95rem', fontWeight: 500, color: 'var(--text-primary)', marginTop: '0.25rem' }}>
                    {result.detailed_info.message || 'Identity successfully verified with clean logs.'}
                  </p>
                </div>

                {/* Grid Metrics */}
                <div className="details-grid">
                  <div className="detail-item">
                    <span className="detail-label">Infrastructure Provider</span>
                    <div className="detail-value">{result.detailed_info.provider || 'Independent Server'}</div>
                  </div>
                  <div className="detail-item">
                    <span className="detail-label">Risk Category</span>
                    <div className={`detail-value ${result.detailed_info.risk_level?.toLowerCase() === 'high' ? 'danger' : ''}`}>
                      {result.detailed_info.risk_level} Risk
                    </div>
                  </div>
                  <div className="detail-item">
                    <span className="detail-label">DNS Record Verification</span>
                    <div className={`detail-value ${result.detailed_info.dns_active ? 'success' : 'danger'}`}>
                      {result.detailed_info.dns_active ? 'MX Record Found' : 'No Mail Server'}
                    </div>
                  </div>
                  <div className="detail-item">
                    <span className="detail-label">SMTP Connection Response</span>
                    <div className="detail-value">{result.detailed_info.smtp_response || 'TCP Timeout / Blocked'}</div>
                  </div>
                  <div className="detail-item">
                    <span className="detail-label">Lifecycle State</span>
                    <div className={`detail-value ${result.lifecycle_state === 'ACTIVE' ? 'success' : (result.lifecycle_state === 'INVALID' || result.lifecycle_state === 'ABANDONED' ? 'danger' : '')}`}>
                      {result.lifecycle_state}
                    </div>
                  </div>
                  <div className="detail-item">
                    <span className="detail-label">Domain/Identity Age</span>
                    <div className="detail-value">
                      {result.detailed_info.domain_age_years > 0 ? `${result.detailed_info.domain_age_years.toFixed(1)} Years` : 'Zero Day'}
                    </div>
                  </div>
                </div>

                {/* Trust/Risk Factors */}
                {result.engagement?.factors && result.engagement.factors.length > 0 && (
                  <div className="factors-section">
                    <h3 className="factors-title">Validation Factors</h3>
                    <div className="factors-list">
                      {result.engagement.factors.map((factor, idx) => {
                        const isPositive = factor.startsWith('+');
                        return (
                          <div 
                            key={idx} 
                            className={`factor-pill ${isPositive ? 'positive' : 'negative'}`}
                          >
                            <span>{isPositive ? '✓' : '✗'}</span>
                            <span>{factor.substring(4)}</span>
                          </div>
                        );
                      })}
                    </div>
                  </div>
                )}
              </div>
            )}
          </section>
        </main>
      ) : (
        /* Tab 2 - Interactive Documentation & Architectural Features */
        <main className="docs-layout">
          
          {/* Section A: Global Lifecycle States */}
          <section className="glass-card">
            <h2 className="docs-section-title">🧭 Core Identity Lifecycle States</h2>
            <p className="card-desc">
              The engine structures and maps analyzed emails into six highly distinct, operational states, utilizing multi-signal metrics to flag identity freshness and security risks.
            </p>
            <div className="lifecycle-grid">
              <div className="lifecycle-card" style={{ borderLeft: '4px solid var(--color-verified)' }}>
                <div className="lifecycle-header">
                  <span className="lifecycle-name" style={{ color: 'var(--color-verified)' }}>ACTIVE</span>
                  <span className="badge-tag" style={{ border: 'none', background: 'hsla(150, 84%, 44%, 0.1)' }}>Tier 1</span>
                </div>
                <p className="lifecycle-desc">
                  Completely clean user mailbox. DNS/MX records resolve perfectly, SMTP server accepts socket delivery with zero soft-fail or reputation flags.
                </p>
              </div>

              <div className="lifecycle-card" style={{ borderLeft: '4px solid var(--primary)' }}>
                <div className="lifecycle-header">
                  <span className="lifecycle-name" style={{ color: 'var(--primary)' }}>CATCH-ALL</span>
                  <span className="badge-tag" style={{ border: 'none', background: 'hsla(190, 100%, 50%, 0.1)' }}>Tier 2</span>
                </div>
                <p className="lifecycle-desc">
                  The domain is configured to accept all incoming mail. Identified via double-probe junk simulation on active enterprise servers.
                </p>
              </div>

              <div className="lifecycle-card" style={{ borderLeft: '4px solid var(--color-suspicious)' }}>
                <div className="lifecycle-header">
                  <span className="lifecycle-name" style={{ color: 'var(--color-suspicious)' }}>STALE</span>
                  <span className="badge-tag" style={{ border: 'none', background: 'hsla(38, 92%, 50%, 0.1)' }}>Tier 2</span>
                </div>
                <p className="lifecycle-desc">
                  Verification data successfully cached, but has exceeded the standard tiered age threshold (30 days), marking it for background re-validation.
                </p>
              </div>

              <div className="lifecycle-card" style={{ borderLeft: '4px solid var(--color-invalid)' }}>
                <div className="lifecycle-header">
                  <span className="lifecycle-name" style={{ color: 'var(--color-invalid)' }}>ABANDONED</span>
                  <span className="badge-tag" style={{ border: 'none', background: 'hsla(350, 78%, 52%, 0.1)' }}>Tier 3</span>
                </div>
                <p className="lifecycle-desc">
                  Address was previously active and verified in SQLite telemetry, but now returns persistent hard failures (SMTP 550 Mailbox Unreachable).
                </p>
              </div>

              <div className="lifecycle-card" style={{ borderLeft: '4px solid var(--color-disposable)' }}>
                <div className="lifecycle-header">
                  <span className="lifecycle-name" style={{ color: 'var(--color-disposable)' }}>FULL</span>
                  <span className="badge-tag" style={{ border: 'none', background: 'hsla(265, 83%, 58%, 0.1)' }}>Tier 1</span>
                </div>
                <p className="lifecycle-desc">
                  Mailbox is currently unreachable because the recipient quota has been exceeded. Explicitly captured via SMTP status code `5.2.2`.
                </p>
              </div>

              <div className="lifecycle-card" style={{ borderLeft: '1px solid var(--border-glow)' }}>
                <div className="lifecycle-header">
                  <span className="lifecycle-name" style={{ color: 'var(--text-muted)' }}>INVALID</span>
                  <span className="badge-tag" style={{ border: 'none', background: 'hsla(0, 0%, 100%, 0.05)' }}>Tier 4</span>
                </div>
                <p className="lifecycle-desc">
                  Fatal syntax validation failure, domain lacks valid MX records, or the host returns a permanent hard-bounce sequence during initial dial.
                </p>
              </div>
            </div>
          </section>

          {/* Section B: Engine Processing Pipeline */}
          <section className="glass-card">
            <h2 className="docs-section-title">📡 Validation Processing Pipeline</h2>
            <p className="card-desc">
              Every email checked by our web service passes through five highly optimized backend checkpoints.
            </p>
            <div className="tracker-steps-docs">
              <div className="tracker-step completed">
                <div className="step-indicator">1</div>
                <div>
                  <h4 style={{ color: 'var(--text-primary)' }}>RFC Syntax Check</h4>
                  <p style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>O(1) strict check for structures, lengths, and double-dots.</p>
                </div>
              </div>
              <div className="tracker-step completed">
                <div className="step-indicator">2</div>
                <div>
                  <h4 style={{ color: 'var(--text-primary)' }}>DNS/MX Lookup</h4>
                  <p style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>Real-time concurrent domain resolutions for active mail exchanges.</p>
                </div>
              </div>
              <div className="tracker-step completed">
                <div className="step-indicator">3</div>
                <div>
                  <h4 style={{ color: 'var(--text-primary)' }}>SMTP Dial & Handshake</h4>
                  <p style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>Raw socket handshake with greylist filters and timeout safeguards.</p>
                </div>
              </div>
              <div className="tracker-step completed">
                <div className="step-indicator">4</div>
                <div>
                  <h4 style={{ color: 'var(--text-primary)' }}>OSINT & Reputation</h4>
                  <p style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>RDAP age lookups, subdomain reputation cascading, and fraud alerts.</p>
                </div>
              </div>
            </div>
          </section>

          {/* Section C: Advanced Intelligence Systems */}
          <section className="glass-card features-table-card">
            <h2 className="docs-section-title">🧠 Advanced Platform Core Capabilities</h2>
            <table className="features-table">
              <thead>
                <tr>
                  <th>Capability</th>
                  <th>Methodology</th>
                  <th>Impact</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td><strong>MX Hub Fingerprinting</strong></td>
                  <td>Matches domain resolution against known underlying disposable MX mail signature tables.</td>
                  <td>Blocks rotating disposable domains instantly by recognizing their core mail servers.</td>
                </tr>
                <tr>
                  <td><strong>Zero-Day Heuristics</strong></td>
                  <td>Scans domain entropy, new TLD structures (.xyz, .tk), and name patterns for disposable keywords.</td>
                  <td>Catches throwaway domains immediately before they register active MX signatures.</td>
                </tr>
                <tr>
                  <td><strong>Autonomous Learning Pump</strong></td>
                  <td>Asynchronous worker scans validation history, evaluates suspects, and promotes newly found patterns.</td>
                  <td>Self-improving infrastructure without requiring continuous manual cataloging or rebuilds.</td>
                </tr>
                <tr>
                  <td><strong>Identity Age Conflict</strong></td>
                  <td>Compares telemetry age against RDAP domain creation date in real-time.</td>
                  <td>Triggers immediate critical alerts for potential domain hijack or phishing impersonations.</td>
                </tr>
                <tr>
                  <td><strong>Port 25 Resiliency</strong></td>
                  <td>Engine is built with timeout fallbacks, custom network block flags, and dynamic greylist filters.</td>
                  <td>Ensures server responsiveness and prevents thread locks under restricted cloud host firewalls.</td>
                </tr>
              </tbody>
            </table>
          </section>
        </main>
      )}

      {/* Footer Branding */}
      <footer className="app-footer">
        <span className="footer-branding">
          Engineered by <strong>Sanjana Mahajan</strong> &copy; {new Date().getFullYear()}
        </span>
        <div className="footer-links">
          <a href="https://github.com/sanjana" target="_blank" rel="noopener noreferrer" className="footer-link">GitHub</a>
          <a href="#" className="footer-link">Portfolio Website</a>
          <a href="#" className="footer-link">Systems Status</a>
        </div>
      </footer>
    </div>
  );
}

export default App;
