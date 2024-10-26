import React, { useState } from 'react';
import { api } from '../services/api';

export function ConnectAccounts() {
  const [twitchUsername, setTwitchUsername] = useState('');
  const [discordUsername, setDiscordUsername] = useState('');

  const handleConnectTwitch = async () => {
    try {
      await api.connectTwitch(twitchUsername);
      alert('Twitch account connected successfully');
    } catch (error) {
      alert('Error connecting Twitch account');
    }
  };

  const handleConnectDiscord = async () => {
    try {
      await api.connectDiscord(discordUsername);
      alert('Discord account connected successfully');
    } catch (error) {
      alert('Error connecting Discord account');
    }
  };

  return (
    <div className="mt-8">
      <h2 className="text-xl mb-4">Connect Accounts</h2>
      <div className="flex flex-col gap-4">
        <div>
          <input
            type="text"
            placeholder="Twitch Username"
            value={twitchUsername}
            onChange={(e) => setTwitchUsername(e.target.value)}
            className="p-2 border rounded mr-2"
          />
          <button onClick={handleConnectTwitch} className="p-2 bg-purple-500 text-white rounded">
            Connect Twitch
          </button>
        </div>
        <div>
          <input
            type="text"
            placeholder="Discord Username"
            value={discordUsername}
            onChange={(e) => setDiscordUsername(e.target.value)}
            className="p-2 border rounded mr-2"
          />
          <button onClick={handleConnectDiscord} className="p-2 bg-blue-500 text-white rounded">
            Connect Discord
          </button>
        </div>
      </div>
    </div>
  );
}
