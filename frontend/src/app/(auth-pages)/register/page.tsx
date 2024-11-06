"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { api } from "../../../services/api";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { CircleX } from "lucide-react";

export default function Register() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      setError(null);
      await api.register(username, password);
      router.push("/login");
    } catch (error: any) {
      console.error("Registration error:", error);
      setError(error.message || "Registration failed");
    }
  };

  return (
    <div className="flex w-full flex-col bg-purple-400/10 md:bg-transparent text-center items-center rounded-xl py-8 px-4 justify-center">
      <div className="flex flex-col items-center justify-center gap-2 mb-8">
        <h1 className="text-3xl font-semibold">Register</h1>
        <p className="text-sm">Start your Journey with airDate now.</p>
      </div>
      <form
        onSubmit={handleSubmit}
        className="flex flex-col w-full md:w-3/4 lg:w-1/2 gap-4"
      >
        <Input
          type="text"
          placeholder="Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          className="p-4 h-12 w-full border rounded"
        />
        <Input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          className="p-4 h-12 w-full border rounded"
        />
        {error && (
          <Badge
            variant={"destructive"}
            className="h-8 hover:bg-destructive/40"
          >
            <CircleX className="mr-2 w-4 h-4" />
            {error}
          </Badge>
        )}
        <Button type="submit" className="p-2 text-white">
          Register
        </Button>
      </form>

      <div className=" mt-10 flex flex-col items-center">
        Already have an account?
        <Button onClick={() => router.push("/login")} variant={"linkHover2"}>
          Login
        </Button>
      </div>
    </div>
  );
}
