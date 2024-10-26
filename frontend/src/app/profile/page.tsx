'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { api } from '../../services/api';
import { ConnectAccounts } from '../../components/ConnectAccounts';
import { GameSearch } from '../../components/GameSearch';

export default function Profile() {
  const [user, setUser] = useState(null);
  const router = useRouter();

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const data = await api.getProfile();
        setUser(data);
      } catch (error) {
        console.error('Error:', error);
        router.push('/login');
      }
    };

    fetchProfile();
  }, [router]);

  if (!user) return <div>Loading...</div>;

  return (
    <div className="flex min-h-screen flex-col items-center justify-center p-24">
      <h1 className="text-2xl mb-4">Profile</h1>
      <p>Username: {user.username}</p>
      <p>User ID: {user.id}</p>
      
      <ConnectAccounts />
      <GameSearch />
    </div>
  );
}
