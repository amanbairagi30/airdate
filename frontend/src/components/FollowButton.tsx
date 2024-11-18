'use client';

import React, { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { api } from '@/services/api';
import { Loader2 } from 'lucide-react';

interface FollowButtonProps {
  targetUsername: string;
  isPrivate: boolean;
  initialFollowState: boolean;
  initialFollowersCount: number;
  currentUsername: string;
  onFollowStateChange?: (isFollowing: boolean) => void;
}

type FollowState = 'not_following' | 'following' | 'requested' | 'self';

const FollowButton: React.FC<FollowButtonProps> = ({
  targetUsername,
  isPrivate,
  initialFollowState,
  initialFollowersCount,
  currentUsername,
  onFollowStateChange
}) => {
  const [followState, setFollowState] = useState<FollowState>('not_following');
  const [followersCount, setFollowersCount] = useState(initialFollowersCount);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    setFollowState(initialFollowState ? 'following' : 'not_following');
    setIsLoading(false);
  }, [initialFollowState]);

  const handleFollow = async () => {
    if (isLoading || followState === 'self' || targetUsername === currentUsername) return;
    
    setIsLoading(true);
    try {
      console.log("Attempting to follow:", targetUsername);
      const response = await api.followUser(targetUsername);
      setFollowState(response.followState);
      
      if (response.followState === 'following') {
        setFollowersCount(prev => prev + 1);
        onFollowStateChange?.(true);
      }
    } catch (error) {
      console.error('Failed to follow user:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleUnfollow = async () => {
    if (isLoading || followState === 'self') return;
    
    setIsLoading(true);
    try {
      const response = await api.unfollowUser(targetUsername);
      setFollowState('not_following');
      setFollowersCount(prev => prev - 1);
      onFollowStateChange?.(false);
    } catch (error) {
      console.error('Failed to unfollow user:', error);
    } finally {
      setIsLoading(false);
    }
  };

  if (followState === 'self') {
    return null;
  }

  const getButtonProps = () => {
    const baseProps = {
      disabled: isLoading,
      size: "sm" as const,
      className: "rounded-full min-w-[100px]"
    };

    switch (followState) {
      case 'following':
        return {
          ...baseProps,
          onClick: handleUnfollow,
          variant: 'destructive' as const,
          children: isLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : 'Unfollow'
        };
      case 'requested':
        return {
          ...baseProps,
          onClick: handleUnfollow,
          variant: 'secondary' as const,
          children: isLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : 'Requested'
        };
      default:
        return {
          ...baseProps,
          onClick: handleFollow,
          variant: 'default' as const,
          children: isLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : 'Follow'
        };
    }
  };

  return <Button {...getButtonProps()} />;
};

export default FollowButton; 