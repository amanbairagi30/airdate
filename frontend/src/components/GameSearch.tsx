import React, { useState } from 'react';
import { api } from '../services/api';
export function GameSearch() {
  const [query, setQuery] = useState('');
  const [games, setGames] = useState([]);
  const handleSearch = async () => {
    try {
      const results = await api.searchGames(query);
      setGames(results);
    } catch (error) {
      alert('Error searching games');
    }
  };

  return (
    <div className="mt-8">
      <h2 className="text-xl mb-4">Search Games</h2>
      <div className="flex gap-2 mb-4">
        <input
          type="text"
          placeholder="Search games..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          className="p-2 border rounded flex-grow"
        />
        <button onClick={handleSearch} className="p-2 bg-green-500 text-white rounded">
          Search
        </button>
      </div>
      <ul className="list-disc pl-5">
        {games.map((game: any) => (
          <li key={game.id}>{game.name}</li>
        ))}
      </ul>
    </div>
  );
}
