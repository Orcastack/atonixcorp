import React from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { useEnterprise } from '../context/EnterpriseContext';

/**
 * ProtectedRoute — guards authenticated routes.
 * Optional `requiredRoles` prop (string[]) restricts access by role.
 * Roles come from EnterpriseContext.currentUserRole.
 * If roles haven't loaded yet (null), access is allowed (fail-open during load).
 */
const ProtectedRoute = ({ children, requiredRoles, requiredPermission }) => {
  const { isAuthenticated, loading, user } = useAuth();
  const { currentUserRole, hasPermission } = useEnterprise();

  if (loading) {
    return (
      <div style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        minHeight: '100vh',
        fontFamily: 'var(--font-family)',
        fontSize: 'var(--font-size-base)',
        color: 'var(--color-silver-dark)',
        background: 'var(--color-silver-very-light)',
      }}>
        Loading...
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  if (!user?.email_verified) {
    return <Navigate to="/verify-email" replace />;
  }

  // Role guard — only enforce when requiredRoles provided AND role is resolved
  if (
    requiredRoles &&
    requiredRoles.length > 0 &&
    currentUserRole !== null &&
    !requiredRoles.includes(currentUserRole)
  ) {
    return (
      <div style={{
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'center',
        alignItems: 'center',
        minHeight: '60vh',
        fontFamily: 'var(--font-family)',
        gap: '12px',
      }}>
        <div style={{ fontSize: 'var(--font-size-xl)', fontWeight: 600, color: 'var(--color-midnight)' }}>
          Access Restricted
        </div>
        <div style={{ fontSize: 'var(--font-size-sm)', color: 'var(--color-silver-dark)' }}>
          Your current role does not have permission to view this page.
        </div>
      </div>
    );
  }

  if (requiredPermission && !hasPermission(requiredPermission)) {
    return (
      <div style={{
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'center',
        alignItems: 'center',
        minHeight: '60vh',
        fontFamily: 'var(--font-family)',
        gap: '12px',
      }}>
        <div style={{ fontSize: 'var(--font-size-xl)', fontWeight: 600, color: 'var(--color-midnight)' }}>
          Access Restricted
        </div>
        <div style={{ fontSize: 'var(--font-size-sm)', color: 'var(--color-silver-dark)' }}>
          Your current accounting access does not permit this page.
        </div>
      </div>
    );
  }

  return children;
};

export default ProtectedRoute;
