import React, { useState, useEffect } from 'react';
import { api } from '../services/api';

export function GameSearch() {
  const [searchQuery, setSearchQuery] = useState('');
  const [searchResults, setSearchResults] = useState<string[]>([]);
  const [connectedGames, setConnectedGames] = useState<string[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [selectedGame, setSelectedGame] = useState<string | null>(null);
  const [gameUsername, setGameUsername] = useState('');
  const [gameId, setGameId] = useState('');
  const [showConnectForm, setShowConnectForm] = useState(false);

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

  const handleSearch = async () => {
    if (!searchQuery.trim()) {
      setError('Please enter a game name to search');
      return;
    }
    
    setLoading(true);
    setError(null);
    
    try {
      const results = await api.searchGames(searchQuery);
      setSearchResults(results || []);
    } catch (error: any) {
      setError(error.message || 'Error searching for games');
      setSearchResults([]);
    } finally {
      setLoading(false);
    }
  };

  const handleInitiateConnect = (gameName: string) => {
    setSelectedGame(gameName);
    setShowConnectForm(true);
    setGameUsername('');
    setGameId('');
    setError(null);
  };

  const handleConnectGame = async () => {
    if (!selectedGame || (!gameUsername.trim() && !gameId.trim())) {
      setError('Please provide either a Game Username or Game ID');
      return;
    }

    try {
      await api.connectGame(selectedGame, {
        username: gameUsername.trim(),
        gameId: gameId.trim()
      });
      const profile = await api.getProfile();
      setConnectedGames(profile.connectedGames || []);
      setError(null);
      // Reset form
      setShowConnectForm(false);
      setSelectedGame(null);
      setGameUsername('');
      setGameId('');
      setSearchResults([]);
      setSearchQuery('');
    } catch (error: any) {
      setError(error.message || 'Error connecting game');
    }
  };

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
          onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
          className="p-2 border rounded flex-1"
        />
        <button 
          onClick={handleSearch}
          disabled={loading}
          className={`p-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition ${
            loading ? 'opacity-50 cursor-not-allowed' : ''
          }`}
        >
          {loading ? 'Searching...' : 'Search'}
        </button>
      </div>

      {error && (
        <div className="mt-2 text-red-500 text-sm">
          {error}
        </div>
      )}

      {/* Connect Game Form */}
      {showConnectForm && selectedGame && (
        <div className="mt-4 p-4 bg-gray-50 rounded-lg">
          <h3 className="text-lg font-medium mb-3">Connect {selectedGame}</h3>
          <p className="text-sm text-gray-600 mb-3">
            Please provide either a Game Username or Game ID (at least one is required)
          </p>
          <div className="space-y-3">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Game Username (Optional)
              </label>
              <input
                type="text"
                value={gameUsername}
                onChange={(e) => setGameUsername(e.target.value)}
                className="w-full p-2 border rounded"
                placeholder="Enter your game username"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Game ID (Optional)
              </label>
              <input
                type="text"
                value={gameId}
                onChange={(e) => setGameId(e.target.value)}
                className="w-full p-2 border rounded"
                placeholder="Enter your game ID"
              />
            </div>
            <div className="flex gap-2">
              <button
                onClick={handleConnectGame}
                className="flex-1 p-2 bg-blue-600 text-white rounded hover:bg-blue-700"
              >
                Connect
              </button>
              <button
                onClick={() => {
                  setShowConnectForm(false);
                  setSelectedGame(null);
                  setError(null);
                }}
                className="flex-1 p-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300"
              >
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Search Results */}
      {searchResults && searchResults.length > 0 && !showConnectForm && (
        <div className="mt-4">
          <h3 className="text-lg font-medium mb-2">Search Results</h3>
          <div className="space-y-2">
            {searchResults.map((game, index) => (
              <div
                key={index}
                className="flex items-center justify-between p-2 bg-gray-50 rounded"
              >
                <span className="text-gray-800">{game}</span>
                {!connectedGames.includes(game) && (
                  <button
                    onClick={() => handleInitiateConnect(game)}
                    className="text-blue-600 hover:text-blue-800 transition-colors"
                  >
                    Connect
                  </button>
                )}
              </div>
            ))}
          </div>
        </div>
      )}

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
    </div>
  );
}
