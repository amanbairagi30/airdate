import React, { useState, useEffect } from 'react';
import { api } from '../services/api';

export function GameSearch() {
  const [searchQuery, setSearchQuery] = useState('');
  const [connectedGames, setConnectedGames] = useState<string[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    // Fetch user's connected games on component mount
    const fetchConnectedGames = async () => {
      try {
        const profile = await api.getProfile();
        setConnectedGames(profile.connectedGames || []);
      } catch (error) {
        console.error('Error fetching connected games:', error);
      }
    };
    fetchConnectedGames();
  }, []);

  const handleDisconnectGame = async (gameName: string) => {
    try {
      await api.disconnectGame(gameName);
      // Refresh the connected games list
      const profile = await api.getProfile();
      setConnectedGames(profile.connectedGames || []);
    } catch (error: any) {
      setError(error.message);
    }
  };

  return (
    <div className="w-full max-w-md">
      <h2 className="text-xl font-semibold mb-4">Search Games</h2>
      <div className="flex gap-2 mb-4">
        <input
          type="text"
          placeholder="Search popular games..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="p-2 border rounded flex-1"
        />
        <button className="p-2 bg-blue-600 text-white rounded">
          Search
        </button>
      </div>

      {/* Connected Games Section */}
      {connectedGames.length > 0 && (
        <div className="mt-4">
          <h3 className="text-lg font-medium mb-2">Connected Games</h3>
          <div className="space-y-2">
            {connectedGames.map((game, index) => (
              <div
                key={index}
                className="flex items-center justify-between p-2 bg-blue-50 rounded"
              >
                <span className="text-blue-800">{game}</span>
                <button
                  onClick={() => handleDisconnectGame(game)}
                  className="text-red-500 hover:text-red-700 transition-colors"
                  title="Disconnect game"
                >
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
            ))}
          </div>
        </div>
      )}

      {error && (
        <div className="mt-2 text-red-500 text-sm">
          {error}
        </div>
      )}
    </div>
  );
}
