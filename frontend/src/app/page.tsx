"use client";

import Landing from "@/components/landing";
import { useEffect, useState } from "react";
import { auth } from "../services/auth";

export default function Home() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  console.log(isAuthenticated);

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
