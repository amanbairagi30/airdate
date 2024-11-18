export interface User {
    username: string;
    isPrivate: boolean;
    isFollowing: boolean;
    followersCount: number;
    followingCount: number;
    twitchUsername?: string;
    discordUsername?: string;
    instagramHandle?: string;
    youtubeChannel?: string;
    connectedGames: string[];
} 