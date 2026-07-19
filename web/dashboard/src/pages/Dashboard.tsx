import React, { useEffect, useState } from 'react';
import { DashboardService } from '../services/api';
import { CheckCircle, AlertTriangle, ShieldAlert } from 'lucide-react';

interface Metrics {
    total_artifacts: number;
    total_verified: number;
    total_signatures: number;
}

interface Artifact {
    id: string;
    digest: string;
    repository: string;
    first_seen: string;
    signatures_count: number;
    slsa_level: number;
    has_sbom: boolean;
}

export const Dashboard: React.FC = () => {
    const [metrics, setMetrics] = useState<Metrics | null>(null);
    const [artifacts, setArtifacts] = useState<Artifact[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const metricsData = await DashboardService.getMetrics();
                const artifactsData = await DashboardService.getArtifacts();
                setMetrics(metricsData);
                setArtifacts(artifactsData.artifacts || []);
            } catch (err) {
                console.error("Failed to load dashboard data:", err);
            } finally {
                setLoading(false);
            }
        };
        fetchData();
    }, []);

    if (loading) return <div>Loading insights...</div>;

    return (
        <div>
            <h1 style={{ marginBottom: '24px' }}>Overview</h1>
            
            <div className="metrics-grid">
                <div className="glass-card metric-card">
                    <div className="metric-header">
                        Total Artifacts
                        <DatabaseIcon />
                    </div>
                    <div className="metric-value">{metrics?.total_artifacts || 0}</div>
                </div>
                
                <div className="glass-card metric-card">
                    <div className="metric-header">
                        Verified Signatures
                        <CheckCircle size={20} color="var(--success)" />
                    </div>
                    <div className="metric-value">{metrics?.total_signatures || 0}</div>
                </div>

                <div className="glass-card metric-card">
                    <div className="metric-header">
                        Fully Compliant
                        <ShieldAlert size={20} color="var(--accent-primary)" />
                    </div>
                    <div className="metric-value">{metrics?.total_verified || 0}</div>
                </div>
            </div>

            <h2 style={{ marginTop: '48px', marginBottom: '24px' }}>Recent Ingestions</h2>
            <div className="glass-panel" style={{ overflow: 'hidden' }}>
                <table className="data-table">
                    <thead>
                        <tr>
                            <th>Repository</th>
                            <th>Digest</th>
                            <th>SLSA Level</th>
                            <th>SBOM</th>
                            <th>Signatures</th>
                            <th>Status</th>
                        </tr>
                    </thead>
                    <tbody>
                        {artifacts.map(art => (
                            <tr key={art.id}>
                                <td style={{ fontWeight: 500 }}>{art.repository || 'unknown'}</td>
                                <td style={{ fontFamily: 'monospace', color: 'var(--text-muted)' }}>
                                    {art.digest.substring(0, 15)}...
                                </td>
                                <td>
                                    {art.slsa_level > 0 
                                        ? <span className="badge badge-success">L{art.slsa_level}</span>
                                        : <span className="badge badge-warning">None</span>}
                                </td>
                                <td>
                                    {art.has_sbom 
                                        ? <CheckCircle size={16} color="var(--success)"/> 
                                        : <AlertTriangle size={16} color="var(--danger)"/>}
                                </td>
                                <td>{art.signatures_count}</td>
                                <td>
                                    {art.signatures_count > 0 && art.slsa_level >= 3 && art.has_sbom 
                                        ? <span className="badge badge-success">Admittable</span>
                                        : <span className="badge badge-danger">Non-Compliant</span>}
                                </td>
                            </tr>
                        ))}
                        {artifacts.length === 0 && (
                            <tr>
                                <td colSpan={6} style={{ textAlign: 'center', padding: '32px', color: 'var(--text-muted)' }}>
                                    No artifacts found in the system.
                                </td>
                            </tr>
                        )}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

const DatabaseIcon = () => (
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <ellipse cx="12" cy="5" rx="9" ry="3"></ellipse>
        <path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"></path>
        <path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"></path>
    </svg>
);
