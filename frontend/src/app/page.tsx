"use client";

import Landing from "@/components/landing";
import { SelectTheme } from "@/components/theme-toggle";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import { useEffect, useState } from "react";
import { auth } from "../services/auth";

export default function Home() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
    setIsAuthenticated(auth.isAuthenticated());
  }, []);

  return (
    <>
      <main>
        <Landing />
      </main>
    </>
  );
}
