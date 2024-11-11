import { NavBar } from "@/components/NavBar";
import { UserFeed } from "@/components/UserFeed";
import React from "react";

export default function AllUsersFeed() {
  return (
    <section className="max-w-7xl mx-auto">
      <div className="p-4">
        <NavBar />
        <div className="my-8">
          <div className="text-3xl font-semibold">
            All <span className="font-serifItalic">Users</span>
          </div>
          <p className="text-sm mt-1">
            Currently , you are viewing the feed of all users in the airdate.
          </p>
        </div>
        <UserFeed position="feed" />
      </div>
    </section>
  );
}
