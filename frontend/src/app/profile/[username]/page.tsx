'use client';

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { api } from '../../../services/api';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { LockIcon } from '@/app/icons/icon';
import {
  DiscordIcon,
  InstaIcon,
  TwitchIcon,
  YoutubeIcon,
} from '@/app/icons/icon';

interface GameConnection {
  name: string;
  username?: string;
  gameId?: string;
}

interface UserProfile {
  username: string;
  twitchUsername?: string;
  discordUsername?: string;
  instagramHandle?: string;
  youtubeChannel?: string;
  connectedGames: GameConnection[];
  isPrivate: boolean;
  followersCount: number;
  followingCount: number;
  isFollowing: boolean;
}

export default function UserProfileView() {
  const params = useParams();
  const router = useRouter();
  const [profile, setProfile] = useState<UserProfile>({
    username: '',
    connectedGames: [],
    isPrivate: false,
    followersCount: 0,
    followingCount: 0,
    isFollowing: false
  });
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const data = await api.getUserProfile(params.username as string);
        setProfile(data);
      } catch (error) {
        setError('Failed to load profile');
        console.error('Profile error:', error);
      }
    };

    if (params.username) {
      fetchProfile();
    }
  }, [params.username]);

  if (error) {
    return (
      <div className='flex min-h-screen flex-col items-center justify-center'>
        <div className='text-red-500 mb-4'>{error}</div>
        <Button onClick={() => router.back()}>Go Back</Button>
      </div>
    );
  }

  if (!profile) {
    return (
      <div className='flex min-h-screen flex-col items-center justify-center'>
        <div className='animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500'></div>
      </div>
    );
  }

  if (profile.isPrivate) {
    return (
      <div className='flex min-h-screen flex-col items-center justify-center'>
        <div className='flex flex-col items-center gap-4'>
          <LockIcon className='w-16 h-16' />
          <h1 className='text-2xl font-bold'>This Account is Private</h1>
          <p className='text-gray-500'>Follow this account to see their profile</p>
          <Button onClick={() => router.back()}>Go Back</Button>
        </div>
      </div>
    );
  }

  return (
    <div className='container mx-auto px-4 py-8'>
      <div className='flex flex-col rounded-2xl border bg-accent/30'>
        <div className='h-32 flex items-center gap-4 justify-start px-6'>
          <Avatar className='w-20 h-20 rounded-2xl'>
            <AvatarImage src='https://github.com/shadcn.png' alt='@shadcn' />
            <AvatarFallback>{profile?.username?.[0]}</AvatarFallback>
          </Avatar>
          <div className='flex flex-col'>
            <div className='text-3xl font-semibold'>@{profile?.username}</div>
            <div className='flex gap-4 text-sm text-gray-500'>
              <span>{profile?.followersCount || 0} followers</span>
              <span>{profile?.followingCount || 0} following</span>
              <span>{profile?.connectedGames?.length || 0} games</span>
            </div>
          </div>
        </div>

        {/* Connected Games */}
        <div className='border-t px-6 py-4'>
          <h2 className='text-xl font-semibold mb-4'>Connected Games</h2>
          {profile?.connectedGames && profile.connectedGames.length > 0 ? (
            <div className='grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4'>
              {profile.connectedGames.map((game, index) => (
                <div key={index} className='flex flex-col gap-1'>
                  <Badge variant='secondary' className='p-2 mb-1'>
                    {game.name}
                  </Badge>
                  {(game.username || game.gameId) && (
                    <div className='text-sm text-gray-500 px-2'>
                      {game.username && (
                        <div>Username: {game.username}</div>
                      )}
                      {game.gameId && (
                        <div>Game ID: {game.gameId}</div>
                      )}
                    </div>
                  )}
                </div>
              ))}
            </div>
          ) : (
            <p className='text-gray-500'>No games connected yet</p>
          )}
        </div>

        {/* Social Links */}
        <div className='border-t px-6 py-4'>
          <h2 className='text-xl font-semibold mb-4'>Connected Accounts</h2>
          <div className='grid grid-cols-1 md:grid-cols-2 gap-4'>
            {profile?.twitchUsername && (
              <div className='flex items-center gap-2'>
                <TwitchIcon className='w-5 h-5' />
                <a
                  href={`https://twitch.tv/${profile.twitchUsername}`}
                  target='_blank'
                  rel='noopener noreferrer'
                  className='hover:underline'
                >
                  {profile.twitchUsername}
                </a>
              </div>
            )}
            {profile?.discordUsername && (
              <div className='flex items-center gap-2'>
                <DiscordIcon className='w-5 h-5' />
                <span>{profile.discordUsername}</span>
              </div>
            )}
            {profile?.instagramHandle && (
              <div className='flex items-center gap-2'>
                <InstaIcon className='w-5 h-5' />
                <a
                  href={`https://instagram.com/${profile.instagramHandle}`}
                  target='_blank'
                  rel='noopener noreferrer'
                  className='hover:underline'
                >
                  {profile.instagramHandle}
                </a>
              </div>
            )}
            {profile?.youtubeChannel && (
              <div className='flex items-center gap-2'>
                <YoutubeIcon className='w-5 h-5' />
                <a
                  href={profile.youtubeChannel}
                  target='_blank'
                  rel='noopener noreferrer'
                  className='hover:underline'
                >
                  {profile.youtubeChannel.split('/').pop()}
                </a>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
} 