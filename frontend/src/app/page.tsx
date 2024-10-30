"use client";
import { useEffect, useState } from 'react';
import { UserFeed } from '../components/UserFeed';
import Link from 'next/link';
import { auth } from '../services/auth';

export default function Home() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
    setIsAuthenticated(auth.isAuthenticated());
  }, []);

  return (
    <main className="flex min-h-screen flex-col items-center p-24">
      <div className="w-full max-w-4xl">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-4xl font-bold">Welcome to Airdate</h1>
          
          {!isAuthenticated && (
            <div className="flex gap-4">
              <Link 
                href="/login"
                className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
              >
                Login
              </Link>
              <Link 
                href="/register"
                className="px-6 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
              >
                Register
              </Link>
            </div>
          )}
        </div>

        <div className="bg-white/5 backdrop-blur-lg rounded-xl p-6 mb-8">
          <h2 className="text-xl font-semibold mb-4">Find Gaming Partners</h2>
          <p className="text-gray-400">
            Connect with fellow gamers, share your favorite games, and find the perfect gaming partners.
            Join our community to start collaborating and making new friends!
          </p>
        </div>

        <UserFeed />
      </div>
    </main>
  );
}
