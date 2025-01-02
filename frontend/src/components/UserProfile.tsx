// import React, { useState, useEffect } from "react";
// import { Box, VStack, Text, Avatar } from "@chakra-ui/react";
// import FollowButton from "./FollowButton";
// import { User } from "../types";

// interface UserProfileProps {
//   username: string;
// }

// const UserProfile: React.FC<UserProfileProps> = ({ username }) => {
//   const [profile, setProfile] = useState<User | null>(null);
//   const [isLoading, setIsLoading] = useState(true);

//   useEffect(() => {
//     const fetchProfile = async () => {
//       try {
//         const response = await fetch(`/api/profile/${username}`);
//         const data = await response.json();
//         setProfile(data);
//       } catch (error) {
//         console.error("Error fetching profile:", error);
//       } finally {
//         setIsLoading(false);
//       }
//     };

//     fetchProfile();
//   }, [username]);

//   if (isLoading || !profile) {
//     return <div>Loading...</div>;
//   }

//   return (
//     <Box p={4}>
//       <VStack spacing={4} align="start">
//         <Avatar size="xl" name={profile.username} />
//         <Text fontSize="2xl" fontWeight="bold">
//           {profile.username}
//         </Text>

//         {/* Follow Button */}
//         <FollowButton
//           targetUsername={profile.username}
//           isPrivate={profile.isPrivate}
//           initialFollowState={profile.isFollowing}
//           initialFollowersCount={profile.followersCount}
//           onFollowStateChange={(isFollowing) => {
//             setProfile((prev) =>
//               prev
//                 ? {
//                     ...prev,
//                     isFollowing,
//                     followersCount: isFollowing
//                       ? prev.followersCount + 1
//                       : prev.followersCount - 1,
//                   }
//                 : null
//             );
//           }}
//         />

//         {/* Profile Stats */}
//         <Box>
//           <Text>Followers: {profile.followersCount}</Text>
//           <Text>Following: {profile.followingCount}</Text>
//         </Box>

//         {/* Connected Accounts */}
//         {profile.twitchUsername && (
//           <Text>Twitch: {profile.twitchUsername}</Text>
//         )}
//         {profile.discordUsername && (
//           <Text>Discord: {profile.discordUsername}</Text>
//         )}
//         {profile.instagramHandle && (
//           <Text>Instagram: {profile.instagramHandle}</Text>
//         )}
//         {profile.youtubeChannel && (
//           <Text>YouTube: {profile.youtubeChannel}</Text>
//         )}

//         {/* Connected Games */}
//         {profile.connectedGames && profile.connectedGames.length > 0 && (
//           <Box>
//             <Text fontWeight="bold">Connected Games:</Text>
//             {profile.connectedGames.map((game, index) => (
//               <Text key={index}>{game}</Text>
//             ))}
//           </Box>
//         )}
//       </VStack>
//     </Box>
//   );
// };

// export default UserProfile;
