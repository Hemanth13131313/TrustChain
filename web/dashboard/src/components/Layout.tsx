import React from 'react';
import { Outlet, NavLink } from 'react-router-dom';
import { ShieldCheck, LayoutDashboard, Database, Settings } from 'lucide-react';

export const Layout: React.FC = () => {
  return (
    <div className="app-container">
      <aside className="sidebar">
        <div className="sidebar-logo">
          <ShieldCheck size={28} color="#60a5fa" />
          TRUSTCHAIN
        </div>
        
        <nav className="nav-menu">
          <NavLink to="/" className={({ isActive }) => `nav-item ${isActive ? 'active' : ''}`}>
            <LayoutDashboard size={20} />
            Dashboard
          </NavLink>
          <NavLink to="/artifacts" className={({ isActive }) => `nav-item ${isActive ? 'active' : ''}`}>
            <Database size={20} />
            Artifact Explorer
          </NavLink>
          <NavLink to="/settings" className={({ isActive }) => `nav-item ${isActive ? 'active' : ''}`}>
            <Settings size={20} />
            Settings
          </NavLink>
        </nav>
      </aside>

      <main className="main-content">
        <header className="topbar">
          <div style={{ color: 'var(--text-muted)' }}>Security Enforcement Control Plane</div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
            <div className="badge badge-success">System Healthy</div>
          </div>
        </header>

        <div className="page-content">
          <Outlet />
        </div>
      </main>
    </div>
  );
};
