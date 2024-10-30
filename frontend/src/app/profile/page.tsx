'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { api } from '../../services/api';
import { ConnectAccounts } from '../../components/ConnectAccounts';
import { GameSearch } from '../../components/GameSearch';

export default function Profile() {
  const [user, setUser] = useState<any>(null);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const token = localStorage.getItem('token');
        if (!token) {
          router.push('/login');
          return;
        }

        const data = await api.getProfile();
        setUser(data);
        setError(null);
      } catch (error: any) {
        console.error('Profile error:', error);
        setError(error.message);
        if (error.message.includes('authentication')) {
          router.push('/login');
        }
      }
    };

    fetchProfile();
  }, [router]);

  if (error) {
    return (
      <div className="flex min-h-screen flex-col items-center justify-center">
        <div className="text-red-500 mb-4">{error}</div>
        <button 
          onClick={() => router.push('/login')} 
          className="p-2 bg-blue-500 text-white rounded"
        >
          Back to Login
        </button>
      </div>
    );
  }

  if (!user) {
    return (
      <div className="flex min-h-screen flex-col items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  return (
    <div className="flex min-h-screen flex-col items-center justify-center p-24">
      <h1 className="text-2xl mb-4">Profile</h1>
      <div className="mb-8">
        <p className="mb-2">Username: {user.username}</p>
        <p className="mb-2">Twitch: {user.twitchUsername || 'Not connected'}</p>
        <p className="mb-2">Discord: {user.discordUsername || 'Not connected'}</p>
      </div>
      
      <ConnectAccounts />
      <GameSearch />

      <button 
        onClick={() => {
          localStorage.removeItem('token');
          localStorage.removeItem('username');
          router.push('/login');
        }}
        className="mt-8 p-2 bg-red-500 text-white rounded"
      >
        Logout
      </button>
    </div>
  );
}
