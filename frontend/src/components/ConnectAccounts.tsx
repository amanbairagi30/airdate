import React, { useState, useEffect } from "react";
import { api } from "../services/api";
import {
  DiscordIcon,
  InstaIcon,
  TwitchIcon,
  YoutubeIcon,
} from "@/app/icons/icon";
import { ChevronDownIcon } from "lucide-react";
import { Input } from "./ui/input";
import { Button } from "./ui/button";

export function ConnectAccounts() {
  const [twitchUsername, setTwitchUsername] = useState("");
  const [discordUsername, setDiscordUsername] = useState("");
  const [instagramHandle, setInstagramHandle] = useState("");
  const [youtubeChannel, setYoutubeChannel] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [connectedAccounts, setConnectedAccounts] = useState({
    twitch: false,
    discord: false,
    instagram: false,
    youtube: false,
  });

  useEffect(() => {
    // Fetch profile to check connected accounts
    const fetchProfile = async () => {
      try {
        const profile = await api.getProfile();
        setConnectedAccounts({
          twitch: !!profile.twitchUsername,
          discord: !!profile.discordUsername,
          instagram: !!profile.instagramHandle,
          youtube: !!profile.youtubeChannel,
        });
      } catch (error) {
        console.error("Error fetching profile:", error);
      }
    };
    fetchProfile();
  }, []);

  const refreshConnections = async () => {
    try {
      const profile = await api.getProfile();
      setConnectedAccounts({
        twitch: !!profile.twitchUsername,
        discord: !!profile.discordUsername,
        instagram: !!profile.instagramHandle,
        youtube: !!profile.youtubeChannel,
      });
    } catch (error) {
      console.error("Error fetching profile:", error);
    }
  };

  const handleConnect = async (
    type: "twitch" | "discord" | "instagram" | "youtube"
  ) => {
    try {
      setError(null);
      setSuccess(null);

      switch (type) {
        case "twitch":
          await api.connectTwitch(twitchUsername);
          setSuccess("Twitch account connected!");
          break;
        case "discord":
          await api.connectDiscord(discordUsername);
          setSuccess("Discord account connected!");
          break;
        case "instagram":
          await api.connectInstagram(instagramHandle);
          setSuccess("Instagram account connected!");
          break;
        case "youtube":
          await api.connectYoutube(youtubeChannel);
          setSuccess("YouTube channel connected!");
          break;
      }

      await refreshConnections(); // Refresh after connecting
    } catch (error: any) {
      setError(error.message || "Failed to connect account");
    }
  };

  const handleDisconnect = async (
    type: "twitch" | "discord" | "instagram" | "youtube"
  ) => {
    try {
      setError(null);
      setSuccess(null);

      switch (type) {
        case "twitch":
          await api.disconnectTwitch();
          setSuccess("Twitch account disconnected!");
          break;
        case "discord":
          await api.disconnectDiscord();
          setSuccess("Discord account disconnected!");
          break;
        case "instagram":
          await api.disconnectInstagram();
          setSuccess("Instagram account disconnected!");
          break;
        case "youtube":
          await api.disconnectYoutube();
          setSuccess("YouTube channel disconnected!");
          break;
      }

      await refreshConnections(); // Refresh after disconnecting
    } catch (error: any) {
      setError(error.message || "Failed to disconnect account");
    }
  };

  const handleDiscount = (platform: string) => {
    window.open(getDiscountLink(platform), "_blank");
  };

  const getDiscountLink = (platform: string) => {
    switch (platform) {
      case "youtube":
        return "https://www.youtube.com/premium";
      case "instagram":
        return "https://business.instagram.com";
      default:
        return "#";
    }
  };

  return (
    <section className="w-full mb-20">
      <div className="mt-10">
        <h1 className="text-2xl font-semibold">
          Social <span className="font-serifItalic font-normal">Links</span>
        </h1>
      </div>

      {/* links to other platform */}
      <div className="my-6 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div className="min-h-[5rem] max-h-fit border-2 flex flex-col rounded-xl">
          <div className="text-white rounded-t-xl flex items-center bg-gradient-to-t from-indigo-700 via-indigo-600 to-indigo-500 justify-between px-4 py-6">
            <div className="flex items-center gap-2">
              <DiscordIcon className="w-6 h-6" /> / @{discordUsername}
            </div>
            <ChevronDownIcon className="cursor-pointer" />
          </div>
          <div className="p-4 flex-1">
            {/* Discord Section */}
            <div className="">
              {connectedAccounts.discord ? (
                <button
                  onClick={() => handleDisconnect("discord")}
                  className="p-2 bg-red-600 text-white rounded w-full hover:bg-red-700 transition"
                >
                  Disconnect Discord
                </button>
              ) : (
                <>
                  <Input
                    type="text"
                    placeholder="Discord Username"
                    value={discordUsername}
                    onChange={(e) => setDiscordUsername(e.target.value)}
                    className="focus-visible:border-indigo-500 h-12 border rounded w-full"
                  />
                  <Button
                    onClick={() => handleConnect("discord")}
                    variant={"default"}
                    className="mt-4 p-2 from-indigo-600 to-indigo-500 rounded w-full"
                  >
                    Connect Discord
                  </Button>
                </>
              )}
            </div>
          </div>
        </div>
        <div className="min-h-[5rem] max-h-fit border-2 flex flex-col rounded-xl">
          <div className="text-white rounded-t-xl flex items-center bg-gradient-to-t from-purple-800 to-purple-600 justify-between px-4 py-6">
            <div className="flex items-center gap-2">
              <TwitchIcon className="w-6 h-6" /> / @{twitchUsername}
            </div>
            <ChevronDownIcon className="cursor-pointer" />
          </div>
          <div className="p-4 flex-1">
            <div className="">
              {connectedAccounts.twitch ? (
                <button
                  onClick={() => handleDisconnect("twitch")}
                  className="p-2 bg-red-600 text-white rounded w-full hover:bg-red-700 transition"
                >
                  Disconnect Twitch
                </button>
              ) : (
                <>
                  <Input
                    type="text"
                    placeholder="Twitch Username"
                    value={twitchUsername}
                    onChange={(e) => setTwitchUsername(e.target.value)}
                    className="focus-visible:border-indigo-500 h-12 border rounded w-full"
                  />
                  <Button
                    onClick={() => handleConnect("twitch")}
                    variant={"default"}
                    className="mt-4 p-2 from-purple-600 to-purple-500 rounded w-full"
                  >
                    Connect Twitch
                  </Button>
                </>
              )}
            </div>
          </div>
        </div>
        <div className="min-h-[5rem] max-h-fit border-2 flex flex-col rounded-xl">
          <div className="text-white rounded-t-xl flex items-center bg-gradient-to-t from-red-700 via-red-600 to-red-500 justify-between px-4 py-6">
            <div className="flex items-center gap-2">
              <YoutubeIcon className="w-6 h-6" /> / @{youtubeChannel}
            </div>
            <ChevronDownIcon className="cursor-pointer" />
          </div>
          <div className="p-4 flex-1">
            <div className="">
              {connectedAccounts.youtube ? (
                <button
                  onClick={() => handleDisconnect("youtube")}
                  className="p-2 bg-red-600 text-white rounded w-full hover:bg-red-700 transition"
                >
                  Disconnect Youtube
                </button>
              ) : (
                <>
                  <Input
                    type="text"
                    placeholder="Youtube Username"
                    value={youtubeChannel}
                    onChange={(e) => setYoutubeChannel(e.target.value)}
                    className="focus-visible:border-red-500 h-12 border rounded w-full"
                  />
                  <Button
                    onClick={() => handleConnect("youtube")}
                    variant={"default"}
                    className="mt-4 p-2 from-red-600 to-red-500 rounded w-full"
                  >
                    Connect Youtube
                  </Button>
                </>
              )}
            </div>
          </div>
        </div>
        <div className="min-h-[5rem] max-h-fit border-2 flex flex-col rounded-xl">
          <div className="text-white rounded-t-xl flex items-center bg-gradient-to-t from-yellow-700 to-yellow-600 justify-between px-4 py-6">
            <div className="flex items-center gap-2">
              <InstaIcon className="w-6 h-6" /> / @{instagramHandle}
            </div>
            <ChevronDownIcon className="cursor-pointer" />
          </div>
          <div className="p-4 flex-1">
            <div className="">
              {connectedAccounts.instagram ? (
                <button
                  onClick={() => handleDisconnect("instagram")}
                  className="p-2 bg-red-600 text-white rounded w-full hover:bg-red-700 transition"
                >
                  Disconnect Instagram
                </button>
              ) : (
                <>
                  <Input
                    type="text"
                    placeholder="Instagram Username"
                    value={instagramHandle}
                    onChange={(e) => setInstagramHandle(e.target.value)}
                    className="focus-visible:border-yellow-500 h-12 border rounded w-full"
                  />
                  <Button
                    onClick={() => handleConnect("instagram")}
                    variant={"default"}
                    className="mt-4 p-2 from-yellow-700 to-yellow-500 rounded w-full"
                  >
                    Connect Instagram
                  </Button>
                </>
              )}
            </div>
          </div>
        </div>
      </div>
    </section>
    // <div className="w-full max-w-md space-y-4">
    //   <h2 className="text-xl font-bold mb-4">Connect Your Accounts</h2>

    //   {error && <div className="text-red-500 p-2 bg-red-50 rounded">{error}</div>}
    //   {success && <div className="text-green-500 p-2 bg-green-50 rounded">{success}</div>}

    //   <div className="space-y-4">
    //     {/* Twitch Section */}
    //     <div className="p-4 border rounded-lg bg-white shadow-sm">
    //       {connectedAccounts.twitch ? (
    //         <button
    //           onClick={() => handleDisconnect('twitch')}
    //           className="p-2 bg-red-600 text-white rounded w-full hover:bg-red-700 transition"
    //         >
    //           Disconnect Twitch
    //         </button>
    //       ) : (
    //         <>
    //           <input
    //             type="text"
    //             placeholder="Twitch Username"
    //             value={twitchUsername}
    //             onChange={(e) => setTwitchUsername(e.target.value)}
    //             className="p-2 border rounded w-full"
    //           />
    //           <button
    //             onClick={() => handleConnect('twitch')}
    //             className="mt-2 p-2 bg-purple-600 text-white rounded w-full"
    //           >
    //             Connect Twitch
    //           </button>
    //         </>
    //       )}
    //     </div>

    //     {/* Instagram Section */}
    //     <div className="p-4 border rounded-lg bg-white shadow-sm">
    //       {connectedAccounts.instagram ? (
    //         <button
    //           onClick={() => handleDisconnect('instagram')}
    //           className="p-2 bg-red-600 text-white rounded w-full hover:bg-red-700 transition"
    //         >
    //           Disconnect Instagram
    //         </button>
    //       ) : (
    //         <>
    //           <input
    //             type="text"
    //             placeholder="Instagram Handle"
    //             value={instagramHandle}
    //             onChange={(e) => setInstagramHandle(e.target.value)}
    //             className="p-2 border rounded w-full focus:ring-2 focus:ring-pink-300"
    //           />
    //           <div className="flex gap-2">
    //             <button
    //               onClick={() => handleConnect('instagram')}
    //               className="mt-2 p-2 bg-pink-600 text-white rounded flex-1 hover:bg-pink-700 transition"
    //             >
    //               Connect Instagram
    //             </button>
    //             <button
    //               onClick={() => handleDiscount('instagram')}
    //               className="mt-2 p-2 bg-pink-100 text-pink-600 rounded flex-1 hover:bg-pink-200 transition"
    //             >
    //               Get Business Account
    //             </button>
    //           </div>
    //         </>
    //       )}
    //     </div>

    //     {/* YouTube Section */}
    //     <div className="p-4 border rounded-lg bg-white shadow-sm">
    //       {connectedAccounts.youtube ? (
    //         <button
    //           onClick={() => handleDisconnect('youtube')}
    //           className="p-2 bg-red-600 text-white rounded w-full hover:bg-red-700 transition"
    //         >
    //           Disconnect YouTube
    //         </button>
    //       ) : (
    //         <>
    //           <input
    //             type="text"
    //             placeholder="YouTube Channel"
    //             value={youtubeChannel}
    //             onChange={(e) => setYoutubeChannel(e.target.value)}
    //             className="p-2 border rounded w-full focus:ring-2 focus:ring-red-300"
    //           />
    //           <div className="flex gap-2">
    //             <button
    //               onClick={() => handleConnect('youtube')}
    //               className="mt-2 p-2 bg-red-600 text-white rounded flex-1 hover:bg-red-700 transition"
    //             >
    //               Connect YouTube
    //             </button>
    //             <button
    //               onClick={() => handleDiscount('youtube')}
    //               className="mt-2 p-2 bg-red-100 text-red-600 rounded flex-1 hover:bg-red-200 transition"
    //             >
    //               Get Premium
    //             </button>
    //           </div>
    //         </>
    //       )}
    //     </div>
    //   </div>
    // </div>
  );
}
