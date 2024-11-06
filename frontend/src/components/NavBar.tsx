"use client";
import React, { useState, useEffect } from "react";
import Link from "next/link";
import { auth } from "../services/auth";
import { useRouter } from "next/navigation";
import { Button } from "./ui/button";
import { SelectTheme } from "./theme-toggle";

export function NavBar() {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean | null>(null);
  const router = useRouter();

  useEffect(() => {
    setIsAuthenticated(auth.isAuthenticated());
  }, []);

  const handleLogout = () => {
    auth.logout();
    setIsAuthenticated(false);
    router.push("/login");
  };

  if (isAuthenticated === null) {
    return null;
  }

  return (
    <>
      <header>
        <nav className="flex items-center justify-between gap-2 p-3 rounded-xl ">
          <div className="flex flex-1 font-medium items-center gap-4">
            <Button variant={"linkHover2"} className="cursor-pointer">
              Products
            </Button>
            <Button variant={"linkHover2"} className="cursor-pointer">
              Services
            </Button>
            <Button variant={"linkHover2"} className="cursor-pointer">
              Themes
            </Button>
          </div>
          <div className="flex flex-1 text-2xl justify-center items-center font-secondary font-extrabold">
            <span className="">Pixel</span>
            <span className="text-purple-400">&</span>
            <span className="">Chill</span>
          </div>
          <div className="flex flex-1 items-center justify-end gap-2">
            <SelectTheme />
            <Button className="font-extrabold">Get Started</Button>
          </div>
        </nav>
      </header>
      {/* <nav className="bg-white shadow-lg">
        <div className="max-w-7xl mx-auto px-4">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <Link href="/" className="text-xl font-bold text-gray-800">
                Airdate
              </Link>
            </div>

            <div className="flex items-center space-x-4">
              {isAuthenticated ? (
                <>
                  <Link
                    href="/profile"
                    className="text-gray-600 hover:text-gray-900"
                  >
                    Profile
                  </Link>
                  <button
                    onClick={handleLogout}
                    className="text-red-600 hover:text-red-800"
                  >
                    Logout
                  </button>
                </>
              ) : (
                <>
                  <Link
                    href="/login"
                    className="text-gray-600 hover:text-gray-900"
                  >
                    Login
                  </Link>
                  <Link
                    href="/register"
                    className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
                  >
                    Register
                  </Link>
                </>
              )}
            </div>
          </div>
        </div>
      </nav> */}
    </>
  );
}
