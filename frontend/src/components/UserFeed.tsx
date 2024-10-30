import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import { api } from '../services/api';

interface UserProfile {
  id: number;
  username: string;
  twitchUsername?: string;
  discordUsername?: string;
  instagramHandle?: string;
  youtubeChannel?: string;
  favoriteGames?: string;
  connectedGames: string[];
}

export function UserFeed() {
  const [users, setUsers] = useState<UserProfile[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const response = await api.getAllUsers();
        setUsers(response || []);
      } catch (error: unknown) {
        const err = error as Error;
        setError(err.message || 'Failed to load users');
        setUsers([]);
      } finally {
        setLoading(false);
      }
    };

    fetchUsers();
  }, []);

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-[200px]">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-red-500 text-center p-4">
        {error}
      </div>
    );
  }

  if (!users || users.length === 0) {
    return (
      <div className="text-center p-4 text-gray-500">
        No users found
      </div>
    );
  }

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      {users.map((user: UserProfile) => (
        <div key={user.id} className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
          <div className="flex items-center space-x-4 mb-4">
            <div className="w-12 h-12 bg-gray-200 rounded-full flex items-center justify-center">
              <span className="text-xl font-bold">{user.username[0].toUpperCase()}</span>
            </div>
            <div>
              <Link href={`/profile/${user.username}`} className="font-bold hover:underline">
                {user.username}
              </Link>
              <div className="text-sm text-gray-500">
                {user.connectedGames?.length || 0} games connected
              </div>
            </div>
          </div>

          <div className="space-y-2">
            {user.twitchUsername && (
              <div className="flex items-center space-x-2 text-purple-600">
                <svg className="w-5 h-5" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M11.571 4.714h1.715v5.143H11.57zm4.715 0H18v5.143h-1.714zM6 0L1.714 4.286v15.428h5.143V24l4.286-4.286h3.428L22.286 12V0zm14.571 11.143l-3.428 3.428h-3.429l-3 3v-3H6.857V1.714h13.714z"/>
                </svg>
                <span>{user.twitchUsername}</span>
              </div>
            )}
            
            {user.discordUsername && (
              <div className="flex items-center space-x-2 text-indigo-600">
                <svg className="w-5 h-5" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M20.317 4.37a19.791 19.791 0 00-4.885-1.515.074.074 0 00-.079.037c-.21.375-.444.864-.608 1.25a18.27 18.27 0 00-5.487 0 12.64 12.64 0 00-.617-1.25.077.077 0 00-.079-.037A19.736 19.736 0 003.677 4.37a.07.07 0 00-.032.027C.533 9.046-.32 13.58.099 18.057a.082.082 0 00.031.057 19.9 19.9 0 005.993 3.03.078.078 0 00.084-.028c.462-.63.874-1.295 1.226-1.994a.076.076 0 00-.041-.106 13.107 13.107 0 01-1.872-.892.077.077 0 01-.008-.128 10.2 10.2 0 00.372-.292.074.074 0 01.077-.01c3.928 1.793 8.18 1.793 12.062 0a.074.074 0 01.078.01c.12.098.246.198.373.292a.077.077 0 01-.006.127 12.299 12.299 0 01-1.873.892.077.077 0 00-.041.107c.36.698.772 1.362 1.225 1.993a.076.076 0 00.084.028 19.839 19.839 0 006.002-3.03.077.077 0 00.032-.054c.5-5.177-.838-9.674-3.549-13.66a.061.061 0 00-.031-.03zM8.02 15.33c-1.183 0-2.157-1.085-2.157-2.419 0-1.333.956-2.419 2.157-2.419 1.21 0 2.176 1.096 2.157 2.42 0 1.333-.956 2.418-2.157 2.418zm7.975 0c-1.183 0-2.157-1.085-2.157-2.419 0-1.333.955-2.419 2.157-2.419 1.21 0 2.176 1.096 2.157 2.42 0 1.333-.946 2.418-2.157 2.418z"/>
                </svg>
                <span>{user.discordUsername}</span>
              </div>
            )}

            {user.favoriteGames && (
              <div className="mt-3">
                <div className="text-sm font-semibold mb-1">Favorite Games:</div>
                <div className="flex flex-wrap gap-2">
                  {user.favoriteGames.split(',').map((game: string, index: number) => (
                    <span 
                      key={index}
                      className="px-2 py-1 bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-100 rounded text-sm"
                    >
                      {game.trim()}
                    </span>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
      ))}
    </div>
  );
} 