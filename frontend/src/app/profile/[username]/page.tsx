'use client';

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { api } from '../../../services/api';
import { ApiError } from '../../../types/errors';
import { auth } from '../../../services/auth';

interface UserProfile {
  username: string;
  twitchUsername?: string;
  discordUsername?: string;
  instagramHandle?: string;
  youtubeChannel?: string;
  connectedGames: string[];
  isPrivate: boolean;
  followersCount: number;
  followingCount: number;
  isFollowing: boolean;
}

export default function UserProfile() {
  const params = useParams();
  const router = useRouter();
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isFollowing, setIsFollowing] = useState(false);

  const isOwnProfile = auth.getUsername() === params.username;

  useEffect(() => {
    const fetchProfile = async () => {
      if (!params.username) return;

      try {
        setIsLoading(true);
        const data = await api.getUserProfile(params.username as string);
        setProfile(data);
        setIsFollowing(data.isFollowing);
        setError(null);
      } catch (error) {
        const apiError = error as ApiError;
        setError(apiError.message || 'Failed to load profile');
        console.error('Profile error:', error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchProfile();
  }, [params.username]);

  const handleFollow = async () => {
    try {
      if (!profile) return;
      
      if (isFollowing) {
        await api.unfollowUser(profile.username);
      } else {
        await api.followUser(profile.username);
      }
      
      setIsFollowing(!isFollowing);
      // Refresh profile to get updated counts
      const updatedProfile = await api.getUserProfile(profile.username);
      setProfile(updatedProfile);
    } catch (error) {
      const apiError = error as ApiError;
      setError(apiError.message || 'Failed to update follow status');
    }
  };

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex min-h-screen flex-col items-center justify-center p-4">
        <div className="bg-red-50 border border-red-200 rounded-lg p-6 text-center">
          <h2 className="text-red-800 text-xl font-semibold mb-2">Error Loading Profile</h2>
          <p className="text-red-600 mb-4">{error}</p>
          <button
            onClick={() => router.push('/')}
            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            Back to Home
          </button>
        </div>
      </div>
    );
  }

  if (!profile) return null;

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <div className="bg-white rounded-lg shadow-lg p-6">
        {/* Profile Header */}
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center space-x-4">
            <div className="w-20 h-20 bg-gray-200 rounded-full flex items-center justify-center">
              <span className="text-3xl font-bold text-gray-600">
                {profile.username[0].toUpperCase()}
              </span>
            </div>
            <div>
              <h1 className="text-2xl font-bold">{profile.username}</h1>
              {profile.isPrivate && !isOwnProfile && !isFollowing && (
                <span className="text-sm text-gray-500">Private Account</span>
              )}
            </div>
          </div>
          {!isOwnProfile && (
            <button
              onClick={handleFollow}
              className={`px-4 py-2 rounded ${
                isFollowing
                  ? 'bg-gray-200 text-gray-800 hover:bg-gray-300'
                  : 'bg-blue-600 text-white hover:bg-blue-700'
              }`}
            >
              {isFollowing ? 'Unfollow' : 'Follow'}
            </button>
          )}
        </div>

        {/* Stats */}
        <div className="flex space-x-8 mb-6 border-b pb-4">
          <div className="text-center">
            <div className="font-bold">{profile.followersCount}</div>
            <div className="text-gray-600">Followers</div>
          </div>
          <div className="text-center">
            <div className="font-bold">{profile.followingCount}</div>
            <div className="text-gray-600">Following</div>
          </div>
          <div className="text-center">
            <div className="font-bold">{profile.connectedGames?.length || 0}</div>
            <div className="text-gray-600">Games</div>
          </div>
        </div>

        {/* Profile Content - Only show if public or own profile or following */}
        {(!profile.isPrivate || isOwnProfile || isFollowing) && (
          <>
            {/* Connected Accounts */}
            <div className="mb-6">
              <h2 className="text-xl font-semibold mb-4">Connected Accounts</h2>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {profile.twitchUsername && (
                  <div className="flex items-center space-x-2">
                    <span className="text-purple-600">Twitch:</span>
                    <a href={`https://twitch.tv/${profile.twitchUsername}`} target="_blank" rel="noopener noreferrer" 
                       className="text-blue-600 hover:underline">
                      {profile.twitchUsername}
                    </a>
                  </div>
                )}
                {profile.discordUsername && (
                  <div className="flex items-center space-x-2">
                    <span className="text-indigo-600">Discord:</span>
                    <span>{profile.discordUsername}</span>
                  </div>
                )}
                {profile.instagramHandle && (
                  <div className="flex items-center space-x-2">
                    <span className="text-pink-600">Instagram:</span>
                    <a href={`https://instagram.com/${profile.instagramHandle}`} target="_blank" rel="noopener noreferrer"
                       className="text-blue-600 hover:underline">
                      {profile.instagramHandle}
                    </a>
                  </div>
                )}
                {profile.youtubeChannel && (
                  <div className="flex items-center space-x-2">
                    <span className="text-red-600">YouTube:</span>
                    <a href={profile.youtubeChannel} target="_blank" rel="noopener noreferrer"
                       className="text-blue-600 hover:underline">
                      {profile.youtubeChannel.split('/').pop()}
                    </a>
                  </div>
                )}
              </div>
            </div>

            {/* Connected Games */}
            {profile.connectedGames && profile.connectedGames.length > 0 && (
              <div>
                <h2 className="text-xl font-semibold mb-4">Connected Games</h2>
                <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
                  {profile.connectedGames.map((game, index) => (
                    <div key={index} className="bg-gray-50 rounded-lg p-3">
                      {game}
                    </div>
                  ))}
                </div>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
} 