import React from 'react';
import AtonixCorpLogo from '../branding/LedgoraLogo';

export const LogoMark = ({ size = 32, variant = 'white' }) => {
  return <AtonixCorpLogo variant={variant} size={size} withText={false} />;
};

export default LogoMark;
