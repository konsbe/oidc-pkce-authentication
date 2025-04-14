import React, { useContext } from 'react';
import { render, screen } from '@testing-library/react';
import { AuthProvider, AuthContext } from './AuthProvider';
import Keycloak from 'keycloak-js';

// Mock Keycloak globally
jest.mock('keycloak-js', () => {
  return jest.fn().mockImplementation(() => ({
    authenticated: true,
    token: 'fake-token',
    tokenParsed: { sub: '123', name: 'Test User' },
    isTokenExpired: jest.fn(() => false),
    login: jest.fn(),
    logout: jest.fn(),
    updateToken: jest.fn(() => Promise.resolve(true)),
    hasRealmRole: jest.fn(() => true),
    hasResourceRole: jest.fn(() => false),
    init: jest.fn(() => Promise.resolve(true)),
  }));
});

describe('AuthProvider', () => {
  test('renders loading when authData is not provided', () => {
    render(<AuthProvider>{<div>Should not see me</div>}</AuthProvider>);
    expect(screen.getByText(/loading/i)).toBeInTheDocument();
  });

  test('renders children when authData is provided', () => {
    const mockKeycloak = new Keycloak();
    render(
      <AuthProvider authData={mockKeycloak}>
        <div>Authenticated UI</div>
      </AuthProvider>
    );
    expect(screen.getByText('Authenticated UI')).toBeInTheDocument();
  });

  test('provides context value to consumer', () => {
    const mockKeycloak = new Keycloak();

    const ConsumerComponent = () => {
      const ctx = useContext(AuthContext);
      return <div>Username: {ctx.tokenParsed.name}</div>;
    };

    render(
      <AuthProvider authData={mockKeycloak}>
        <ConsumerComponent />
      </AuthProvider>
    );

    expect(screen.getByText('Username: Test User')).toBeInTheDocument();
  });
});
