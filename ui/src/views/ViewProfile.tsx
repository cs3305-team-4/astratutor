import React from 'react';
import { useParams } from 'react-router';
import { Layout } from 'antd';

import { AccountType } from '../api/definitions';

import { Profile } from '../components/Profile';

export function ViewProfile(): React.ReactElement {
  const { uuid }: { uuid: string } = useParams();

  return (
    <Layout style={{ backgroundColor: 'white', padding: '2rem' }}>
      <Profile uuid={uuid} type={AccountType.Tutor} />
    </Layout>
  );
}
