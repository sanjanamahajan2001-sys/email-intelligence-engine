import React, { useState, useEffect } from 'react';

function App() {
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

      {/* Main Grid Dashboard */}
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
