"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { api } from "../../../services/api";
import { auth } from "../../../services/auth";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { CircleX } from "lucide-react";

export default function Login() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!username.trim() || !password.trim()) {
      setError("Username and password are required");
      return;
    }

    try {
      setError(null);
      setLoading(true);

      const data = await api.login(username, password);

      if (data.token && data.username) {
        auth.login(data.token, data.username);
        router.push("/profile");
      } else {
        throw new Error("Invalid response from server");
      }
    } catch (error) {
      console.error("Login error:", error);
      setError(
        (error as string) || "Login failed. Please check your credentials."
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex bg-purple-400/10 md:bg-transparent text-center  flex-col items-center rounded-xl py-8 justify-center">
      <div className="flex flex-col items-center justify-center gap-2 mb-8">
        <h1 className="text-3xl font-semibold">Login</h1>
        <p className="text-sm ">Welcome back , Good to see you back !</p>
      </div>

      <form
        onSubmit={handleSubmit}
        className="flex flex-col w-[90%] md:w-3/4 lg:w-1/2 gap-4"
      >
        <Input
          type="text"
          placeholder="Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          className="p-4 h-12 w-full border rounded"
          disabled={loading}
        />
        <Input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          className="p-4 h-12 w-full border rounded"
          disabled={loading}
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
        <Button
          type="submit"
          className="p-2 text-white disabled:opacity-50"
          disabled={loading}
        >
          {loading ? "Logging in..." : "Login"}
        </Button>
      </form>
      <div className=" mt-10 flex flex-col items-center">
        Don&apos;t have an account?
        <Button
          onClick={() => router.push("/register")}
          variant={"linkHover2"}
          disabled={loading}
        >
          Register
        </Button>
      </div>
    </div>
  );
}
