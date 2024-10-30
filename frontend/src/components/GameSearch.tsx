import React, { useState } from 'react';
import { api } from '../services/api';

export function GameSearch() {
  const [query, setQuery] = useState('');
  const [games, setGames] = useState<string[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [selectedGame, setSelectedGame] = useState<string | null>(null);
  const [gameId, setGameId] = useState('');

  const handleSearch = async () => {
    if (!query.trim()) {
      setError('Please enter a search term');
      return;
    }

    try {
      setLoading(true);
      setError(null);
      const results = await api.searchGames(query);
      setGames(results);
      if (results.length === 0) {
        setError('No games found');
      }
    } catch (err: any) {
      console.error('Search error:', err);
      setError(err.message || 'Failed to search games');
      setGames([]);
    } finally {
      setLoading(false);
    }
  };

  const handleGameSelect = (game: string) => {
    setSelectedGame(game);
    setGameId('');
    setError(null);
    setSuccess(null);
  };

  const handleConnectGame = async () => {
    if (!selectedGame || !gameId.trim()) {
      setError('Please enter your game ID');
      return;
    }

    try {
      await api.connectGame(selectedGame, gameId);
      setSuccess(`Successfully connected to ${selectedGame}`);
      setSelectedGame(null);
      setGameId('');
    } catch (err: any) {
      setError(err.message || 'Failed to connect game');
    }
  };

  return (
    <div className="w-full max-w-md mt-8">
      <h2 className="text-xl font-bold mb-4">Search Games</h2>
      
      <div className="flex gap-2">
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
          placeholder="Search popular games..."
          className="p-2 border rounded flex-1 focus:ring-2 focus:ring-blue-300"
        />
        <button
          onClick={handleSearch}
          disabled={loading}
          className="p-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition disabled:opacity-50"
        >
          {loading ? 'Searching...' : 'Search'}
        </button>
      </div>

      {error && (
        <div className="mt-2 text-red-500 p-2 bg-red-50 rounded">
          {error}
        </div>
      )}

      {success && (
        <div className="mt-2 text-green-500 p-2 bg-green-50 rounded">
          {success}
        </div>
      )}

      {selectedGame ? (
        <div className="mt-4 p-4 border rounded-lg bg-white">
          <h3 className="font-semibold mb-2">Connect to {selectedGame}</h3>
          <input
            type="text"
            value={gameId}
            onChange={(e) => setGameId(e.target.value)}
            placeholder="Enter your game ID/username"
            className="p-2 border rounded w-full mb-2"
          />
          <div className="flex gap-2">
            <button
              onClick={handleConnectGame}
              className="p-2 bg-green-600 text-white rounded flex-1 hover:bg-green-700"
            >
              Connect
            </button>
            <button
              onClick={() => setSelectedGame(null)}
              className="p-2 bg-gray-200 text-gray-700 rounded flex-1 hover:bg-gray-300"
            >
              Cancel
            </button>
          </div>
        </div>
      ) : (
        games.length > 0 && (
          <div className="mt-4 space-y-2">
            {games.map((game, index) => (
              <div
                key={index}
                onClick={() => handleGameSelect(game)}
                className="p-3 border rounded-lg bg-white hover:bg-gray-50 transition cursor-pointer shadow-sm"
              >
                {game}
              </div>
            ))}
          </div>
        )
      )}
    </div>
  );
}
