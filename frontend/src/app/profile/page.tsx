"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { api } from "../../services/api";
import { ConnectAccounts } from "../../components/ConnectAccounts";
import { GameSearch } from "../../components/GameSearch";
import { ApiError } from "../../types/errors";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import {
  DiscordIcon,
  InstaIcon,
  LockIcon,
  OpenLockICon,
  TwitchIcon,
  YoutubeIcon,
} from "../icons/icon";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { ChevronDownIcon } from "lucide-react";

interface UserProfile {
  username: string;
  twitchUsername?: string;
  discordUsername?: string;
  instagramHandle?: string;
  youtubeChannel?: string;
  connectedGames: string[];
  isPrivate: boolean;
}

export default function Profile() {
  const [user, setUser] = useState<UserProfile | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isPrivate, setIsPrivate] = useState(false);
  const router = useRouter();

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const token = localStorage.getItem("token");
        if (!token) {
          router.push("/login");
          return;
        }

        const data = await api.getProfile();
        setUser(data);
        setIsPrivate(data.isPrivate || false);
        setError(null);
      } catch (error) {
        const apiError = error as ApiError;
        console.error("Profile error:", apiError);
        setError(apiError.message);
        if (apiError.message.includes("authentication")) {
          router.push("/login");
        }
      }
    };

    fetchProfile();
  }, [router]);

  const handlePrivacyToggle = async () => {
    try {
      await api.updatePrivacySettings(!isPrivate);
      setIsPrivate(!isPrivate);
    } catch (error) {
      const apiError = error as ApiError;
      console.error("Privacy update error:", apiError);
      setError(apiError.message);
    }
  };

  if (error) {
    return (
      <div className="flex min-h-screen flex-col items-center justify-center">
        <div className="text-red-500 mb-4">{error}</div>
        <button
          onClick={() => router.push("/login")}
          className="p-2 bg-blue-500 text-white rounded"
        >
          Back to Login
        </button>
      </div>
    );
  }

  if (!user) {
    return (
      <div className="flex min-h-screen flex-col items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  return (
    <section className="">
      <div className="my-4">
        <h1 className="text-2xl font-semibold">
          User <span className="font-serifItalic font-normal">Profile</span>
        </h1>
      </div>
      <div className="flex flex-col rounded-2xl border bg-accent/30">
        <div className="h-32  flex items-center gap-4 justify-start px-6">
          <Avatar className="w-20 rounded-2xl h-20">
            <AvatarImage src="https://github.com/shadcn.png" alt="@shadcn" />
            <AvatarFallback>CN</AvatarFallback>
          </Avatar>
          <div className="flex flex-col">
            <div className="text-3xl font-semibold">@{user?.username}</div>
            <div>{user.connectedGames.length} connected games</div>
          </div>
        </div>
        <div className="border-t px-6 py-3">
          <Badge
            onClick={handlePrivacyToggle}
            className={`${
              isPrivate ? "bg-yellow-800 hover:bg-yellow-800" : ""
            } h-6 cursor-pointer py-3`}
          >
            {isPrivate ? (
              <div className="flex items-center gap-1">
                <LockIcon /> Private Account
              </div>
            ) : (
              <div className="flex items-center gap-1">
                <OpenLockICon /> Public Account
              </div>
            )}
          </Badge>
        </div>
      </div>

      <div className="flex min-h-screen flex-col">
        <ConnectAccounts />
        <GameSearch />

        <button
          onClick={() => {
            localStorage.removeItem("token");
            localStorage.removeItem("username");
            router.push("/login");
          }}
          className="mt-8 p-2 bg-red-500 text-white rounded"
        >
          Logout
        </button>
      </div>
    </section>
  );
}
